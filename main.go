package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	log.Info().Msg("Starting")
	err := loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
}
