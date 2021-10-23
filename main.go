package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kelseyhightower/envconfig"
	"github.com/kofoworola/passbasetest/integrations/fixer"

	conversionpb "github.com/kofoworola/passbasetest/proto/v1/conversion"
	projectpb "github.com/kofoworola/passbasetest/proto/v1/project"

	"github.com/kofoworola/passbasetest/services/conversion"
	"github.com/kofoworola/passbasetest/services/project"
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

	go func() {
		if err := setupGrpcGateway(&cfg); err != nil {
			log.Fatal(err)
		}
	}()

	// TODO handle graceful shutdown
	intercaptor := serverInterceptor{storage}
	serv := grpc.NewServer(grpc.UnaryInterceptor(intercaptor.Interceptor()))
	reflection.Register(serv)
	projectpb.RegisterProjectServiceServer(serv, project.New(storage))
	conversionpb.RegisterConversionServiceServer(serv, conversion.New(fixer.New(cfg.Fixer)))
	if err := serv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func setupGrpcGateway(config *Config) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := projectpb.RegisterProjectServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		"localhost:"+config.Port,
		opts,
	); err != nil {
		return err
	}

	if err := conversionpb.RegisterConversionServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		"localhost:"+config.Port,
		opts,
	); err != nil {
		return err
	}

	return http.ListenAndServe(":"+config.RestPort, mux)
}
