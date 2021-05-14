package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"

	logger "github.com/sirupsen/logrus"
)

const (
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel Level = iota
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logger.TextFormatter{})

	logger.AddHook(&ContextHook{})
	// Report caller
	//logger.SetReportCaller(true)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(logger.InfoLevel)
}

// Level is the log level
type Level uint32

// ContextHook is a logrus hook
type ContextHook struct{}

// Levels implementation
func (hook *ContextHook) Levels() []logger.Level {
	return logger.AllLevels
}

// Fire implementation
func (hook *ContextHook) Fire(entry *logger.Entry) error {
	if pc, file, line, ok := runtime.Caller(10); ok {
		funcName := runtime.FuncForPC(pc).Name()
		entry.Message = fmt.Sprintf("%s:%v:%s %s", path.Base(file), line, path.Base(funcName), entry.Message)
	}
	return nil
}

func resolvePointers(args []interface{}) []interface{} {
	outArgs := make([]interface{}, 0, len(args))
	for idx := range args {
		arg := args[idx]
		if arg == nil {
			outArgs = append(outArgs, arg)
			continue
		}
		outArgs = append(outArgs, printFriendly(arg))
	}
	return outArgs
}

func printFriendly(i interface{}) interface{} {
	value := reflect.ValueOf(i)
	if value.Type().Kind() == reflect.Struct || value.Type().Kind() == reflect.Map || value.Type().Kind() == reflect.Slice {
		ba, _ := json.Marshal(value.Interface())
		return string(ba)
	}
	if value.Type().Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		value = reflect.Indirect(value)
		return printFriendly(value.Interface())
	}
	return i
}

// IsLevelEnabled checks the given level is enabled
func IsLevelEnabled(level Level) bool {
	return logger.IsLevelEnabled(logger.Level(level))
}

// Debug logs at debug level
func Debug(args ...interface{}) {
	if IsLevelEnabled(DebugLevel) {
		rArgs := resolvePointers(args)
		logger.Debug(rArgs...)
	}
}

// Debugf logs at debug level
func Debugf(format string, args ...interface{}) {
	if IsLevelEnabled(DebugLevel) {
		rArgs := resolvePointers(args)
		logger.Debugf(format, rArgs...)
	}
}

// Error logs at error level
func Error(args ...interface{}) {
	if IsLevelEnabled(ErrorLevel) {
		rArgs := resolvePointers(args)
		logger.Error(rArgs...)
	}
}

// Errorf logs at error level
func Errorf(format string, args ...interface{}) {
	if IsLevelEnabled(ErrorLevel) {
		rArgs := resolvePointers(args)
		logger.Errorf(format, rArgs...)
	}
}

// Info logs at info level
func Info(args ...interface{}) {
	if IsLevelEnabled(InfoLevel) {
		rArgs := resolvePointers(args)
		logger.Info(rArgs...)
	}
}

// Infof logs at info level
func Infof(format string, args ...interface{}) {
	if IsLevelEnabled(InfoLevel) {
		rArgs := resolvePointers(args)
		logger.Infof(format, rArgs...)
	}
}

// Warn logs at warn level
func Warn(args ...interface{}) {
	if IsLevelEnabled(WarnLevel) {
		rArgs := resolvePointers(args)
		logger.Warn(rArgs...)
	}
}

// Warnf logs at warn level
func Warnf(format string, args ...interface{}) {
	if IsLevelEnabled(WarnLevel) {
		rArgs := resolvePointers(args)
		logger.Warnf(format, rArgs...)
	}
}
