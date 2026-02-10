package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/steveyegge/fastbeads/internal/types"
	"gopkg.in/yaml.v3"
)

type JSONLToFilesResult struct {
	Total       int
	Written     int
	Collisions  int
	Skipped     int
	Rewritten   int
	JSONLBackup string
	Manifest    Manifest
}

type JSONLToFilesOptions struct {
	DryRun   bool
	Force    bool
	Backup   bool
	TempDir  string
	JSONLIn  string
	FilesOut string
}

// ConvertJSONLToFiles converts a JSONL file into a directory of YAML files.
func ConvertJSONLToFiles(opts JSONLToFilesOptions) (*JSONLToFilesResult, error) {
	if opts.JSONLIn == "" || opts.FilesOut == "" {
		return nil, fmt.Errorf("missing input or output path")
	}

	in, err := os.Open(opts.JSONLIn)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	decoder := json.NewDecoder(in)
	latest := make(map[string]*types.Issue)
	for {
		var issue types.Issue
		if err := decoder.Decode(&issue); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return nil, err
		}
		issue.SetDefaults()
		if err := issue.EnsureIdentity(); err != nil {
			return nil, err
		}
		latest[issue.ID] = &issue
	}

	issues := make([]*types.Issue, 0, len(latest))
	for _, issue := range latest {
		issues = append(issues, issue)
	}
	sort.Slice(issues, func(i, j int) bool { return issues[i].ID < issues[j].ID })

	result := &JSONLToFilesResult{Total: len(issues)}

	idToUUID := make(map[string]string)
	for _, issue := range issues {
		idToUUID[issue.ID] = issue.UUID
	}

	// Rewrite dependencies to UUIDs where possible.
	for _, issue := range issues {
		changed := false
		for _, dep := range issue.Dependencies {
			if dep.IssueID != "" {
				if uuid, ok := idToUUID[dep.IssueID]; ok && uuid != dep.IssueID {
					dep.IssueID = uuid
					changed = true
				}
			}
			if dep.DependsOnID != "" {
				if uuid, ok := idToUUID[dep.DependsOnID]; ok && uuid != dep.DependsOnID {
					dep.DependsOnID = uuid
					changed = true
				}
			}
		}
		if changed {
			result.Rewritten++
		}
	}
	if opts.DryRun {
		result.Manifest = BuildManifest(issues)
		return result, nil
	}

	if err := os.MkdirAll(filepath.Dir(opts.FilesOut), 0755); err != nil {
		return nil, err
	}

	if fi, err := os.Stat(opts.FilesOut); err == nil && fi.IsDir() {
		entries, _ := os.ReadDir(opts.FilesOut)
		if len(entries) > 0 && !opts.Force {
			return nil, fmt.Errorf("issues dir not empty: %s (use --force to overwrite)", opts.FilesOut)
		}
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	tmpDir := opts.TempDir
	if tmpDir == "" {
		tmpDir = filepath.Join(filepath.Dir(opts.FilesOut), ".issues.tmp")
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, err
	}

	seenNames := make(map[string]bool)
	for _, issue := range issues {
		name := sanitizeFilename(issue.ID) + ".yaml"
		if seenNames[name] {
			result.Collisions++
			name = fmt.Sprintf("%s~%d.yaml", sanitizeFilename(issue.ID), result.Collisions)
		}
		seenNames[name] = true
		path := filepath.Join(tmpDir, name)
		data, err := yaml.Marshal(issue)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
		result.Written++
	}

	if opts.Backup {
		backupPath := opts.JSONLIn + ".bak"
		input, err := os.ReadFile(opts.JSONLIn)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(backupPath, input, 0644); err != nil {
			return nil, err
		}
		result.JSONLBackup = backupPath
	}

	if err := os.RemoveAll(opts.FilesOut); err != nil {
		return nil, err
	}
	if err := os.Rename(tmpDir, opts.FilesOut); err != nil {
		return nil, err
	}
	result.Manifest = BuildManifest(issues)
	return result, nil
}

func sanitizeFilename(s string) string {
	return strings.ReplaceAll(s, string(os.PathSeparator), "_")
}
