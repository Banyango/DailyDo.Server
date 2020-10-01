package model

import (
	"gopkg.in/guregu/null.v4"
	"time"
)

// Day model
type Day struct {
	ID           string      `json:"id" db:"id"`
	UserID       string      `json:"user" db:"user_id"`
	Summary      null.String `json:"summary" db:"summary"`
	ParentTaskID string      `json:"todoTopLevel" db:"parent_task_id"`
	Date         time.Time   `json:"date" db:"date"`
}
