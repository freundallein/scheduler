package domain

import "time"

type State int

const (
	Pending    State = 0
	Processing State = 1
	Succeeded  State = 2
	Failed     State = 3
)

type Task struct {
	ID        string
	ClaimID   string
	State     int
	ExecuteAt time.Time
	UpdatedAt time.Time
	Deadline  time.Time
	Payload   map[string]interface{}
	Result    string
}

type Service interface {
	// Public interface
	Set(*Task) error
	Get(id string) (*Task, error)

	// Private interface
	Issue() ([]*Task, error)
	Succeed(id, claimID, result string) error
	Fail(id, claimID, reason string) error
}

type Gateway interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	ClaimPending(amount int) ([]*Task, error)
	MarkAsSucceeded(id, claimID, result string) error
	MarkAsFailed(id, claimID, reason string) error
}
