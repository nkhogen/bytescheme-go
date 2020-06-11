package service

import (
	"net/http"
)

// Authorizer for requests
type Authorizer func(*http.Request, interface{}) error

// Authorize is the implemented method for runtime.Authorizer
func (authorizer Authorizer) Authorize(r *http.Request, i interface{}) error {
	return authorizer(r, i)
}
