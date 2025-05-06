package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LogrusWriter struct {
	Logger *logrus.Logger
	Level  logrus.Level
	Pkg    string
}

func (w *LogrusWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	switch w.Level {
	case logrus.DebugLevel:
		w.Logger.WithField("pkg", w.Pkg).Debug(msg)
	case logrus.InfoLevel:
		w.Logger.WithField("pkg", w.Pkg).Info(msg)
	case logrus.WarnLevel:
		w.Logger.WithField("pkg", w.Pkg).Warn(msg)
	case logrus.ErrorLevel:
		w.Logger.WithField("pkg", w.Pkg).Error(msg)
	default:
		w.Logger.WithField("pkg", w.Pkg).Info(msg)
	}
	return len(p), nil
}

var (
	LoggerInstance = logrus.New()
)

func Instance() *logrus.Logger {
	return LoggerInstance
}

func LoggerWriter(pkg string, level logrus.Level) *LogrusWriter {
	return &LogrusWriter{Logger: LoggerInstance, Level: level, Pkg: pkg}
}

func Logger(c context.Context) *logrus.Entry {
	return LoggerInstance.WithContext(c)
}
