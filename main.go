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
func CmdServe(ctx *cli.Context) error {
	commandWrapper(ctx, func() {
		parrot.Println("Serving gadget for REST Apis on port", settings.RestPort())
		serve()
	})
	return nil
}
*/
