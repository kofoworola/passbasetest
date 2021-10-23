package init

import (
	"context"

	pb "github.com/kofoworola/passbasetest/proto/v1/project"
	"github.com/kofoworola/passbasetest/storage"
)

type Service struct {
	pb.UnimplementedInitServiceServer

	storage createProjectStorage
}

type createProjectStorage interface {
	CreateProject(context.Context, string, string) (*storage.Project, error)
}
