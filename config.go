package main

import (
	"github.com/kofoworola/passbasetest/integrations/fixer"
	"github.com/kofoworola/passbasetest/storage/postgres"
)

type Config struct {
	Port     string `default:"8080"`
	RestPort string `default:"8081"`

	Postgres postgres.Config
	Fixer    fixer.Config
}
