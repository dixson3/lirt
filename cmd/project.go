package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dixson3/lirt/internal/client"
	"github.com/spf13/cobra"
)

var (
	projectNameFlag  string
	projectDescFlag  string
	projectStateFlag string
	projectLeadFlag  string
	projectPriorityFlag string
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  `Create, view, edit, and manage Linear projects.`,
}

// projectListCmd represents the project list command
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `List all projects with their state, priority, and lead.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Check cache first
		cacheKey := "projects"
		var projects interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &projects); err == nil && found {
				return formatter.Output(projects)
			}
		}

		// Fetch from API
		projects, err = apiClient.ListProjects(getContext())
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}

		// Cache results
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, projects)
		}

		return formatter.Output(projects)
	},
}

// projectViewCmd represents the project view command
var projectViewCmd = &cobra.Command{
	Use:   "view <project-id>",
	Short: "View project details",
	Long:  `View detailed information about a specific project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		projectID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("project-%s", projectID)
		var project interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &project); err == nil && found {
				return formatter.Output(project)
			}
		}

		// Fetch from API
		project, err = apiClient.GetProject(getContext(), projectID)
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, project)
		}

		return formatter.Output(project)
	},
}

// projectIssuesCmd represents the project issues command
var projectIssuesCmd = &cobra.Command{
	Use:   "issues <project-id>",
	Short: "List project issues",
	Long:  `List all issues in a specific project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		projectID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("project-issues-%s", projectID)
		var issues interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &issues); err == nil && found {
				return formatter.Output(issues)
			}
		}

		// Fetch from API
		issues, err = apiClient.ListProjectIssues(getContext(), projectID)
		if err != nil {
			return fmt.Errorf("failed to list project issues: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, issues)
		}

		return formatter.Output(issues)
	},
}

// projectCreateCmd represents the project create command
var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long: `Create a new project in Linear.

Project states: backlog, planned, started, paused, completed, canceled

Examples:
  lirt project create --name "Q1 Initiative"
  lirt project create --name "Migration" --description "Database migration" --state planned`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Validate required flags
		if projectNameFlag == "" {
			return fmt.Errorf("--name is required")
		}

		// Build input
		input := &client.CreateProjectInput{
			Name: projectNameFlag,
		}

		if projectDescFlag != "" {
			input.Description = &projectDescFlag
		}

		if projectStateFlag != "" {
			// Validate state
			validStates := []string{"backlog", "planned", "started", "paused", "completed", "canceled"}
			valid := false
			for _, s := range validStates {
				if projectStateFlag == s {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid state: %s (must be one of: %s)", projectStateFlag, strings.Join(validStates, ", "))
			}
			input.State = &projectStateFlag
		}

		if projectPriorityFlag != "" {
			priority, err := parsePriority(projectPriorityFlag)
			if err != nil {
				return err
			}
			input.Priority = &priority
		}

		if projectLeadFlag != "" {
			input.LeadID = &projectLeadFlag
		}

		// Create project
		project, err := apiClient.CreateProject(getContext(), input)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Created project %s\n", project.Name)
			fmt.Printf("  %s\n", project.URL)
		}

		return formatter.Output(project)
	},
}

// projectEditCmd represents the project edit command
var projectEditCmd = &cobra.Command{
	Use:   "edit <project-id>",
	Short: "Edit a project",
	Long:  `Edit an existing project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		projectID := args[0]

		// Build input
		input := &client.UpdateProjectInput{}

		if projectNameFlag != "" {
			input.Name = &projectNameFlag
		}

		if projectDescFlag != "" {
			input.Description = &projectDescFlag
		}

		if projectStateFlag != "" {
			// Validate state
			validStates := []string{"backlog", "planned", "started", "paused", "completed", "canceled"}
			valid := false
			for _, s := range validStates {
				if projectStateFlag == s {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid state: %s (must be one of: %s)", projectStateFlag, strings.Join(validStates, ", "))
			}
			input.State = &projectStateFlag
		}

		if projectPriorityFlag != "" {
			priority, err := parsePriority(projectPriorityFlag)
			if err != nil {
				return err
			}
			input.Priority = &priority
		}

		if projectLeadFlag != "" {
			input.LeadID = &projectLeadFlag
		}

		// Update project
		if err := apiClient.UpdateProject(getContext(), projectID, input); err != nil {
			return fmt.Errorf("failed to update project: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Updated project %s\n", projectID)
		}

		return nil
	},
}

// projectArchiveCmd represents the project archive command
var projectArchiveCmd = &cobra.Command{
	Use:   "archive <project-id>",
	Short: "Archive a project",
	Long:  `Archive a project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		projectID := args[0]

		// Confirm
		if !quietFlag {
			fmt.Printf("Archive project %s? (y/N): ", projectID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Archive project
		if err := apiClient.ArchiveProject(getContext(), projectID); err != nil {
			return fmt.Errorf("failed to archive project: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Archived project %s\n", projectID)
		}

		return nil
	},
}

// projectDeleteCmd represents the project delete command
var projectDeleteCmd = &cobra.Command{
	Use:   "delete <project-id>",
	Short: "Delete a project",
	Long:  `Permanently delete a project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		projectID := args[0]

		// Confirm
		if !quietFlag {
			fmt.Printf("Delete project %s? This cannot be undone. (y/N): ", projectID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Delete project
		if err := apiClient.DeleteProject(getContext(), projectID); err != nil {
			return fmt.Errorf("failed to delete project: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Deleted project %s\n", projectID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	// Add subcommands
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectViewCmd)
	projectCmd.AddCommand(projectIssuesCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectEditCmd)
	projectCmd.AddCommand(projectArchiveCmd)
	projectCmd.AddCommand(projectDeleteCmd)

	// Flags for project create
	projectCreateCmd.Flags().StringVar(&projectNameFlag, "name", "", "Project name (required)")
	projectCreateCmd.Flags().StringVar(&projectDescFlag, "description", "", "Project description")
	projectCreateCmd.Flags().StringVar(&projectStateFlag, "state", "", "Project state (backlog, planned, started, paused, completed, canceled)")
	projectCreateCmd.Flags().StringVar(&projectPriorityFlag, "priority", "", "Priority (0-4 or urgent/high/medium/low/none)")
	projectCreateCmd.Flags().StringVar(&projectLeadFlag, "lead", "", "Lead user ID")

	// Flags for project edit
	projectEditCmd.Flags().StringVar(&projectNameFlag, "name", "", "Project name")
	projectEditCmd.Flags().StringVar(&projectDescFlag, "description", "", "Project description")
	projectEditCmd.Flags().StringVar(&projectStateFlag, "state", "", "Project state (backlog, planned, started, paused, completed, canceled)")
	projectEditCmd.Flags().StringVar(&projectPriorityFlag, "priority", "", "Priority (0-4 or urgent/high/medium/low/none)")
	projectEditCmd.Flags().StringVar(&projectLeadFlag, "lead", "", "Lead user ID")
}
