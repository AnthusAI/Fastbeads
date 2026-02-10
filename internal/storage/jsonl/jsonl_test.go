package jsonl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/steveyegge/fastbeads/internal/types"
)

func TestStoreCreatePersistsJSONL(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "issues.jsonl")

	store := New(jsonlPath)
	if err := store.LoadFromJSONL(); err != nil {
		t.Fatalf("LoadFromJSONL error: %v", err)
	}

	issue := &types.Issue{Title: "Test issue", Status: types.StatusOpen, IssueType: types.TypeTask}
	if err := store.CreateIssue(t.Context(), issue, "tester"); err != nil {
		t.Fatalf("CreateIssue error: %v", err)
	}

	if _, err := os.Stat(jsonlPath); err != nil {
		t.Fatalf("expected JSONL file to exist: %v", err)
	}
}
