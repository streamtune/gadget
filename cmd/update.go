package cmd

import (
	"github.com/spf13/cobra"
)

// updateCmd represents the output command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update",
	Long:  `Update command`,
	Run: func(cmd *cobra.Command, args []string) {
		commandWrapper(args, func() {
			Parrot.Debug("Update command invoked")

			err := Commands.Update()

			if err != nil {
				Parrot.Warn("Update", err)
				return
			}
		})
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
