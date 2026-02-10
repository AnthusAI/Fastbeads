package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/steveyegge/fastbeads/internal/beads"
	"github.com/steveyegge/fastbeads/internal/storage/convert"
	"gopkg.in/yaml.v3"
)

var migrateStorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Migrate storage format (jsonl <-> files)",
	Run: func(cmd *cobra.Command, _ []string) {
		to, _ := cmd.Flags().GetString("to")
		from, _ := cmd.Flags().GetString("from")
		dest, _ := cmd.Flags().GetString("dest")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")

		beadsDir := beads.FindBeadsDir()
		if beadsDir == "" {
			FatalErrorWithHint("no .beads directory found", "run 'fbd init' first")
		}

		switch to {
		case "files":
			if from == "" {
				from = filepath.Join(beadsDir, "issues.jsonl")
			}
			if dest == "" {
				dest = filepath.Join(beadsDir, "issues")
			}
			opts := convert.JSONLToFilesOptions{
				DryRun:   dryRun,
				Force:    force,
				Backup:   !dryRun,
				JSONLIn:  from,
				FilesOut: dest,
			}
			result, err := convert.ConvertJSONLToFiles(opts)
			if err != nil {
				FatalError(err.Error())
			}
			if !dryRun {
				if err := updateConfigStorage(beadsDir, "files"); err != nil {
					FatalError(err.Error())
				}
			}
			if jsonOutput {
				outputJSON(map[string]interface{}{
					"total":        result.Total,
					"written":      result.Written,
					"collisions":   result.Collisions,
					"backup_jsonl": result.JSONLBackup,
					"rewritten":    result.Rewritten,
				})
				return
			}
			fmt.Printf("Converted %d issues to %s\n", result.Written, dest)
			if result.Rewritten > 0 {
				fmt.Printf("Rewrote dependencies for %d issue(s)\n", result.Rewritten)
			}
			if result.JSONLBackup != "" {
				fmt.Printf("Backed up JSONL to %s\n", result.JSONLBackup)
			}
		case "jsonl":
			if from == "" {
				from = filepath.Join(beadsDir, "issues")
			}
			if dest == "" {
				dest = filepath.Join(beadsDir, "issues.jsonl")
			}
			opts := convert.FilesToJSONLOptions{
				DryRun:   dryRun,
				Force:    force,
				FilesIn:  from,
				JSONLOut: dest,
			}
			result, err := convert.ConvertFilesToJSONL(opts)
			if err != nil {
				FatalError(err.Error())
			}
			if !dryRun {
				if err := updateConfigStorage(beadsDir, "jsonl"); err != nil {
					FatalError(err.Error())
				}
			}
			if jsonOutput {
				outputJSON(map[string]interface{}{
					"total":   result.Total,
					"written": result.Written,
				})
				return
			}
			fmt.Printf("Converted %d issues to %s\n", result.Written, dest)
		default:
			FatalErrorWithHint("unsupported migration target", "use --to=files or --to=jsonl")
		}
	},
}

func updateConfigStorage(beadsDir string, storage string) error {
	path := filepath.Join(beadsDir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if cfg == nil {
		cfg = map[string]interface{}{}
	}
	cfg["storage"] = storage
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

func init() {
	migrateStorageCmd.Flags().String("to", "files", "Target storage format (files|jsonl)")
	migrateStorageCmd.Flags().String("from", "", "Source path (default: issues.jsonl for files, issues/ for jsonl)")
	migrateStorageCmd.Flags().String("dest", "", "Destination path (default: issues/ for files, issues.jsonl for jsonl)")
	migrateStorageCmd.Flags().Bool("dry-run", false, "Preview conversion without writing")
	migrateStorageCmd.Flags().Bool("force", false, "Overwrite existing destination directory")
	migrateCmd.AddCommand(migrateStorageCmd)
}
