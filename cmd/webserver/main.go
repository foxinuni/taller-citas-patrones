//go:build wireinject
// +build wireinject

package main

import (
	"os"

	"github.com/foxinuni/citas/core"
	"github.com/foxinuni/citas/core/controllers"
	"github.com/foxinuni/citas/core/stores"
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	address core.ListenAddress = "localhost:8080"
	relpath stores.FsPath      = "data"
)

var webserverSet = wire.NewSet(core.NewSistemaCitas, controllers.NewCitaController, stores.NewInFsCitaStore)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func main() {
	webserver, err := bootstrapRESP(address, relpath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to bootstrap HTTP server!")
	}

	log.Info().Msgf("HTTP server is now listening on %s", address)
	if err := webserver.Listen(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed to start!")
	}
}

func bootstrapRESP(_ core.ListenAddress, _ stores.FsPath) (*core.SistemaCitas, error) {
	wire.Build(webserverSet)
	return &core.SistemaCitas{}, nil
}
