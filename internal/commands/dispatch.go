package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jasonmassey/devdash-cli-go/internal/api"
	"github.com/jasonmassey/devdash-cli-go/internal/resolve"
	"github.com/spf13/cobra"
)

func init() {
	dispatchCmd.Flags().Int("priority", 2, "Job priority: 0-4")
	dispatchCmd.Flags().String("worker", "", "Execution backend: docker, e2b, railway")
	rootCmd.AddCommand(dispatchCmd)
}

var dispatchCmd = &cobra.Command{
	Use:   "dispatch <id>",
	Short: "Dispatch a bead for execution",
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

		// Fetch bead for prompt
		beadData, err := client.Get("/beads/" + uuid + "?projectId=" + pid)
		if err != nil {
			return err
		}
		var bead api.Bead
		json.Unmarshal(beadData, &bead)

		// Build prompt: preInstructions > description > subject
		prompt := bead.PreInstructions
		if prompt == "" {
			prompt = bead.Description
		}
		if prompt == "" {
			prompt = bead.Subject
		}

		body := map[string]interface{}{
			"beadId":    uuid,
			"projectId": pid,
			"prompt":    prompt,
		}

		if cmd.Flags().Changed("priority") {
			p, _ := cmd.Flags().GetInt("priority")
			body["priority"] = p
		}
		if worker, _ := cmd.Flags().GetString("worker"); worker != "" {
			body["workerType"] = worker
		}

		data, err := client.Post("/jobs", body)
		if err != nil {
			return err
		}

		var job api.Job
		json.Unmarshal(data, &job)

		fmt.Printf("Job queued: %s\n", job.ID)
		fmt.Printf("  Bead:   %s — %s\n", shortID(uuid), bead.Subject)
		fmt.Printf("  Status: %s\n", job.Status)
		if job.WorkerType != "" {
			fmt.Printf("  Worker: %s\n", job.WorkerType)
		}
		return nil
	},
}
