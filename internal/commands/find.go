package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(findCmd)
}

var findCmd = &cobra.Command{
	Use:   "find <uuid>",
	Short: "Look up a bead by full UUID across all projects",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}

		uuid := args[0]
		// Require full UUID format
		if len(uuid) != 36 || strings.Count(uuid, "-") != 4 {
			return fmt.Errorf("find requires a full UUID (36 characters with dashes)")
		}

		data, err := client.Get("/beads/" + uuid)
		if err != nil {
			return err
		}

		var raw json.RawMessage
		if err := json.Unmarshal(data, &raw); err != nil {
			fmt.Println(string(data))
			return nil
		}

		out, _ := json.MarshalIndent(raw, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}
