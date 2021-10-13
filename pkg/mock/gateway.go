package mock

import (
	domain "github.com/freundallein/scheduler/pkg"
)

type Gateway struct {
	CreateFn          func(task *domain.Task) (*domain.Task, error)
	FindByIDFn        func(id string) (*domain.Task, error)
	ClaimPendingFn    func(amount int) ([]*domain.Task, error)
	MarkAsSucceededFn func(id, claimID, result string) error
	MarkAsFailedFn    func(id, claimID, reason string) error
}

func (m *Gateway) Create(task *domain.Task) (*domain.Task, error) {
	if m.CreateFn == nil {
		panic("Gateway.CreateFn is not implemented")
	}
	return m.CreateFn(task)
}
func (m *Gateway) FindByID(id string) (*domain.Task, error) {
	if m.CreateFn == nil {
		panic("Gateway.FindByIDFn is not implemented")
	}
	return m.FindByIDFn(id)
}
func (m *Gateway) ClaimPending(amount int) ([]*domain.Task, error) {
	if m.CreateFn == nil {
		panic("Gateway.ClaimPendingFn is not implemented")
	}
	return m.ClaimPendingFn(amount)
}
func (m *Gateway) MarkAsSucceeded(id, claimID, result string) error {
	if m.CreateFn == nil {
		panic("Gateway.MarkAsSucceededFn is not implemented")
	}
	return m.MarkAsSucceededFn(id, claimID, result)
}
func (m *Gateway) MarkAsFailed(id, claimID, reason string) error {
	if m.CreateFn == nil {
		panic("Gateway.MarkAsFailedFn is not implemented")
	}
	return m.MarkAsFailedFn(id, claimID, reason)
}
