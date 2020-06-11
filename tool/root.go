package tool

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tool ...",
		Short: "CLI for tools",
	}
)

func setupToolCommands() {
	replaceMiddlewareCmd := &cobra.Command{
		Use:   "replace-middleware",
		Short: "Replaces middleware",
		Run:   replaceMiddleware,
	}
	replaceMiddlewareCmd.Flags().StringVarP(&processorCmdParams.filepath, "file", "f", processorCmdParams.filepath, "Middleware filepath")
	replaceMiddlewareCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(replaceMiddlewareCmd)

	processSwaggerCmd := &cobra.Command{
		Use:   "process-swagger",
		Short: "Post process swagger",
		Run:   processSwagger,
	}
	processSwaggerCmd.Flags().StringVarP(&processorCmdParams.filepath, "file", "f", processorCmdParams.filepath, "Swagger filepath")
	processSwaggerCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(processSwaggerCmd)
}

func replaceMiddleware(cmd *cobra.Command, args []string) {
	err := ReplaceGlobalMiddlewareFunc(processorCmdParams.filepath)
	if err != nil {
		panic(err)
	}
}

func processSwagger(cmd *cobra.Command, args []string) {
	err := ProcessSwagger(processorCmdParams.filepath)
	if err != nil {
		panic(err)
	}
}

// Execute is the entry for the command
func Execute() {
	setupToolCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
