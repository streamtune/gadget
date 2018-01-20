package cmd

import (
	"github.com/spf13/cobra"
)

var noList int
var liList string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List",
	Long:  `List command`,
	Run: func(cmd *cobra.Command, args []string) {
		commandWrapper(args, func() {
			Parrot.Debug("List command invoked")

			if cmd.Flag("no").Changed {
				err := Commands.ListByNumber(noList)

				if err != nil {
					Parrot.Warn("List", err)
					return
				}
			} else {
				if cmd.Flag("li").Changed {
					err := Commands.ListByName(liList)
					if err != nil {
						Parrot.Warn("List", err)
						return
					}
				} else {
					err := Commands.List()
					if err != nil {
						Parrot.Warn("List", err)
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
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&noList, "no", "n", 10, "Get a limit number of images")
	listCmd.Flags().StringVarP(&liList, "li", "l", "", "Get images with name like")
}
