//go:build !dolt

package doctor

// MigrationValidationResult provides machine-parseable migration validation output.
// This stub exists when Dolt support is not enabled.
type MigrationValidationResult struct {
	Phase          string   `json:"phase"`
	Ready          bool     `json:"ready"`
	Backend        string   `json:"backend"`
	JSONLCount     int      `json:"jsonl_count"`
	SQLiteCount    int      `json:"sqlite_count"`
	DoltCount      int      `json:"dolt_count"`
	MissingInDB    []string `json:"missing_in_db"`
	MissingInJSONL []string `json:"missing_in_jsonl"`
	Errors         []string `json:"errors"`
	Warnings       []string `json:"warnings"`
	JSONLValid     bool     `json:"jsonl_valid"`
	JSONLMalformed int      `json:"jsonl_malformed"`
	DoltHealthy    bool     `json:"dolt_healthy"`
	DoltLocked     bool     `json:"dolt_locked"`
	SchemaValid    bool     `json:"schema_valid"`
	RecommendedFix string   `json:"recommended_fix"`
}

// CheckMigrationReadiness returns N/A when Dolt is not enabled.
func CheckMigrationReadiness(path string) (DoctorCheck, MigrationValidationResult) {
	return DoctorCheck{
			Name:     "Migration Readiness",
			Status:   StatusOK,
			Message:  "N/A (requires dolt build tag)",
			Category: CategoryMaintenance,
		}, MigrationValidationResult{
			Phase:   "pre-migration",
			Ready:   false,
			Backend: "unknown",
			Errors:  []string{"Dolt migration requires dolt build tag"},
		}
}

// CheckMigrationCompletion returns N/A when Dolt is not enabled.
func CheckMigrationCompletion(path string) (DoctorCheck, MigrationValidationResult) {
	return DoctorCheck{
			Name:     "Migration Completion",
			Status:   StatusOK,
			Message:  "N/A (requires dolt build tag)",
			Category: CategoryMaintenance,
		}, MigrationValidationResult{
			Phase:   "post-migration",
			Ready:   false,
			Backend: "unknown",
			Errors:  []string{"Dolt migration requires dolt build tag"},
		}
}

// CheckDoltLocks returns N/A when Dolt is not enabled.
func CheckDoltLocks(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Locks",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryMaintenance,
	}
}
