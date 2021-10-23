package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/kofoworola/passbasetest/storage"
)

var ignoreOpt = cmpopts.IgnoreFields(storage.Project{}, "ID", "Created", "Updated")

func TestCreateProject(t *testing.T) {
	expected := &storage.Project{
		Name:   "test",
		ApiKey: uuid.New().String(),
	}
	project, err := _testStorage.CreateProject(context.Background(), expected.Name, expected.ApiKey)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expected, project, ignoreOpt); diff != "" {
		t.Fatal(diff)
	}
}

func TestGetProject(t *testing.T) {
	project, err := _testStorage.CreateProject(context.Background(), "test", uuid.New().String())
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Success", func(t *testing.T) {
		gotten, err := _testStorage.GetProjectByApiKey(context.Background(), project.ApiKey)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(project, gotten); diff != "" {
			t.Fatal(diff)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		_, err := _testStorage.GetProjectByApiKey(context.Background(), "wrong-key")
		if err != sql.ErrNoRows {
			t.Fatalf("expected %v gotten %v", sql.ErrNoRows, err)
		}

	})
}
