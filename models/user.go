package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"unique"`
	Name      string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
