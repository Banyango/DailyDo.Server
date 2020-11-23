package model

import (
	"gopkg.in/guregu/null.v4"
)

// Parent model
type Task struct {
	ID        string      `json:"id" db:"id"`
	Type      string      `json:"type" db:"discriminator"`
	UserID    string      `json:"user" db:"user_id"`
	Text      null.String `json:"text" db:"text"`
	TaskID    null.String `json:"parent" db:"task_id"`
	Completed bool        `json:"completed" db:"completed"`
	Order     null.String `json:"order" db:"task_order"`
}

//TaskQuery - task queries
type TaskQuery struct {
	Limit  int
	Offset int
	SortBy string
	Type   string
}
