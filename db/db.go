package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // for postgres
}

func InitDB(config DBConfig) error {
	var err error

	connCtr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	DB, err = sql.Open("postgres", connCtr)

	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	return CreateTables()
}

func Close() error {
	return DB.Close()
}

func CreateTables() error {
	usersTable := `
  CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        first_name TEXT,
        last_name TEXT
    );
  `

	_, err := DB.Exec(usersTable)
	return err
}
