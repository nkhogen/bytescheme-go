// Code generated by go-swagger; DO NOT EDIT.

package controller

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetControllerParams creates a new GetControllerParams object
// no default values defined in spec.
func NewGetControllerParams() GetControllerParams {

	return GetControllerParams{}
}

// GetControllerParams contains all the bound params for the get controller operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetController
type GetControllerParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*API key
	  Required: true
	  In: header
	*/
	Authorization string
	/*
	  Required: true
	  In: path
	*/
	ControllerID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetControllerParams() beforehand.
func (o *GetControllerParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindAuthorization(r.Header[http.CanonicalHeaderKey("Authorization")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	rControllerID, rhkControllerID, _ := route.Params.GetOK("controllerId")
	if err := o.bindControllerID(rControllerID, rhkControllerID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAuthorization binds and validates parameter Authorization from header.
func (o *GetControllerParams) bindAuthorization(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("Authorization", "header")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("Authorization", "header", raw); err != nil {
		return err
	}

	o.Authorization = raw

	return nil
}

// bindControllerID binds and validates parameter ControllerID from path.
func (o *GetControllerParams) bindControllerID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ControllerID = raw

	return nil
}