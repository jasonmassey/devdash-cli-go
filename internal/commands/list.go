package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/jasonmassey/devdash-cli-go/internal/api"
	"github.com/jasonmassey/devdash-cli-go/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().String("status", "", "Filter by status: pending, in_progress, completed")
	listCmd.Flags().String("since", "", "Filter by updatedAt (Nh, Nd, Nw, or YYYY-MM-DD)")
	listCmd.Flags().String("parent", "", "Filter by parent bead ID")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := requireProject()
		if err != nil {
			return err
		}

		beads, err := api.FetchAll[api.Bead](client, "/beads?projectId="+pid)
		if err != nil {
			return err
		}

		statusFilter, _ := cmd.Flags().GetString("status")
		since, _ := cmd.Flags().GetString("since")
		parent, _ := cmd.Flags().GetString("parent")

		var sinceFilter string
		if since != "" {
			sinceFilter, err = output.FormatSinceISO(since)
			if err != nil {
				return err
			}
		}

		// Build set of completed IDs for blocked detection
		completedIDs := make(map[string]bool)
		for _, b := range beads {
			if b.Status == "completed" {
				completedIDs[b.ID] = true
			}
		}

		var filtered []api.Bead
		for _, b := range beads {
			if statusFilter != "" && b.Status != statusFilter {
				continue
			}
			if sinceFilter != "" && b.UpdatedAt.Format("2006-01-02T15:04:05.000Z") < sinceFilter {
				continue
			}
			if parent != "" && b.ParentBeadID != parent {
				continue
			}
			filtered = append(filtered, b)
		}

		// Sort by priority
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].Priority < filtered[j].Priority
		})

		if len(filtered) == 0 {
			fmt.Fprintln(os.Stderr, "No issues found.")
			return nil
		}

		for _, b := range filtered {
			fmt.Println(output.FormatListLine(b))
		}
		return nil
	},
}
