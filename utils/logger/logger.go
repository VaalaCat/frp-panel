package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type LogrusWriter struct {
	Logger *logrus.Logger
	Level  logrus.Level
}

func (w *LogrusWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	switch w.Level {
	case logrus.DebugLevel:
		w.Logger.WithField("pkg", "frp").Debug(msg)
	case logrus.InfoLevel:
		w.Logger.WithField("pkg", "frp").Info(msg)
	case logrus.WarnLevel:
		w.Logger.WithField("pkg", "frp").Warn(msg)
	case logrus.ErrorLevel:
		w.Logger.WithField("pkg", "frp").Error(msg)
	default:
		w.Logger.WithField("pkg", "frp").Info(msg)
	}
	return len(p), nil
}

var (
	logger = &LogrusWriter{
		Logger: logrus.New(),
		Level:  logrus.InfoLevel,
	}
)

func Instance() *logrus.Logger {
	return logger.Logger
}

func Logger(c context.Context) *logrus.Entry {
	return logger.Logger.WithContext(c)
}
