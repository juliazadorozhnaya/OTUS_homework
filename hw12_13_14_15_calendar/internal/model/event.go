package model

import (
	"time"
)

type IEvent interface {
	GetID() string
	GetTitle() string
	GetDescription() string
	GetBeginning() time.Time
	GetFinish() time.Time
	GetNotification() time.Time
	GetUserID() string
}

type Event struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Beginning    time.Time `json:"beginning"`
	Finish       time.Time `json:"finish"`
	Notification time.Time `json:"notification"`
	UserID       string    `json:"userId"`
}

func (event *Event) GetID() string {
	return event.ID
}

func (event *Event) GetTitle() string {
	return event.Title
}

func (event *Event) GetDescription() string {
	return event.Description
}

func (event *Event) GetBeginning() time.Time {
	return event.Beginning
}

func (event *Event) GetFinish() time.Time {
	return event.Finish
}

func (event *Event) GetNotification() time.Time {
	return event.Notification
}

func (event *Event) GetUserID() string {
	return event.UserID
}
