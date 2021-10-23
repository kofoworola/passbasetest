package postgres

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kofoworola/passbasetest/storage"
)

func TestCreateProject(t *testing.T) {
	expected := &storage.Project{
		Name:   "test",
		ApiKey: "api-key",
	}
	project, err := _testStorage.CreateProject(context.Background(), expected.Name, expected.ApiKey)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expected, project, cmpopts.IgnoreFields(storage.Project{}, "ID", "Created", "Updated")); diff != "" {
		t.Fatal(diff)
	}
}
