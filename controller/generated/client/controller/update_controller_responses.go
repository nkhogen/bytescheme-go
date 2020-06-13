// Code generated by go-swagger; DO NOT EDIT.

package controller

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"bytescheme/controller/generated/models"
)

// UpdateControllerReader is a Reader for the UpdateController structure.
type UpdateControllerReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateControllerReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateControllerOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUpdateControllerDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateControllerOK creates a UpdateControllerOK with default headers values
func NewUpdateControllerOK() *UpdateControllerOK {
	return &UpdateControllerOK{}
}

/*UpdateControllerOK handles this case with default header values.

UpdateControllerResponse is the response for controller update
*/
type UpdateControllerOK struct {
	Payload *models.Controller
}

func (o *UpdateControllerOK) Error() string {
	return fmt.Sprintf("[PUT /v1/controllers/{controllerId}][%d] updateControllerOK  %+v", 200, o.Payload)
}

func (o *UpdateControllerOK) GetPayload() *models.Controller {
	return o.Payload
}

func (o *UpdateControllerOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Controller)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateControllerDefault creates a UpdateControllerDefault with default headers values
func NewUpdateControllerDefault(code int) *UpdateControllerDefault {
	return &UpdateControllerDefault{
		_statusCode: code,
	}
}

/*UpdateControllerDefault handles this case with default header values.

APIErrorResponse is all API errors
*/
type UpdateControllerDefault struct {
	_statusCode int

	Payload *models.APIError
}

// Code gets the status code for the update controller default response
func (o *UpdateControllerDefault) Code() int {
	return o._statusCode
}

func (o *UpdateControllerDefault) Error() string {
	return fmt.Sprintf("[PUT /v1/controllers/{controllerId}][%d] UpdateController default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateControllerDefault) GetPayload() *models.APIError {
	return o.Payload
}

func (o *UpdateControllerDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}