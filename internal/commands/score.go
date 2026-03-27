package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jasonmassey/devdash-cli-go/internal/resolve"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scoreCmd)
}

// ScoreResult holds scoring output from the API.
type ScoreResult struct {
	BeadID            string `json:"beadId"`
	ComplexityScore   int    `json:"complexityScore"`
	AutomabilityScore int    `json:"automabilityScore"`
	AutomabilityGrade string `json:"automabilityGrade"`
	Factors           []string `json:"factors"`
}

// BulkScoreResult holds bulk scoring output.
type BulkScoreResult struct {
	Scored []ScoreResult `json:"scored"`
}

var scoreCmd = &cobra.Command{
	Use:   "score [<id>]",
	Short: "Score beads for automability",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := requireProject()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			// Bulk score all unscored beads
			data, err := client.Post("/projects/"+pid+"/beads/score", nil)
			if err != nil {
				return err
			}

			var result BulkScoreResult
			if err := json.Unmarshal(data, &result); err != nil {
				// Fallback: try as array
				var scored []ScoreResult
				if err2 := json.Unmarshal(data, &scored); err2 == nil {
					result.Scored = scored
				} else {
					out, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
					fmt.Println(string(out))
					return nil
				}
			}

			fmt.Printf("Scored %d beads:\n", len(result.Scored))
			for _, s := range result.Scored {
				fmt.Printf("  %s  %s (%d/100)  complexity=%d\n",
					shortID(s.BeadID), s.AutomabilityGrade, s.AutomabilityScore, s.ComplexityScore)
			}
			return nil
		}

		// Single bead score
		uuid, err := resolve.IDWithFetch(args[0], client, pid)
		if err != nil {
			return err
		}

		data, err := client.Post("/projects/"+pid+"/beads/"+uuid+"/score", nil)
		if err != nil {
			return err
		}

		var raw json.RawMessage
		json.Unmarshal(data, &raw)
		out, _ := json.MarshalIndent(raw, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}
