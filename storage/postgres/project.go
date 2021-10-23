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

func (s *Storage) GetProjectByApiKey(ctx context.Context, apiKey string) (*storage.Project, error) {
	stmt := `SELECT * FROM project WHERE api_key = $1`
	var project storage.Project
	if err := s.db.GetContext(ctx, &project, stmt, apiKey); err != nil {
		return nil, err
	}
	return &project, nil

}
