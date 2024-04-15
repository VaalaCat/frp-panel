package watcher

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
)

type Client interface {
	Run()
	Stop()
}

type client struct {
	s gocron.Scheduler
}

func NewClient(f func(clientID, clientSecret string) error, clientID, clientSecret string) Client {
	s, err := gocron.NewScheduler()
	if err != nil {
		logrus.WithError(err).Fatalf("create scheduler error")
	}

	_, err = s.NewJob(
		gocron.DurationJob(time.Second*30),
		gocron.NewTask(f, clientID, clientSecret),
	)
	if err != nil {
		logrus.WithError(err).Fatalf("create job error")
	}
	return &client{
		s: s,
	}
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
