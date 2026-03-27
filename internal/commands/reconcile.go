package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	reconcileCmd.Flags().Bool("dry-run", true, "Preview findings only (default)")
	reconcileCmd.Flags().Bool("auto-fix", false, "Apply fixes automatically")
	reconcileCmd.Flags().Bool("json", false, "Output raw JSON")
	rootCmd.AddCommand(reconcileCmd)
}

// ReconcileResult holds the API response for reconcile-tasks.
type ReconcileResult struct {
	Findings []ReconcileFinding `json:"findings"`
	Fixed    int                `json:"fixed"`
}

// ReconcileFinding is a single finding from reconciliation.
type ReconcileFinding struct {
	Type          string   `json:"type"`
	Severity      string   `json:"severity"`
	Subject       string   `json:"subject"`
	Reason        string   `json:"reason"`
	RelatedBeadID string   `json:"relatedBeadId,omitempty"`
	BeadID        string   `json:"beadId,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

var reconcileCmd = &cobra.Command{
	Use:   "reconcile-tasks",
	Short: "Audit and fix backlog inconsistencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := requireProject()
		if err != nil {
			return err
		}

		autoFix, _ := cmd.Flags().GetBool("auto-fix")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		body := map[string]interface{}{
			"projectId": pid,
			"autoFix":   autoFix,
		}

		data, err := client.Post("/jobs/reconcile-tasks", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			var raw json.RawMessage
			json.Unmarshal(data, &raw)
			out, _ := json.MarshalIndent(raw, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		var result ReconcileResult
		if err := json.Unmarshal(data, &result); err != nil {
			// Fallback to raw output
			fmt.Println(string(data))
			return nil
		}

		fmt.Println("Auditing backlog...")
		fmt.Println()

		// Group by type
		typeCounts := make(map[string]int)
		for _, f := range result.Findings {
			typeCounts[f.Type]++
		}

		for t, c := range typeCounts {
			fmt.Printf("  %s: %d\n", t, c)
		}
		fmt.Println()

		// Print findings
		for _, f := range result.Findings {
			related := ""
			if f.RelatedBeadID != "" {
				related = fmt.Sprintf(" [related: %s]", shortID(f.RelatedBeadID))
			}
			fmt.Printf("[%s] %s\n  → %s%s\n\n", f.Severity, f.Subject, f.Reason, related)
		}

		if !autoFix && len(result.Findings) > 0 {
			fmt.Println("Run with --auto-fix to fix dependency issues automatically")
		}

		if autoFix {
			fmt.Printf("Fixed %d issues.\n", result.Fixed)
		}

		return nil
	},
}
