package scheduler

import (
	"context"

	domain "github.com/freundallein/scheduler/pkg"
	"github.com/google/uuid"
)

// Service implements a domain.Scheduler and domain.Worker.
type Service struct {
	taskGateway domain.Gateway
}

// New returns domain.Scheduler & domain.Worker implementation.
func New(taskGateway domain.Gateway) *Service {
	return &Service{
		taskGateway: taskGateway,
	}
}

// Set allows to enqueue task.
func (svc *Service) Set(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return svc.taskGateway.Create(ctx, task)
}

// Get allows to poll a task state.
func (svc *Service) Get(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return svc.taskGateway.FindByID(ctx, id)
}

// Claim gives a task to worker.
func (svc *Service) Claim(ctx context.Context, amount int) ([]*domain.Task, error) {
	return svc.taskGateway.ClaimPending(ctx, amount)
}

// Succeed marks a task as done.
func (svc *Service) Succeed(ctx context.Context, id, claimID uuid.UUID, result map[string]interface{}) error {
	return svc.taskGateway.MarkAsSucceeded(ctx, id, claimID, result)
}

// Fail marks a task as failed.
func (svc *Service) Fail(ctx context.Context, id, claimID uuid.UUID, reason string) error {
	return svc.taskGateway.MarkAsFailed(ctx, id, claimID, reason)
}
