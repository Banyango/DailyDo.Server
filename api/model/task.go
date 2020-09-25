package model

import (
	"gopkg.in/guregu/null.v4"
)

// Task model
type Task struct {
	ID        string         `json:"id" db:"id"`
	Type      string         `json:"type" db:"type"`
	UserID    string         `json:"user" db:"user_id"`
	Text      null.String `json:"text" db:"text"`
	TaskID    null.String `json:"task" db:"task_id"`
	Completed bool           `json:"completed" db:"completed"`
	Order     int            `json:"order" db:"task_order"`
}

//TaskQuery - task queries
type TaskQuery struct {
	Limit  int
	Offset int
	SortBy string
	Type   string
}