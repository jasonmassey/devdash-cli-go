package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jasonmassey/devdash-cli-go/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	projectCreateCmd.Flags().String("name", "", "Project name (required)")
	projectCreateCmd.Flags().String("repo", "", "GitHub repo (owner/repo format)")
	projectCreateCmd.Flags().String("description", "", "Project description")

	projectDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	rootCmd.AddCommand(projectCmd)
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}

		body := map[string]string{"name": name}
		if repo, _ := cmd.Flags().GetString("repo"); repo != "" {
			body["githubRepo"] = repo
		}
		if desc, _ := cmd.Flags().GetString("description"); desc != "" {
			body["description"] = desc
		}

		data, err := client.Post("/projects", body)
		if err != nil {
			return err
		}

		var project api.Project
		json.Unmarshal(data, &project)
		fmt.Printf("Created project %s: %s\n", project.ID, project.Name)
		return nil
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}

		data, err := client.Get("/projects")
		if err != nil {
			return err
		}

		var projects []api.Project
		json.Unmarshal(data, &projects)

		for _, p := range projects {
			repo := ""
			if p.GithubRepo != "" {
				repo = fmt.Sprintf(" (%s)", p.GithubRepo)
			}
			fmt.Printf("%s  %s%s\n", shortID(p.ID), p.Name, repo)
		}
		return nil
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete <project-id>",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuth(); err != nil {
			return err
		}

		_, err := client.Delete("/projects/" + args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Deleted project: %s\n", args[0])
		return nil
	},
}
