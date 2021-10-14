package apiserv

import (
	"context"
	"time"

	domain "github.com/freundallein/scheduler/pkg"
	"github.com/google/uuid"
)

type Scheduler struct {
	sch domain.Scheduler
}

type SetParams struct {
	CorrelationID string `json:"corrID"`
}

// Set accepts task that should be executed.
// curl -H "Content-Type: application/json" -X POST -d '{"jsonrpc": "2.0", "method": "Scheduler.Set", "params":[{"corrID":"123"}], "id": "1"}' http://0.0.0.0:8000/rpc/v0
func (handler *Scheduler) Set(params *SetParams, result *map[string]interface{}) error {
	payload := map[string]interface{}{
		"corrID": params.CorrelationID,
	}
	task := &domain.Task{
		ID:        uuid.New(),
		ExecuteAt: time.Now(),
		Deadline:  time.Now().Add(time.Hour),
		Payload:   payload,
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

type GetParams struct {
	ID string `json:"id"`
}

// Get should be used for task state polling.
// curl -H "Content-Type: application/json" -X POST -d \
//  '{"jsonrpc": "2.0", "method": "Scheduler.Get", "params":[{"id":"123"}], "id": "1"}' \
//  http://0.0.0.0:8000/rpc/v0
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
