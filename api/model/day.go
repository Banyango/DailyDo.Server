package model

import (
	"encoding/json"
	"gopkg.in/guregu/null.v4"
	"time"
)

// Day model
type Day struct {
	ID           string      `json:"id" db:"id"`
	UserID       string      `json:"user" db:"user_id"`
	Summary      null.String `json:"summary" db:"summary"`
	ParentTaskID string      `json:"todoTopLevel" db:"parent_task_id"`
	Date         time.Time   `db:"date"`
}

func (d *Day) MarshalJSON() ([]byte, error) {
	type Alias Day
	return json.Marshal(&struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(d),
		Date: d.Date.Format("Mon Jan _2"),
	})
}