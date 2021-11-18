package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// State describes task states.
type State string

const (
	// StatePending means, that task was created and is waiting for processing.
	StatePending State = "pending"
	// StateProcessing means, that task is processing by a worker.
	StateProcessing State = "processing"
	// StateSucceeded means, that task was successfully processed.
	StateSucceeded State = "succeeded"
	// StateFailed means, that we got failure during processing.
	StateFailed State = "failed"
)

// Task describes a work unit.
type Task struct {
	// ID is a task identifier.
	ID uuid.UUID `json:"id"`
	// ClaimID used for worker identification and result linearisation.
	ClaimID *uuid.UUID `json:"claimId,omitempty"`
	// State describes current task's state.
	State State `json:"state"`
	// ExecuteAt allows scheduler to define when to execute a task.
	ExecuteAt time.Time `json:"executeAt"`
	// Deadline  allows scheduler to define when task becomes stale.
	Deadline time.Time `json:"deadline"`
	// Payload describes the task itself.
	Payload map[string]interface{} `json:"payload"`
	// Result shows the result of a task processing.
	Result map[string]interface{} `json:"result,omitempty"`
	// Meta used for service information.
	Meta map[string]interface{} `json:"-"`
	// CreatedAt shows when task was created.
	CreatedAt time.Time `json:"createdAt"`
	// DoneAt shows when task was succeeded.
	DoneAt sql.NullTime `json:"doneAt,omitempty"`
}

// Scheduler used for task planning and polling.
type Scheduler interface {
	// Set allows to enqueue task.
	Set(ctx context.Context, task *Task) (*Task, error)
	// Get allows to poll a task state.
	Get(ctx context.Context, id uuid.UUID) (*Task, error)
}

// Worker used for task processing.
type Worker interface {
	// Claim gives a task to worker.
	Claim(ctx context.Context, amount int) ([]*Task, error)
	// Succeed marks a task as done.
	Succeed(ctx context.Context, id, claimID uuid.UUID, result map[string]interface{}) error
	// Fail marks a task as failed.
	Fail(ctx context.Context, id, claimID uuid.UUID, reason string) error
}

// Supervisor is used for storage maintenance.
type Supervisor interface {
	// DeleteStaleTasks cleans storage from stale tasks.
	DeleteStaleTasks(ctx context.Context, staleHours int) error
	// TODO: delete or move tasks with N attempts
}

// Gateway describes database access to a task.
type Gateway interface {
	// Create makes record with new task.
	Create(ctx context.Context, task *Task) (*Task, error)
	// FindByID allows to poll a task state.
	FindByID(ctx context.Context, id uuid.UUID) (*Task, error)
	// ClaimPending used for locking tasks.
	ClaimPending(ctx context.Context, amount int) ([]*Task, error)
	// MarkAsSucceeded marks a task as successfully processed.
	MarkAsSucceeded(ctx context.Context, id, claimID uuid.UUID, result map[string]interface{}) error
	// MarkAsFailed marks a task as failed.
	MarkAsFailed(ctx context.Context, id, claimID uuid.UUID, reason string) error
	// DeleteStaleTasks removes stale tasks.
	DeleteStaleTasks(ctx context.Context, staleHours int) (int64, error)
}
