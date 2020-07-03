package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/go-kit/kit/log"

	"github.com/sirupsen/logrus"
)

const (
	//WarnLevel warn leve logger
	WarnLevel = logrus.WarnLevel
	//DebugLevel debug level logger
	DebugLevel = logrus.DebugLevel
	//InfoLevel logger info level
	InfoLevel = logrus.InfoLevel
	//ErrorLevel error level logger
	ErrorLevel = logrus.ErrorLevel
)

var once sync.Once
var l *logrus.Logger
var e *logrus.Entry

var (
	bufferPool *sync.Pool

	// qualified package name, cached at first use
	logrusPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Used for caller information initialisation
	callerInitOnce sync.Once
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

func init() {
	once.Do(func() {
		l = logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				f = getCaller()
				s := strings.Split(f.Function, ".")
				funcname := s[len(s)-1]
				return funcname, fmt.Sprintf("%s:%d", f.File, f.Line)
			},
		})
		l.SetReportCaller(true)
		l.SetOutput(os.Stdout)
		e = logrus.NewEntry(l)
	})
	minimumCallerDepth = 1
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage && pkg != "github.com/sirupsen/logrus" {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

//SetLevel ...
func SetLevel(level logrus.Level) {
	l.SetLevel(level)
	e = logrus.NewEntry(l)
}

//SetOutput set output for logger
func SetOutput(out io.Writer) {
	l.SetOutput(out)
	e = logrus.NewEntry(l)
}

//Info log info
func Info(args ...interface{}) {
	e.Info(args...)
}

//Debug ...
func Debug(args ...interface{}) {
	e.Debug(args...)
}

//Warn ...
func Warn(args ...interface{}) {
	e.Debug(args...)
}

//Error log error
func Error(args ...interface{}) {
	e.Error(args...)
}

type loggerGoKit struct{}

//NewLoggerGokit ...
func NewLoggerGokit() log.Logger {
	return loggerGoKit{}
}

func (l loggerGoKit) Log(keyvals ...interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Error when logging")
		}
	}()
	keyvalMap := make(map[string]interface{})
	for i := 0; i < len(keyvals); i = i + 2 {
		keyvalMap[keyvals[i].(string)] = keyvals[i+1]
	}
	logrus.Debug(keyvalMap)
	return err
}
