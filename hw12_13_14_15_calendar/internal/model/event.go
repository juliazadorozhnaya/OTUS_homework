package model

import (
	"time"
)

type Event struct {
	ID           string
	Title        string
	Description  string
	Beginning    time.Time
	Finish       time.Time
	Notification time.Time
	UserID       string
}
