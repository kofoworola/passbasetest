package postgres

import (
	"context"
	"fmt"

	"github.com/kofoworola/passbasetest/storage"
)

const createProjectSQL = `
	INSERT INTO project(
		name,
		api_key
	) VALUES (
		:name,
		:api_key
	) returning *
`

func (s *Storage) CreateProject(ctx context.Context, name, apiKey string) (*storage.Project, error) {
	input := storage.Project{
		Name:   name,
		ApiKey: apiKey,
	}
	var project storage.Project

	stmt, err := s.db.PrepareNamedContext(ctx, createProjectSQL)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}

	if err := stmt.Get(&project, input); err != nil {
		return nil, fmt.Errorf("error inserting into table: %w", err)
	}
	return &project, nil
}
