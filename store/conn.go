package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// Connect establishes a connection to the PostgreSQL database using the provided configuration.
// It builds a PostgreSQL connection string from the individual parameters.
func Connect(host string, port int, user string, password string, dbname string) error {
	// Build PostgreSQL connection string
	// Format: postgres://user:password@host:port/dbname?sslmode=disable
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	// conn, err := sql.Open("postgres", connectionString)
	// if err != nil {
	// 	return fmt.Errorf("error opening database connection: %w", err)
	// }
	//
	// err = conn.Ping()
	// if err != nil {
	// 	conn.Close()
	// 	return fmt.Errorf("error pinging the database: %w", err)
	// }

	conn := sqlx.MustConnect("postgres", connectionString)

	db = conn
	return nil
}

// GetDBConnection returns the current database connection.
func GetDBConnection() *sqlx.DB {
	return db
}
