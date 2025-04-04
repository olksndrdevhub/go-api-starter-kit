package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type DBInterface interface {
	Init() error
	Close() error
	GetDB() *sql.DB
}

type DBConfig struct {
	Type     string // "sqlite3" or "postgres"
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // for postgres
	FilePath string // for sqlite
}

type SQLiteDB struct {
	db   *sql.DB
	path string
}

func NewSQLiteDB(config DBConfig) *SQLiteDB {
	return &SQLiteDB{
		path: config.FilePath,
	}
}

func (s *SQLiteDB) Init() error {
	var err error

	s.db, err = sql.Open("sqlite3", s.path)
	if err != nil {
		return err
	}

	err = s.db.Ping()
	if err != nil {
		return err
	}

	return s.CreateTables()
}

func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

func (s *SQLiteDB) GetDB() *sql.DB {
	return s.db
}

func (s *SQLiteDB) CreateTables() error {
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

	_, err := s.db.Exec(usersTable)
	return err
}

type PostgresDB struct {
	db     *sql.DB
	config DBConfig
}

func NewPostgresDB(config DBConfig) *PostgresDB {
	return &PostgresDB{
		config: config,
	}
}

func (p *PostgresDB) Init() error {
	var err error

	connStr := "host=" + p.config.Host + " port=" + p.config.Port + " user=" + p.config.User + " password=" + p.config.Password + " dbname=" + p.config.DBName + " sslmode=" + p.config.SSLMode
	p.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = p.db.Ping()
	if err != nil {
		return err
	}

	return p.CreateTables()
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) GetDB() *sql.DB {
	return p.db
}

func (p *PostgresDB) CreateTables() error {
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

	_, err := p.db.Exec(usersTable)
	return err
}

func NewDB(config DBConfig) error {
	var db DBInterface

	switch config.Type {
	case "sqlite":
		db = NewSQLiteDB(config)
	case "postgres":
		db = NewPostgresDB(config)
	default:
		return fmt.Errorf("unknown database type: %s", config.Type)
	}

	if err := db.Init(); err != nil {
		return err
	}

	DB = db.GetDB()
	return nil
}

func SetDB(db *sql.DB) {
	DB = db
}
