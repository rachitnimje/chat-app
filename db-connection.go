package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// when we use _ blank import, we import the package only for the side effects and not the exported functions
// when the program executes, the init() method of the package executes, and it registers itself with database/sql
// thus registering the postgresql driver. The code will work directly with the sql interface and not pq package

const (
	host     = "localhost"
	port     = 5432
	dbname   = "yap_up"
	user     = "postgres"
	password = "cricket360"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("error connecting to db: ", err)
	}

	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			fmt.Println("error closing db connection: ", err)
		}
	}(db)

	fmt.Println("connected successfully")

	query := "drop table chats"
	_, err = db.Query(query)
	if err != nil {
		fmt.Println("Error dropping table chats.")
	}
}
