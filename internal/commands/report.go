package commands

import (
	"fmt"

	"github.com/jasonmassey/devdash-cli-go/internal/api"
	"github.com/jasonmassey/devdash-cli-go/internal/resolve"
	"github.com/spf13/cobra"
)

func init() {
	reportCmd.Flags().String("status", "", "Status: code_complete, committed, pushed, error (required)")
	reportCmd.Flags().String("summary", "", "Progress summary")
	reportCmd.Flags().Int("files-changed", 0, "Number of files changed")
	reportCmd.Flags().String("branch", "", "Git branch name")
	reportCmd.Flags().String("commit", "", "Git commit SHA")
	reportCmd.Flags().String("error", "", "Error message (when status=error)")
	rootCmd.AddCommand(reportCmd)
}

var reportCmd = &cobra.Command{
	Use:   "report <id>",
	Short: "Report progress on an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := requireProject()
		if err != nil {
			return err
		}

		status, _ := cmd.Flags().GetString("status")
		if status == "" {
			return fmt.Errorf("--status is required (code_complete, committed, pushed, error)")
		}

		uuid, err := resolve.IDWithFetch(args[0], client, pid)
		if err != nil {
			return err
		}

		req := api.ReportRequest{
			Status: status,
		}

		if v, _ := cmd.Flags().GetString("summary"); v != "" {
			req.Summary = v
		}
		if cmd.Flags().Changed("files-changed") {
			v, _ := cmd.Flags().GetInt("files-changed")
			req.FilesChanged = &v
		}
		if v, _ := cmd.Flags().GetString("branch"); v != "" {
			req.Branch = v
		}
		if v, _ := cmd.Flags().GetString("commit"); v != "" {
			req.Commit = v
		}
		if v, _ := cmd.Flags().GetString("error"); v != "" {
			req.Error = v
		}

		_, err = client.Post("/beads/"+uuid+"/report", req)
		if err != nil {
			// Fire-and-forget: report errors to stderr but don't fail
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: report failed: %v\n", err)
			return nil
		}

		fmt.Printf("Report submitted: %s for %s\n", status, uuid)
		return nil
	},
}
