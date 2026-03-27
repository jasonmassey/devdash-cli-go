package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	adminResetUserCmd.Flags().Bool("confirm", false, "Skip confirmation prompt")
	adminCmd.AddCommand(adminResetUserCmd)
	rootCmd.AddCommand(adminCmd)
}

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin commands (requires ADMIN_SECRET)",
}

var adminResetUserCmd = &cobra.Command{
	Use:   "reset-user <user-id>",
	Short: "Reset a user's data",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := getAdminSecret()
		if secret == "" {
			return fmt.Errorf("admin secret not found — set ADMIN_SECRET env var or create ~/.config/dev-dash/admin-secret")
		}

		if cfg == nil {
			return fmt.Errorf("configuration not loaded")
		}

		url := cfg.APIURL + "/api/admin/reset-user/" + args[0]

		req, err := http.NewRequest("POST", url, bytes.NewReader([]byte("{}")))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-admin-secret", secret)

		httpClient := &http.Client{Timeout: 30 * time.Second}
		resp, err := httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
		}

		var raw json.RawMessage
		json.Unmarshal(body, &raw)
		out, _ := json.MarshalIndent(raw, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func getAdminSecret() string {
	if s := os.Getenv("ADMIN_SECRET"); s != "" {
		return s
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(filepath.Join(home, ".config", "dev-dash", "admin-secret"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
