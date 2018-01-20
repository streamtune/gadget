package cmd

import (
	"github.com/spf13/cobra"
)

// reviveCmd represents the output command
var reviveCmd = &cobra.Command{
	Use:   "revive",
	Short: "Revive",
	Long:  `Revive command`,
	Run: func(cmd *cobra.Command, args []string) {
		commandWrapper(args, func() {
			Parrot.Debug("Revive command invoked")

			err := Commands.Revive()

			if err != nil {
				Parrot.Warn("Revive", err)
				return
			}
		})
	},
}

func init() {
	RootCmd.AddCommand(reviveCmd)
}
