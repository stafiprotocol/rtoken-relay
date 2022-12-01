package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	// Log a message at the given level with context key/value pairs
	Trace(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}

type log struct {
	entry *logrus.Entry
}

func NewLog(field ...interface{}) Logger {
	return &log{logrus.WithFields(transField(field))}
}

func (l *log) Trace(msg string, ctx ...interface{}) {
	l.entry.WithFields(transField(ctx)).Trace(msg)
}

func (l *log) Debug(msg string, ctx ...interface{}) {
	l.entry.WithFields(transField(ctx)).Debug(msg)
}

func (l *log) Info(msg string, ctx ...interface{}) {
	l.entry.WithFields(transField(ctx)).Info(msg)
}

func (l *log) Warn(msg string, ctx ...interface{}) {
	l.entry.WithFields(transField(ctx)).Warn(msg)
}

func (l *log) Error(msg string, ctx ...interface{}) {
	l.entry.WithFields(transField(ctx)).Error(msg)
}

func transField(datas []interface{}) logrus.Fields {
	field := make(logrus.Fields)
	for i := 0; i < len(datas); i += 2 {
		key := fmt.Sprintf("%v", datas[i])
		if i+1 < len(datas) {
			field[key] = datas[i+1]
		} else {
			field["field"] = datas[i]
		}
	}
	return field
}
