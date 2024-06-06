package model

import (
	"time"
)

// IEvent интерфейс для структуры Event, предоставляющий методы доступа к полям.
type IEvent interface {
	GetID() string
	GetTitle() string
	GetDescription() string
	GetBeginning() time.Time
	GetFinish() time.Time
	GetNotification() time.Time
	GetUserID() string
}

// Event структура, представляющая событие.
type Event struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Beginning    time.Time `json:"beginning"`
	Finish       time.Time `json:"finish"`
	Notification time.Time `json:"notification"`
	UserID       string    `json:"userId"`
}

// GetID возвращает ID события.
func (event *Event) GetID() string {
	return event.ID
}

// GetTitle возвращает заголовок события.
func (event *Event) GetTitle() string {
	return event.Title
}

// GetDescription возвращает описание события.
func (event *Event) GetDescription() string {
	return event.Description
}

// GetBeginning возвращает время начала события.
func (event *Event) GetBeginning() time.Time {
	return event.Beginning
}

// GetFinish возвращает время окончания события.
func (event *Event) GetFinish() time.Time {
	return event.Finish
}

// GetNotification возвращает время уведомления о событии.
func (event *Event) GetNotification() time.Time {
	return event.Notification
}

// GetUserID возвращает ID пользователя, связанного с событием.
func (event *Event) GetUserID() string {
	return event.UserID
}
