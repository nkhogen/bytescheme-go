// Code generated by go-swagger; DO NOT EDIT.

package controller

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"bytescheme/controller/generated/models"
)

// GetControllerOKCode is the HTTP code returned for type GetControllerOK
const GetControllerOKCode int = 200

/*GetControllerOK GetControllerResponse is the response model

swagger:response getControllerOK
*/
type GetControllerOK struct {

	/*
	  In: Body
	*/
	Payload *models.Controller `json:"body,omitempty"`
}

// NewGetControllerOK creates GetControllerOK with default headers values
func NewGetControllerOK() *GetControllerOK {

	return &GetControllerOK{}
}

// WithPayload adds the payload to the get controller o k response
func (o *GetControllerOK) WithPayload(payload *models.Controller) *GetControllerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get controller o k response
func (o *GetControllerOK) SetPayload(payload *models.Controller) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetControllerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetControllerDefault APIErrorResponse is all API errors

swagger:response getControllerDefault
*/
type GetControllerDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.APIError `json:"body,omitempty"`
}

// NewGetControllerDefault creates GetControllerDefault with default headers values
func NewGetControllerDefault(code int) *GetControllerDefault {
	if code <= 0 {
		code = 500
	}

	return &GetControllerDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get controller default response
func (o *GetControllerDefault) WithStatusCode(code int) *GetControllerDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get controller default response
func (o *GetControllerDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get controller default response
func (o *GetControllerDefault) WithPayload(payload *models.APIError) *GetControllerDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get controller default response
func (o *GetControllerDefault) SetPayload(payload *models.APIError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetControllerDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}