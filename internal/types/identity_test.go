package types

import "testing"

func TestEnsureIdentity(t *testing.T) {
	issue := &Issue{ID: "bd-123"}
	if err := issue.EnsureIdentity(); err != nil {
		t.Fatalf("EnsureIdentity error: %v", err)
	}
	if issue.UUID == "" {
		t.Fatalf("expected UUID to be set")
	}
	if issue.DisplayID != "bd-123" {
		t.Fatalf("expected DisplayID to default to ID, got %q", issue.DisplayID)
	}
}

func TestEnsureIdentityPreservesDisplayID(t *testing.T) {
	issue := &Issue{ID: "bd-123", DisplayID: "bd-123"}
	if err := issue.EnsureIdentity(); err != nil {
		t.Fatalf("EnsureIdentity error: %v", err)
	}
	if issue.DisplayID != "bd-123" {
		t.Fatalf("expected DisplayID to remain unchanged")
	}
}
