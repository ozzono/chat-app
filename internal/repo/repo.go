package repo

import (
	"chat-app/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepo(dsn string) (*Repo, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Room{},
		&models.Message{},
	)
	if err != nil {
		return nil, err
	}

	return &Repo{DB: db}, nil
}

func (r *Repo) CreateRoom(name string) error {
	return r.DB.Create(&models.Room{Name: name}).Error
}

func (r *Repo) GetRooms() ([]models.Room, error) {
	var rooms models.Rooms
	err := r.DB.Find(&rooms).Error
	return rooms, err
}

func (r *Repo) CreateMessage(msg models.Message) error {
	return r.DB.Create(&msg).Error
}

func (r *Repo) GetMessages(room string) ([]models.Message, error) {
	var msgs models.Messages
	err := r.DB.Find(&msgs, "room = ?", room).Error
	return msgs, err
}
