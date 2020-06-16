package cli

import (
	"bytescheme/controller/generated/client/controller"
	"bytescheme/controller/generated/models"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

func setupControllerCommand() {
	controllerCmd := &cobra.Command{
		Use:   "controller",
		Short: "Controller commands",
	}
	rootCmd.AddCommand(controllerCmd)

	controllerCmdGet := &cobra.Command{
		Use:   "get",
		Short: "Gets a controller",
		Run:   controllerCommandGet,
	}
	// -h already registered for help
	controllerCmdGet.Flags().StringVarP(&controllerCmdParams.host, "server", "s", controllerCmdParams.host, "Host to be connected")
	controllerCmdGet.Flags().IntVarP(&controllerCmdParams.port, "port", "n", controllerCmdParams.port, "Port to be connected")
	controllerCmdGet.Flags().StringVarP(&controllerCmdParams.apiKey, "apikey", "a", controllerCmdParams.apiKey, "API key for the service access")
	controllerCmdGet.Flags().StringVarP(&controllerCmdParams.controllerID, "controller", "i", controllerCmdParams.controllerID, "Controller ID")
	controllerCmdGet.MarkFlagRequired("controller")
	controllerCmd.AddCommand(controllerCmdGet)

	controllerCmdSet := &cobra.Command{
		Use:   "set",
		Short: "Sets a controller",
		Run:   controllerCommandSet,
	}

	controllerCmdSet.Flags().StringVarP(&controllerCmdParams.host, "server", "s", controllerCmdParams.host, "Host to be connected")
	controllerCmdSet.Flags().IntVarP(&controllerCmdParams.port, "port", "n", controllerCmdParams.port, "Port to be connected")
	controllerCmdSet.Flags().StringVarP(&controllerCmdParams.apiKey, "apikey", "a", controllerCmdParams.apiKey, "API key for the service access")
	controllerCmdSet.Flags().StringVarP(&controllerCmdParams.controllerID, "controller", "i", controllerCmdParams.controllerID, "Controller ID")
	controllerCmdSet.Flags().IntVarP(&controllerCmdParams.pinID, "pin", "p", controllerCmdParams.pinID, "Pin ID")
	controllerCmdSet.Flags().BoolVarP(&controllerCmdParams.pinHigh, "enable", "e", controllerCmdParams.pinHigh, "Pin enable value")
	controllerCmdSet.MarkFlagRequired("controller")
	controllerCmdSet.MarkFlagRequired("pin")
	controllerCmdSet.MarkFlagRequired("high")
	controllerCmd.AddCommand(controllerCmdSet)

}

func getControllerClient() controller.ClientService {
	server := fmt.Sprintf("%s:%d", controllerCmdParams.host, controllerCmdParams.port)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	transport := httptransport.NewWithClient(server, "", []string{controllerCmdParams.scheme}, client)
	return controller.New(transport, strfmt.Default)
}

func getControllerAuth() runtime.ClientAuthInfoWriter {
	return httptransport.APIKeyAuth("Authorization", "header", controllerCmdParams.apiKey)
}

func controllerCommandGet(cmd *cobra.Command, args []string) {
	client := getControllerClient()
	authParam := getControllerAuth()

	params := controller.NewGetControllerParams()
	params.ControllerID = controllerCmdParams.controllerID

	ok, err := client.GetController(params, authParam)
	if err != nil {
		fmt.Printf("Error in getting the controller. Error: %s\n", err.Error())
		return
	}
	ba, _ := json.MarshalIndent(ok.Payload, " ", " ")
	fmt.Printf("Controller \n%s\n", string(ba))
}

func controllerCommandSet(cmd *cobra.Command, args []string) {
	client := getControllerClient()
	authParam := getControllerAuth()

	params := controller.NewGetControllerParams()
	params.ControllerID = controllerCmdParams.controllerID

	ok, err := client.GetController(params, authParam)
	if err != nil {
		fmt.Printf("Error in getting the controller. Error: %s\n", err.Error())
		return
	}
	cntlr := ok.Payload
	pins := cntlr.Pins
	for idx := range pins {
		pin := pins[idx]
		if int(*pin.ID) == controllerCmdParams.pinID {
			pin.Value = models.PinValueLow
			if controllerCmdParams.pinHigh {
				pin.Value = models.PinValueHigh
			}
			updateParams := controller.NewUpdateControllerParams()
			updateParams.ControllerID = *cntlr.ID
			updateParams.Payload = cntlr
			updateOk, err := client.UpdateController(updateParams, authParam)
			if err != nil {
				fmt.Printf("Error in updating the controller. Error: %s\n", err.Error())
			} else {
				ba, _ := json.MarshalIndent(updateOk.Payload, " ", " ")
				fmt.Printf("Controller \n%s\n", string(ba))
			}
			return
		}
	}
	fmt.Printf("Unrecognized pin %d\n", controllerCmdParams.pinID)
}
