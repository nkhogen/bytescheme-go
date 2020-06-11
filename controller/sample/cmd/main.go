package main

import (
	"bytescheme/controller/generated/client/controller"
	"bytescheme/controller/generated/models"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

const (
	Host         = "bytescheme.mynetgear.com"
	Port         = 443
	Scheme       = "https"
	APIKey       = "test"
	ControllerID = "bfd8dd0a-10db-4782-86ec-b27f52d6362c"
)

func getControllerClient() controller.ClientService {
	server := fmt.Sprintf("%s:%d", Host, Port)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	transport := httptransport.NewWithClient(server, "", []string{Scheme}, client)
	return controller.New(transport, strfmt.Default)
}

func getControllerAuth() runtime.ClientAuthInfoWriter {
	return httptransport.APIKeyAuth("Authorization", "header", APIKey)
}

func getController() (*models.Controller, error) {
	client := getControllerClient()
	authParam := getControllerAuth()

	params := controller.NewGetControllerParams()
	params.ControllerID = ControllerID

	ok, err := client.GetController(params, authParam)
	if err != nil {
		fmt.Printf("Error in getting the controller. Error: %s\n", err.Error())
		return nil, err
	}
	return ok.Payload, nil
}

func updateController(cntlr *models.Controller) error {
	client := getControllerClient()
	authParam := getControllerAuth()
	updateParams := controller.NewUpdateControllerParams()
	updateParams.ControllerID = *cntlr.ID
	updateParams.Payload = cntlr
	_, err := client.UpdateController(updateParams, authParam)
	return err
}

func sleep(sec int) {
	time.Sleep(time.Duration(int(time.Second) * sec))
}

func main() {

	cntlr, err := getController()
	if err != nil {
		panic(err)
	}
	pin := cntlr.Pins[0]

	sleepTime := 1
	for {
		pin.Value = models.PinValueHigh
		updateController(cntlr)
		sleep(sleepTime)
		pin.Value = models.PinValueLow
		updateController(cntlr)
		sleep(sleepTime)
	}
}
