// Package jsonl implements a JSONL file-backed storage backend.
package jsonl

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/steveyegge/fastbeads/internal/storage"
	"github.com/steveyegge/fastbeads/internal/storage/memory"
	"github.com/steveyegge/fastbeads/internal/types"
)

// Store is a JSONL-backed storage built on the in-memory backend.
// It persists by rewriting a snapshot JSONL file after mutations.
type Store struct {
	*memory.MemoryStorage
	jsonlPath string
}

// New creates a JSONL-backed store for the given jsonlPath.
func New(jsonlPath string) *Store {
	return &Store{
		MemoryStorage: memory.New(jsonlPath),
		jsonlPath:     jsonlPath,
	}
}

// LoadFromJSONL loads issues from jsonlPath into memory.
func (s *Store) LoadFromJSONL() error {
	if _, err := os.Stat(s.jsonlPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	f, err := os.Open(s.jsonlPath)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	seen := make(map[string]*types.Issue)
	for {
		var issue types.Issue
		if err := dec.Decode(&issue); err != nil {
			if errors.Is(err, os.ErrClosed) || errors.Is(err, context.Canceled) {
				return err
			}
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return err
		}
		issue.SetDefaults()
		if err := issue.EnsureIdentity(); err != nil {
			return err
		}
		seen[issue.ID] = &issue
	}

	issues := make([]*types.Issue, 0, len(seen))
	for _, issue := range seen {
		issues = append(issues, issue)
	}
	sort.Slice(issues, func(i, j int) bool { return issues[i].ID < issues[j].ID })
	normalizeDepsFromUUID(issues)
	return s.MemoryStorage.LoadFromIssues(issues)
}

// Path returns the JSONL file path.
func (s *Store) Path() string {
	return s.jsonlPath
}

// CreateIssue creates and persists a new issue.
func (s *Store) CreateIssue(ctx context.Context, issue *types.Issue, actor string) error {
	if err := s.MemoryStorage.CreateIssue(ctx, issue, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// CreateIssues creates and persists multiple issues.
func (s *Store) CreateIssues(ctx context.Context, issues []*types.Issue, actor string) error {
	if err := s.MemoryStorage.CreateIssues(ctx, issues, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// CreateIssuesWithFullOptions persists after creation.
func (s *Store) CreateIssuesWithFullOptions(ctx context.Context, issues []*types.Issue, actor string, opts storage.BatchCreateOptions) error {
	if err := s.MemoryStorage.CreateIssuesWithFullOptions(ctx, issues, actor, opts); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// UpdateIssue updates and persists the issue.
func (s *Store) UpdateIssue(ctx context.Context, id string, updates map[string]interface{}, actor string) error {
	if err := s.MemoryStorage.UpdateIssue(ctx, id, updates, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// CloseIssue updates and persists the issue.
func (s *Store) CloseIssue(ctx context.Context, id string, reason string, actor string, session string) error {
	if err := s.MemoryStorage.CloseIssue(ctx, id, reason, actor, session); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// ClaimIssue updates and persists the issue.
func (s *Store) ClaimIssue(ctx context.Context, id string, actor string) error {
	if err := s.MemoryStorage.ClaimIssue(ctx, id, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// DeleteIssue removes issue and persists.
func (s *Store) DeleteIssue(ctx context.Context, id string) error {
	if err := s.MemoryStorage.DeleteIssue(ctx, id); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// AddDependency persists after mutation.
func (s *Store) AddDependency(ctx context.Context, dep *types.Dependency, actor string) error {
	if err := s.MemoryStorage.AddDependency(ctx, dep, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// RemoveDependency persists after mutation.
func (s *Store) RemoveDependency(ctx context.Context, issueID, dependsOnID string, actor string) error {
	if err := s.MemoryStorage.RemoveDependency(ctx, issueID, dependsOnID, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// AddLabel persists after mutation.
func (s *Store) AddLabel(ctx context.Context, issueID, label, actor string) error {
	if err := s.MemoryStorage.AddLabel(ctx, issueID, label, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// RemoveLabel persists after mutation.
func (s *Store) RemoveLabel(ctx context.Context, issueID, label, actor string) error {
	if err := s.MemoryStorage.RemoveLabel(ctx, issueID, label, actor); err != nil {
		return err
	}
	return s.writeSnapshot()
}

// AddIssueComment persists after mutation.
func (s *Store) AddIssueComment(ctx context.Context, issueID, author, text string) (*types.Comment, error) {
	comment, err := s.MemoryStorage.AddIssueComment(ctx, issueID, author, text)
	if err != nil {
		return nil, err
	}
	if err := s.writeSnapshot(); err != nil {
		return nil, err
	}
	return comment, nil
}

// ImportIssueComment persists after mutation.
func (s *Store) ImportIssueComment(ctx context.Context, issueID, author, text string, createdAt time.Time) (*types.Comment, error) {
	comment, err := s.MemoryStorage.ImportIssueComment(ctx, issueID, author, text, createdAt)
	if err != nil {
		return nil, err
	}
	if err := s.writeSnapshot(); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *Store) writeSnapshot() error {
	if err := os.MkdirAll(filepath.Dir(s.jsonlPath), 0755); err != nil {
		return err
	}
	lock := flock.New(s.jsonlPath + ".lock")
	if err := lock.Lock(); err != nil {
		return err
	}
	defer func() { _ = lock.Unlock() }()

	issues := s.MemoryStorage.GetAllIssues()
	return writeJSONLAtomic(s.jsonlPath, rewriteDepsToUUID(issues))
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

func normalizeDepsFromUUID(issues []*types.Issue) {
	uuidToID := make(map[string]string)
	for _, issue := range issues {
		if issue == nil || issue.UUID == "" || issue.ID == "" {
			continue
		}
		uuidToID[issue.UUID] = issue.ID
	}

	for _, issue := range issues {
		if issue == nil {
			continue
		}
		for _, dep := range issue.Dependencies {
			if dep == nil {
				continue
			}
			if id, ok := uuidToID[dep.IssueID]; ok {
				dep.IssueID = id
			}
			if id, ok := uuidToID[dep.DependsOnID]; ok {
				dep.DependsOnID = id
			}
		}
	}
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
