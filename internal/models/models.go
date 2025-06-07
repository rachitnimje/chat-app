package models

import (
	"gorm.io/gorm"
)

type User struct {
	//Id        uint      `gorm:"primaryKey;autoIncrement"`
	gorm.Model
	Username string `gorm:"unique"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
	//CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Message struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	RoomID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID"`
	Room    Room   `gorm:"foreignKey:RoomID"`
}

type Room struct {
	gorm.Model
	Name        string `json:"name"`
	Description string
	CreatedBy   uint
	Creator     User `gorm:"foreignKey:CreatedBy"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Message{}, &Room{})
}
