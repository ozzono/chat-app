package models

import (
	"fmt"
	"time"
)

type Session struct {
	Room     string
	Nickname string
}

type Message struct {
	Room      string    `json:"room"       gorm:"room"`
	Nickname  string    `json:"nickname"   gorm:"nickname"  binding:"required"`
	Timestamp time.Time `json:"timestamp"  gorm:"timestamp"`
	Content   string    `json:"content"    gorm:"content"   binding:"required"`
}

type Rooms []Room
type Messages []Message

func ToMap[T Custom](input []T) map[string]T {
	valueMap := make(map[string]T)
	for _, value := range input {
		valueMap[value.GetID()] = value
	}
	return valueMap
}

type Custom interface {
	GetID() string
}

func (r Room) GetID() string {
	return r.ID
}

func (m Message) Fmt() string {
	return fmt.Sprintf("[%s] %s: %s", m.Timestamp.Format("2006-01-02 15:04:05"), m.Nickname, m.Content)
}
