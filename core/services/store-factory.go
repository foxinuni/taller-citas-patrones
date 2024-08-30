package services

import (
	"context"

	"github.com/foxinuni/citas/core/stores"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

/* StoreFactory - abstract factory to create new stores */
type StoreFactory interface {
	NewCitaStore() (stores.CitaStore, error)
}

/* InFsStoreConfigProvider - interface to provide configuration for InFsStoreFactory */
type InFsStoreConfigProvider interface {
	GetDataPath() string
}

type InFsStoreFactory struct {
	config InFsStoreConfigProvider
}

func NewInFsStoreFactory(config InFsStoreConfigProvider) StoreFactory {
	log.Info().Msgf("Using in-fs store (data path: %q)", config.GetDataPath())
	return &InFsStoreFactory{config: config}
}

func (f *InFsStoreFactory) NewCitaStore() (stores.CitaStore, error) {
	return stores.NewInFsCitaStore(f.config.GetDataPath())
}

/* PostgresStoreConfigProvider - interface to provide configuration for PostgresStoreFactory */
type PostgresStoreConfigProvider interface {
	GetConnString() string
	GetMigrationPath() string
}

type PostgresStoreFactory struct {
	config PostgresStoreConfigProvider
}

func NewPostgresStoreFactory(config PostgresStoreConfigProvider) StoreFactory {
	log.Info().Msgf("Using postgres store (database url: %q, migrations: %q)", config.GetConnString(), config.GetMigrationPath())
	return &PostgresStoreFactory{config: config}
}

func (f *PostgresStoreFactory) NewCitaStore() (stores.CitaStore, error) {
	pool, err := pgxpool.Connect(context.Background(), f.config.GetConnString())
	if err != nil {
		return nil, err
	}

	// perform migrations
	migrations, err := migrate.New(f.config.GetMigrationPath(), f.config.GetConnString())
	if err != nil {
		return nil, err
	}
	defer migrations.Close()

	log.Info().Msg("Applying migrations to the database...")
	if err := migrations.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return nil, err
		}

		log.Info().Msg("Database is up to date! No migrations to apply.")
	}

	return stores.NewPostgresCitaStore(pool), nil
}
