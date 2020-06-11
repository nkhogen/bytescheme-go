package util

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

const (
	// XRequestIDKey is the request ID key
	XRequestIDKey = ContextKey("x-request-id")
	// XPrincipalKey is the security principal key
	XPrincipalKey = ContextKey("x-principal")
)

// ContextKey is the context key
type ContextKey string

// GetStringContextValue returns string context value
func GetStringContextValue(ctx context.Context, key ContextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

// GetIntContextValue returns int context value
func GetIntContextValue(ctx context.Context, key ContextKey) int {
	value, ok := ctx.Value(key).(int)
	if !ok {
		return 0
	}
	return value
}

// Convert converts from one object to another compatible object.
func Convert(from interface{}, to interface{}) error {
	data, err := ConvertToJSON(from)
	if err != nil {
		return err
	}
	return ConvertFromJSON(data, to)
}

// ConvertFromJSON converts from JSON string to an object which can be a protobuf type.
func ConvertFromJSON(jsonData []byte, to interface{}) error {
	var err error
	toMsg, ok := to.(proto.Message)
	if ok {
		err = jsonpb.UnmarshalString(string(jsonData), toMsg)
	} else {
		err = json.Unmarshal(jsonData, to)
	}
	if err != nil {
		return fmt.Errorf("Unable to convert from JSON to object. Error: %s", err.Error())
	}
	return nil
}

// ConvertToJSON converts an object which can be a protobuf type to a JSON string.
func ConvertToJSON(from interface{}) ([]byte, error) {
	var data []byte
	var err error
	fromMsg, ok := from.(proto.Message)
	if ok {
		marshaller := jsonpb.Marshaler{}
		jstr, err := marshaller.MarshalToString(fromMsg)
		if err != nil {
			return nil, fmt.Errorf("Unable to convert object to JSON. Error: %s", err.Error())
		}
		data = []byte(jstr)
	} else {
		data, err = json.Marshal(from)
		if err != nil {
			return nil, fmt.Errorf("Unable to convert from JSON to object. Error: %s", err.Error())
		}
	}
	return data, nil
}

// ResolveFilepath resolves the relative file path from bytescheme folder
func ResolveFilepath(relFilepath string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "../../", relFilepath)
}
