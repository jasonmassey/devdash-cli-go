package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	tokenCmd.AddCommand(tokenCreateCmd)
	tokenCmd.AddCommand(tokenListCmd)
	tokenCmd.AddCommand(tokenRevokeCmd)
	rootCmd.AddCommand(tokenCmd)
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage API tokens",
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}

		data, err := client.Post("/auth/tokens", map[string]string{
			"name": args[0],
		})
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

var tokenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}

		data, err := client.Get("/auth/tokens")
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

var tokenRevokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Revoke an API token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}

		_, err := client.Delete("/auth/tokens/" + args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Revoked token: %s\n", args[0])
		return nil
	},
}
