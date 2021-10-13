package scheduler

import (
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

func (svc *Service) Set(task *domain.Task) (*domain.Task, error) {
	return svc.taskGateway.Create(task)
}

func (svc *Service) Get(id string) (*domain.Task, error) {
	return svc.taskGateway.FindByID(id)
}

func (svc *Service) Issue(amount int) ([]*domain.Task, error) {
	return svc.taskGateway.ClaimPending(amount)
}

func (svc *Service) Succeed(id, claimID, result string) error {
	return svc.taskGateway.MarkAsSucceeded(id, claimID, result)
}

func (svc *Service) Fail(id, claimID, reason string) error {
	return svc.taskGateway.MarkAsFailed(id, claimID, reason)
}
