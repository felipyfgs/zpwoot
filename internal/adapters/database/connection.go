package database

import (
	"database/sql"
	"fmt"
	"time"

	"zpwoot/internal/adapters/config"
	"zpwoot/internal/adapters/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database wraps sqlx.DB with additional functionality
type Database struct {
	*sqlx.DB
	config *config.Config
	logger *logger.Logger
}

// New creates a new database connection
func New(cfg *config.Config, log *logger.Logger) (*Database, error) {
	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Database connection established successfully")

	return &Database{
		DB:     db,
		config: cfg,
		logger: log,
	}, nil
}

// NewFromAppConfig creates a database connection from app config
func NewFromAppConfig(cfg *config.Config, log *logger.Logger) (*Database, error) {
	return New(cfg, log)
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

// Begin starts a new transaction
func (db *Database) Begin() (*sql.Tx, error) {
	return db.DB.Begin()
}

// Exec executes a query without returning any rows
func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

// Query executes a query that returns rows
func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (db *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

// Health checks database health
func (db *Database) Health() error {
	return db.Ping()
}
