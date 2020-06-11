// Code generated by go-swagger; DO NOT EDIT.

package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteStoreKeysHandlerFunc turns a function with the right signature into a delete store keys handler
type DeleteStoreKeysHandlerFunc func(DeleteStoreKeysParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteStoreKeysHandlerFunc) Handle(params DeleteStoreKeysParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// DeleteStoreKeysHandler interface for that can handle valid delete store keys params
type DeleteStoreKeysHandler interface {
	Handle(DeleteStoreKeysParams, interface{}) middleware.Responder
}

// NewDeleteStoreKeys creates a new http.Handler for the delete store keys operation
func NewDeleteStoreKeys(ctx *middleware.Context, handler DeleteStoreKeysHandler) *DeleteStoreKeys {
	return &DeleteStoreKeys{Context: ctx, Handler: handler}
}

/*DeleteStoreKeys swagger:route DELETE /v1/store/keys/{key} Store deleteStoreKeys

Delete a key or keys

Delete a key or keys

*/
type DeleteStoreKeys struct {
	Context *middleware.Context
	Handler DeleteStoreKeysHandler
}

func (o *DeleteStoreKeys) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteStoreKeysParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
