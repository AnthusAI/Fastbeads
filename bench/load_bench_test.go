//go:build bench

package bench

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/steveyegge/fastbeads/internal/storage/files"
	"github.com/steveyegge/fastbeads/internal/storage/jsonl"
)

func BenchmarkLoadJSONL(b *testing.B) {
	count := fixtureCount()
	jsonlPath := fixtureJSONLPath(b, count)
	// Warm OS cache to approximate hot load.
	{
		store := jsonl.New(jsonlPath)
		if err := store.LoadFromJSONL(); err != nil {
			b.Fatalf("LoadFromJSONL warm: %v", err)
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store := jsonl.New(jsonlPath)
		if err := store.LoadFromJSONL(); err != nil {
			b.Fatalf("LoadFromJSONL: %v", err)
		}
	}
}

func BenchmarkLoadFiles(b *testing.B) {
	count := fixtureCount()
	issuesDir := fixtureFilesDir(b, count)
	// Warm OS cache to approximate hot load.
	{
		store := files.New(issuesDir)
		if err := store.LoadFromDir(); err != nil {
			b.Fatalf("LoadFromDir warm: %v", err)
		}
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store := files.New(issuesDir)
		if err := store.LoadFromDir(); err != nil {
			b.Fatalf("LoadFromDir: %v", err)
		}
	}
}

func fixtureJSONLPath(b *testing.B, count int) string {
	b.Helper()
	dir := filepath.Join(os.TempDir(), "beads-bench-cache", "fixtures")
	path := filepath.Join(dir, "issues-"+itoa(count)+".jsonl")
	if _, err := os.Stat(path); err != nil {
		b.Fatalf("fixture missing: %s (run scripts/gen_fixture.go)", path)
	}
	return path
}

func fixtureFilesDir(b *testing.B, count int) string {
	b.Helper()
	dir := filepath.Join(os.TempDir(), "beads-bench-cache", "fixtures")
	path := filepath.Join(dir, "issues")
	if _, err := os.Stat(path); err != nil {
		b.Fatalf("fixture dir missing: %s (run migration on JSONL fixture)", path)
	}
	return path
}

func fixtureCount() int {
	if v := os.Getenv("BEADS_FIXTURE_COUNT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 1000
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
