// Code generated by go-swagger; DO NOT EDIT.

package controller

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ListControllersHandlerFunc turns a function with the right signature into a list controllers handler
type ListControllersHandlerFunc func(ListControllersParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ListControllersHandlerFunc) Handle(params ListControllersParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ListControllersHandler interface for that can handle valid list controllers params
type ListControllersHandler interface {
	Handle(ListControllersParams, interface{}) middleware.Responder
}

// NewListControllers creates a new http.Handler for the list controllers operation
func NewListControllers(ctx *middleware.Context, handler ListControllersHandler) *ListControllers {
	return &ListControllers{Context: ctx, Handler: handler}
}

/*ListControllers swagger:route GET /v1/controllers Controller listControllers

List all controllers

List all controllers

*/
type ListControllers struct {
	Context *middleware.Context
	Handler ListControllersHandler
}

func (o *ListControllers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListControllersParams()

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
