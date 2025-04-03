package db

import (
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// when we use _ blank import, we import the package only for the side effects and not the exported functions
// when the program executes, the init() method of the package executes, and it registers itself with database/sql
// thus registering the postgresql driver. The code will work directly with the sql interface and not pq package

var DB *gorm.DB

const (
	host     = "localhost"
	port     = 5432
	dbname   = "yap_up"
	user     = "postgres"
	password = "cricket360"
)

func ConnectDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, dbname, user, password)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %v", err)
	}

	DB = db
	fmt.Println("successfully connected to db")
	return nil
}
