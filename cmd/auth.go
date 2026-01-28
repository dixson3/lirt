package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/dixson3/lirt/internal/client"
	"github.com/dixson3/lirt/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	authAPIKeyFlag string
	authProfileFlag string
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication and credentials",
	Long:  `Manage authentication and credentials for Linear API access.`,
}

// authLoginCmd represents the auth login command
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Linear API",
	Long: `Authenticate with Linear API using a Personal API Key.

The API key will be validated by calling the Linear API, and if valid,
stored in ~/.config/lirt/credentials with secure file permissions (0600).

Examples:
  # Interactive login (prompts for API key)
  lirt auth login

  # Non-interactive login
  lirt auth login --api-key lin_api_xxxxx...

  # Login to named profile
  lirt auth login --profile work`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(authProfileFlag)
		apiKey := authAPIKeyFlag

		// If no API key provided, prompt for it
		if apiKey == "" {
			fmt.Print("Enter your Linear API key: ")
			keyBytes, err := term.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			if err != nil {
				return fmt.Errorf("failed to read API key: %w", err)
			}
			apiKey = strings.TrimSpace(string(keyBytes))
		}

		if apiKey == "" {
			return fmt.Errorf("API key is required")
		}

		// Validate API key by calling viewer query
		if !quietFlag {
			fmt.Println("Validating API key...")
		}

		testClient, err := client.New(apiKey)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		viewer, err := testClient.GetViewer(getContext())
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Check if profile exists and confirm overwrite
		if _, err := config.LoadAPIKey(profile); err == nil {
			if !quietFlag {
				fmt.Printf("Profile '%s' already exists. Overwrite? (y/N): ", profile)
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					fmt.Println("Cancelled.")
					return nil
				}
			}
		}

		// Save API key to credentials file
		if err := config.SaveAPIKey(profile, apiKey); err != nil {
			return fmt.Errorf("failed to save API key: %w", err)
		}

		// Save workspace name to config file
		if err := config.SaveConfigValue(profile, "workspace", viewer.Organization.Name); err != nil {
			return fmt.Errorf("failed to save workspace name: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Authenticated as %s (%s)\n", viewer.Name, viewer.Email)
			fmt.Printf("✓ Workspace: %s\n", viewer.Organization.Name)
			fmt.Printf("✓ Saved to profile '%s'\n", profile)
		}

		return nil
	},
}

// authStatusCmd represents the auth status command
var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Display the current authentication state including profile, workspace, and user information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(authProfileFlag)

		cfg, err := config.LoadConfig(profile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.APIKey == "" {
			fmt.Println("Not authenticated")
			fmt.Printf("Run 'lirt auth login' to set up credentials\n")
			return fmt.Errorf("not authenticated")
		}

		// Get viewer info
		apiClient, err := client.New(cfg.APIKey)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		viewer, err := apiClient.GetViewer(getContext())
		if err != nil {
			return fmt.Errorf("failed to get viewer: %w", err)
		}

		// Display status
		fmt.Printf("Profile:    %s\n", profile)
		fmt.Printf("Workspace:  %s\n", viewer.Organization.Name)
		fmt.Printf("User:       %s (%s)\n", viewer.Name, viewer.Email)
		fmt.Printf("Key prefix: %s\n", client.MaskAPIKey(cfg.APIKey))

		return nil
	},
}

// authTokenCmd represents the auth token command
var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print API key to stdout",
	Long: `Print the API key to stdout for piping to other tools.

Warning: The API key will be visible in your terminal.

Example:
  # Use with curl
  curl -H "Authorization: Bearer $(lirt auth token)" https://api.linear.app/graphql`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(authProfileFlag)

		apiKey, err := config.LoadAPIKey(profile)
		if err != nil {
			return fmt.Errorf("no API key found: %w", err)
		}

		fmt.Println(apiKey)
		return nil
	},
}

// authLogoutCmd represents the auth logout command
var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove credentials for a profile",
	Long:  `Remove a profile from both the credentials and config files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(authProfileFlag)

		// Check if profile exists
		if _, err := config.LoadAPIKey(profile); err != nil {
			return fmt.Errorf("profile '%s' not found", profile)
		}

		// Confirm deletion
		if !quietFlag {
			fmt.Printf("Remove profile '%s'? (y/N): ", profile)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Delete profile
		if err := config.DeleteProfile(profile); err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Profile '%s' removed\n", profile)
		}

		return nil
	},
}

// authListCmd represents the auth list command
var authListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	Long:  `List all configured profiles with their workspace names and API key prefixes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := config.ListProfiles()
		if err != nil {
			return fmt.Errorf("failed to list profiles: %w", err)
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles configured")
			fmt.Println("Run 'lirt auth login' to set up credentials")
			return nil
		}

		// Build table data
		type profileRow struct {
			Profile   string
			Workspace string
			KeyPrefix string
		}

		rows := []profileRow{}
		for profile, workspace := range profiles {
			apiKey, err := config.LoadAPIKey(profile)
			keyPrefix := "—"
			if err == nil {
				keyPrefix = client.MaskAPIKey(apiKey)
			}

			rows = append(rows, profileRow{
				Profile:   profile,
				Workspace: workspace,
				KeyPrefix: keyPrefix,
			})
		}

		// Output using formatter
		return formatter.Output(rows)
	},
}

// authSwitchCmd represents the auth switch command
var authSwitchCmd = &cobra.Command{
	Use:   "switch <profile>",
	Short: "Switch to a different profile",
	Long: `Print an export command to switch to a different profile.

Example:
  # Bash/Zsh
  eval "$(lirt auth switch work)"

  # Fish
  lirt auth switch work | source`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := args[0]

		// Verify profile exists
		if _, err := config.LoadAPIKey(profile); err != nil {
			return fmt.Errorf("profile '%s' not found", profile)
		}

		// Print export command
		fmt.Printf("export LIRT_PROFILE=%s\n", profile)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Add subcommands
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authTokenCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authListCmd)
	authCmd.AddCommand(authSwitchCmd)

	// Flags for auth login
	authLoginCmd.Flags().StringVar(&authAPIKeyFlag, "api-key", "", "API key (non-interactive)")
	authLoginCmd.Flags().StringVar(&authProfileFlag, "profile", "", "Profile name (default: default)")

	// Flags for other commands
	authStatusCmd.Flags().StringVar(&authProfileFlag, "profile", "", "Profile name")
	authTokenCmd.Flags().StringVar(&authProfileFlag, "profile", "", "Profile name")
	authLogoutCmd.Flags().StringVar(&authProfileFlag, "profile", "", "Profile name")
}
