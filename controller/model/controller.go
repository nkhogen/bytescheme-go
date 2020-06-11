package model

import (
	gmodels "bytescheme/controller/generated/models"
	"context"
)

// Processor for a controller
type Processor interface {
	GetController(context.Context, string) (*gmodels.Controller, error)
	SyncController(context.Context, *gmodels.Controller) (*gmodels.Controller, error)
	GetConfig() *ProcessorConfig
	Close() error
}

// ProcessorConfig is the config for a processor
type ProcessorConfig struct {
	Host       string              `json:"host,omitempty"`
	Port       int                 `json:"port,omitempty"`
	APIKey     string              `json:"apiKey,omitempty"`
	Controller *gmodels.Controller `json:"controller,omitempty"`
	Version    int64               `json:"version"`
}
