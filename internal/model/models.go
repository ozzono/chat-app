package model

import "time"

type Session struct {
	Room     string
	Nickname string
}

type Message struct {
	Room      string
	Nickname  string
	Timestamp time.Time
	Msg       string
}
