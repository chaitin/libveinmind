package log

import (
	"github.com/sirupsen/logrus"
)

type logrusCore struct {
	l *logrus.Logger
}

func (c *logrusCore) Enabled(l Level) bool {
	return c.l.IsLevelEnabled(logrus.Level(l))
}

func (c *logrusCore) Do(l Log) {
	entry := c.l.WithTime(l.Time)
	if l.Fields != nil {
		entry = entry.WithFields(logrus.Fields(l.Fields))
	}
	entry.Log(logrus.Level(l.Level), l.Msg)
}

// NewLogrus specifies a logrus logger as the underlying Core
// implementation. The plugin log will goes into a logrus
// logger instead of the default logger.
func NewLogrus(l *logrus.Logger) *Logger {
	return New(&logrusCore{l: l})
}
