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
	initiativeNameFlag string
	initiativeDescFlag string
)

// initiativeCmd represents the initiative command
var initiativeCmd = &cobra.Command{
	Use:   "initiative",
	Short: "Manage initiatives",
	Long:  `Create, view, edit, and manage Linear initiatives.`,
}

// initiativeListCmd represents the initiative list command
var initiativeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all initiatives",
	Long:  `List all initiatives.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Check cache first
		cacheKey := "initiatives"
		var initiatives interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &initiatives); err == nil && found {
				return formatter.Output(initiatives)
			}
		}

		// Fetch from API
		initiatives, err = apiClient.ListInitiatives(getContext())
		if err != nil {
			return fmt.Errorf("failed to list initiatives: %w", err)
		}

		// Cache results
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, initiatives)
		}

		return formatter.Output(initiatives)
	},
}

// initiativeViewCmd represents the initiative view command
var initiativeViewCmd = &cobra.Command{
	Use:   "view <initiative-id>",
	Short: "View initiative details",
	Long:  `View detailed information about a specific initiative.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		initiativeID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("initiative-%s", initiativeID)
		var initiative interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &initiative); err == nil && found {
				return formatter.Output(initiative)
			}
		}

		// Fetch from API
		initiative, err = apiClient.GetInitiative(getContext(), initiativeID)
		if err != nil {
			return fmt.Errorf("failed to get initiative: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, initiative)
		}

		return formatter.Output(initiative)
	},
}

// initiativeProjectsCmd represents the initiative projects command
var initiativeProjectsCmd = &cobra.Command{
	Use:   "projects <initiative-id>",
	Short: "List initiative projects",
	Long:  `List all projects in a specific initiative.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		initiativeID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("initiative-projects-%s", initiativeID)
		var projects interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &projects); err == nil && found {
				return formatter.Output(projects)
			}
		}

		// Fetch from API
		projects, err = apiClient.ListInitiativeProjects(getContext(), initiativeID)
		if err != nil {
			return fmt.Errorf("failed to list initiative projects: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, projects)
		}

		return formatter.Output(projects)
	},
}

// initiativeCreateCmd represents the initiative create command
var initiativeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new initiative",
	Long: `Create a new initiative in Linear.

Examples:
  lirt initiative create --name "2024 Strategy"
  lirt initiative create --name "Product Expansion" --description "Expand into new markets"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Validate required flags
		if initiativeNameFlag == "" {
			return fmt.Errorf("--name is required")
		}

		// Build input
		input := &client.CreateInitiativeInput{
			Name: initiativeNameFlag,
		}

		if initiativeDescFlag != "" {
			input.Description = &initiativeDescFlag
		}

		// Create initiative
		initiative, err := apiClient.CreateInitiative(getContext(), input)
		if err != nil {
			return fmt.Errorf("failed to create initiative: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Created initiative %s\n", initiative.Name)
		}

		return formatter.Output(initiative)
	},
}

// initiativeEditCmd represents the initiative edit command
var initiativeEditCmd = &cobra.Command{
	Use:   "edit <initiative-id>",
	Short: "Edit an initiative",
	Long:  `Edit an existing initiative.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		initiativeID := args[0]

		// Build input
		input := &client.UpdateInitiativeInput{}

		if initiativeNameFlag != "" {
			input.Name = &initiativeNameFlag
		}

		if initiativeDescFlag != "" {
			input.Description = &initiativeDescFlag
		}

		// Update initiative
		if err := apiClient.UpdateInitiative(getContext(), initiativeID, input); err != nil {
			return fmt.Errorf("failed to update initiative: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Updated initiative %s\n", initiativeID)
		}

		return nil
	},
}

// initiativeArchiveCmd represents the initiative archive command
var initiativeArchiveCmd = &cobra.Command{
	Use:   "archive <initiative-id>",
	Short: "Archive an initiative",
	Long:  `Archive an initiative.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		initiativeID := args[0]

		// Confirm
		if !quietFlag {
			fmt.Printf("Archive initiative %s? (y/N): ", initiativeID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Archive initiative
		if err := apiClient.ArchiveInitiative(getContext(), initiativeID); err != nil {
			return fmt.Errorf("failed to archive initiative: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Archived initiative %s\n", initiativeID)
		}

		return nil
	},
}

// initiativeDeleteCmd represents the initiative delete command
var initiativeDeleteCmd = &cobra.Command{
	Use:   "delete <initiative-id>",
	Short: "Delete an initiative",
	Long:  `Permanently delete an initiative.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		initiativeID := args[0]

		// Confirm
		if !quietFlag {
			fmt.Printf("Delete initiative %s? This cannot be undone. (y/N): ", initiativeID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Delete initiative
		if err := apiClient.DeleteInitiative(getContext(), initiativeID); err != nil {
			return fmt.Errorf("failed to delete initiative: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Deleted initiative %s\n", initiativeID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initiativeCmd)

	// Add subcommands
	initiativeCmd.AddCommand(initiativeListCmd)
	initiativeCmd.AddCommand(initiativeViewCmd)
	initiativeCmd.AddCommand(initiativeProjectsCmd)
	initiativeCmd.AddCommand(initiativeCreateCmd)
	initiativeCmd.AddCommand(initiativeEditCmd)
	initiativeCmd.AddCommand(initiativeArchiveCmd)
	initiativeCmd.AddCommand(initiativeDeleteCmd)

	// Flags for initiative create
	initiativeCreateCmd.Flags().StringVar(&initiativeNameFlag, "name", "", "Initiative name (required)")
	initiativeCreateCmd.Flags().StringVar(&initiativeDescFlag, "description", "", "Initiative description")

	// Flags for initiative edit
	initiativeEditCmd.Flags().StringVar(&initiativeNameFlag, "name", "", "Initiative name")
	initiativeEditCmd.Flags().StringVar(&initiativeDescFlag, "description", "", "Initiative description")
}
