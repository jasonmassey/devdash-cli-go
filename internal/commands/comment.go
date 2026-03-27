package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jasonmassey/devdash-cli-go/internal/api"
	"github.com/jasonmassey/devdash-cli-go/internal/resolve"
	"github.com/spf13/cobra"
)

func init() {
	commentCmd.Flags().String("body", "", "Comment body (required)")
	rootCmd.AddCommand(commentCmd)
	rootCmd.AddCommand(commentsCmd)
}

var commentCmd = &cobra.Command{
	Use:   "comment <id>",
	Short: "Add a comment to an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := requireProject()
		if err != nil {
			return err
		}

		body, _ := cmd.Flags().GetString("body")
		if body == "" {
			return fmt.Errorf("--body is required")
		}

		uuid, err := resolve.IDWithFetch(args[0], client, pid)
		if err != nil {
			return err
		}

		_, err = client.Post("/beads/"+uuid+"/comments", api.CommentRequest{Body: body})
		return err
	},
}

var commentsCmd = &cobra.Command{
	Use:   "comments <id>",
	Short: "List comments on an issue",
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

		data, err := client.Get("/beads/" + uuid + "/comments?projectId=" + pid)
		if err != nil {
			return err
		}

		// Pretty-print JSON array
		var comments []json.RawMessage
		if err := json.Unmarshal(data, &comments); err != nil {
			// Just output raw
			fmt.Println(string(data))
			return nil
		}

		out, _ := json.MarshalIndent(comments, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}
