package cmd

import (
	"github.com/spf13/cobra"
)

var idLabels string
var taLabels string

// labelsCmd represents the labels command
var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Labels",
	Long:  `Labels command`,
	Run: func(cmd *cobra.Command, args []string) {
		commandWrapper(args, func() {
			Parrot.Debug("Labels command invoked")

			if cmd.Flag("id").Changed {
				err := Commands.LabelsById(idLabels)

				if err != nil {
					Parrot.Warn("Labels", err)
					return
				}
			} else {
				if cmd.Flag("li").Changed {
					err := Commands.LabelsByTag(taLabels)
					if err != nil {
						Parrot.Warn("Labels", err)
						return
					}
				} else {
					err := Commands.Labels()
					if err != nil {
						Parrot.Warn("Labels", err)
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
	RootCmd.AddCommand(labelsCmd)

	labelsCmd.Flags().StringVarP(&idLabels, "id", "i", "", "Get labels by id")
	labelsCmd.Flags().StringVarP(&taLabels, "li", "l", "", "Get labels by tag like")
}
