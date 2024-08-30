package main

import "github.com/foxinuni/citas/core/services"

var _ services.InFsStoreConfigProvider = (*ApplicationConfig)(nil)
var _ services.PostgresStoreConfigProvider = (*ApplicationConfig)(nil)

type ApplicationConfig struct {
	ListenAddress  string `default:"0.0.0.0:8080" env:"LISTEN_ADDRESS"`
	UseFsStore     bool   `default:"true" env:"USE_FS_STORE"`
	DataPath       string `default:"data" env:"DATA_PATH"`
	MigrationsPath string `default:"migrations" env:"MIGRATIONS_SRC"`
	DatabaseURL    string `env:"DATABASE_URL"`
}

func (s *ApplicationConfig) GetDataPath() string {
	return s.DataPath
}

func (s *ApplicationConfig) GetConnString() string {
	return s.DatabaseURL
}

func (s *ApplicationConfig) GetListenAddress() string {
	return s.ListenAddress
}

func (s *ApplicationConfig) GetUseFsStore() bool {
	return s.UseFsStore
}

func (s *ApplicationConfig) GetMigrationPath() string {
	return s.MigrationsPath
}
