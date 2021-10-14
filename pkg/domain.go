package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type State string

const (
	StatePending    State = "pending"
	StateProcessing State = "processing"
	StateSucceeded  State = "succeeded"
	StateFailed     State = "failed"
)

type Task struct {
	ID        uuid.UUID              `json:"id"`
	ClaimID   uuid.UUID              `json:"claimId,omitempty"`
	State     State                  `json:"state"`
	ExecuteAt time.Time              `json:"executeAt"`
	Deadline  time.Time              `json:"deadline"`
	Payload   map[string]interface{} `json:"payload"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Meta      map[string]interface{} `json:"-"`
}

type Scheduler interface {
	// Public interface
	Set(ctx context.Context, task *Task) (*Task, error)
	Get(ctx context.Context, id string) (*Task, error)

	// Private interface
	Issue(ctx context.Context, amount int) ([]*Task, error)
	Succeed(ctx context.Context, id, claimID, result string) error
	Fail(ctx context.Context, id, claimID, reason string) error
}

type Gateway interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	FindByID(ctx context.Context, id string) (*Task, error)
	ClaimPending(ctx context.Context, amount int) ([]*Task, error)
	MarkAsSucceeded(ctx context.Context, id, claimID, result string) error
	MarkAsFailed(ctx context.Context, id, claimID, reason string) error
}
