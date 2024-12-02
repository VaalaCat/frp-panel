package watcher

import (
	"context"
	"time"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/go-co-op/gocron/v2"
)

type Client interface {
	Run()
	Stop()
	AddDurationTask(time.Duration, any, ...any) error
	AddCronTask(string, any, ...any) error
}

type client struct {
	s gocron.Scheduler
}

func NewClient() Client {
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Fatalf("create scheduler error")
	}
	return &client{
		s: s,
	}
}

func (c *client) AddDurationTask(duration time.Duration, function any, parameters ...any) error {
	_, err := c.s.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(function, parameters...),
	)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Fatalf("create task error")
	}
	return err
}

func (c *client) AddCronTask(cron string, function any, parameters ...any) error {
	_, err := c.s.NewJob(
		gocron.CronJob(cron, true),
		gocron.NewTask(function, parameters...),
	)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Fatalf("create task error")
	}
	return err
}

func (c *client) Run() {
	ctx := context.Background()
	logger.Logger(ctx).Infof("start to run scheduler, interval: 30s")
	c.s.Start()
}

func (c *client) Stop() {
	if err := c.s.Shutdown(); err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("shutdown scheduler error")
	}
}
