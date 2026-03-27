package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jasonmassey/devdash-cli-go/internal/resolve"
	"github.com/spf13/cobra"
)

func init() {
	activityCmd.Flags().Int("limit", 0, "Maximum number of results")
	rootCmd.AddCommand(activityCmd)
}

var activityCmd = &cobra.Command{
	Use:   "activity [<id>]",
	Short: "View activity log",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := requireProject()
		if err != nil {
			return err
		}

		path := "/activity?projectId=" + pid

		if len(args) > 0 {
			uuid, err := resolve.IDWithFetch(args[0], client, pid)
			if err != nil {
				return err
			}
			path = "/activity?beadId=" + uuid
		}

		if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
			path += fmt.Sprintf("&limit=%d", limit)
		}

		data, err := client.Get(path)
		if err != nil {
			return err
		}

		var activity json.RawMessage
		if err := json.Unmarshal(data, &activity); err != nil {
			fmt.Println(string(data))
			return nil
		}

		out, _ := json.MarshalIndent(activity, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}
