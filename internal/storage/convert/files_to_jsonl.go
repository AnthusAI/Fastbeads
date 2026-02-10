package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/steveyegge/fastbeads/internal/types"
	"gopkg.in/yaml.v3"
)

type FilesToJSONLResult struct {
	Total    int
	Written  int
	Manifest Manifest
}

type FilesToJSONLOptions struct {
	DryRun   bool
	Force    bool
	FilesIn  string
	JSONLOut string
}

// ConvertFilesToJSONL exports a directory of YAML issues into a JSONL file.
func ConvertFilesToJSONL(opts FilesToJSONLOptions) (*FilesToJSONLResult, error) {
	if opts.FilesIn == "" || opts.JSONLOut == "" {
		return nil, fmt.Errorf("missing input or output path")
	}

	var issues []*types.Issue
	err := filepath.WalkDir(opts.FilesIn, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".yaml") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var issue types.Issue
		if err := yaml.Unmarshal(data, &issue); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if issue.ID == "" {
			issue.ID = idFromFilename(d.Name())
		}
		issue.SetDefaults()
		if err := issue.EnsureIdentity(); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		issues = append(issues, &issue)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(issues, func(i, j int) bool { return issues[i].ID < issues[j].ID })
	issues = rewriteDepsToUUID(issues)

	result := &FilesToJSONLResult{
		Total:    len(issues),
		Written:  len(issues),
		Manifest: BuildManifest(issues),
	}
	if opts.DryRun {
		return result, nil
	}

	if err := os.MkdirAll(filepath.Dir(opts.JSONLOut), 0755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(opts.JSONLOut); err == nil && !opts.Force {
		return nil, fmt.Errorf("jsonl path already exists: %s (use --force to overwrite)", opts.JSONLOut)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if err := writeJSONLAtomic(opts.JSONLOut, issues); err != nil {
		return nil, err
	}
	return result, nil
}

func idFromFilename(name string) string {
	base := strings.TrimSuffix(name, ".yaml")
	if idx := strings.Index(base, "~"); idx > 0 {
		return base[:idx]
	}
	return base
}

func rewriteDepsToUUID(issues []*types.Issue) []*types.Issue {
	idToUUID := make(map[string]string)
	for _, issue := range issues {
		if issue == nil || issue.ID == "" || issue.UUID == "" {
			continue
		}
		idToUUID[issue.ID] = issue.UUID
	}

	rewritten := make([]*types.Issue, 0, len(issues))
	for _, issue := range issues {
		if issue == nil {
			continue
		}
		issueCopy := *issue
		if len(issue.Dependencies) > 0 {
			deps := make([]*types.Dependency, len(issue.Dependencies))
			for i, dep := range issue.Dependencies {
				if dep == nil {
					continue
				}
				depCopy := *dep
				if uuid, ok := idToUUID[depCopy.IssueID]; ok {
					depCopy.IssueID = uuid
				}
				if uuid, ok := idToUUID[depCopy.DependsOnID]; ok {
					depCopy.DependsOnID = uuid
				}
				deps[i] = &depCopy
			}
			issueCopy.Dependencies = deps
		}
		rewritten = append(rewritten, &issueCopy)
	}
	return rewritten
}

func writeJSONLAtomic(path string, issues []*types.Issue) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	tmpFile, err := os.CreateTemp(dir, base+".tmp.*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if err := tmpFile.Chmod(0644); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return err
	}

	enc := json.NewEncoder(tmpFile)
	for _, issue := range issues {
		if err := enc.Encode(issue); err != nil {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)
			return err
		}
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}
