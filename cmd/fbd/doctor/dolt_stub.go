//go:build !dolt

package doctor

import "fmt"

// DoltPerfMetrics is a stub for non-dolt builds.
type DoltPerfMetrics struct{}

// CheckDoltConnection returns N/A when Dolt is not enabled.
func CheckDoltConnection(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Connection",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryCore,
	}
}

// CheckDoltSchema returns N/A when Dolt is not enabled.
func CheckDoltSchema(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Schema",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryCore,
	}
}

// CheckDoltIssueCount returns N/A when Dolt is not enabled.
func CheckDoltIssueCount(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt-JSONL Sync",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryData,
	}
}

// CheckDoltStatus returns N/A when Dolt is not enabled.
func CheckDoltStatus(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Status",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryData,
	}
}

// RunDoltPerformanceDiagnostics returns an error when Dolt is not enabled.
func RunDoltPerformanceDiagnostics(path string, enableProfiling bool) (*DoltPerfMetrics, error) {
	return nil, fmt.Errorf("dolt performance diagnostics require dolt build tag")
}

// PrintDoltPerfReport is a no-op when Dolt is not enabled.
func PrintDoltPerfReport(metrics *DoltPerfMetrics) {
}

// CheckDoltPerformance returns N/A when Dolt is not enabled.
func CheckDoltPerformance(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Performance",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryPerformance,
	}
}

// CompareDoltModes returns an error when Dolt is not enabled.
func CompareDoltModes(path string) error {
	return fmt.Errorf("dolt mode comparison requires dolt build tag")
}
