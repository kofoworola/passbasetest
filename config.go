package main

import "github.com/kofoworola/passbasetest/storage/postgres"

type Config struct {
	Port     string `default:"3030"`
	Postgres postgres.Config
}
