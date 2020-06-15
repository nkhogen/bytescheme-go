package cli

import (
	"bytescheme/controller/operation"
	"bytescheme/controller/service"
	"bytescheme/controller/shared"
	"fmt"

	"github.com/spf13/cobra"
)

func setupServiceCommand() {
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Runs the service",
	}
	rootCmd.AddCommand(serviceCmd)

	serviceCmdStart := &cobra.Command{
		Use:   "start",
		Short: "Starts the service",
		Run:   serviceCommandStart,
	}
	// -h already registered for help
	serviceCmdStart.Flags().StringVarP(&serviceCmdParams.host, "host", "s", serviceCmdParams.host, "Host to be connected")
	serviceCmdStart.Flags().IntVarP(&serviceCmdParams.port, "port", "n", serviceCmdParams.port, "Port to be connected")
	serviceCmd.AddCommand(serviceCmdStart)
}

func serviceCommandStart(cmd *cobra.Command, args []string) {
	shared.InitStore()
	registry, err := operation.NewRegistry()
	if err != nil {
		panic(err)
	}
	service.InitTimer(registry)
	svc, err := service.NewService(serviceCmdParams.host, serviceCmdParams.port, registry)
	if err != nil {
		panic(err)
	}
	err = svc.Serve()
	if err != nil {
		fmt.Printf("Error occurred while starting the service. Error: %s\n", err.Error())
	}
}
