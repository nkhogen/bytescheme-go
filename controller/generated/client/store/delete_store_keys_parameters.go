// Code generated by go-swagger; DO NOT EDIT.

package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewDeleteStoreKeysParams creates a new DeleteStoreKeysParams object
// with the default values initialized.
func NewDeleteStoreKeysParams() *DeleteStoreKeysParams {
	var ()
	return &DeleteStoreKeysParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteStoreKeysParamsWithTimeout creates a new DeleteStoreKeysParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteStoreKeysParamsWithTimeout(timeout time.Duration) *DeleteStoreKeysParams {
	var ()
	return &DeleteStoreKeysParams{

		timeout: timeout,
	}
}

// NewDeleteStoreKeysParamsWithContext creates a new DeleteStoreKeysParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteStoreKeysParamsWithContext(ctx context.Context) *DeleteStoreKeysParams {
	var ()
	return &DeleteStoreKeysParams{

		Context: ctx,
	}
}

// NewDeleteStoreKeysParamsWithHTTPClient creates a new DeleteStoreKeysParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteStoreKeysParamsWithHTTPClient(client *http.Client) *DeleteStoreKeysParams {
	var ()
	return &DeleteStoreKeysParams{
		HTTPClient: client,
	}
}

/*DeleteStoreKeysParams contains all the parameters to send to the API endpoint
for the delete store keys operation typically these are written to a http.Request
*/
type DeleteStoreKeysParams struct {

	/*Authorization
	  API key

	*/
	Authorization string
	/*Key
	  Key of the value

	*/
	Key string
	/*Prefix
	  Set it to true if the key is a prefix

	*/
	Prefix *bool

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete store keys params
func (o *DeleteStoreKeysParams) WithTimeout(timeout time.Duration) *DeleteStoreKeysParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete store keys params
func (o *DeleteStoreKeysParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete store keys params
func (o *DeleteStoreKeysParams) WithContext(ctx context.Context) *DeleteStoreKeysParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete store keys params
func (o *DeleteStoreKeysParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete store keys params
func (o *DeleteStoreKeysParams) WithHTTPClient(client *http.Client) *DeleteStoreKeysParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete store keys params
func (o *DeleteStoreKeysParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAuthorization adds the authorization to the delete store keys params
func (o *DeleteStoreKeysParams) WithAuthorization(authorization string) *DeleteStoreKeysParams {
	o.SetAuthorization(authorization)
	return o
}

// SetAuthorization adds the authorization to the delete store keys params
func (o *DeleteStoreKeysParams) SetAuthorization(authorization string) {
	o.Authorization = authorization
}

// WithKey adds the key to the delete store keys params
func (o *DeleteStoreKeysParams) WithKey(key string) *DeleteStoreKeysParams {
	o.SetKey(key)
	return o
}

// SetKey adds the key to the delete store keys params
func (o *DeleteStoreKeysParams) SetKey(key string) {
	o.Key = key
}

// WithPrefix adds the prefix to the delete store keys params
func (o *DeleteStoreKeysParams) WithPrefix(prefix *bool) *DeleteStoreKeysParams {
	o.SetPrefix(prefix)
	return o
}

// SetPrefix adds the prefix to the delete store keys params
func (o *DeleteStoreKeysParams) SetPrefix(prefix *bool) {
	o.Prefix = prefix
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteStoreKeysParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// header param Authorization
	if err := r.SetHeaderParam("Authorization", o.Authorization); err != nil {
		return err
	}

	// path param key
	if err := r.SetPathParam("key", o.Key); err != nil {
		return err
	}

	if o.Prefix != nil {

		// query param prefix
		var qrPrefix bool
		if o.Prefix != nil {
			qrPrefix = *o.Prefix
		}
		qPrefix := swag.FormatBool(qrPrefix)
		if qPrefix != "" {
			if err := r.SetQueryParam("prefix", qPrefix); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}