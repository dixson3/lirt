package cmd

import (
	"fmt"

	"github.com/dixson3/lirt/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify lirt configuration settings.`,
}

// configListCmd represents the config list command
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config values",
	Long:  `List all configuration values for the current profile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(profileFlag)
		cfg, err := config.LoadConfig(profile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Build config map
		configMap := map[string]interface{}{
			"profile":   profile,
			"workspace": cfg.Workspace,
			"team":      cfg.Team,
			"format":    cfg.Format,
		}

		return formatter.Output(configMap)
	},
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Long:  `Get the value of a specific configuration key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(profileFlag)
		cfg, err := config.LoadConfig(profile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		key := args[0]
		var value string

		switch key {
		case "workspace":
			value = cfg.Workspace
		case "team":
			value = cfg.Team
		case "format":
			value = cfg.Format
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}

		if value == "" {
			return fmt.Errorf("config key %s is not set", key)
		}

		fmt.Println(value)
		return nil
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long:  `Set the value of a specific configuration key.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(profileFlag)
		key := args[0]
		value := args[1]

		// Validate key
		validKeys := []string{"workspace", "team", "format"}
		valid := false
		for _, k := range validKeys {
			if key == k {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid config key: %s (valid keys: workspace, team, format)", key)
		}

		// Save config value
		if err := config.SaveConfigValue(profile, key, value); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Set %s = %s\n", key, value)
		}

		return nil
	},
}

// configUnsetCmd represents the config unset command
var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Unset a config value",
	Long:  `Remove a configuration key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile := config.GetProfile(profileFlag)
		key := args[0]

		// Save empty value to unset
		if err := config.SaveConfigValue(profile, key, ""); err != nil {
			return fmt.Errorf("failed to unset config: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Unset %s\n", key)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Add subcommands
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
}
