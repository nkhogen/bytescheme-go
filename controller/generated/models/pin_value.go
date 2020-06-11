// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// PinValue PinValue is the type for the pin value
//
// swagger:model PinValue
type PinValue string

const (

	// PinValueHigh captures enum value "High"
	PinValueHigh PinValue = "High"

	// PinValueLow captures enum value "Low"
	PinValueLow PinValue = "Low"
)

// for schema
var pinValueEnum []interface{}

func init() {
	var res []PinValue
	if err := json.Unmarshal([]byte(`["High","Low"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		pinValueEnum = append(pinValueEnum, v)
	}
}

func (m PinValue) validatePinValueEnum(path, location string, value PinValue) error {
	if err := validate.Enum(path, location, value, pinValueEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this pin value
func (m PinValue) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validatePinValueEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
