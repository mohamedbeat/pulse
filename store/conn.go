package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

// Connect establishes a connection to the SQLite database using the provided file path.
// The dbpath parameter should be the path to the SQLite database file.
func Connect(dbpath string) error {
	// SQLite connection string is just the file path
	// You can add parameters like ?cache=shared&mode=rwc for additional options
	connectionString := fmt.Sprintf("file:%s?cache=shared&_fk=1", dbpath)

	// sqlx.MustConnect will panic on error, so we'll use sqlx.Connect instead
	conn, err := sqlx.Connect("sqlite3", connectionString)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Test the connection
	err = conn.Ping()
	if err != nil {
		conn.Close()
		return fmt.Errorf("error pinging the database: %w", err)
	}

	// Set reasonable connection limits for SQLite
	conn.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	conn.SetMaxIdleConns(1)

	db = conn
	return nil
}

// GetDBConnection returns the current database connection.
func GetDBConnection() *sqlx.DB {
	return db
}

// Optional: Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

