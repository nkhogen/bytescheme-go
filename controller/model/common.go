package model

import (
	gmodels "bytescheme/controller/generated/models"
)

// ServiceError is the error returned for a service
type ServiceError struct {
	*gmodels.APIError
}

// Code is compatible with swagger generated error
func (err *ServiceError) Code() int {
	return int(err.Status)
}

func (err *ServiceError) Error() string {
	return err.Message
}

// NewServiceError returns the instance of ServiceError
func NewServiceError(status int, err error) *ServiceError {
	if serviceErr, ok := err.(*ServiceError); ok {
		return serviceErr
	}
	return &ServiceError{
		&gmodels.APIError{
			Status:  int32(status),
			Message: err.Error(),
		},
	}
}
