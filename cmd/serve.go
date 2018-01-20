package cmd

import (
	"github.com/spf13/cobra"

	rest "github.com/streamtune/gadget/rest"
)

// serveCmd represents the output command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve",
	Long:  `Serve command`,
	Run: func(cmd *cobra.Command, args []string) {
		commandWrapper(args, func() {
			Parrot.Debug("Serve command invoked")

			restServer := rest.NewRest(*Parrot, *Configuration, *Repository, *Commands, *Utilities)

			restServer.Serve()
		})
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
