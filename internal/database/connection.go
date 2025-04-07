package database

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rachitnimje/chat-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// when we use _ blank import, we import the package only for the side effects and not the exported functions
// when the program executes, the init() method of the package executes, and it registers itself with database/sql
// thus registering the postgresql driver. The code will work directly with the sql interface and not pq package

var DB *gorm.DB

func ConnectDB(host string, port int, user string, password string, dbname string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, dbname, user, password)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %v", err)
	}

	DB = db
	log.Println("successfully connected to database")

	// migrate the database
	err = models.Migrate(DB)
	if err != nil {
		log.Fatal("error migrating the database: ", err)
	}
	return nil
}
