package main

import (
	"fmt"
	"os"

	"github.com/streamtune/gadget/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

/*
var parrot = quant.NewParrot("gadget")

var settings = Settings{}
var repository = Repository{}

//var pathUtils = quant.NewPathUtils()
var commands = Commands{}

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

	parrot.Debug("Parrot will talk a lot!")
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
			Subcommands: []cli.Command{
				{
					Name:   "no",
					Usage:  "Get a limit number of images",
					Action: CmdListByNumber,
				},
				{
					Name:   "li",
					Usage:  "Get images with name like",
					Action: CmdListByName,
				},
				{
					Name:   "all",
					Usage:  "Get info of the images",
					Action: CmdList,
				},
			},
		},
		{
			Name:    "info",
			Aliases: []string{"in"},
			Usage:   "info the images",
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
func CmdTest(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		parrot.Info("Info")
		parrot.Debug("Debug", "is", "a", "nice", "thing")
		parrot.Warn("Attenction", "please", nil)
		parrot.Error("This is an error", "my friend", "!")
		Parse()
	})
	return nil
}

func CmdServe(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		parrot.Println("Serving gadget for REST Apis on port", settings.RestPort())
		serve()
	})
	return nil
}

func CmdRevive(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		err := commands.Revive()

		if err != nil {
			parrot.Error("Revive", err)
			panic(err)
		}
	})
	return nil
}

func CmdUpdate(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		err := commands.Update()

		if err != nil {
			parrot.Warn("Update", err)
		}

	})
	return nil
}

func CmdList(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		commands.List()
	})

	return nil
}

func CmdListByNumber(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		co, err := intFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}
		commands.ListByNumber(co)
	})

	return nil
}

func CmdListByName(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		name, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}
		commands.ListByName(name)
	})
	return nil
}

func CmdLabels(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		commands.Labels()
	})
	return nil
}

func CmdLabelsById(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		commands.LabelsById(id)
	})
	return nil
}

func CmdLabelsByTag(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		commands.LabelsByTag(id)
	})
	return nil
}

// Infos
func CmdInfo(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		commands.Info()
	})
	return nil
}

func CmdInfoById(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		commands.InfoById(id)
	})
	return nil
}

func CmdInfoByTag(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		commands.InfoByTag(id)
	})
	return nil
}

// Volumes
func CmdVolumes(ctx *cli.Context) error {
	commandWrapper(ctx, func() {

		commands.Volumes()
	})
	return nil
}

func CmdVolumesById(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}
		commands.VolumesById(id)
	})
	return nil
}

func CmdVolumesByTag(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		id, err := stringFromArguments(ctx)
		if err != nil {
			parrot.Error("Error...", err)
			return
		}

		commands.VolumesByTag(id)
	})
	return nil
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
*/
