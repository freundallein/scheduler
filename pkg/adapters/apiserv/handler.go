package apiserv

import (
	"context"
	"fmt"
	"strconv"
	"time"

	domain "github.com/freundallein/scheduler/pkg"
	"github.com/google/uuid"
)

// Scheduler is a JSON RPC handler.
type Scheduler struct {
	sch domain.Scheduler
}

// SetParams describes input params for Set procedure.
type SetParams struct {
	ID            uuid.UUID `json:"id"`
	CorrelationID string    `json:"corrID"`
	ExecuteAt     time.Time `json:"executeAt"`
	Deadline      time.Time `json:"deadline"`
	Payload       map[string]interface{}
}

// Set accepts task that should be executed.
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Scheduler.Set", "params":[{"corrID":"123","id":"bd954d5e-2b11-49a8-be81-2a53e25a9dc3", "executeAt":"2021-10-14T18:32:11+03:00","deadline":"2021-11-14T18:32:11+03:00","payload": {"type":"parse", "source": "example.com"}}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Scheduler) Set(params *SetParams, result *map[string]interface{}) error {
	task := &domain.Task{
		ID:        params.ID,
		ExecuteAt: params.ExecuteAt,
		Deadline:  params.Deadline,
		Payload:   params.Payload,
		Meta: map[string]interface{}{
			"corrID": params.CorrelationID,
		},
	}
	ctx := context.Background()
	task, err := handler.sch.Set(ctx, task)
	if err != nil {
		return err
	}
	*result = map[string]interface{}{
		"id": task.ID,
	}
	return nil
}

// GetParams describes input params for Get procedure.
type GetParams struct {
	ID uuid.UUID `json:"id"`
}

// Get should be used for task state polling.
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Scheduler.Get", "params":[{"id":"bd954d5e-2b11-49a8-be81-2a53e25a9dc3"}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Scheduler) Get(params *GetParams, result *map[string]interface{}) error {
	ctx := context.Background()
	task, err := handler.sch.Get(ctx, params.ID)
	if err != nil {
		return err
	}
	*result = map[string]interface{}{
		"task": task,
		"meta": task.Meta,
	}
	return nil
}

// ClaimParams describes input params for Claim procedure.
type ClaimParams struct {
	Amount string `json:"amount"`
}

// Claim is for claiming one task or more for processing.
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Scheduler.Claim", "params":[{"amount":"3"}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Scheduler) Claim(params *ClaimParams, result *map[string]interface{}) error {
	ctx := context.Background()
	amount, err := strconv.Atoi(params.Amount)
	if err != nil {
		return err
	}
	if amount >= 100 { // Hardcoded batch size
		return fmt.Errorf("amount should be under 100")
	}
	tasks, err := handler.sch.Claim(ctx, amount)
	if err != nil {
		return err
	}
	*result = map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	}
	return nil
}

// SucceedParams describes input params for Succeed procedure.
type SucceedParams struct {
	ID      uuid.UUID              `json:"id"`
	ClaimID uuid.UUID              `json:"claimID"`
	Result  map[string]interface{} `json:"result"`
}

// Succeed marks task as done.
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Scheduler.Succeed", "params":[{"id":"bd954d5e-2b11-49a8-be81-2a53e25a9dc3","claimID":"f5dca270-be27-45aa-ae3a-6e5a600dd965","result": {"data": "job is done"}}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Scheduler) Succeed(params *SucceedParams, result *map[string]interface{}) error {
	ctx := context.Background()

	err := handler.sch.Succeed(ctx, params.ID, params.ClaimID, params.Result)
	if err != nil {
		return err
	}
	*result = map[string]interface{}{
		"message": "success",
	}
	return nil
}

// FailParams describes input params for Fail procedure.
type FailParams struct {
	ID      uuid.UUID `json:"id"`
	ClaimID uuid.UUID `json:"claimID"`
	Reason  string    `json:"reason"`
}

// Fail marks task as failed.
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Scheduler.Fail", "params":[{"id":"bd954d5e-2b11-49a8-be81-2a53e25a9dc3","claimID":"09cd1033-2e13-4ff4-9e7d-35f4c58359ef","reason": "there was no one at home"}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Scheduler) Fail(params *FailParams, result *map[string]interface{}) error {
	if params.Reason == "" {
		return fmt.Errorf("reason should not be empty")
	}
	ctx := context.Background()
	err := handler.sch.Fail(ctx, params.ID, params.ClaimID, params.Reason)
	if err != nil {
		return err
	}
	*result = map[string]interface{}{
		"message": "success",
	}
	return nil
}
