package model

import "time"

// Post Staging model
type PostStaging struct {
	ID       string    `json:"id"`
	Name     string    `json:"title"`
	URL      string    `json:"url"`
	UserID   string    `json:"user"`
	PostDate time.Time `json:"date"`
	Imported bool      `json:"imported"`
	Quality  int       `json:"quality"`
}
