package main

import (
	"os"

	"github.com/cristalhq/aconfig"
	"github.com/foxinuni/citas/core"
	"github.com/foxinuni/citas/core/controllers"
	"github.com/foxinuni/citas/core/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var config ApplicationConfig

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	configLoader := aconfig.LoaderFor(&config, aconfig.Config{})
	if err := configLoader.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration!")
	}
}

func main() {
	// Store factory setup
	var storeFactory services.StoreFactory
	if config.GetUseFsStore() {
		storeFactory = services.NewInFsStoreFactory(&config)
	} else {
		storeFactory = services.NewPostgresStoreFactory(&config)
	}

	// Store setup
	citasStore, err := storeFactory.NewCitaStore()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create cita store!")
	}

	// Controllers setup
	citasController := controllers.NewCitaController(citasStore)

	// Web server setup
	webserver := core.NewSistemaCitas(config.GetListenAddress(), citasController)

	// Bootstrap the web server
	log.Info().Msgf("HTTP server is now listening on %s", config.GetListenAddress())
	if err := webserver.Listen(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed to start!")
	}
}
