package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/kofoworola/passbasetest/storage/postgres"
)

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error loading config: %w", err)
	}

	storage, err := postgres.New(&cfg.Postgres)
	if err != nil {
		log.Fatalf("error loading storage")
	}
}
