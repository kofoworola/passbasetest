package main

import (
	"context"
	"database/sql"

	"github.com/kofoworola/passbasetest/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var ignorePaths = []string{
	"/project.ProjectService/CreateProject",
}

type serverInterceptor struct {
	storage projectStorage
}

type projectStorage interface {
	GetProjectByApiKey(context.Context, string) (*storage.Project, error)
}

func (i *serverInterceptor) Interceptor() grpc.UnaryServerInterceptor {
	return grpc.UnaryServerInterceptor(func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		// ignore the ignorables
		found := false
		for _, p := range ignorePaths {
			if info.FullMethod == p {
				found = true
				break
			}
		}
		if found {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "error completing request")
		}

		authHeader, ok := md["authorization"]
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid auth data")
		}

		apiKey := authHeader[0]

		_, err := i.storage.GetProjectByApiKey(ctx, apiKey)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, status.Error(codes.PermissionDenied, "unauthorized")
			}
			return nil, status.Error(codes.Internal, "error completing request")
		}

		return handler(ctx, req)
	})
}
