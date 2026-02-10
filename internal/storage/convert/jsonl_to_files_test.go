package convert

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/steveyegge/fastbeads/internal/types"
)

func TestConvertJSONLToFilesRewritesDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "issues.jsonl")
	outDir := filepath.Join(tmpDir, "issues")

	a := &types.Issue{ID: "bd-a", Title: "A", Status: types.StatusOpen, IssueType: types.TypeTask}
	b := &types.Issue{ID: "bd-b", Title: "B", Status: types.StatusOpen, IssueType: types.TypeTask}
	a.Dependencies = []*types.Dependency{
		{IssueID: "bd-a", DependsOnID: "bd-b", Type: types.DepBlocks},
	}

	data := []byte(mustJSONL(t, a) + mustJSONL(t, b))
	if err := os.WriteFile(jsonlPath, data, 0644); err != nil {
		t.Fatalf("write jsonl: %v", err)
	}

	result, err := ConvertJSONLToFiles(JSONLToFilesOptions{
		JSONLIn:  jsonlPath,
		FilesOut: outDir,
		Backup:   false,
	})
	if err != nil {
		t.Fatalf("ConvertJSONLToFiles: %v", err)
	}
	if result.Rewritten == 0 {
		t.Fatalf("expected dependency rewrite")
	}
}

func mustJSONL(t *testing.T, issue *types.Issue) string {
	t.Helper()
	issue.SetDefaults()
	if err := issue.EnsureIdentity(); err != nil {
		t.Fatalf("EnsureIdentity: %v", err)
	}
	b, err := json.Marshal(issue)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b) + "\n"
}
