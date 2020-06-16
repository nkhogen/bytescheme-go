package operation

import (
	"bytescheme/common/log"
	"bytescheme/common/service"
	"bytescheme/common/util"
	cntlr "bytescheme/controller/generated/client/controller"
	gmodels "bytescheme/controller/generated/models"
	"bytescheme/controller/model"
	"context"
	"fmt"
	"sync"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	rpio "github.com/stianeikeland/go-rpio/v4"
)

const (
	// IDMultiplier is the multiplier for attached devices
	IDMultiplier int = 100
	// SetPowerStatusFormat is the message format to set power status on the connected device
	SetPowerStatusFormat string = "SET %d %d"
	// GetPowerStatusFormat is the message format to get power status from the connected device
	GetPowerStatusFormat string = "GET %d"
)

// LocalProcessor means the RPI is local
type LocalProcessor struct {
	rwLock      *sync.RWMutex
	config      *model.ProcessorConfig
	eventServer *service.EventServer
}

// RemoteProcessor means the RPI is remote
type RemoteProcessor struct {
	client cntlr.ClientService
	config *model.ProcessorConfig
}

// ResolvePin resolve into client and pin ID
func ResolvePin(id int) (int, int) {
	clientID := id / IDMultiplier
	pinID := id % IDMultiplier
	return clientID, pinID
}

// NewProcessor creates an instance of either a local or a remote processor with the given config
func NewProcessor(config *model.ProcessorConfig) (model.Processor, error) {
	if config.Controller != nil && config.Controller.ID != nil {
		return NewLocalProcessor(config)
	}
	return NewRemoteProcessor(config)
}

// NewLocalProcessor creates an instance of a local processor
func NewLocalProcessor(config *model.ProcessorConfig) (model.Processor, error) {
	err := rpio.Open()
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	processor := &LocalProcessor{
		rwLock: &sync.RWMutex{},
		config: config,
	}
	if config.Host != "" && config.Port > 0 {
		eventServer, err := service.NewEventServer(config.Host, config.Port, func(clientID int) error {
			// Sync the client with the cached data
			_, err := processor.SyncController(context.Background(), config.Controller)
			return err
		})
		if err != nil {
			return nil, model.NewServiceError(500, err)
		}
		processor.eventServer = eventServer
	}
	util.ShutdownHandler.RegisterCloseable(processor)
	return processor, nil
}

// NewRemoteProcessor creates an instance of remote processor
func NewRemoteProcessor(config *model.ProcessorConfig) (model.Processor, error) {
	// create the transport
	transport := httptransport.New(fmt.Sprintf("%s:%d", config.Host, config.Port), "", []string{"http"})
	// create the API client, with the transport
	client := cntlr.New(transport, strfmt.Default)
	processor := &RemoteProcessor{
		client: client,
		config: config,
	}
	return processor, nil
}

// Close closes the processor
func (localProcessor *LocalProcessor) Close() error {
	rpio.Close()
	eventServer := localProcessor.eventServer
	if eventServer != nil {
		return eventServer.Close()
	}
	return nil
}

// Close closes the processor
func (remoteProcessor *RemoteProcessor) Close() error {
	return nil
}

// GetConfig returns the processor config
func (localProcessor *LocalProcessor) GetConfig() *model.ProcessorConfig {
	return localProcessor.config
}

// GetConfig returns the processor config
func (remoteProcessor *RemoteProcessor) GetConfig() *model.ProcessorConfig {
	return remoteProcessor.config
}

// WriteClientPin writes message to the client pin
func (localProcessor *LocalProcessor) WriteClientPin(clientID, pinID int, value gmodels.PinValue) (gmodels.PinValue, error) {
	pinValue := 0
	if value == gmodels.PinValueHigh {
		pinValue = 1
	}
	message := fmt.Sprintf(SetPowerStatusFormat, pinID, pinValue)
	message, err := localProcessor.eventServer.Send(clientID, message)
	if err != nil {
		return gmodels.PinValueLow, err
	}
	if message == "TRUE" {
		return gmodels.PinValueHigh, nil
	}
	return gmodels.PinValueLow, nil
}

// ReadClientPin reads message from the client pin
func (localProcessor *LocalProcessor) ReadClientPin(clientID, pinID int) (gmodels.PinValue, error) {
	message := fmt.Sprintf(GetPowerStatusFormat, pinID)
	message, err := localProcessor.eventServer.Send(clientID, message)
	if err != nil {
		return gmodels.PinValueLow, err
	}
	if message == "TRUE" {
		return gmodels.PinValueHigh, nil
	}
	return gmodels.PinValueLow, nil
}

// SyncController syncs the given local controller
func (localProcessor *LocalProcessor) SyncController(ctx context.Context, controller *gmodels.Controller) (*gmodels.Controller, error) {
	localProcessor.rwLock.Lock()
	defer localProcessor.rwLock.Unlock()
	cntlr := localProcessor.config.Controller
	inPinMap := map[int]*gmodels.Pin{}
	for idx := range controller.Pins {
		inPin := controller.Pins[idx]
		inPinMap[int(*inPin.ID)] = inPin
	}
	// Sync all pins
	for idx := range cntlr.Pins {
		pin := cntlr.Pins[idx]
		pinID := int(*pin.ID)
		inPin, ok := inPinMap[pinID]
		clientID, pinID := ResolvePin(pinID)
		if !ok {
			if clientID > 0 {
				pinValue, err := localProcessor.ReadClientPin(clientID, pinID)
				if err == nil {
					pin.Value = pinValue
				} else {
					log.Errorf("Error occurred in reading pin %d for client %d. Error: %s", pinID, clientID, err.Error())
				}
				continue
			}
			rpin := rpio.Pin(pinID)
			in := rpin.Read()
			if in == rpio.High {
				pin.Value = gmodels.PinValueHigh
			} else {
				pin.Value = gmodels.PinValueLow
			}
			continue
		}
		if clientID > 0 {
			if pin.Mode == gmodels.PinModeOutput {
				pinValue, err := localProcessor.WriteClientPin(clientID, pinID, inPin.Value)
				if err == nil {
					pin.Value = pinValue
				} else {
					log.Errorf("Error occurred in setting pin %d for client %d. Error: %s", pinID, clientID, err.Error())
				}
				continue
			}
			pinValue, err := localProcessor.ReadClientPin(clientID, pinID)
			if err == nil {
				pin.Value = pinValue
			} else {
				log.Errorf("Error occurred in reading pin %d for client %d. Error: %s", pinID, clientID, err.Error())
			}
			continue
		}
		rpin := rpio.Pin(pinID)
		if pin.Mode == gmodels.PinModeOutput {
			rpin.Output()
			log.Infof("Setting pin %d for controller %s to %s", pinID, *controller.ID, inPin.Value)
			if inPin.Value == gmodels.PinValueHigh {
				rpin.High()
			} else {
				rpin.Low()
			}
			pin.Value = inPin.Value
			continue
		}
		log.Infof("Reading pin %d for controller %s", pinID, *controller.ID)
		rpin.Input()
		in := rpin.Read()
		if in == rpio.High {
			pin.Value = gmodels.PinValueHigh
		} else {
			pin.Value = gmodels.PinValueLow
		}
	}
	controller = &gmodels.Controller{}
	err := util.Convert(cntlr, controller)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	return controller, nil
}

// SyncController syncs the given remote controller
func (remoteProcessor *RemoteProcessor) SyncController(ctx context.Context, controller *gmodels.Controller) (*gmodels.Controller, error) {
	param := cntlr.NewUpdateControllerParamsWithContext(ctx)
	param.Payload = controller
	authorization := remoteProcessor.GetAuthorization()
	ok, err := remoteProcessor.client.UpdateController(param, authorization)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	return ok.Payload, nil
}

// GetController returns the given controller from the local processor
func (localProcessor *LocalProcessor) GetController(ctx context.Context, controllerID string) (*gmodels.Controller, error) {
	cntlr := localProcessor.config.Controller
	if controllerID != *cntlr.ID {
		return nil, fmt.Errorf("Invalid controller ID received by local processor")
	}
	controller := &gmodels.Controller{}
	err := util.Convert(cntlr, controller)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	return controller, nil
}

// GetController returns the given controller from a remote processor
func (remoteProcessor *RemoteProcessor) GetController(ctx context.Context, controllerID string) (*gmodels.Controller, error) {
	param := cntlr.NewGetControllerParamsWithContext(ctx)
	param.ControllerID = controllerID
	authorization := remoteProcessor.GetAuthorization()
	ok, err := remoteProcessor.client.GetController(param, authorization)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	return ok.Payload, nil
}

// GetAuthorization returns the authorization
func (remoteProcessor *RemoteProcessor) GetAuthorization() runtime.ClientAuthInfoWriter {
	return httptransport.APIKeyAuth("Authorization", "header", remoteProcessor.config.APIKey)
}
