// Package plugin/log provides a common log system that is
// based on the plugin/service.
//
// Unlike the scenario of those in logrus and zap, the log
// in each plugins should be transferred back to the host.
// And the bottleneck should be the case of IPC that the
// plugin waits for the host's loggerBase result. So there're
// several related control parameter while creating loggers.
package log

import (
	"fmt"
	"time"
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// Fields that will be marshaled and used while logging.
type Fields map[string]interface{}

// Log is a entry created by the logger and is serializable
// between host and plugin.
type Log struct {
	Time   time.Time `json:time`
	Level  Level     `json:level`
	Fields Fields    `json:fields,omitempty`
	Msg    string    `json:msg`
}

// Core is the core object that is sugared by the logger but
// performs the actual loggerBase behaviour.
type Core interface {
	Enabled(level Level) bool
	Do(Log)
}

type logger interface {
	enabled(Level) bool
	log(time.Time, Level, string)
}

type loggerBase struct {
	l logger
}

func (b *loggerBase) logf(n Level, msg string, obj ...interface{}) {
	if b.l.enabled(n) {
		b.l.log(time.Now(), n, fmt.Sprintf(msg, obj...))
	}
}

func (b *loggerBase) Panicf(msg string, obj ...interface{}) {
	b.logf(PanicLevel, msg, obj...)
}

func (b *loggerBase) Fatalf(msg string, obj ...interface{}) {
	b.logf(FatalLevel, msg, obj...)
}

func (b *loggerBase) Errorf(msg string, obj ...interface{}) {
	b.logf(ErrorLevel, msg, obj...)
}

func (b *loggerBase) Warnf(msg string, obj ...interface{}) {
	b.logf(WarnLevel, msg, obj...)
}

func (b *loggerBase) Warningf(msg string, obj ...interface{}) {
	b.logf(WarnLevel, msg, obj...)
}

func (b *loggerBase) Infof(msg string, obj ...interface{}) {
	b.logf(InfoLevel, msg, obj...)
}

func (b *loggerBase) Debugf(msg string, obj ...interface{}) {
	b.logf(DebugLevel, msg, obj...)
}

func (b *loggerBase) Tracef(msg string, obj ...interface{}) {
	b.logf(TraceLevel, msg, obj...)
}

func (b *loggerBase) logln(n Level, obj ...interface{}) {
	if b.l.enabled(n) {
		b.l.log(time.Now(), n, fmt.Sprintln(obj...))
	}
}

func (b *loggerBase) Panicln(obj ...interface{}) {
	b.logln(PanicLevel, obj...)
}

func (b *loggerBase) Fatalln(obj ...interface{}) {
	b.logln(FatalLevel, obj...)
}

func (b *loggerBase) Errorln(obj ...interface{}) {
	b.logln(ErrorLevel, obj...)
}

func (b *loggerBase) Warnln(obj ...interface{}) {
	b.logln(WarnLevel, obj...)
}

func (b *loggerBase) Warningln(obj ...interface{}) {
	b.logln(WarnLevel, obj...)
}

func (b *loggerBase) Infoln(obj ...interface{}) {
	b.logln(InfoLevel, obj...)
}

func (b *loggerBase) Debugln(obj ...interface{}) {
	b.logln(DebugLevel, obj...)
}

func (b *loggerBase) Traceln(obj ...interface{}) {
	b.logln(TraceLevel, obj...)
}

func (b *loggerBase) log(n Level, obj ...interface{}) {
	if b.l.enabled(n) {
		b.l.log(time.Now(), n, fmt.Sprint(obj...))
	}
}

func (b *loggerBase) Panic(obj ...interface{}) {
	b.log(PanicLevel, obj...)
}

func (b *loggerBase) Fatal(obj ...interface{}) {
	b.log(FatalLevel, obj...)
}

func (b *loggerBase) Error(obj ...interface{}) {
	b.log(ErrorLevel, obj...)
}

func (b *loggerBase) Warn(obj ...interface{}) {
	b.log(WarnLevel, obj...)
}

func (b *loggerBase) Warning(obj ...interface{}) {
	b.log(WarnLevel, obj...)
}

func (b *loggerBase) Info(obj ...interface{}) {
	b.log(InfoLevel, obj...)
}

func (b *loggerBase) Debug(obj ...interface{}) {
	b.log(DebugLevel, obj...)
}

func (b *loggerBase) Trace(obj ...interface{}) {
	b.log(TraceLevel, obj...)
}

// Logger is the configured and sugared logger.
type Logger struct {
	loggerBase
	core Core
}

func (l *Logger) enabled(n Level) bool {
	return l.core.Enabled(n)
}

func (l *Logger) log(t time.Time, n Level, msg string) {
	l.core.Do(Log{
		Time:  t,
		Level: n,
		Msg:   msg,
	})
}

func New(core Core) *Logger {
	result := &Logger{
		core: core,
	}
	result.loggerBase.l = result
	return result
}

// Entry is the decorated logging entry.
type Entry struct {
	loggerBase
	l      *Logger
	fields Fields
}

func (e *Entry) enabled(n Level) bool {
	return e.l.core.Enabled(n)
}

func (e *Entry) log(t time.Time, n Level, msg string) {
	e.l.core.Do(Log{
		Time:   t,
		Level:  n,
		Fields: e.fields,
		Msg:    msg,
	})
}

func (l *Logger) WithFields(f Fields) *Entry {
	result := &Entry{
		l:      l,
		fields: f,
	}
	result.loggerBase.l = result
	return result
}
