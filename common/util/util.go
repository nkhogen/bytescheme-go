package util

import (
	"bytescheme/common/log"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
)

const (
	// XRequestIDKey is the request ID key
	XRequestIDKey = ContextKey("x-request-id")
	// XPrincipalKey is the security principal key
	XPrincipalKey = ContextKey("x-principal")
)

var (
	// ShutdownHandler is the shutdown hook
	ShutdownHandler *shutDownHandler
)

type shutDownHandler struct {
	lock       *sync.Mutex
	closeables map[uintptr]Closeable
	ctx        context.Context
	cancel     context.CancelFunc
}

// Closeable has a Close method
type Closeable interface {
	Close() error
}

// ContextKey is the context key
type ContextKey string

func init() {
	ShutdownHandler = newShutDownHandler()
}

func newShutDownHandler() *shutDownHandler {
	handler := &shutDownHandler{
		lock:       &sync.Mutex{},
		closeables: map[uintptr]Closeable{},
	}
	handler.ctx, handler.cancel = context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		handler.cancel()
		handler.lock.Lock()
		defer handler.lock.Unlock()
		for _, closeable := range handler.closeables {
			closeable.Close()
		}
	}()
	return handler
}

// Context returns the cancellable context in the shutdown handler
func (handler *shutDownHandler) Context() context.Context {
	return handler.ctx
}

// RegisterCloseable registers closeables for shut down hook
func (handler *shutDownHandler) RegisterCloseable(closeable Closeable) error {
	if closeable == nil {
		return fmt.Errorf("Invalid closeable")
	}
	value := reflect.ValueOf(closeable)
	handler.lock.Lock()
	defer handler.lock.Unlock()
	handler.closeables[value.Pointer()] = closeable
	return nil
}

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

// MonitorJob invokes a monitoring at a regular interval until timeout
func MonitorJob(ctx context.Context, checkInterval time.Duration, timeout time.Duration, job func(context.Context) (bool, error)) error {
	timer := time.After(timeout)
	ticker := time.Tick(checkInterval)
	done := make(chan bool, 1)
	isContinue, err := job(ctx)
	if err != nil {
		log.Errorf("Error in executing job")
	}
	if !isContinue {
		return err
	}
	done <- true
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer:
			return fmt.Errorf("Timed out")
		case <-ticker:
			select {
			case <-done:
				isContinue, err = job(ctx)
				done <- true
				if err != nil {
					log.Errorf("Error in executing job")
				}
				if !isContinue {
					return err
				}
			}
		}
	}
}


func GetUUID() string {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")
	return id
}