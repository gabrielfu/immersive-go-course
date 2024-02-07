package kron

import (
	"time"

	"github.com/robfig/cron/v3"
)

var parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

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

func (j Job) String() string {
	return `Job<"` + j.Command + `" @ ` + j.Time.String() + `>`
}

type Scheduler struct {
	JobsDue chan Job
	specs   []*Spec
	stop    chan struct{}
}

func NewScheduler(size int) *Scheduler {
	return &Scheduler{
		JobsDue: make(chan Job, size),
		specs:   []*Spec{},
		stop:    make(chan struct{}),
	}
}

// Add adds a new spec to the scheduler
func (s *Scheduler) Add(schedule cron.Schedule, command string) error {
	next := schedule.Next(time.Now())
	spec := &Spec{
		Schedule: schedule,
		Next:     next,
		Command:  command,
	}
	s.specs = append(s.specs, spec)
	go s.start(spec)
	return nil
}

func (s *Scheduler) start(spec *Spec) {
	ticker := time.NewTicker(time.Until(spec.Next))
	for {
		select {
		case <-s.stop:
			ticker.Stop()
			return
		default:
		}

		<-ticker.C
		s.JobsDue <- Job{
			Time:    spec.Next,
			Command: spec.Command,
		}
		spec.Next = spec.Schedule.Next(time.Now())
		ticker.Reset(time.Until(spec.Next))
	}
}

// Stop stops the scheduler
// but does not stop any jobs that are currently running
func (s *Scheduler) Stop() {
	s.stop <- struct{}{}
}
