package project

import (
	"context"
	"log"

	"github.com/google/uuid"
	pb "github.com/kofoworola/passbasetest/proto/v1/project"
	"github.com/kofoworola/passbasetest/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedProjectServiceServer

	storage createProjectStorage
}

func New(storage createProjectStorage) *Service {
	return &Service{storage: storage}
}

type createProjectStorage interface {
	CreateProject(context.Context, string, string) (*storage.Project, error)
}

func (s *Service) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	project, err := s.storage.CreateProject(ctx, req.GetName(), uuid.New().String())
	if err != nil {
		log.Printf("error creating project: %v", err)
		return nil, status.Error(codes.Internal, "error creating project")
	}
	return &pb.CreateProjectResponse{
		ApiKey: &project.ApiKey,
		Name:   &project.Name,
	}, nil
}
