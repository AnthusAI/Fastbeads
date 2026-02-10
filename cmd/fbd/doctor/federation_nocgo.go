//go:build !dolt

package doctor

// CheckFederationRemotesAPI returns N/A when Dolt is not enabled.
func CheckFederationRemotesAPI(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Federation remotesapi",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryFederation,
	}
}

// CheckFederationPeerConnectivity returns N/A when Dolt is not enabled.
func CheckFederationPeerConnectivity(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Peer Connectivity",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryFederation,
	}
}

// CheckFederationSyncStaleness returns N/A when Dolt is not enabled.
func CheckFederationSyncStaleness(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Sync Staleness",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryFederation,
	}
}

// CheckFederationConflicts returns N/A when Dolt is not enabled.
func CheckFederationConflicts(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Federation Conflicts",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryFederation,
	}
}

// CheckDoltServerModeMismatch returns N/A when Dolt is not enabled.
func CheckDoltServerModeMismatch(path string) DoctorCheck {
	return DoctorCheck{
		Name:     "Dolt Mode",
		Status:   StatusOK,
		Message:  "N/A (requires dolt build tag)",
		Category: CategoryFederation,
	}
}
