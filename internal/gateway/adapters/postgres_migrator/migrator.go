package postgres_migrator

import (
	"embed"
	"fmt"
	"github.com/cardio-analyst/backend/internal/gateway/ports/migrator"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	log "github.com/sirupsen/logrus"
)

const migrationsDir = "migrations"

// migrations source
//
//go:embed migrations
var embedFS embed.FS

var _ migrator.Migrator = (*PostgresMigrator)(nil)

type PostgresMigrator struct {
	driver migration.Driver
}

func NewPostgresMigrator(dsn string) (*PostgresMigrator, error) {
	dbDriver, err := postgres.New(dsn)
	if err != nil {
		return nil, fmt.Errorf("migrator driver initialization failed: %w", err)
	}

	return &PostgresMigrator{
		driver: dbDriver,
	}, nil
}

func (m *PostgresMigrator) Migrate() error {
	embedSource := &migration.EmbedMigrationSource{
		EmbedFS: embedFS,
		Dir:     migrationsDir,
	}

	// run all up migrations
	applied, err := migration.Migrate(m.driver, embedSource, migration.Up, 0)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Infof("migrations applied: %v", applied)
	return nil
}

func (m *PostgresMigrator) Close() error {
	return m.driver.Close()
}
