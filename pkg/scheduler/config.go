package scheduler

import "github.com/prometheus/client_golang/prometheus"

// Option is used to configure Service.
type Option func(service *Service)

// WithTasksEnqueued configures Service to use counter metrics.
func WithTasksEnqueued(counter prometheus.Counter) Option {
	return func(s *Service) {
		s.tasksEnqueued = counter
	}
}

// WithTaskRequestPolled configures Service to use counter metrics.
func WithTaskRequestPolled(counter prometheus.Counter) Option {
	return func(s *Service) {
		s.taskRequestPolled = counter
	}
}

// WithTasksClaimed configures Service to use counter metrics.
func WithTasksClaimed(counter prometheus.Counter) Option {
	return func(s *Service) {
		s.tasksClaimed = counter
	}
}

// WithTasksSucceeded configures Service to use counter metrics.
func WithTasksSucceeded(counter prometheus.Counter) Option {
	return func(s *Service) {
		s.tasksSucceeded = counter
	}
}

// WithTasksFailed configures Service to use counter metrics.
func WithTasksFailed(counter prometheus.Counter) Option {
	return func(s *Service) {
		s.tasksFailed = counter
	}
}

// SupervisorOption is used to configure Worker.
type SupervisorOption func(service *Supervisor)

// WithStaleTasksDeleted configures Supervisor to use counter metrics.
func WithStaleTasksDeleted(counter prometheus.Counter) SupervisorOption {
	return func(s *Supervisor) {
		s.staleTasksDeletedCounter = counter
	}
}
