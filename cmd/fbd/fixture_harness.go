package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/steveyegge/fastbeads/internal/storage/convert"
)

var fixtureHarnessCmd = &cobra.Command{
	Use:   "fixture-harness",
	Short: "Validate JSONL->files conversion against fixture manifests",
	Run: func(cmd *cobra.Command, _ []string) {
		count, _ := cmd.Flags().GetInt("count")
		outDir, _ := cmd.Flags().GetString("out")

		jsonlPath := filepath.Join(outDir, fmt.Sprintf("issues-%d.jsonl", count))
		manifestPath := filepath.Join(outDir, fmt.Sprintf("issues-%d.manifest.json", count))
		filesDir := filepath.Join(outDir, "issues")

		if _, err := os.Stat(jsonlPath); err != nil {
			FatalError(fmt.Sprintf("missing fixture: %s", jsonlPath))
		}
		if _, err := os.Stat(manifestPath); err != nil {
			FatalError(fmt.Sprintf("missing manifest: %s", manifestPath))
		}

		data, err := os.ReadFile(manifestPath)
		if err != nil {
			FatalError(err.Error())
		}
		var expected convert.Manifest
		if err := json.Unmarshal(data, &expected); err != nil {
			FatalError(err.Error())
		}

		result, err := convert.ConvertJSONLToFiles(convert.JSONLToFilesOptions{
			JSONLIn:  jsonlPath,
			FilesOut: filesDir,
			Force:    true,
			Backup:   false,
		})
		if err != nil {
			FatalError(err.Error())
		}

		if !manifestsEqual(expected, result.Manifest) {
			FatalError("manifest mismatch after conversion")
		}

		fmt.Printf("Fixture harness ok: %d issues\n", result.Manifest.TotalIssues)
	},
}

func manifestsEqual(a, b convert.Manifest) bool {
	if a.TotalIssues != b.TotalIssues || a.Deps != b.Deps || a.Tombstones != b.Tombstones {
		return false
	}
	if len(a.ByStatus) != len(b.ByStatus) || len(a.ByType) != len(b.ByType) {
		return false
	}
	for k, v := range a.ByStatus {
		if b.ByStatus[k] != v {
			return false
		}
	}
	for k, v := range a.ByType {
		if b.ByType[k] != v {
			return false
		}
	}
	return true
}

func init() {
	fixtureHarnessCmd.Flags().Int("count", 1000, "Fixture issue count")
	fixtureHarnessCmd.Flags().String("out", filepath.Join(os.TempDir(), "beads-bench-cache", "fixtures"), "Fixture directory")
	rootCmd.AddCommand(fixtureHarnessCmd)
}
