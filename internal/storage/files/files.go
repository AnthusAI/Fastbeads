// Package files implements a YAML-per-issue storage backend.
package files

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/steveyegge/fastbeads/internal/storage"
	"github.com/steveyegge/fastbeads/internal/storage/memory"
	"github.com/steveyegge/fastbeads/internal/types"
	"gopkg.in/yaml.v3"
)

const (
	issueFileExt = ".yaml"
	lockExt      = ".lock"
)

// Store is a file-backed storage built on the in-memory backend.
// It persists each issue as a YAML file under issuesDir.
type Store struct {
	*memory.MemoryStorage
	issuesDir string
}

// New creates a file-backed store using the provided issues directory.
func New(issuesDir string) *Store {
	return &Store{
		MemoryStorage: memory.New(""),
		issuesDir:     issuesDir,
	}
}

// LoadFromDir loads all YAML issues from issuesDir into memory.
func (s *Store) LoadFromDir() error {
	if err := os.MkdirAll(s.issuesDir, 0755); err != nil {
		return err
	}

	var issues []*types.Issue
	err := filepath.WalkDir(s.issuesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), issueFileExt) {
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
		return err
	}

	// Stable order ensures deterministic counters.
	sort.Slice(issues, func(i, j int) bool { return issues[i].ID < issues[j].ID })
	normalizeDepsFromUUID(issues)
	return s.MemoryStorage.LoadFromIssues(issues)
}

// Path returns the issues directory.
func (s *Store) Path() string {
	return s.issuesDir
}

// CreateIssue creates and persists a new issue.
func (s *Store) CreateIssue(ctx context.Context, issue *types.Issue, actor string) error {
	if err := s.MemoryStorage.CreateIssue(ctx, issue, actor); err != nil {
		return err
	}
	return s.writeIssueFile(issue.ID)
}

// CreateIssues creates and persists multiple issues.
func (s *Store) CreateIssues(ctx context.Context, issues []*types.Issue, actor string) error {
	if err := s.MemoryStorage.CreateIssues(ctx, issues, actor); err != nil {
		return err
	}
	for _, issue := range issues {
		if err := s.writeIssueFile(issue.ID); err != nil {
			return err
		}
	}
	return nil
}

// CreateIssuesWithFullOptions persists each issue after creation.
func (s *Store) CreateIssuesWithFullOptions(ctx context.Context, issues []*types.Issue, actor string, opts storage.BatchCreateOptions) error {
	if err := s.MemoryStorage.CreateIssuesWithFullOptions(ctx, issues, actor, opts); err != nil {
		return err
	}
	for _, issue := range issues {
		if err := s.writeIssueFile(issue.ID); err != nil {
			return err
		}
	}
	return nil
}

// UpdateIssue updates and persists the issue.
func (s *Store) UpdateIssue(ctx context.Context, id string, updates map[string]interface{}, actor string) error {
	if err := s.MemoryStorage.UpdateIssue(ctx, id, updates, actor); err != nil {
		return err
	}
	return s.writeIssueFile(id)
}

// CloseIssue updates and persists the issue.
func (s *Store) CloseIssue(ctx context.Context, id string, reason string, actor string, session string) error {
	if err := s.MemoryStorage.CloseIssue(ctx, id, reason, actor, session); err != nil {
		return err
	}
	return s.writeIssueFile(id)
}

// ClaimIssue updates and persists the issue.
func (s *Store) ClaimIssue(ctx context.Context, id string, actor string) error {
	if err := s.MemoryStorage.ClaimIssue(ctx, id, actor); err != nil {
		return err
	}
	return s.writeIssueFile(id)
}

// DeleteIssue removes issue and file.
func (s *Store) DeleteIssue(ctx context.Context, id string) error {
	if err := s.MemoryStorage.DeleteIssue(ctx, id); err != nil {
		return err
	}
	return s.deleteIssueFile(id)
}

// AddDependency persists after mutation.
func (s *Store) AddDependency(ctx context.Context, dep *types.Dependency, actor string) error {
	if err := s.MemoryStorage.AddDependency(ctx, dep, actor); err != nil {
		return err
	}
	return s.writeIssueFile(dep.IssueID)
}

// RemoveDependency persists after mutation.
func (s *Store) RemoveDependency(ctx context.Context, issueID, dependsOnID string, actor string) error {
	if err := s.MemoryStorage.RemoveDependency(ctx, issueID, dependsOnID, actor); err != nil {
		return err
	}
	return s.writeIssueFile(issueID)
}

// AddLabel persists after mutation.
func (s *Store) AddLabel(ctx context.Context, issueID, label, actor string) error {
	if err := s.MemoryStorage.AddLabel(ctx, issueID, label, actor); err != nil {
		return err
	}
	return s.writeIssueFile(issueID)
}

// RemoveLabel persists after mutation.
func (s *Store) RemoveLabel(ctx context.Context, issueID, label, actor string) error {
	if err := s.MemoryStorage.RemoveLabel(ctx, issueID, label, actor); err != nil {
		return err
	}
	return s.writeIssueFile(issueID)
}

// AddIssueComment persists after mutation.
func (s *Store) AddIssueComment(ctx context.Context, issueID, author, text string) (*types.Comment, error) {
	comment, err := s.MemoryStorage.AddIssueComment(ctx, issueID, author, text)
	if err != nil {
		return nil, err
	}
	if err := s.writeIssueFile(issueID); err != nil {
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
	if err := s.writeIssueFile(issueID); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *Store) writeIssueFile(issueID string) error {
	issue, err := s.issueForWrite(issueID)
	if err != nil {
		return err
	}
	if issue == nil {
		return fmt.Errorf("issue not found: %s", issueID)
	}
	if err := os.MkdirAll(s.issuesDir, 0755); err != nil {
		return err
	}
	path := s.issueFilePath(issueID)
	lock, err := lockFile(path)
	if err != nil {
		return err
	}
	defer func() { _ = lock.Unlock() }()

	data, err := yaml.Marshal(issue)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(s.issuesDir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()
	if err := tmpFile.Chmod(0644); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return err
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

func (s *Store) issueForWrite(issueID string) (*types.Issue, error) {
	issue, err := s.fullIssue(issueID)
	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, fmt.Errorf("issue not found: %s", issueID)
	}

	idToUUID := make(map[string]string)
	for _, existing := range s.MemoryStorage.GetAllIssues() {
		if existing != nil && existing.ID != "" && existing.UUID != "" {
			idToUUID[existing.ID] = existing.UUID
		}
	}
	return rewriteDepsToUUID(issue, idToUUID), nil
}

func (s *Store) deleteIssueFile(issueID string) error {
	path := s.issueFilePath(issueID)
	lock, err := lockFile(path)
	if err != nil {
		return err
	}
	defer func() { _ = lock.Unlock() }()

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func (s *Store) fullIssue(issueID string) (*types.Issue, error) {
	issue, err := s.MemoryStorage.GetIssue(context.Background(), issueID)
	if err != nil || issue == nil {
		return nil, err
	}
	comments, err := s.MemoryStorage.GetIssueComments(context.Background(), issueID)
	if err != nil {
		return nil, err
	}
	issue.Comments = comments
	return issue, nil
}

func (s *Store) issueFilePath(issueID string) string {
	filename := issueFilename(issueID)
	return filepath.Join(s.issuesDir, filename)
}

func issueFilename(issueID string) string {
	return sanitizeFilename(issueID) + issueFileExt
}

func idFromFilename(name string) string {
	base := strings.TrimSuffix(name, issueFileExt)
	if idx := strings.Index(base, "~"); idx > 0 {
		return base[:idx]
	}
	return base
}

func sanitizeFilename(s string) string {
	return strings.ReplaceAll(s, string(os.PathSeparator), "_")
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

func rewriteDepsToUUID(issue *types.Issue, idToUUID map[string]string) *types.Issue {
	if issue == nil {
		return nil
	}
	issueCopy := *issue
	if len(issue.Dependencies) == 0 {
		return &issueCopy
	}

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
	return &issueCopy
}

func lockFile(path string) (*flock.Flock, error) {
	lock := flock.New(path + lockExt)
	if err := lock.Lock(); err != nil {
		return nil, err
	}
	return lock, nil
}
