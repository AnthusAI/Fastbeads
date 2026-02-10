package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/steveyegge/fastbeads/internal/beads"
	"github.com/steveyegge/fastbeads/internal/config"
	"github.com/steveyegge/fastbeads/internal/debug"
	"github.com/steveyegge/fastbeads/internal/storage"
	"github.com/steveyegge/fastbeads/internal/storage/factory"
	"github.com/steveyegge/fastbeads/internal/storage/files"
	"github.com/steveyegge/fastbeads/internal/storage/jsonl"
	"github.com/steveyegge/fastbeads/internal/storage/memory"
	"github.com/steveyegge/fastbeads/internal/syncbranch"
)

// ensureDirectMode makes sure the CLI is operating in direct-storage mode.
func ensureDirectMode(_ string) error {
	return ensureStoreActive()
}

// fallbackToDirectMode ensures a local store is ready.
// With the daemon removed, this simply ensures the store is active.
func fallbackToDirectMode(_ string) error {
	return ensureStoreActive()
}

// ensureStoreActive guarantees that a storage backend is initialized and tracked.
// Uses the factory to respect metadata.json backend configuration (SQLite, Dolt embedded, or Dolt server).
func ensureStoreActive() error {
	lockStore()
	active := isStoreActive() && getStore() != nil
	unlockStore()
	if active {
		return nil
	}

	// Find the .beads directory
	beadsDir := beads.FindBeadsDir()
	if beadsDir == "" {
		return fmt.Errorf("no beads database found.\n" +
			"Hint: run 'fbd init' to create a database in the current directory,\n" +
			"      or use 'fbd --no-db' for JSONL-only mode")
	}

	// GH#1349: Ensure sync branch worktree exists if configured.
	// This must happen before any JSONL operations to fix fresh clone scenario
	// where findJSONLPath would otherwise fall back to main's stale JSONL.
	if _, err := syncbranch.EnsureWorktree(context.Background()); err != nil {
		// Log warning but don't fail - operations can still work with main's JSONL
		// This allows graceful degradation if worktree creation fails
		debug.Logf("Warning: could not ensure sync worktree: %v", err)
	}

	// Use factory to create the appropriate backend (SQLite, Dolt embedded, or Dolt server)
	// based on metadata.json configuration
	store, err := factory.NewFromConfig(getRootContext(), beadsDir)
	if err != nil {
		return fmt.Errorf("failed to open storage: %w", err)
	}

	// Update the database path for compatibility with code that expects it
	if dbPath := beads.FindDatabasePath(); dbPath != "" {
		setDBPath(dbPath)
	}

	applyConfigToStore(beadsDir, store)

	lockStore()
	setStore(store)
	setStoreActive(true)
	unlockStore()

	if isAutoImportEnabled() && ShouldImportJSONL(rootCtx, store) {
		autoImportIfNewer()
	}

	return nil
}

func applyConfigToStore(beadsDir string, store storage.Storage) {
	if store == nil {
		return
	}

	// Apply custom types/statuses from config.yaml for file/jsonl backends.
	if customTypes := config.GetCustomTypesFromYAML(); len(customTypes) > 0 {
		_ = store.SetConfig(rootCtx, "types.custom", strings.Join(customTypes, ","))
	}
	if customStatuses := config.GetCustomStatusesFromYAML(); len(customStatuses) > 0 {
		_ = store.SetConfig(rootCtx, "status.custom", strings.Join(customStatuses, ","))
	}

	prefix := config.GetString("issue-prefix")
	if prefix == "" {
		if mem := memoryFromStore(store); mem != nil {
			detected, err := detectPrefix(beadsDir, mem)
			if err == nil {
				prefix = detected
			}
		}
	}
	if prefix != "" {
		_ = store.SetConfig(rootCtx, "issue_prefix", prefix)
	}
}

func memoryFromStore(store storage.Storage) *memory.MemoryStorage {
	switch s := store.(type) {
	case *files.Store:
		return s.MemoryStorage
	case *jsonl.Store:
		return s.MemoryStorage
	default:
		return nil
	}
}
