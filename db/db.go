package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const DBName = "database.db"

func InitDB() error {
	var err error

	_, err = os.Stat(DBName)
	if os.IsNotExist(err) {
		fmt.Println("Creating database file...")
		file, err := os.Create(DBName)
		if err != nil {
			return err
		}
		file.Close()
	}

	DB, err = sql.Open("sqlite3", "./"+DBName)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	err = CreateTables()
	if err != nil {
		return err
	}

	return nil
}

func CreateTables() error {
	usersTable := `
  create table if not exists users (
    id integer primary key autoincrement,
    email text not null unique,
    password text not null,
    first_name text,
    last_name text,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp
    );
  `

	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return err
	}

	return nil
}

func CloseDBConn() error {
	if DB != nil {
		return DB.Close()
	}

	return nil
}
