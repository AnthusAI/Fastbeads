package storage_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/steveyegge/fastbeads/internal/storage/convert"
	"github.com/steveyegge/fastbeads/internal/storage/files"
	"github.com/steveyegge/fastbeads/internal/storage/jsonl"
	"github.com/steveyegge/fastbeads/internal/types"
)

func BenchmarkLoadJSONL(b *testing.B) {
	for _, count := range []int{1000, 5000, 10000} {
		b.Run(fmt.Sprintf("issues=%d", count), func(b *testing.B) {
			jsonlPath, _ := prepareFixtureFiles(b, count)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				store := jsonl.New(jsonlPath)
				if err := store.LoadFromJSONL(); err != nil {
					b.Fatalf("LoadFromJSONL failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkLoadFiles(b *testing.B) {
	for _, count := range []int{1000, 5000, 10000} {
		b.Run(fmt.Sprintf("issues=%d", count), func(b *testing.B) {
			_, issuesDir := prepareFixtureFiles(b, count)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				store := files.New(issuesDir)
				if err := store.LoadFromDir(); err != nil {
					b.Fatalf("LoadFromDir failed: %v", err)
				}
			}
		})
	}
}

func prepareFixtureFiles(b *testing.B, count int) (string, string) {
	b.Helper()

	dir := b.TempDir()
	jsonlPath := filepath.Join(dir, "issues.jsonl")
	issuesDir := filepath.Join(dir, "issues")

	generateJSONLFixture(b, jsonlPath, count)
	if _, err := convert.ConvertJSONLToFiles(convert.JSONLToFilesOptions{
		JSONLIn:  jsonlPath,
		FilesOut: issuesDir,
		Force:    true,
		Backup:   false,
	}); err != nil {
		b.Fatalf("ConvertJSONLToFiles failed: %v", err)
	}

	return jsonlPath, issuesDir
}

func generateJSONLFixture(b *testing.B, jsonlPath string, count int) {
	b.Helper()

	f, err := os.Create(jsonlPath)
	if err != nil {
		b.Fatalf("create JSONL: %v", err)
	}
	defer f.Close()

	rng := rand.New(rand.NewSource(42))
	now := time.Now().UTC()

	epicCount := int(float64(count) * 0.1)
	if epicCount < 1 {
		epicCount = 1
	}

	issues := make([]*types.Issue, 0, count)
	for i := 0; i < count; i++ {
		issueType := types.TypeTask
		if i < epicCount {
			issueType = types.TypeEpic
		}
		issue := &types.Issue{
			ID:          fmt.Sprintf("fbd-%06d", i+1),
			Title:       fmt.Sprintf("Fixture issue %d", i+1),
			Description: "Synthetic fixture data for load benchmarks.",
			Status:      types.StatusOpen,
			Priority:    rng.Intn(5),
			IssueType:   issueType,
			CreatedAt:   now.Add(-time.Duration(rng.Intn(720)) * time.Hour),
			UpdatedAt:   now,
			Labels:      []string{"benchmark"},
		}
		if err := issue.EnsureIdentity(); err != nil {
			b.Fatalf("ensure identity: %v", err)
		}
		if rng.Float64() < 0.2 {
			issue.Comments = []*types.Comment{
				{Author: "fixture", Text: "Synthetic comment", CreatedAt: now},
			}
		}
		issues = append(issues, issue)
	}

	for i := epicCount; i < len(issues); i++ {
		issue := issues[i]
		parent := issues[rng.Intn(epicCount)]
		issue.Dependencies = append(issue.Dependencies, &types.Dependency{
			IssueID:     issue.ID,
			DependsOnID: parent.ID,
			Type:        types.DepParentChild,
			CreatedAt:   now,
			CreatedBy:   "fixture",
		})
		for d := 0; d < 2; d++ {
			target := issues[rng.Intn(len(issues))].ID
			if target == issue.ID {
				continue
			}
			issue.Dependencies = append(issue.Dependencies, &types.Dependency{
				IssueID:     issue.ID,
				DependsOnID: target,
				Type:        types.DepBlocks,
				CreatedAt:   now,
				CreatedBy:   "fixture",
			})
		}
	}

	enc := json.NewEncoder(f)
	for _, issue := range issues {
		if err := enc.Encode(issue); err != nil {
			b.Fatalf("encode JSONL: %v", err)
		}
	}
}
