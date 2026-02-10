package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/steveyegge/fastbeads/internal/types"
	"gopkg.in/yaml.v3"
)

func TestStoreCreatePersistsYAML(t *testing.T) {
	tmpDir := t.TempDir()
	issuesDir := filepath.Join(tmpDir, "issues")

	store := New(issuesDir)
	if err := store.LoadFromDir(); err != nil {
		t.Fatalf("LoadFromDir error: %v", err)
	}

	issue := &types.Issue{Title: "Test issue", Status: types.StatusOpen, IssueType: types.TypeTask}
	if err := store.CreateIssue(t.Context(), issue, "tester"); err != nil {
		t.Fatalf("CreateIssue error: %v", err)
	}

	path := filepath.Join(issuesDir, issue.ID+issueFileExt)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read issue file: %v", err)
	}
	var got types.Issue
	if err := yaml.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal yaml: %v", err)
	}
	if got.UUID == "" {
		t.Fatalf("expected UUID to be set")
	}
	if got.DisplayID != issue.ID {
		t.Fatalf("expected DisplayID to match ID")
	}
}
