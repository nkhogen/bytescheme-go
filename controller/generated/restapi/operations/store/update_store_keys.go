// Code generated by go-swagger; DO NOT EDIT.

package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateStoreKeysHandlerFunc turns a function with the right signature into a update store keys handler
type UpdateStoreKeysHandlerFunc func(UpdateStoreKeysParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateStoreKeysHandlerFunc) Handle(params UpdateStoreKeysParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// UpdateStoreKeysHandler interface for that can handle valid update store keys params
type UpdateStoreKeysHandler interface {
	Handle(UpdateStoreKeysParams, interface{}) middleware.Responder
}

// NewUpdateStoreKeys creates a new http.Handler for the update store keys operation
func NewUpdateStoreKeys(ctx *middleware.Context, handler UpdateStoreKeysHandler) *UpdateStoreKeys {
	return &UpdateStoreKeys{Context: ctx, Handler: handler}
}

/*UpdateStoreKeys swagger:route PUT /v1/store/keys Store updateStoreKeys

Save a key value pair

Save a key value pair

*/
type UpdateStoreKeys struct {
	Context *middleware.Context
	Handler UpdateStoreKeysHandler
}

func (o *UpdateStoreKeys) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateStoreKeysParams()

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
