package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := setupGrpcGateway(ctx, &cfg); err != nil {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	intercaptor := serverInterceptor{storage}
	serv := grpc.NewServer(grpc.UnaryInterceptor(intercaptor.Interceptor()))

	go func() {
		<-done
		fmt.Println("stopping service...")
		cancel()
		serv.GracefulStop()
	}()

	reflection.Register(serv)
	projectpb.RegisterProjectServiceServer(serv, project.New(storage))
	conversionpb.RegisterConversionServiceServer(serv, conversion.New(fixer.New(cfg.Fixer)))
	if err := serv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func setupGrpcGateway(ctx context.Context, config *Config) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := projectpb.RegisterProjectServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:"+config.Port,
		opts,
	); err != nil {
		return err
	}

	if err := conversionpb.RegisterConversionServiceHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:"+config.Port,
		opts,
	); err != nil {
		return err
	}

	return http.ListenAndServe(":"+config.RestPort, mux)
}
