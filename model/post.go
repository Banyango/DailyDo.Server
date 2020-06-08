package model

import "time"

type Post struct {
	Id     string    `json:"id"`
	Name   string    `json:"name"`
	Url    string    `json:"url"`
	UserId string    `json:"user"`
	Date   time.Time `json:"date"`
}

type PostQuery struct {
	Limit    int
	Offset   int
	Category string
	SortBy   string
}
