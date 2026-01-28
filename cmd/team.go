package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage teams",
	Long:  `List and view Linear teams.`,
}

// teamListCmd represents the team list command
var teamListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all teams",
	Long:  `List all teams with their keys, names, and descriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		// Check cache first
		cacheKey := "teams"
		var teams interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &teams); err == nil && found {
				return formatter.Output(teams)
			}
		}

		// Fetch from API
		teams, err = client.ListTeams(getContext())
		if err != nil {
			return fmt.Errorf("failed to list teams: %w", err)
		}

		// Cache the results
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, teams)
		}

		return formatter.Output(teams)
	},
}

// teamViewCmd represents the team view command
var teamViewCmd = &cobra.Command{
	Use:   "view <key-or-id>",
	Short: "View team details",
	Long:  `View detailed information about a specific team.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("team view: not yet implemented")
	},
}

// teamMembersCmd represents the team members command
var teamMembersCmd = &cobra.Command{
	Use:   "members <key-or-id>",
	Short: "List team members",
	Long:  `List all members of a specific team.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("team members: not yet implemented")
	},
}

// teamStatesCmd represents the team states command
var teamStatesCmd = &cobra.Command{
	Use:   "states <key-or-id>",
	Short: "List team workflow states",
	Long:  `List workflow states for a specific team.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("team states: not yet implemented")
	},
}

// teamLabelsCmd represents the team labels command
var teamLabelsCmd = &cobra.Command{
	Use:   "labels <key-or-id>",
	Short: "List team labels",
	Long:  `List labels for a specific team.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("team labels: not yet implemented")
	},
}

// teamCyclesCmd represents the team cycles command
var teamCyclesCmd = &cobra.Command{
	Use:   "cycles <key-or-id>",
	Short: "List team cycles",
	Long:  `List cycles for a specific team.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("team cycles: not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(teamCmd)

	// Add subcommands
	teamCmd.AddCommand(teamListCmd)
	teamCmd.AddCommand(teamViewCmd)
	teamCmd.AddCommand(teamMembersCmd)
	teamCmd.AddCommand(teamStatesCmd)
	teamCmd.AddCommand(teamLabelsCmd)
	teamCmd.AddCommand(teamCyclesCmd)
}
