package doctor

import (
	"github.com/steveyegge/fastbeads/internal/configfile"
	"github.com/steveyegge/fastbeads/internal/storage/factory"
)

// GetBackend returns the configured backend type from configuration.
// It checks config.yaml first (storage-backend key), then falls back to metadata.json.
// Returns "sqlite" (default) or "dolt".
// hq-3446fc.17: Use factory.GetBackendFromConfig for consistent backend detection.
func GetBackend(beadsDir string) string {
	return factory.GetBackendFromConfig(beadsDir)
}

// IsDoltBackend returns true if the configured backend is Dolt.
func IsDoltBackend(beadsDir string) bool {
	return GetBackend(beadsDir) == configfile.BackendDolt
}
