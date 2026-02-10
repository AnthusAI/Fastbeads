package convert

import (
	"encoding/json"
	"os"

	"github.com/steveyegge/fastbeads/internal/types"
)

type Manifest struct {
	TotalIssues int            `json:"total_issues"`
	ByStatus    map[string]int `json:"by_status"`
	ByType      map[string]int `json:"by_type"`
	Deps        int            `json:"dependencies"`
	Tombstones  int            `json:"tombstones"`
}

func BuildManifest(issues []*types.Issue) Manifest {
	m := Manifest{
		ByStatus: make(map[string]int),
		ByType:   make(map[string]int),
	}
	for _, issue := range issues {
		m.TotalIssues++
		m.ByStatus[string(issue.Status)]++
		m.ByType[string(issue.IssueType)]++
		if issue.Status == types.StatusTombstone {
			m.Tombstones++
		}
		m.Deps += len(issue.Dependencies)
	}
	return m
}

func WriteManifest(path string, manifest Manifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
