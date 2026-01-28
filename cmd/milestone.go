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
	milestoneProjectFlag    string
	milestoneNameFlag       string
	milestoneDescFlag       string
	milestoneTargetDateFlag string
)

// milestoneCmd represents the milestone command
var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Manage milestones",
	Long:  `Create, view, edit, and manage project milestones.`,
}

// milestoneListCmd represents the milestone list command
var milestoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List milestones",
	Long:  `List milestones, optionally filtered by project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Check cache first
		cacheKey := fmt.Sprintf("milestones-%s", milestoneProjectFlag)
		var milestones interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &milestones); err == nil && found {
				return formatter.Output(milestones)
			}
		}

		// Fetch from API
		milestones, err = apiClient.ListMilestones(getContext(), milestoneProjectFlag)
		if err != nil {
			return fmt.Errorf("failed to list milestones: %w", err)
		}

		// Cache results
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, milestones)
		}

		return formatter.Output(milestones)
	},
}

// milestoneViewCmd represents the milestone view command
var milestoneViewCmd = &cobra.Command{
	Use:   "view <milestone-id>",
	Short: "View milestone details",
	Long:  `View detailed information about a specific milestone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		milestoneID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("milestone-%s", milestoneID)
		var milestone interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &milestone); err == nil && found {
				return formatter.Output(milestone)
			}
		}

		// Fetch from API
		milestone, err = apiClient.GetMilestone(getContext(), milestoneID)
		if err != nil {
			return fmt.Errorf("failed to get milestone: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, milestone)
		}

		return formatter.Output(milestone)
	},
}

// milestoneIssuesCmd represents the milestone issues command
var milestoneIssuesCmd = &cobra.Command{
	Use:   "issues <milestone-id>",
	Short: "List milestone issues",
	Long:  `List all issues in a specific milestone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		milestoneID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("milestone-issues-%s", milestoneID)
		var issues interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &issues); err == nil && found {
				return formatter.Output(issues)
			}
		}

		// Fetch from API
		issues, err = apiClient.ListMilestoneIssues(getContext(), milestoneID)
		if err != nil {
			return fmt.Errorf("failed to list milestone issues: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, issues)
		}

		return formatter.Output(issues)
	},
}

// milestoneCreateCmd represents the milestone create command
var milestoneCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new milestone",
	Long: `Create a new milestone in a project.

Examples:
  lirt milestone create --project <id> --name "Beta Release"
  lirt milestone create --project <id> --name "Q1 Goals" --target-date 2024-03-31`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Validate required flags
		if milestoneProjectFlag == "" {
			return fmt.Errorf("--project is required")
		}
		if milestoneNameFlag == "" {
			return fmt.Errorf("--name is required")
		}

		// Build input
		input := &client.CreateMilestoneInput{
			ProjectID: milestoneProjectFlag,
			Name:      milestoneNameFlag,
		}

		if milestoneDescFlag != "" {
			input.Description = &milestoneDescFlag
		}

		if milestoneTargetDateFlag != "" {
			input.TargetDate = &milestoneTargetDateFlag
		}

		// Create milestone
		milestone, err := apiClient.CreateMilestone(getContext(), input)
		if err != nil {
			return fmt.Errorf("failed to create milestone: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Created milestone %s\n", milestone.Name)
		}

		return formatter.Output(milestone)
	},
}

// milestoneEditCmd represents the milestone edit command
var milestoneEditCmd = &cobra.Command{
	Use:   "edit <milestone-id>",
	Short: "Edit a milestone",
	Long:  `Edit an existing milestone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		milestoneID := args[0]

		// Build input
		input := &client.UpdateMilestoneInput{}

		if milestoneNameFlag != "" {
			input.Name = &milestoneNameFlag
		}

		if milestoneDescFlag != "" {
			input.Description = &milestoneDescFlag
		}

		if milestoneTargetDateFlag != "" {
			input.TargetDate = &milestoneTargetDateFlag
		}

		// Update milestone
		if err := apiClient.UpdateMilestone(getContext(), milestoneID, input); err != nil {
			return fmt.Errorf("failed to update milestone: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Updated milestone %s\n", milestoneID)
		}

		return nil
	},
}

// milestoneDeleteCmd represents the milestone delete command
var milestoneDeleteCmd = &cobra.Command{
	Use:   "delete <milestone-id>",
	Short: "Delete a milestone",
	Long:  `Permanently delete a milestone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		milestoneID := args[0]

		// Confirm
		if !quietFlag {
			fmt.Printf("Delete milestone %s? This cannot be undone. (y/N): ", milestoneID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Delete milestone
		if err := apiClient.DeleteMilestone(getContext(), milestoneID); err != nil {
			return fmt.Errorf("failed to delete milestone: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Deleted milestone %s\n", milestoneID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(milestoneCmd)

	// Add subcommands
	milestoneCmd.AddCommand(milestoneListCmd)
	milestoneCmd.AddCommand(milestoneViewCmd)
	milestoneCmd.AddCommand(milestoneIssuesCmd)
	milestoneCmd.AddCommand(milestoneCreateCmd)
	milestoneCmd.AddCommand(milestoneEditCmd)
	milestoneCmd.AddCommand(milestoneDeleteCmd)

	// Flags for milestone list
	milestoneListCmd.Flags().StringVar(&milestoneProjectFlag, "project", "", "Filter by project ID")

	// Flags for milestone create
	milestoneCreateCmd.Flags().StringVar(&milestoneProjectFlag, "project", "", "Project ID (required)")
	milestoneCreateCmd.Flags().StringVar(&milestoneNameFlag, "name", "", "Milestone name (required)")
	milestoneCreateCmd.Flags().StringVar(&milestoneDescFlag, "description", "", "Milestone description")
	milestoneCreateCmd.Flags().StringVar(&milestoneTargetDateFlag, "target-date", "", "Target date (YYYY-MM-DD)")

	// Flags for milestone edit
	milestoneEditCmd.Flags().StringVar(&milestoneNameFlag, "name", "", "Milestone name")
	milestoneEditCmd.Flags().StringVar(&milestoneDescFlag, "description", "", "Milestone description")
	milestoneEditCmd.Flags().StringVar(&milestoneTargetDateFlag, "target-date", "", "Target date (YYYY-MM-DD)")
}
