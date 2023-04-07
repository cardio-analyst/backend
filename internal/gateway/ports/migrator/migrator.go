package migrator

type Migrator interface {
	Migrate() error
	Close() error
}
