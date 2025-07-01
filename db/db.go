package db

import (
	"context"
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var Conn *sql.DB //This is a package level variable that holds the database connection and can be used anywhere in the package

// InitDB initilizes the database and creates the tables if they do not exist
func InitDB(ctx context.Context) error {

	var err error

	Conn, err = sql.Open("sqlite3", "file:conversations.db?_foreign_keys=on")

	if err != nil {
		return err
	}

	log.Println("running migrations for connections.db")
	
	sqlBytes, err := os.ReadFile("sql/conversations.sql")
	if err != nil {
		return err               
	}

	_, err = Conn.Exec(string(sqlBytes))

	if err != nil {
		return err
	}

	//prepare all the insert statements once
	if err = PrepareConvoInsertStatement(ctx); err != nil {log.Fatalf("failed to prepare conversation insert statement: %v", err) }
	if err = PrepareProcessedInsertStatement(ctx); err != nil {log.Fatalf("Failed to prepare processed insert statement: %v", err)}

	return err
}

