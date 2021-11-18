package scheduler

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"

	domain "github.com/freundallein/scheduler/pkg"
	"github.com/google/uuid"
)

// Service implements a domain.Scheduler and domain.Worker.
type Service struct {
	taskGateway domain.Gateway

	tasksEnqueued     prometheus.Counter
	taskRequestPolled prometheus.Counter
	tasksClaimed      prometheus.Counter
	tasksSucceeded    prometheus.Counter
	tasksFailed       prometheus.Counter
}

// New returns domain.Scheduler & domain.Worker implementation.
func New(taskGateway domain.Gateway, opts ...Option) *Service {
	svc := &Service{
		taskGateway: taskGateway,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Set allows to enqueue task.
func (svc *Service) Set(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	svc.tasksEnqueued.Inc()
	return svc.taskGateway.Create(ctx, task)
}

// Get allows to poll a task state.
func (svc *Service) Get(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	svc.taskRequestPolled.Inc()
	return svc.taskGateway.FindByID(ctx, id)
}

// Claim gives a task to worker.
func (svc *Service) Claim(ctx context.Context, amount int) ([]*domain.Task, error) {
	tasks, err := svc.taskGateway.ClaimPending(ctx, amount)
	if err != nil {
		return nil, err
	}
	svc.tasksClaimed.Add(float64(len(tasks)))
	return tasks, nil
}

// Succeed marks a task as done.
func (svc *Service) Succeed(ctx context.Context, id, claimID uuid.UUID, result map[string]interface{}) error {
	svc.tasksSucceeded.Inc()
	return svc.taskGateway.MarkAsSucceeded(ctx, id, claimID, result)
}

// Fail marks a task as failed.
func (svc *Service) Fail(ctx context.Context, id, claimID uuid.UUID, reason string) error {
	svc.tasksFailed.Inc()
	return svc.taskGateway.MarkAsFailed(ctx, id, claimID, reason)
}
