package domain

import (
	"time"

	"github.com/google/uuid"
)

type Run struct {
	RunID     uuid.UUID `db:"run_id"`
	AccountID string    `db:"account_id"`
	Initiator string    `db:"initiator"`
	Label     string    `db:"label"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"timestamp"`
}

type RunRepository interface {
	GetRun(id uuid.UUID) (*Run, error)
	GetRuns(id string, limit, offset int) ([]Run, error)
	UpdateRunStatus(r *Run) error
	CreateRun(r *Run) error
}
