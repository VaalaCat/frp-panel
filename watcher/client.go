package watcher

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
)

type Client interface {
	Run()
	Stop()
	AddTask(time.Duration, any, ...any) error
}

type client struct {
	s gocron.Scheduler
}

func NewClient() Client {
	s, err := gocron.NewScheduler()
	if err != nil {
		logrus.WithError(err).Fatalf("create scheduler error")
	}
	return &client{
		s: s,
	}
}

func (c *client) AddTask(duration time.Duration, function any, parameters ...any) error {
	_, err := c.s.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(function, parameters...),
	)
	if err != nil {
		logrus.WithError(err).Fatalf("create task error")
	}
	return err
}

func (c *client) Run() {
	logrus.Infof("start to run scheduler, interval: 30s")
	c.s.Start()
}

func (c *client) Stop() {
	if err := c.s.Shutdown(); err != nil {
		logrus.WithError(err).Errorf("shutdown scheduler error")
	}
}
