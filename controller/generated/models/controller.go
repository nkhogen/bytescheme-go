// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Controller Controller model
//
// swagger:model Controller
type Controller struct {

	// ID
	// Required: true
	// Min Length: 1
	ID *string `json:"id"`

	// pins
	// Required: true
	Pins []*Pin `json:"pins"`

	// description
	Description string `json:"description,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this controller
func (m *Controller) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePins(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Controller) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	if err := validate.MinLength("id", "body", string(*m.ID), 1); err != nil {
		return err
	}

	return nil
}

func (m *Controller) validatePins(formats strfmt.Registry) error {

	if err := validate.Required("pins", "body", m.Pins); err != nil {
		return err
	}

	for i := 0; i < len(m.Pins); i++ {
		if swag.IsZero(m.Pins[i]) { // not required
			continue
		}

		if m.Pins[i] != nil {
			if err := m.Pins[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("pins" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Controller) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Controller) UnmarshalBinary(b []byte) error {
	var res Controller
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
