package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Query metadata and enumerations",
	Long:  `Query workflow states, priorities, labels, and other metadata for scripting.`,
}

// metaStatesCmd represents the meta states command
var metaStatesCmd = &cobra.Command{
	Use:   "states [team-id]",
	Short: "List workflow states",
	Long:  `List workflow states for a specific team or all teams.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Get team ID from arg or flag
		teamID := ""
		if len(args) > 0 {
			teamID = args[0]
		} else if teamFlag != "" {
			resolvedID, err := resolveTeamID(apiClient, teamFlag)
			if err != nil {
				return err
			}
			teamID = resolvedID
		}

		if teamID == "" {
			return fmt.Errorf("team ID or --team flag is required")
		}

		// Check cache
		cacheKey := fmt.Sprintf("states-%s", teamID)
		var states interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &states); err == nil && found {
				return formatter.Output(states)
			}
		}

		// Fetch from API
		states, err = apiClient.ListWorkflowStates(getContext(), teamID)
		if err != nil {
			return fmt.Errorf("failed to list workflow states: %w", err)
		}

		// Cache results
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, states)
		}

		return formatter.Output(states)
	},
}

// metaPrioritiesCmd represents the meta priorities command
var metaPrioritiesCmd = &cobra.Command{
	Use:   "priorities",
	Short: "List priority levels",
	Long:  `List all available priority levels and their numeric values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Priority levels are static
		priorities := []map[string]interface{}{
			{"value": 0, "name": "none", "label": "No Priority"},
			{"value": 1, "name": "urgent", "label": "Urgent"},
			{"value": 2, "name": "high", "label": "High"},
			{"value": 3, "name": "medium", "label": "Medium"},
			{"value": 4, "name": "low", "label": "Low"},
		}

		return formatter.Output(priorities)
	},
}

// metaLabelsCmd represents the meta labels command
var metaLabelsCmd = &cobra.Command{
	Use:   "labels [team-id]",
	Short: "List labels",
	Long:  `List labels for a specific team or all teams.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("meta labels: not yet implemented")
	},
}

// metaCyclesCmd represents the meta cycles command
var metaCyclesCmd = &cobra.Command{
	Use:   "cycles [team-id]",
	Short: "List cycles",
	Long:  `List cycles for a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("meta cycles: not yet implemented")
	},
}

// metaIssueTypesCmd represents the meta issue-types command
var metaIssueTypesCmd = &cobra.Command{
	Use:   "issue-types",
	Short: "List issue types",
	Long:  `List all available issue types.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Issue types are typically static
		issueTypes := []map[string]interface{}{
			{"id": "issue", "name": "Issue"},
			{"id": "bug", "name": "Bug"},
			{"id": "feature", "name": "Feature"},
			{"id": "improvement", "name": "Improvement"},
		}

		return formatter.Output(issueTypes)
	},
}

func init() {
	rootCmd.AddCommand(metaCmd)

	// Add subcommands
	metaCmd.AddCommand(metaStatesCmd)
	metaCmd.AddCommand(metaPrioritiesCmd)
	metaCmd.AddCommand(metaLabelsCmd)
	metaCmd.AddCommand(metaCyclesCmd)
	metaCmd.AddCommand(metaIssueTypesCmd)
}
