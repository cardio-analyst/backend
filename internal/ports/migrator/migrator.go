package migrator

// Migrator TODO
type Migrator interface {
	// Migrate TODO
	Migrate() error
	// Close TODO
	Close() error
}
