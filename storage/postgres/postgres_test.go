package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/kelseyhightower/envconfig"
)

var _testStorage *Storage

func TestMain(m *testing.M) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error setting up config: %v", err)
	}

	stg, err := New(&cfg)
	if err != nil {
		log.Fatalf("error setting up storage: %v", err)
	}

	_testStorage = stg
	// first drop db
	//	if err := _testStorage.migrate.Drop(); err != nil {
	//		log.Fatalf("error dropping db: %v", err)
	//	}
	if err := _testStorage.Migrate(); err != nil {
		log.Fatalf("error migrating: %v", err)
	}

	code := m.Run()
	if err := _testStorage.migrate.Drop(); err != nil {
		log.Printf("error dropping db: %v", err)
	}

	os.Exit(code)
}
