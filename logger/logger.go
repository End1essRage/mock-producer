package logger

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"
)

// Global logger entry
var Log = logrus.WithFields(logrus.Fields{
	"serviceName": "mock-producer",
	"logtype":     "logtypeT",
})

type MsgLogger struct {
	fields watermill.LogFields
}

func (l MsgLogger) Error(msg string, err error, fields watermill.LogFields) {
	MapFields(fields).WithError(err).Error(msg)
}
func (l MsgLogger) Info(msg string, fields watermill.LogFields) {
	MapFields(fields).Info(msg)
}
func (l MsgLogger) Debug(msg string, fields watermill.LogFields) {
	MapFields(fields).Debug(msg)
}
func (l MsgLogger) Trace(msg string, fields watermill.LogFields) {
	MapFields(fields).Trace(msg)
}
func (l MsgLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	l.fields = fields
	return l
}

func MapFields(fields watermill.LogFields) *logrus.Entry {
	logrusFields := make(logrus.Fields)
	for k, v := range fields {
		logrusFields[k] = v
	}

	return Log.WithField("caller", "watermill").WithFields(logrusFields)
}
