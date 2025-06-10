package models

import (
	"gorm.io/gorm"
	"time"
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

//type Room struct {
//	gorm.Model
//	Name        string `gorm:"not null"`
//	Description string
//	CreatedBy   uint
//	Creator     User `gorm:"foreignKey:CreatedBy"`
//}

type Room struct {
	gorm.Model `json:"-"`     // Ignore gorm.Model fields for default JSON marshalling
	ID         uint           `json:"id"`                  // Explicitly map to lowercase 'id'
	CreatedAt  time.Time      `json:"createdAt"`           // Explicitly map to camelCase 'createdAt'
	UpdatedAt  time.Time      `json:"updatedAt"`           // Explicitly map to camelCase 'updatedAt'
	DeletedAt  gorm.DeletedAt `json:"deletedAt,omitempty"` // Explicitly map and handle omitempty

	Name        string `json:"name" gorm:"not null"`          // Explicitly map to lowercase 'name'
	Description string `json:"description"`                   // Explicitly map to lowercase 'description'
	CreatedBy   uint   `json:"createdBy"`                     // Explicitly map to camelCase 'createdBy'
	Creator     User   `gorm:"foreignKey:CreatedBy" json:"-"` // Ignore Creator field for JSON if not needed on frontend
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Message{}, &Room{})
}
