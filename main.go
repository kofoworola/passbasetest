package main

import (
	"errors"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	"github.com/kelseyhightower/envconfig"
	"github.com/kofoworola/passbasetest/integrations/fixer"
	conversionpb "github.com/kofoworola/passbasetest/proto/v1/conversion"
	"github.com/kofoworola/passbasetest/proto/v1/project"
	"github.com/kofoworola/passbasetest/services/conversion"
	"github.com/kofoworola/passbasetest/services/initproject"
	"github.com/kofoworola/passbasetest/storage/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	storage, err := postgres.New(&cfg.Postgres)
	if err != nil {
		log.Fatalf("error loading storage")
	}
	if err := storage.Migrate(); err != nil && errors.Is(migrate.ErrNoChange, err) {
		log.Fatalf("error migrating: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("error starting listener: %v", err)
	}

	// TODO handle graceful shutdown
	serv := grpc.NewServer()
	reflection.Register(serv)
	project.RegisterInitServiceServer(serv, initproject.New(storage))
	conversionpb.RegisterConversionServiceServer(serv, conversion.New(fixer.New(cfg.Fixer)))
	if err := serv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
