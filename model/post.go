package model

import "time"

// Post model
type Post struct {
	ID       string    `json:"id"`
	Name     string    `json:"title"`
	URL      string    `json:"url"`
	UserID   string    `json:"user"`
	PostDate time.Time `json:"date"`
}

//PostQuery - post queries
type PostQuery struct {
	Limit    int
	Offset   int
	Category string
	SortBy   string
}
