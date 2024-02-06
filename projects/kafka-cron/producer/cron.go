package main

import (
	"time"

	"github.com/robfig/cron/v3"
)

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

func ParseCronSchedule(scheduleSpec string) (cron.Schedule, error) {
	return parser.Parse(scheduleSpec)
}

// Spec represents a scheduled job
type Spec struct {
	Schedule cron.Schedule
	Next     time.Time
	Command  string
}

// Job represents a job to be executed at a specific time
type Job struct {
	Time    time.Time `json:"time"`
	Command string    `json:"command"`
}

type Scheduler struct {
	Specs   []Spec
	JobsDue chan Job
	stop    chan struct{}
}

func NewScheduler(size int) *Scheduler {
	return &Scheduler{
		JobsDue: make(chan Job, size),
	}
}

// Add adds a new spec to the scheduler
func (s *Scheduler) Add(schedule cron.Schedule, command string) error {
	next := schedule.Next(time.Now())
	s.Specs = append(s.Specs, Spec{
		Schedule: schedule,
		Next:     next,
		Command:  command,
	})
	return nil
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	go func() {
		for {
			select {
			case <-s.stop:
				return
			default:
			}

			for i, spec := range s.Specs {
				if spec.Next.Before(time.Now()) {
					s.JobsDue <- Job{
						Time:    spec.Next,
						Command: spec.Command,
					}
					s.Specs[i].Next = spec.Schedule.Next(time.Now())
				}
			}
		}
	}()
}

// Stop stops the scheduler
// but does not stop any jobs that are currently running
func (s *Scheduler) Stop() {
	s.stop <- struct{}{}
}
