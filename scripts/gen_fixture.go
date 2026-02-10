package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/steveyegge/fastbeads/internal/storage/convert"
	"github.com/steveyegge/fastbeads/internal/types"
	"gopkg.in/yaml.v3"
)

type vocab struct {
	Titles       []string `yaml:"titles"`
	Descriptions []string `yaml:"descriptions"`
	Labels       []string `yaml:"labels"`
}

func main() {
	var (
		seed        = flag.Int64("seed", 42, "random seed")
		count       = flag.Int("count", 5000, "number of issues")
		epicRatio   = flag.Float64("epic-ratio", 0.1, "fraction of issues that are epics")
		depDensity  = flag.Int("dep-density", 2, "average dependencies per issue")
		commentRate = flag.Float64("comment-rate", 0.2, "fraction of issues with comments")
		outDir      = flag.String("out", "/tmp/beads-bench-cache/fixtures", "output directory")
	)
	flag.Parse()

	rng := rand.New(rand.NewSource(*seed))

	vocabPath := filepath.Join("scripts", "fixtures", "vocab.yaml")
	data, err := os.ReadFile(vocabPath)
	if err != nil {
		panic(err)
	}
	var v vocab
	if err := yaml.Unmarshal(data, &v); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		panic(err)
	}
	jsonlPath := filepath.Join(*outDir, fmt.Sprintf("issues-%d.jsonl", *count))
	manifestPath := filepath.Join(*outDir, fmt.Sprintf("issues-%d.manifest.json", *count))
	f, err := os.Create(jsonlPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	now := time.Now().UTC()
	issues := make([]*types.Issue, 0, *count)

	epicCount := int(float64(*count) * *epicRatio)
	for i := 0; i < *count; i++ {
		id := fmt.Sprintf("fbd-%06d", i+1)
		issueType := types.TypeTask
		if i < epicCount {
			issueType = types.TypeEpic
		}
		issue := &types.Issue{
			ID:          id,
			Title:       v.Titles[i%len(v.Titles)],
			Description: v.Descriptions[i%len(v.Descriptions)],
			Status:      types.StatusOpen,
			Priority:    rng.Intn(5),
			IssueType:   issueType,
			CreatedAt:   now.Add(-time.Duration(rng.Intn(720)) * time.Hour),
			UpdatedAt:   now,
			Labels:      []string{v.Labels[i%len(v.Labels)]},
		}
		if err := issue.EnsureIdentity(); err != nil {
			panic(err)
		}
		if rng.Float64() < *commentRate {
			issue.Comments = []*types.Comment{
				{Author: "fixture", Text: "Synthetic comment", CreatedAt: now},
			}
		}
		issues = append(issues, issue)
	}

	// Add dependencies: epics as parents, cross-links between tasks.
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
		for d := 0; d < *depDensity; d++ {
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
			panic(err)
		}
	}

	manifest := convert.BuildManifest(issues)
	if err := convert.WriteManifest(manifestPath, manifest); err != nil {
		panic(err)
	}

	fmt.Printf("Wrote %d issues to %s\n", len(issues), jsonlPath)
	fmt.Printf("Wrote manifest to %s\n", manifestPath)
}
