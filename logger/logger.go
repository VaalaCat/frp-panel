package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

func Instance() *logrus.Logger {
	return logger
}

func Logger(c context.Context) *logrus.Entry {
	return logger.WithContext(c)
}
