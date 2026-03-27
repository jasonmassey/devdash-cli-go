package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jasonmassey/devdash-cli-go/internal/api"
	"github.com/jasonmassey/devdash-cli-go/internal/resolve"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

// AnalyzeResult holds the analysis output from a completed job.
type AnalyzeResult struct {
	EstimatedComplexity string           `json:"estimatedComplexity"`
	AffectedFiles       json.RawMessage  `json:"affectedFiles"`
	AffectedModules     json.RawMessage  `json:"affectedModules"`
	ShouldSubdivide     bool             `json:"shouldSubdivide"`
	Reasoning           string           `json:"reasoning"`
	AgentInstructions   string           `json:"agentInstructions"`
	Subtasks            []AnalyzeSubtask `json:"subtasks"`
}

// AnalyzeSubtask is a subtask created by analysis.
type AnalyzeSubtask struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze <id>",
	Short: "Trigger sandbox analysis for an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := requireProject()
		if err != nil {
			return err
		}

		uuid, err := resolve.IDWithFetch(args[0], client, pid)
		if err != nil {
			return err
		}

		// Fetch bead for subject
		beadData, err := client.Get("/beads/" + uuid + "?projectId=" + pid)
		if err != nil {
			return err
		}
		var bead api.Bead
		json.Unmarshal(beadData, &bead)

		fmt.Printf("Analyzing: %s\n\n", bead.Subject)

		// Trigger analysis
		data, err := client.Post("/jobs/analyze", map[string]string{
			"beadId":    uuid,
			"projectId": pid,
		})
		if err != nil {
			return fmt.Errorf("failed to start analysis: %w", err)
		}

		var job api.Job
		if err := json.Unmarshal(data, &job); err != nil {
			return err
		}

		// Poll for completion (5s interval, 300s timeout)
		timeout := time.After(300 * time.Second)
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				return fmt.Errorf("analysis timed out after 5 minutes")
			case <-ticker.C:
				statusData, err := client.Get("/jobs/" + job.ID)
				if err != nil {
					continue
				}
				var current api.Job
				if err := json.Unmarshal(statusData, &current); err != nil {
					continue
				}

				switch current.Status {
				case "completed":
					return printAnalyzeResult(current, pid)
				case "failed":
					errMsg := current.Error
					if current.FailureAnalysis != nil {
						errMsg = current.FailureAnalysis.Summary
					}
					return fmt.Errorf("analysis failed: %s", errMsg)
				default:
					fmt.Fprintf(os.Stderr, "Status: %s...\n", current.Status)
				}
			}
		}
	},
}

func printAnalyzeResult(job api.Job, pid string) error {
	if job.Result == nil {
		fmt.Println("Analysis complete (no result data)")
		return nil
	}

	resultBytes, err := json.Marshal(job.Result)
	if err != nil {
		return err
	}

	var result AnalyzeResult
	if err := json.Unmarshal(resultBytes, &result); err != nil {
		// Fallback: print raw result
		out, _ := json.MarshalIndent(job.Result, "", "  ")
		fmt.Println(string(out))
		return nil
	}

	fmt.Printf("Complexity:  %s\n", result.EstimatedComplexity)
	fmt.Printf("Files:       %s\n", string(result.AffectedFiles))
	fmt.Printf("Modules:     %s\n", string(result.AffectedModules))
	fmt.Printf("Subdivide:   %v\n", result.ShouldSubdivide)

	if result.Reasoning != "" {
		fmt.Printf("\n## Reasoning\n%s\n", result.Reasoning)
	}
	if result.AgentInstructions != "" {
		fmt.Printf("\n## Agent Instructions\n%s\n", result.AgentInstructions)
	}

	if len(result.Subtasks) > 0 {
		fmt.Printf("\n## Decomposed into %d subtasks:\n", len(result.Subtasks))
		for _, st := range result.Subtasks {
			fmt.Printf("  - %s %s\n", shortID(st.ID), st.Subject)
		}
	}

	// Check for new children
	childData, err := client.Get("/beads?projectId=" + pid + "&parentId=" + job.BeadID)
	if err == nil {
		var children []api.Bead
		json.Unmarshal(childData, &children)
		if len(children) > 0 && len(result.Subtasks) == 0 {
			fmt.Printf("\n## Decomposed into %d subtasks:\n", len(children))
			for _, c := range children {
				fmt.Printf("  - %s %s\n", shortID(c.ID), c.Subject)
			}
		}
	}

	return nil
}
