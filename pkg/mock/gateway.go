package mock

import (
	"context"

	domain "github.com/freundallein/scheduler/pkg"
)

type Gateway struct {
	CreateFn          func(task *domain.Task) (*domain.Task, error)
	FindByIDFn        func(id string) (*domain.Task, error)
	ClaimPendingFn    func(amount int) ([]*domain.Task, error)
	MarkAsSucceededFn func(id, claimID, result string) error
	MarkAsFailedFn    func(id, claimID, reason string) error
}

func (m *Gateway) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if m.CreateFn == nil {
		panic("Gateway.CreateFn is not implemented")
	}
	return m.CreateFn(task)
}
func (m *Gateway) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	if m.FindByIDFn == nil {
		panic("Gateway.FindByIDFn is not implemented")
	}
	return m.FindByIDFn(id)
}
func (m *Gateway) ClaimPending(ctx context.Context, amount int) ([]*domain.Task, error) {
	if m.ClaimPendingFn == nil {
		panic("Gateway.ClaimPendingFn is not implemented")
	}
	return m.ClaimPendingFn(amount)
}
func (m *Gateway) MarkAsSucceeded(ctx context.Context, id, claimID, result string) error {
	if m.MarkAsSucceededFn == nil {
		panic("Gateway.MarkAsSucceededFn is not implemented")
	}
	return m.MarkAsSucceededFn(id, claimID, result)
}
func (m *Gateway) MarkAsFailed(ctx context.Context, id, claimID, reason string) error {
	if m.MarkAsFailedFn == nil {
		panic("Gateway.MarkAsFailedFn is not implemented")
	}
	return m.MarkAsFailedFn(id, claimID, reason)
}
