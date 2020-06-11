package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "controller ...",
		Short: "CLI for controller",
	}
)

// Execute is the entry for the command
func Execute() {
	setupControllerCommand()
	setupStoreCommand()
	setupServiceCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
