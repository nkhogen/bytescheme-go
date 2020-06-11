package service

import (
	"bytescheme/common/auth"
	"fmt"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

const (
	// PermissionProperty is the vendor extension for path permission
	PermissionProperty = "x-permission"
)

func addPathPermission(method, path string, operation *spec.Operation, pathPerms *[]*auth.PathPermission) error {
	if operation == nil || pathPerms == nil {
		return nil
	}
	extensions := operation.VendorExtensible.Extensions
	value, ok := extensions[PermissionProperty]
	if !ok {
		return nil
	}
	perm, ok := value.(string)
	if !ok {
		return fmt.Errorf("Invalid value %+v for %s", value, PermissionProperty)
	}
	pathPerm := &auth.PathPermission{
		Method:     method,
		Path:       path,
		Permission: perm,
	}
	*pathPerms = append(*pathPerms, pathPerm)
	return nil
}

// GetPathPermissions extracts the path permissions from the swagger document
func GetPathPermissions(doc *loads.Document) ([]*auth.PathPermission, error) {
	pathPerms := []*auth.PathPermission{}
	if doc == nil {
		return pathPerms, fmt.Errorf("Invalid swagger document")
	}
	swaggerSpec := doc.Spec()
	if swaggerSpec == nil {
		return pathPerms, fmt.Errorf("Invalid swagger document")
	}
	for path, pItem := range swaggerSpec.SwaggerProps.Paths.Paths {
		pathProps := pItem.PathItemProps
		err := addPathPermission("GET", path, pathProps.Get, &pathPerms)
		if err != nil {
			return []*auth.PathPermission{}, err
		}
		err = addPathPermission("POST", path, pathProps.Post, &pathPerms)
		if err != nil {
			return []*auth.PathPermission{}, err
		}
		err = addPathPermission("PUT", path, pathProps.Put, &pathPerms)
		if err != nil {
			return []*auth.PathPermission{}, err
		}
		err = addPathPermission("DELETE", path, pathProps.Delete, &pathPerms)
		if err != nil {
			return []*auth.PathPermission{}, err
		}
	}
	return pathPerms, nil
}
