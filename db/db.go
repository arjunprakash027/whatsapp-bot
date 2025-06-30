package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var Conn *sql.DB //This is a package level variable that holds the database connection and can be used anywhere in the package

// InitDB initilizes the database and creates the tables if they do not exist
func InitDB() error {

	var err error

	Conn, err = sql.Open("sqlite3", "file:conversations.db?_foreign_keys=on")

	if err != nil {
		return err
	}

	log.Println("running migrations for connections.db")
	
	sqlBytes, err := os.ReadFile("db/conversations.sql")
	if err != nil {
		return err               
	}

	_, err = Conn.Exec(string(sqlBytes))

	return err
}

