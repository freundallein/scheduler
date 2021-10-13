package domain

import (
	"time"
)

type State string

const (
	StatePending    State = "pending"
	StateProcessing State = "processing"
	StateSucceeded  State = "succeeded"
	StateFailed     State = "failed"
)

type Task struct {
	ID        string                 `json:"id"`
	ClaimID   string                 `json:"-"`
	State     State                  `json:"state"`
	ExecuteAt time.Time              `json:"executeAt"`
	UpdatedAt time.Time              `json:"-"`
	Deadline  time.Time              `json:"deadline"`
	Payload   map[string]interface{} `json:"payload"`
	Meta      map[string]interface{} `json:"-"`
	Result    string                 `json:"result,omitempty"`
}

type Scheduler interface {
	// Public interface
	Set(*Task) (*Task, error)
	Get(id string) (*Task, error)

	// Private interface
	Issue(amount int) ([]*Task, error)
	Succeed(id, claimID, result string) error
	Fail(id, claimID, reason string) error
}

type Gateway interface {
	Create(task *Task) (*Task, error)
	FindByID(id string) (*Task, error)
	ClaimPending(amount int) ([]*Task, error)
	MarkAsSucceeded(id, claimID, result string) error
	MarkAsFailed(id, claimID, reason string) error
}
