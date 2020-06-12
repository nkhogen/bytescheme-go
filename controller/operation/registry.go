package operation

import (
	"bytescheme/common/db"
	"bytescheme/common/util"
	gmodels "bytescheme/controller/generated/models"
	"bytescheme/controller/model"
	"bytescheme/controller/shared"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	// ControllerKeyPrefix is the prefix added to controller ID key
	ControllerKeyPrefix = "controller/"
)

// Registry maintains all the controllers
type Registry struct {
	// Controller ID to processor
	processors map[string]model.Processor
	rwLock     *sync.RWMutex
}

// NewRegistry instantiates a registry
func NewRegistry() (*Registry, error) {
	registry := &Registry{
		processors: map[string]model.Processor{},
		rwLock:     &sync.RWMutex{},
	}
	return registry, nil
}

func (registry *Registry) getProcessorConfig(controllerID string) (*model.ProcessorConfig, error) {
	controllerKey := ControllerKeyPrefix + controllerID
	value, err := shared.Store.Get(controllerKey)
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	if value == nil {
		err = fmt.Errorf("Unrecognized controller %s", controllerID)
		return nil, model.NewServiceError(404, err)
	}
	config := &model.ProcessorConfig{}
	err = util.ConvertFromJSON([]byte(*value), config)
	if err != nil {
		return nil, model.NewServiceError(400, err)
	}
	if config.Controller != nil {
		config.Controller.ID = &controllerID
	}
	return config, nil
}

func (registry *Registry) setProcessorConfig(controllerID string, config *model.ProcessorConfig) error {
	controllerKey := ControllerKeyPrefix + controllerID
	ba, err := util.ConvertToJSON(config)
	if err != nil {
		return model.NewServiceError(500, err)
	}
	config.Version = time.Now().UnixNano()
	err = shared.Store.Set(&db.KeyValue{
		Key:   controllerKey,
		Value: string(ba),
	})
	if err != nil {
		return model.NewServiceError(500, err)
	}
	return nil
}

// Submit submits the callback which is invoked with the target processor
func (registry *Registry) Submit(ctx context.Context, controllerID string, callback func(controllerID string, processor model.Processor) (bool, error)) error {
	registry.rwLock.Lock()
	registry.rwLock.Unlock()
	processor, ok := registry.processors[controllerID]
	if ok {
		config, err := registry.getProcessorConfig(controllerID)
		if err != nil {
			processor.Close()
			delete(registry.processors, controllerID)
			return model.NewServiceError(500, err)
		}
		if config.Version == 0 {
			// Changed by user
			processor.Close()
			delete(registry.processors, controllerID)
			processor, err = NewProcessor(config)
			if err != nil {
				return model.NewServiceError(500, err)
			}
			registry.processors[controllerID] = processor
			if config.Controller != nil {
				// Sync the controller with all the pins
				// Read existing pin values
				cntlr, err := processor.SyncController(ctx, &gmodels.Controller{ID: &controllerID, Pins: []*gmodels.Pin{}})
				if err == nil {
					config.Controller = cntlr
				}
			}
			// Update version
			err = registry.setProcessorConfig(controllerID, config)
			if err != nil {
				return model.NewServiceError(500, err)
			}
		}
	} else {
		config, err := registry.getProcessorConfig(controllerID)
		if err != nil {
			return model.NewServiceError(500, err)
		}
		processor, err = NewProcessor(config)
		if err != nil {
			return model.NewServiceError(500, err)
		}
		registry.processors[controllerID] = processor
	}
	saveConfig, err := callback(controllerID, processor)
	if err != nil {
		return model.NewServiceError(500, err)
	}
	if saveConfig {
		err = registry.setProcessorConfig(controllerID, processor.GetConfig())
	}
	return err
}

// ListControllers returns all the controllers
func (registry *Registry) ListControllers(ctx context.Context) ([]*gmodels.Controller, error) {
	controllers := []*gmodels.Controller{}
	keys, err := shared.Store.GetKeys(ControllerKeyPrefix)
	if err != nil {
		return controllers, model.NewServiceError(500, err)
	}
	for _, key := range keys {
		controllerID := strings.TrimPrefix(key, ControllerKeyPrefix)
		err = registry.Submit(ctx, controllerID, func(controllerID string, processor model.Processor) (bool, error) {
			controller, err := processor.GetController(ctx, controllerID)
			if err != nil {
				return false, err
			}
			controllers = append(controllers, controller)
			return false, nil
		})
		if err != nil {
			return []*gmodels.Controller{}, err
		}
	}
	return controllers, nil
}

// GetController returns the given controller
func (registry *Registry) GetController(ctx context.Context, controllerID string) (*gmodels.Controller, error) {
	var err error
	var controller *gmodels.Controller
	err = registry.Submit(ctx, controllerID, func(controllerID string, processor model.Processor) (bool, error) {
		controller, err = processor.GetController(ctx, controllerID)
		return false, err
	})
	if err != nil {
		return nil, err
	}
	return controller, nil
}

// UpdateController updates the given controller
func (registry *Registry) UpdateController(ctx context.Context, controller *gmodels.Controller) (*gmodels.Controller, error) {
	var err error
	var cntlr *gmodels.Controller
	if controller.ID == nil {
		return cntlr, fmt.Errorf("Invalid controller ID")
	}
	err = registry.Submit(ctx, *controller.ID, func(controllerID string, processor model.Processor) (bool, error) {
		cntlr, err = processor.SyncController(ctx, controller)
		return true, err
	})
	if err != nil {
		return nil, model.NewServiceError(500, err)
	}
	return cntlr, nil
}
