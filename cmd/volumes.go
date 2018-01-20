package cmd

import (
	"github.com/spf13/cobra"
)

var idVolumes string
var taVolumes string

// volumesCmd represents the volumes command
var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "Volumes",
	Long:  `Volumes command`,
	Run: func(cmd *cobra.Command, args []string) {
		commandWrapper(args, func() {
			Parrot.Debug("Volumes command invoked")

			if cmd.Flag("id").Changed {
				err := Commands.VolumesById(idVolumes)

				if err != nil {
					Parrot.Warn("Volumes", err)
					return
				}
			} else {
				if cmd.Flag("li").Changed {
					err := Commands.VolumesByTag(taVolumes)
					if err != nil {
						Parrot.Warn("Volumes", err)
						return
					}
				} else {
					err := Commands.Volumes()
					if err != nil {
						Parrot.Warn("Volumes", err)
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
	RootCmd.AddCommand(volumesCmd)

	volumesCmd.Flags().StringVarP(&idVolumes, "id", "i", "", "Get volumes by id")
	volumesCmd.Flags().StringVarP(&taVolumes, "li", "l", "", "Get volumes by tag like")
}
