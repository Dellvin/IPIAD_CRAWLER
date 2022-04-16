package model

import "time"

type News struct {
	ID          string
	Code        int
	Header      string
	Body        string
	PublishedAt time.Time
	Author      string
	Link        string
	Error       error
}
