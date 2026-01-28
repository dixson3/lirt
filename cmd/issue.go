package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dixson3/lirt/internal/client"
	"github.com/spf13/cobra"
)

var (
	issueTeamFlag      string
	issueStateFlag     string
	issueAssigneeFlag  string
	issueLabelFlag     []string
	issueProjectFlag   string
	issuePriorityFlag  string
	issueMilestoneFlag string
	issueParentFlag    string
	issueSearchFlag    string
	issueTitleFlag     string
	issueDescFlag      string
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Long:  `Create, view, edit, and manage Linear issues.`,
}

// issueListCmd represents the issue list command
var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long:  `List issues with optional filters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Build filters
		filters := &client.IssueFilters{}

		if issueTeamFlag != "" {
			teamID, err := resolveTeamID(apiClient, issueTeamFlag)
			if err != nil {
				return err
			}
			filters.TeamID = &teamID
		}

		if issueStateFlag != "" {
			filters.StateID = &issueStateFlag
		}

		if issueAssigneeFlag != "" {
			filters.AssigneeID = &issueAssigneeFlag
		}

		if issuePriorityFlag != "" {
			priority, err := parsePriority(issuePriorityFlag)
			if err != nil {
				return err
			}
			filters.Priority = &priority
		}

		if issueSearchFlag != "" {
			filters.Search = &issueSearchFlag
		}

		// Check cache first
		cacheKey := fmt.Sprintf("issues-%s-%s-%s-%s-%s", issueTeamFlag, issueStateFlag, issueAssigneeFlag, issuePriorityFlag, issueSearchFlag)
		var issues interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &issues); err == nil && found {
				return formatter.Output(issues)
			}
		}

		// Fetch from API
		issues, err = apiClient.ListIssues(getContext(), filters)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		// Cache results
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, issues)
		}

		return formatter.Output(issues)
	},
}

// issueViewCmd represents the issue view command
var issueViewCmd = &cobra.Command{
	Use:   "view <issue-id>",
	Short: "View issue details",
	Long:  `View detailed information about a specific issue. Accepts issue identifier (e.g., ENG-123) or UUID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Check cache
		cacheKey := fmt.Sprintf("issue-%s", id)
		var issue interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &issue); err == nil && found {
				return formatter.Output(issue)
			}
		}

		// Fetch from API
		issue, err = apiClient.GetIssue(getContext(), id)
		if err != nil {
			return fmt.Errorf("failed to get issue: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, issue)
		}

		return formatter.Output(issue)
	},
}

// issueCreateCmd represents the issue create command
var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long: `Create a new issue in Linear.

Examples:
  lirt issue create --team ENG --title "Fix bug"
  lirt issue create --team ENG --title "New feature" --description "Add support for X" --priority high`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Validate required flags
		if issueTeamFlag == "" {
			return fmt.Errorf("--team is required")
		}
		if issueTitleFlag == "" {
			return fmt.Errorf("--title is required")
		}

		// Resolve team ID
		teamID, err := resolveTeamID(apiClient, issueTeamFlag)
		if err != nil {
			return err
		}

		// Build input
		input := &client.CreateIssueInput{
			TeamID: teamID,
			Title:  issueTitleFlag,
		}

		if issueDescFlag != "" {
			input.Description = &issueDescFlag
		}

		if issuePriorityFlag != "" {
			priority, err := parsePriority(issuePriorityFlag)
			if err != nil {
				return err
			}
			input.Priority = &priority
		}

		if issueStateFlag != "" {
			input.StateID = &issueStateFlag
		}

		if issueAssigneeFlag != "" {
			input.AssigneeID = &issueAssigneeFlag
		}

		if issueProjectFlag != "" {
			input.ProjectID = &issueProjectFlag
		}

		if issueParentFlag != "" {
			parentID, err := apiClient.ResolveIssueID(getContext(), issueParentFlag)
			if err != nil {
				return err
			}
			input.ParentID = &parentID
		}

		// Create issue
		issue, err := apiClient.CreateIssue(getContext(), input)
		if err != nil {
			return fmt.Errorf("failed to create issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Created issue %s: %s\n", issue.Identifier, issue.Title)
			fmt.Printf("  %s\n", issue.URL)
		}

		return formatter.Output(issue)
	},
}

// issueEditCmd represents the issue edit command
var issueEditCmd = &cobra.Command{
	Use:   "edit <issue-id>",
	Short: "Edit an issue",
	Long:  `Edit an existing issue. Accepts issue identifier (e.g., ENG-123) or UUID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Build input
		input := &client.UpdateIssueInput{}

		if issueTitleFlag != "" {
			input.Title = &issueTitleFlag
		}

		if issueDescFlag != "" {
			input.Description = &issueDescFlag
		}

		if issuePriorityFlag != "" {
			priority, err := parsePriority(issuePriorityFlag)
			if err != nil {
				return err
			}
			input.Priority = &priority
		}

		if issueStateFlag != "" {
			input.StateID = &issueStateFlag
		}

		if issueAssigneeFlag != "" {
			input.AssigneeID = &issueAssigneeFlag
		}

		if issueProjectFlag != "" {
			input.ProjectID = &issueProjectFlag
		}

		if issueParentFlag != "" {
			parentID, err := apiClient.ResolveIssueID(getContext(), issueParentFlag)
			if err != nil {
				return err
			}
			input.ParentID = &parentID
		}

		// Update issue
		if err := apiClient.UpdateIssue(getContext(), id, input); err != nil {
			return fmt.Errorf("failed to update issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Updated issue %s\n", args[0])
		}

		return nil
	},
}

// issueCloseCmd represents the issue close command
var issueCloseCmd = &cobra.Command{
	Use:   "close <issue-id>",
	Short: "Close an issue",
	Long:  `Close an issue by transitioning it to the first completed state.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Get current issue to find team
		issue, err := apiClient.GetIssue(getContext(), id)
		if err != nil {
			return err
		}

		// Get workflow states for team
		states, err := apiClient.ListWorkflowStates(getContext(), issue.Team.ID)
		if err != nil {
			return err
		}

		// Find first completed state
		var completedStateID string
		for _, state := range states {
			if state.Type == "completed" {
				completedStateID = state.ID
				break
			}
		}

		if completedStateID == "" {
			return fmt.Errorf("no completed state found for team %s", issue.Team.Key)
		}

		// Update issue to completed state
		input := &client.UpdateIssueInput{
			StateID: &completedStateID,
		}

		if err := apiClient.UpdateIssue(getContext(), id, input); err != nil {
			return fmt.Errorf("failed to close issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Closed issue %s\n", args[0])
		}

		return nil
	},
}

// issueReopenCmd represents the issue reopen command
var issueReopenCmd = &cobra.Command{
	Use:   "reopen <issue-id>",
	Short: "Reopen a closed issue",
	Long:  `Reopen a closed issue by transitioning it to the first unstarted state.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Get current issue to find team
		issue, err := apiClient.GetIssue(getContext(), id)
		if err != nil {
			return err
		}

		// Get workflow states for team
		states, err := apiClient.ListWorkflowStates(getContext(), issue.Team.ID)
		if err != nil {
			return err
		}

		// Find first unstarted state
		var unstartedStateID string
		for _, state := range states {
			if state.Type == "unstarted" {
				unstartedStateID = state.ID
				break
			}
		}

		if unstartedStateID == "" {
			return fmt.Errorf("no unstarted state found for team %s", issue.Team.Key)
		}

		// Update issue to unstarted state
		input := &client.UpdateIssueInput{
			StateID: &unstartedStateID,
		}

		if err := apiClient.UpdateIssue(getContext(), id, input); err != nil {
			return fmt.Errorf("failed to reopen issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Reopened issue %s\n", args[0])
		}

		return nil
	},
}

// issueTransitionCmd represents the issue transition command
var issueTransitionCmd = &cobra.Command{
	Use:   "transition <issue-id> <state>",
	Short: "Transition issue to a specific state",
	Long:  `Transition an issue to a specific workflow state.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Update issue state
		input := &client.UpdateIssueInput{
			StateID: &args[1],
		}

		if err := apiClient.UpdateIssue(getContext(), id, input); err != nil {
			return fmt.Errorf("failed to transition issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Transitioned issue %s to state %s\n", args[0], args[1])
		}

		return nil
	},
}

// issueArchiveCmd represents the issue archive command
var issueArchiveCmd = &cobra.Command{
	Use:   "archive <issue-id>",
	Short: "Archive an issue",
	Long:  `Archive an issue.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Confirm
		if !quietFlag {
			fmt.Printf("Archive issue %s? (y/N): ", args[0])
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Archive issue
		if err := apiClient.ArchiveIssue(getContext(), id); err != nil {
			return fmt.Errorf("failed to archive issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Archived issue %s\n", args[0])
		}

		return nil
	},
}

// issueDeleteCmd represents the issue delete command
var issueDeleteCmd = &cobra.Command{
	Use:   "delete <issue-id>",
	Short: "Delete an issue",
	Long:  `Permanently delete an issue.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Confirm
		if !quietFlag {
			fmt.Printf("Delete issue %s? This cannot be undone. (y/N): ", args[0])
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Delete issue
		if err := apiClient.DeleteIssue(getContext(), id); err != nil {
			return fmt.Errorf("failed to delete issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Deleted issue %s\n", args[0])
		}

		return nil
	},
}

// issueAssignCmd represents the issue assign command
var issueAssignCmd = &cobra.Command{
	Use:   "assign <issue-id> <user-id>",
	Short: "Assign an issue to a user",
	Long:  `Assign an issue to a specific user.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Update assignee
		userID := args[1]
		input := &client.UpdateIssueInput{
			AssigneeID: &userID,
		}

		if err := apiClient.UpdateIssue(getContext(), id, input); err != nil {
			return fmt.Errorf("failed to assign issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Assigned issue %s to %s\n", args[0], args[1])
		}

		return nil
	},
}

// issueUnassignCmd represents the issue unassign command
var issueUnassignCmd = &cobra.Command{
	Use:   "unassign <issue-id>",
	Short: "Unassign an issue",
	Long:  `Remove the assignee from an issue.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		id, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Remove assignee (set to null)
		var nullAssignee *string
		input := &client.UpdateIssueInput{
			AssigneeID: nullAssignee,
		}

		if err := apiClient.UpdateIssue(getContext(), id, input); err != nil {
			return fmt.Errorf("failed to unassign issue: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Unassigned issue %s\n", args[0])
		}

		return nil
	},
}

// Helper function to resolve team key/ID to ID
func resolveTeamID(apiClient *client.Client, teamKeyOrID string) (string, error) {
	// If it looks like a UUID, return as-is
	if len(teamKeyOrID) > 20 {
		return teamKeyOrID, nil
	}

	// Otherwise fetch teams and find by key
	teams, err := apiClient.ListTeams(getContext())
	if err != nil {
		return "", err
	}

	for _, team := range teams {
		if team.Key == teamKeyOrID {
			return team.ID, nil
		}
	}

	return "", fmt.Errorf("team not found: %s", teamKeyOrID)
}

// Helper function to parse priority value
func parsePriority(priority string) (int, error) {
	// Try parsing as number first
	if val, err := strconv.Atoi(priority); err == nil {
		if val < 0 || val > 4 {
			return 0, fmt.Errorf("priority must be 0-4 or urgent/high/medium/low/none")
		}
		return val, nil
	}

	// Parse as name
	switch strings.ToLower(priority) {
	case "urgent":
		return 1, nil
	case "high":
		return 2, nil
	case "medium":
		return 3, nil
	case "low":
		return 4, nil
	case "none", "no priority":
		return 0, nil
	default:
		return 0, fmt.Errorf("invalid priority: %s (must be 0-4 or urgent/high/medium/low/none)", priority)
	}
}

func init() {
	rootCmd.AddCommand(issueCmd)

	// Add subcommands
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueViewCmd)
	issueCmd.AddCommand(issueCreateCmd)
	issueCmd.AddCommand(issueEditCmd)
	issueCmd.AddCommand(issueCloseCmd)
	issueCmd.AddCommand(issueReopenCmd)
	issueCmd.AddCommand(issueTransitionCmd)
	issueCmd.AddCommand(issueArchiveCmd)
	issueCmd.AddCommand(issueDeleteCmd)
	issueCmd.AddCommand(issueAssignCmd)
	issueCmd.AddCommand(issueUnassignCmd)

	// Flags for issue list
	issueListCmd.Flags().StringVar(&issueTeamFlag, "team", "", "Filter by team key or ID")
	issueListCmd.Flags().StringVar(&issueStateFlag, "state", "", "Filter by state ID")
	issueListCmd.Flags().StringVar(&issueAssigneeFlag, "assignee", "", "Filter by assignee ID")
	issueListCmd.Flags().StringSliceVar(&issueLabelFlag, "label", []string{}, "Filter by label IDs")
	issueListCmd.Flags().StringVar(&issueProjectFlag, "project", "", "Filter by project ID")
	issueListCmd.Flags().StringVar(&issuePriorityFlag, "priority", "", "Filter by priority (0-4 or urgent/high/medium/low/none)")
	issueListCmd.Flags().StringVar(&issueMilestoneFlag, "milestone", "", "Filter by milestone ID")
	issueListCmd.Flags().StringVar(&issueSearchFlag, "search", "", "Search issues by text")

	// Flags for issue create
	issueCreateCmd.Flags().StringVar(&issueTeamFlag, "team", "", "Team key or ID (required)")
	issueCreateCmd.Flags().StringVar(&issueTitleFlag, "title", "", "Issue title (required)")
	issueCreateCmd.Flags().StringVar(&issueDescFlag, "description", "", "Issue description")
	issueCreateCmd.Flags().StringVar(&issuePriorityFlag, "priority", "", "Priority (0-4 or urgent/high/medium/low/none)")
	issueCreateCmd.Flags().StringVar(&issueStateFlag, "state", "", "State ID")
	issueCreateCmd.Flags().StringVar(&issueAssigneeFlag, "assignee", "", "Assignee user ID")
	issueCreateCmd.Flags().StringVar(&issueProjectFlag, "project", "", "Project ID")
	issueCreateCmd.Flags().StringVar(&issueParentFlag, "parent", "", "Parent issue ID or identifier")

	// Flags for issue edit
	issueEditCmd.Flags().StringVar(&issueTitleFlag, "title", "", "Issue title")
	issueEditCmd.Flags().StringVar(&issueDescFlag, "description", "", "Issue description")
	issueEditCmd.Flags().StringVar(&issuePriorityFlag, "priority", "", "Priority (0-4 or urgent/high/medium/low/none)")
	issueEditCmd.Flags().StringVar(&issueStateFlag, "state", "", "State ID")
	issueEditCmd.Flags().StringVar(&issueAssigneeFlag, "assignee", "", "Assignee user ID")
	issueEditCmd.Flags().StringVar(&issueProjectFlag, "project", "", "Project ID")
	issueEditCmd.Flags().StringVar(&issueParentFlag, "parent", "", "Parent issue ID or identifier")
}
