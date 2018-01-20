package cmd

import (
	"github.com/spf13/cobra"
)

var idInfo string
var taInfo string

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Info",
	Long:  `Info command`,
	Run: func(cmd *cobra.Command, args []string) {
		commandWrapper(args, func() {
			Parrot.Debug("Info command invoked")

			if cmd.Flag("id").Changed {
				err := Commands.InfoById(idInfo)

				if err != nil {
					Parrot.Warn("Info", err)
					return
				}
			} else {
				if cmd.Flag("li").Changed {
					err := Commands.InfoByTag(taInfo)
					if err != nil {
						Parrot.Warn("Info", err)
						return
					}
				} else {
					err := Commands.Info()
					if err != nil {
						Parrot.Warn("Info", err)
						return
					}
				}
			}

			/*
				var no = cmd.Flag("no").Value
				var li = cmd.Flag("li").Value.String()
				var all = cmd.Flag("all").Value.String()
			*/

			/*
				if id != "" {
					var command, err = Repository.FindById(id)

					if err != nil {
						Parrot.Println("Error retrieving command in the store ("+id+")", err)
						return
					}

					Parrot.Println(command.String())
				} else {
					var commands, err = Repository.GetAllCommands()

					if err != nil {
						Parrot.Println("Error retrieving commands in the store", err)
						return
					}

					for _, c := range commands {
						Parrot.Println(c.String())
					}
				}
			*/

		})
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringVarP(&idInfo, "id", "i", "", "Get info by id")
	infoCmd.Flags().StringVarP(&taInfo, "li", "l", "", "Get info by tag like")
}
