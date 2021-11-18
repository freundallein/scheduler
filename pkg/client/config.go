package client

// SchedulerOption is used to configure Scheduler.
type SchedulerOption func(service *Scheduler)

// WithToken provides scheduler API access token.
func WithToken(token string) SchedulerOption {
	return func(s *Scheduler) {
		s.accessToken = token
	}
}

// WorkerOption is used to configure Worker.
type WorkerOption func(service *Worker)

// WithWorkerToken provides worker API access token.
func WithWorkerToken(token string) WorkerOption {
	return func(s *Worker) {
		s.accessToken = token
	}
}
