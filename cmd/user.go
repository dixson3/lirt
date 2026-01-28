package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  `List, view, and query Linear users.`,
}

// userListCmd represents the user list command
var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `List all users in the workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Check cache first
		cacheKey := "users"
		var users interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &users); err == nil && found {
				return formatter.Output(users)
			}
		}

		// Fetch from API
		users, err = apiClient.ListUsers(getContext())
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		// Cache results
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, users)
		}

		return formatter.Output(users)
	},
}

// userViewCmd represents the user view command
var userViewCmd = &cobra.Command{
	Use:   "view <user-id>",
	Short: "View user details",
	Long:  `View detailed information about a specific user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		userID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("user-%s", userID)
		var user interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &user); err == nil && found {
				return formatter.Output(user)
			}
		}

		// Fetch from API
		user, err = apiClient.GetUser(getContext(), userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, user)
		}

		return formatter.Output(user)
	},
}

// userMeCmd represents the user me command
var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user",
	Long:  `Display information about the currently authenticated user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Fetch viewer (current user)
		viewer, err := apiClient.GetViewer(getContext())
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		return formatter.Output(viewer)
	},
}

// userIssuesCmd represents the user issues command
var userIssuesCmd = &cobra.Command{
	Use:   "issues <user-id>",
	Short: "List user's assigned issues",
	Long:  `List all issues assigned to a specific user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		userID := args[0]

		// Check cache
		cacheKey := fmt.Sprintf("user-issues-%s", userID)
		var issues interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &issues); err == nil && found {
				return formatter.Output(issues)
			}
		}

		// Fetch from API
		issues, err = apiClient.ListUserIssues(getContext(), userID)
		if err != nil {
			return fmt.Errorf("failed to list user issues: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, issues)
		}

		return formatter.Output(issues)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Add subcommands
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userViewCmd)
	userCmd.AddCommand(userMeCmd)
	userCmd.AddCommand(userIssuesCmd)
}
