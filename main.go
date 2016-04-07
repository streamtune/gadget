package main

import (
	"errors"
	//	"fmt"
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
var pathUtils = quant.NewPathUtils()

func initDB() {
	repository.InitDB()
}

func closeDB() {
	repository.CloseDB()
}

func readSettings() {
	settings.LoadSettings()

	if settings.DebugMode() {
		parrot = quant.NewVerboseParrot("gadget")
	}

	parrot.Debug("Parrot is set to talk a lot!")
}

func main() {
	readSettings()
	initDB()

	// -------------------
	app := cli.NewApp()
	app.Name = "gadget"
	app.Usage = "The inspector will be used to inspect docker images"
	app.Version = "0.0.1"
	app.Copyright = "gi4nks - 2016"

	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"up"},
			Usage:   "update the database",
			Action:  CmdUpdate,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list all the images",
			Action:  CmdList,
		},
		{
			Name:    "inspect",
			Aliases: []string{"in"},
			Usage:   "inspect the images",
			Subcommands: []cli.Command{
				{
					Name:   "id",
					Usage:  "Get the info of specific image",
					Action: CmdInfoById,
				},
				{
					Name:   "tag",
					Usage:  "Get the info of specific tag",
					Action: CmdInfoByTag,
				},
				{
					Name:   "all",
					Usage:  "Get info of the images",
					Action: CmdInfo,
				},
			},
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
			Name:    "volumes",
			Aliases: []string{"vo"},
			Usage:   "show me the volumes of images",
			Subcommands: []cli.Command{
				{
					Name:   "id",
					Usage:  "Get the volumes of specific image",
					Action: CmdVolumesById,
				},
				{
					Name:   "tag",
					Usage:  "Get the volumes of specific tag",
					Action: CmdVolumesByTag,
				},
				{
					Name:   "all",
					Usage:  "Get all the volumes",
					Action: CmdVolumes,
				},
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"se"},
			Usage:   "serving gadget for rest apis",
			Action:  CmdServe,
		},
		{
			Name:    "test",
			Aliases: []string{"te"},
			Usage:   "test",
			Action:  CmdTest,
		},
	}

	app.Run(os.Args)
	defer closeDB()
}

// List of functions
func CmdTest(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		parrot.Info("Info")
		parrot.Debug("Debug", "is", "a", "nice", "thing")
		parrot.Warn("Attenction", "please", nil)
		parrot.Error("This is an error", "my friend", "!")
		Parse()
	})
}

func CmdServe(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		parrot.Println("Serving gadget for REST Apis on port", settings.RestPort())
		serve()
	})
}

func CmdRevive(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		parrot.Debug("Reviving gadget will reinitialize all the datas.")

		repository.BackupSchema()
		repository.InitSchema()

		parrot.Println("Gadget revived")
	})
}

func CmdUpdate(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		endpoint := settings.LocalDockerEndpoint()
		client, _ := docker.NewClient(endpoint)

		if settings.UseDockerMachine() {
			endpoint = settings.MachineDockerEndpoint()

			parrot.Debug("Configuring client for tls")

			client, _ = docker.NewTLSClient(endpoint,
				settings.MachineDockerCertFile(),
				settings.MachineDockerKeyFile(),
				settings.MachineDockerCAFile())
		}

		imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})

		var c = 0

		for _, img := range imgs {
			var id = TruncateID(img.ID)

			parrot.Debug("ID is", id)

			if !repository.Exists(id) {
				repository.Put(img)
				c = c + 1
				parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), "added to bucket")
			} else {
				parrot.Debug("["+id+"] - ", strings.Join(img.RepoTags, ", "), " not inserted in bucket because already exists")
			}
		}

		parrot.Println("Added " + strconv.Itoa(c) + " images")
	})
}

func CmdLabels(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		//var ii = ImageLabel{}
		//var iis = [][]string{}

		/*
			for _, img := range repository.GetAll() {

					for _, r := range AsImageLabels(img).Rows() {
						iis = append(iis, r)
					}

			}
		*/

		//parrot.TablePrint(ii.Header(), iis)
	})
}

func CmdLabelsById(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		parrot.Debug("ID", id)
		/*

			var img = repository.Get(id)


			var ii = ImageLabels{}
			var iis = [][]string{}

			for _, r := range AsImageLabels(img).Rows() {
				iis = append(iis, r)
			}

			parrot.TablePrint(ii.Header(), iis)
		*/
	})
}

func CmdLabelsByTag(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		parrot.Debug("ID", id)

		/*
			var img = repository.FindByTag(id)

			var ii = ImageLabels{}
			var iis = [][]string{}

			for _, r := range AsImageLabels(img).Rows() {
				iis = append(iis, r)
			}

			parrot.TablePrint(ii.Header(), iis)
		*/
	})
}

// Volumes
func CmdInfo(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		var images = repository.GetAll()

		var header = Image{}
		var iis = [][]string{}

		for _, i := range images {
			for _, r := range i.Rows() {
				iis = append(iis, r)
			}
		}

		parrot.TablePrint(header.Header(), iis)
	})
}

func CmdInfoById(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		var image = repository.Get(id)

		var header = Image{}
		var iis = [][]string{}

		for _, r := range image.Rows() {
			iis = append(iis, r)
		}

		parrot.TablePrint(header.Header(), iis)
	})
}

func CmdInfoByTag(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		var images = repository.FindByTag(id)

		var header = Image{}
		var iis = [][]string{}

		for _, i := range images {
			for _, r := range i.Rows() {
				iis = append(iis, r)
			}
		}

		parrot.TablePrint(header.Header(), iis)
	})
}

// Volumes
func CmdVolumes(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		/*
			var ii = ImageVolumes{}
			var iis = [][]string{}

			for _, img := range repository.GetAll() {
				for _, r := range AsImageVolumes(img).Rows() {
					iis = append(iis, r)
				}
			}

			parrot.TablePrint(ii.Header(), iis)
		*/
	})
}

func CmdVolumesById(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		var img = repository.Get(id)

		if len(img.Labels) == 0 {
			parrot.Info("[" + id + "] - No labels defined")
		} else {
			parrot.Info("[" + id + "] - " + asJson(img.Labels))
		}
	})
}

func CmdVolumesByTag(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		parrot.Debug("ID", id)
		/*
			var img = repository.FindByTag(id)

			if len(img.Labels) == 0 {
				parrot.Info("[" + id + "] - No labels defined")
			} else {
				parrot.Info("[" + id + "] - " + asJson(img.Labels))
			}
		*/
	})
}

func CmdList(ctx *cli.Context) {
	commandWrapper(ctx, func() {
		var images = repository.GetAll()

		var header = Image{}
		var iis = [][]string{}

		for _, i := range images {
			for _, r := range i.Rows() {
				iis = append(iis, r)
			}
		}

		parrot.TablePrint(header.Header(), iis)
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
