package repo

import (
	"chat-app/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepo(sqliteDB string) (*Repo, error) {
	db, err := gorm.Open(sqlite.Open(sqliteDB), &gorm.Config{})
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

func (r *Repo) AddRoom(name string) error {
	return r.DB.Create(&models.Room{ID: name}).Error
}

func (r *Repo) GetRooms() ([]models.Room, error) {
	var rooms models.Rooms
	err := r.DB.Find(&rooms).Error
	return rooms, err
}

func (r *Repo) AddMessage(msg models.Message) error {
	return r.DB.Create(&msg).Error
}

func (r *Repo) GetMessages(room string) ([]models.Message, error) {
	var msgs []models.Message
	err := r.DB.Find(&msgs, "room = ?", room).Order("timestamp ASC").Error
	return msgs, err
}
