package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// InitDB initilizes the database and creates the tables if they do not exist
func InitDB() error {

	var err error
	var Conn *sql.DB

	Conn, err = sql.Open("sqlite3", "file:conversations.db?_foreign_keys=on")

	if err != nil {
		return err
	}

	log.Println("running migrations for connections.db")
	sqlBytes, _ := os.ReadFile("conversations.sql")
	_, err = Conn.Exec(string(sqlBytes))

	return err
}

