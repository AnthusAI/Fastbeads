//go:build !dolt

package main

import (
	"fmt"
	"os"
)

// handleToDoltMigration is a stub for builds without Dolt support.
func handleToDoltMigration(dryRun bool, autoYes bool) {
	if jsonOutput {
		outputJSON(map[string]interface{}{
			"error":   "dolt_not_available",
			"message": "Dolt backend requires the dolt build tag (and CGO). This binary was built without Dolt support.",
		})
	} else {
		fmt.Fprintf(os.Stderr, "Error: Dolt backend requires the dolt build tag (and CGO)\n")
		fmt.Fprintf(os.Stderr, "This binary was built without Dolt support.\n")
		fmt.Fprintf(os.Stderr, "To use Dolt, rebuild with: CGO_ENABLED=1 go build -tags dolt\n")
	}
	os.Exit(1)
}

// handleToSQLiteMigration is a stub for builds without Dolt support.
func handleToSQLiteMigration(dryRun bool, autoYes bool) {
	if jsonOutput {
		outputJSON(map[string]interface{}{
			"error":   "dolt_not_available",
			"message": "Dolt backend requires the dolt build tag (and CGO). This binary was built without Dolt support.",
		})
	} else {
		fmt.Fprintf(os.Stderr, "Error: Dolt backend requires the dolt build tag (and CGO)\n")
		fmt.Fprintf(os.Stderr, "This binary was built without Dolt support.\n")
		fmt.Fprintf(os.Stderr, "To use Dolt, rebuild with: CGO_ENABLED=1 go build -tags dolt\n")
	}
	os.Exit(1)
}
