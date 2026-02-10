package syncbranch

import (
	"os/exec"
	"strings"
	"testing"
)

func currentBranch(t *testing.T, dir string) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command git rev-parse --abbrev-ref HEAD failed: %v\nOutput: %s", err, output)
	}
	return strings.TrimSpace(string(output))
}
