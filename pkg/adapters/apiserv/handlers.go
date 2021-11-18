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
	svc domain.Scheduler
}

// SetParams describes input params for Set procedure.
type SetParams struct {
	ID        uuid.UUID `json:"id"`
	ExecuteAt time.Time `json:"executeAt"`
	Deadline  time.Time `json:"deadline"`
	Payload   map[string]interface{}
}

// Set accepts task that should be executed.
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Scheduler.Set", "params":[{"id":"bd954d5e-2b11-49a8-be81-2a53e25a9dc3", "executeAt":"2021-10-14T18:32:11+03:00","deadline":"2021-11-14T18:32:11+03:00","payload": {"type":"parse", "source": "example.com"}}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Scheduler) Set(params *SetParams, result *map[string]interface{}) error {
	task := &domain.Task{
		ID:        params.ID,
		ExecuteAt: params.ExecuteAt.UTC(),
		Deadline:  params.Deadline.UTC(),
		Payload:   params.Payload,
		Meta:      map[string]interface{}{},
	}
	ctx := context.Background()
	task, err := handler.svc.Set(ctx, task)
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
	task, err := handler.svc.Get(ctx, params.ID)
	if err != nil {
		return err
	}
	*result = map[string]interface{}{
		"task": task,
		"meta": task.Meta,
	}
	return nil
}

// Worker is a JSON RPC handler.
type Worker struct {
	svc domain.Worker
}

// ClaimParams describes input params for Claim procedure.
type ClaimParams struct {
	Amount string `json:"amount"`
}

// Claim is for claiming one task or more for processing.
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Worker.Claim", "params":[{"amount":"3"}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Worker) Claim(params *ClaimParams, result *map[string]interface{}) error {
	ctx := context.Background()
	amount, err := strconv.Atoi(params.Amount)
	if err != nil {
		return err
	}
	if amount >= 100 { // Hardcoded batch size
		return fmt.Errorf("amount should be under 100")
	}
	tasks, err := handler.svc.Claim(ctx, amount)
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
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Worker.Succeed", "params":[{"id":"bd954d5e-2b11-49a8-be81-2a53e25a9dc3","claimID":"f5dca270-be27-45aa-ae3a-6e5a600dd965","result": {"data": "job is done"}}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Worker) Succeed(params *SucceedParams, result *map[string]interface{}) error {
	ctx := context.Background()

	err := handler.svc.Succeed(ctx, params.ID, params.ClaimID, params.Result)
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
// curl -X POST -H 'Auth: token' -d '{"jsonrpc": "2.0", "method": "Worker.Fail", "params":[{"id":"bd954d5e-2b11-49a8-be81-2a53e25a9dc3","claimID":"032b8d9e-8f73-4a4e-a850-e2ed716099bf","reason": "there was no one at home"}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Worker) Fail(params *FailParams, result *map[string]interface{}) error {
	if params.Reason == "" {
		return fmt.Errorf("reason should not be empty")
	}
	ctx := context.Background()
	err := handler.svc.Fail(ctx, params.ID, params.ClaimID, params.Reason)
	if err != nil {
		return err
	}
	*result = map[string]interface{}{
		"message": "success",
	}
	return nil
}
