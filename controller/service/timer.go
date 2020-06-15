package service

import (
	"bytescheme/common/log"
	"bytescheme/common/service"
	"bytescheme/common/util"
	gmodels "bytescheme/controller/generated/models"
	"bytescheme/controller/operation"
	"bytescheme/controller/shared"
	"context"
)

var (
	// Timer for controller
	Timer *ControllerTimer
)

// ControllerTimer is the controller timer
type ControllerTimer struct {
	registry *operation.Registry
	timer    *service.Timer
}

// InitTimer initiates the controller timer
func InitTimer(registry *operation.Registry) {
	Timer = &ControllerTimer{
		registry: registry,
	}
	Timer.timer = service.NewTimer(shared.Store, Timer.OnEvent)
}

// OnEvent is the callback for events triggered by the timer
func (timer *ControllerTimer) OnEvent(id string, data map[string]interface{}) error {
	controller := &gmodels.Controller{}
	err := util.Convert(data, controller)
	if err != nil {
		return err
	}
	_, err = timer.registry.UpdateController(context.TODO(), controller)
	if err != nil {
		log.Errorf("Error occurred in updating controller from timer. Error: %s\n", err.Error())
		return err
	}
	return nil
}
