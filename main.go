package main

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
	"github.com/gi4nks/quant"
)

var parrot = quant.NewParrot("gadget")

var settings = Settings{}
var repository = Repository{}

func initDB() {
	repository.InitDB()
	repository.InitSchema()
}

func closeDB() {
	repository.CloseDB()
}

func readSettings() {
	settings.LoadSettings()

	if settings.DebugMode() {
		parrot = quant.NewVerboseParrot("gadget")
	}

	parrot.Debug("Parrot is set to talk so much!")
}

func main() {
	readSettings()
	initDB()

	// -------------------
	app := cli.NewApp()
	app.Name = "gadget"
	app.Usage = "the inspector gadget will be used to inspect docker images"
	app.Version = "0.0.1"
	app.Copyright = "gi4nks - 2016"

	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"bu"},
			Usage:   "build",
			Action:  CmdBuild,
		},
		{
			Name:    "list",
			Aliases: []string{"li"},
			Usage:   "list all the images",
			Action:  CmdList,
		},
		{
			Name:    "revive",
			Aliases: []string{"re"},
			Usage:   "revive gadget",
			Action:  CmdRevive,
		},
		{
			Name:    "labels",
			Aliases: []string{"la"},
			Usage:   "show me the labels of images",
			Subcommands: []cli.Command{
				{
					Name:   "id",
					Usage:  "Get the labels of specific image",
					Action: CmdLabelsById,
				},
				{
					Name:   "tag",
					Usage:  "Get the labels of specific tag",
					Action: CmdLabelsByTag,
				},
				{
					Name:   "all",
					Usage:  "Get all the labels",
					Action: CmdLabels,
				},
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"se"},
			Usage:   "serving gadget for rest apis",
			Action:  CmdServe,
		},
	}

	app.Run(os.Args)
	defer closeDB()
}

// List of functions
func CmdServe(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		parrot.Info("==> Serving gadget for REST Apis.")
		serve()
	})
}

func CmdRevive(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		parrot.Info("==> Reviving gadget will reinitialize all the datas.")

		repository.BackupSchema()
		repository.InitSchema()
	})
}

func CmdLabels(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		var images = repository.GetAll()

		for _, img := range images {
			parrot.Info("[" + img.ID + "] - " + strings.Join(img.RepoTags, ", ") + " [" + strconv.Itoa(len(img.Labels)) + "]")
		}
	})
}

func CmdLabelsById(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		var img = repository.FindById(id)

		if len(img.Labels) == 0 {
			parrot.Info("[" + img.ID + "] - No labels defined")
		} else {
			parrot.Info("[" + img.ID + "] - " + asJson(img.Labels))
		}
	})
}

func CmdLabelsByTag(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		var img = repository.FindByTag(id)

		if len(img.Labels) == 0 {
			parrot.Info("[" + id + "] - No labels defined")
		} else {
			parrot.Info("[" + id + "] - " + asJson(img.Labels))
		}
	})
}

func CmdBuild(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		endpoint := "unix:///var/run/docker.sock"
		client, _ := docker.NewClient(endpoint)
		imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
		for _, img := range imgs {
			if !repository.Exists(img.ID) {
				repository.Put(img)
				parrot.Info("[" + img.ID + "] - " + strings.Join(img.RepoTags, ", ") + " Added to bucket")
			} else {
				parrot.Info("[" + img.ID + "] - " + strings.Join(img.RepoTags, ", ") + " Not inserted in bucket because already exists")
			}
		}
	})
}

func CmdList(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		imgs := repository.GetAll()
		for _, img := range imgs {
			parrot.Info("[" + img.ID + "] - " + strings.Join(img.RepoTags, ", ") + asJson(img))

		}
	})
}

func CmdWrapper(ctx *cli.Context) {
}

// ----------------
// Arguments from command string
// ----------------
func stringsFromArguments(ctx *cli.Context) ([]string, error) {
	if !ctx.Args().Present() {
		return nil, errors.New("Value must be provided!")
	}

	str := ctx.Args()

	return str, nil
}

func stringFromArguments(ctx *cli.Context) (string, error) {
	if !ctx.Args().Present() {
		return "", errors.New("Value must be provided!")
	}

	str := ctx.Args()[0]

	return str, nil
}

func intFromArguments(ctx *cli.Context) (int, error) {
	if !ctx.Args().Present() {
		return -1, errors.New("Value must be provided!")
	}

	i, err := strconv.Atoi(ctx.Args()[0])
	if err != nil {
		return -1, err
	}

	return i, nil
}

// -------------------------------
// Cli command wrapper
// -------------------------------
func commandWrapper(ctx *cli.Context, cmd quant.Action0) {
	CmdWrapper(ctx)

	cmd()
}
