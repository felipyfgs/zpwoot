package adapters

import (
	"database/sql"
	"fmt"
	"github.com/zpwoot/internal/infra/database"
)

// DatabaseAdapterImpl implements DatabaseAdapter
type DatabaseAdapterImpl struct {
	db             *database.DB
	dataSourceName string
	migrationsPath string
}

// NewDatabaseAdapter creates a new database adapter
func NewDatabaseAdapter(dataSourceName, migrationsPath string) *DatabaseAdapterImpl {
	return &DatabaseAdapterImpl{
		dataSourceName: dataSourceName,
		migrationsPath: migrationsPath,
	}
}

// Connect establishes database connection
func (d *DatabaseAdapterImpl) Connect() error {
	db, err := database.NewDB(d.dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	d.db = db
	return nil
}

// Close closes database connection
func (d *DatabaseAdapterImpl) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Migrate runs database migrations
func (d *DatabaseAdapterImpl) Migrate() error {
	if d.db == nil {
		return fmt.Errorf("database not connected")
	}
	return d.db.RunMigrations(d.migrationsPath)
}

// Health checks database health
func (d *DatabaseAdapterImpl) Health() error {
	if d.db == nil {
		return fmt.Errorf("database not connected")
	}
	return d.db.Ping()
}

// GetDB returns the database instance
func (d *DatabaseAdapterImpl) GetDB() *sql.DB {
	if d.db == nil {
		return nil
	}
	return d.db.DB
}
