package models

import (
	"chat-app/pkg/queue"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	Room     string
	Nickname string
}

type Message struct {
	Room      string
	Nickname  string
	Timestamp time.Time
	Content   string
}

type Room struct {
	gorm.Model
	Name      string          `gorm:"unique"`
	TaskQueue chan queue.Task `gorm:"-"`
	Worker    *queue.Worker   `gorm:"-"`
}

type Rooms []Room
type Messages []Message

func ToMap[T Custom](input []T) map[string]T {
	valueMap := make(map[string]T)
	for _, value := range input {
		valueMap[value.GetName()] = value
	}
	return valueMap
}

type Custom interface {
	GetName() string
}

func (r Room) GetName() string {
	return r.Name
}

func (m Message) GetName() string {
	return m.Room
}
