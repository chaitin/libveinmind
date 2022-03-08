package log

import (
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/chaitin/libveinmind/go/plugin/service"
)

// Variables related to the initialization of default logger.
var (
	defaultOnce       sync.Once
	defaultClientCore *clientCore
	defaultLogger     *Logger
	defaultError      error
)

// SetDefaultLogger updates the default logger if not hosted.
//
// The logger provided will be ignored if it is being hosted,
// and it is recommended to invoke the program right at the
// start since it is not synchronized.
func SetDefaultLogger(l *Logger) {
	_ = DefaultLogger()
	if defaultClientCore == nil {
		defaultLogger = l
	}
}

// DefaultLogger returns the default global logger for logging.
//
// When the program is hosted, it returns the logger that can
// communicate with the host. Otherwise it returns the logger
// that is equivalent to logrus.StandardLogger().
func DefaultLogger() *Logger {
	defaultOnce.Do(func() {
		hasService := false
		if service.Hosted() {
			ok, err := service.HasNamespace(Namespace)
			if err != nil {
				defaultError = err
			}
			hasService = ok
		}
		if hasService {
			core, err := newClientCore()
			if err != nil {
				defaultError = err
				return
			}
			defaultClientCore = core
			defaultLogger = New(defaultClientCore)
		} else {
			defaultLogger = NewLogrus(logrus.StandardLogger())
		}
	})
	if defaultError != nil {
		panic(defaultError)
	}
	return defaultLogger
}

// Destroy after synchronizing the default logger initialized.
func Destroy() {
	defaultOnce.Do(func() {
		// Disable further initialization of the client core.
		defaultError = xerrors.New("race with destroy default logger")
	})
	if defaultClientCore != nil {
		defaultClientCore.CloseWait()
	}
}

func Panicf(msg string, obj ...interface{}) {
	DefaultLogger().Panicf(msg, obj...)
}

func Fatalf(msg string, obj ...interface{}) {
	DefaultLogger().Fatalf(msg, obj...)
}

func Errorf(msg string, obj ...interface{}) {
	DefaultLogger().Errorf(msg, obj...)
}

func Warnf(msg string, obj ...interface{}) {
	DefaultLogger().Warnf(msg, obj...)
}

func Warningf(msg string, obj ...interface{}) {
	DefaultLogger().Warningf(msg, obj...)
}

func Infof(msg string, obj ...interface{}) {
	DefaultLogger().Infof(msg, obj...)
}

func Debugf(msg string, obj ...interface{}) {
	DefaultLogger().Debugf(msg, obj...)
}

func Tracef(msg string, obj ...interface{}) {
	DefaultLogger().Tracef(msg, obj...)
}

func Panicln(obj ...interface{}) {
	DefaultLogger().Panicln(obj...)
}

func Fatalln(obj ...interface{}) {
	DefaultLogger().Fatalln(obj...)
}

func Errorln(obj ...interface{}) {
	DefaultLogger().Errorln(obj...)
}

func Warnln(obj ...interface{}) {
	DefaultLogger().Warnln(obj...)
}

func Warningln(obj ...interface{}) {
	DefaultLogger().Warningln(obj...)
}

func Infoln(obj ...interface{}) {
	DefaultLogger().Infoln(obj...)
}

func Debugln(obj ...interface{}) {
	DefaultLogger().Debugln(obj...)
}

func Traceln(obj ...interface{}) {
	DefaultLogger().Traceln(obj...)
}

func Panic(obj ...interface{}) {
	DefaultLogger().Panic(obj...)
}

func Fatal(obj ...interface{}) {
	DefaultLogger().Fatal(obj...)
}

func Error(obj ...interface{}) {
	DefaultLogger().Error(obj...)
}

func Warn(obj ...interface{}) {
	DefaultLogger().Warn(obj...)
}

func Warning(obj ...interface{}) {
	DefaultLogger().Warning(obj...)
}

func Info(obj ...interface{}) {
	DefaultLogger().Info(obj...)
}

func Debug(obj ...interface{}) {
	DefaultLogger().Debug(obj...)
}

func Trace(obj ...interface{}) {
	DefaultLogger().Trace(obj...)
}

func WithFields(fields Fields) *Entry {
	return DefaultLogger().WithFields(fields)
}
