package postgres

import (
	"database/sql"
	"log/slog"
	"sync"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Postgres struct {
	DbConn *sql.DB // Capitalized to export the field
}

var once = sync.Once{}

// NewPostgres initializes a singleton instance of Postgres and returns it.
// If connection creation fails, it returns an error.
func NewPostgres() (*Postgres, error) {
	var db *sql.DB
	var err error
	once.Do(func() { db, err = createConnection() })
	if err != nil {
		return nil, err
	}
	return &Postgres{
		DbConn: db, // Use the exported field
	}, nil
}

// createConnection establishes a connection to the PostgreSQL database
// and returns the connection object.
func createConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", viper.GetString("postgresURL"))
	if err != nil {
		return nil, err
	}

	// Set connection pool limits (optional tuning)
	db.SetMaxIdleConns(10)   // maximum number of idle connections
	db.SetMaxOpenConns(50)   // maximum number of open connections
	db.SetConnMaxLifetime(0) // connection lifetime (set to 0 for no limit)

	// Ping the database to ensure it's reachable
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	slog.Info("Connected to postgres", "URL", viper.GetString("postgresURL"))
	return db, nil
}

// Close closes the database connection gracefully
func (p *Postgres) Close() error {
	return p.DbConn.Close() // Use the exported field
}
