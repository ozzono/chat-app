package repo

import (
	"context"

	"chat-app/pkg/queue"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name      string          `gorm:"unique"`
	TaskQueue chan queue.Task `gorm:"-"`
	Worker    *queue.Worker   `gorm:"-"`
}

type Repo struct {
	DB *gorm.DB
}

func NewRepo(dsn string) (*Repo, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Room{})
	if err != nil {
		return nil, err
	}

	return &Repo{DB: db}, nil
}

func (r *Repo) CreateRoom(name string) error {
	room := Room{Name: name, TaskQueue: make(chan queue.Task)}
	room.Worker = queue.NewWorker(name)
	room.Worker.TaskQueue = room.TaskQueue
	go room.Worker.StartWorker(context.Background()) // Pass a valid context
	return r.DB.Create(&room).Error
}

func (r *Repo) GetRooms() ([]Room, error) {
	var rooms []Room
	err := r.DB.Find(&rooms).Error
	return rooms, err
}
