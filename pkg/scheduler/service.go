package scheduler

import (
	"context"

	domain "github.com/freundallein/scheduler/pkg"
)

type Service struct {
	taskGateway domain.Gateway
}

func New(taskGateway domain.Gateway) *Service {
	return &Service{
		taskGateway: taskGateway,
	}
}

func (svc *Service) Set(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return svc.taskGateway.Create(ctx, task)
}

func (svc *Service) Get(ctx context.Context, id string) (*domain.Task, error) {
	return svc.taskGateway.FindByID(ctx, id)
}

func (svc *Service) Issue(ctx context.Context, amount int) ([]*domain.Task, error) {
	return svc.taskGateway.ClaimPending(ctx, amount)
}

func (svc *Service) Succeed(ctx context.Context, id, claimID, result string) error {
	return svc.taskGateway.MarkAsSucceeded(ctx, id, claimID, result)
}

func (svc *Service) Fail(ctx context.Context, id, claimID, reason string) error {
	return svc.taskGateway.MarkAsFailed(ctx, id, claimID, reason)
}
