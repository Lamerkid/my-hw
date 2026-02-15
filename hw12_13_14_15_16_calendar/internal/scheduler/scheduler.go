package scheduler

import (
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
)

type Scheduler struct {
	app      *app.App
	producer rmq.Producer
	logger   Logger
	interval time.Duration
}

func NewScheduler(logg Logger, producer rmq.Producer, app *app.App, interval time.Duration) *Scheduler {
	return &Scheduler{
		app:      app,
		producer: producer,
		logger:   logg,
		interval: interval,
	}
}
