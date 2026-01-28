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
	commentIssueFlag string
	commentBodyFlag  string
	commentFileFlag  string
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage comments",
	Long:  `Add, edit, and delete comments on issues.`,
}

// commentListCmd represents the comment list command
var commentListCmd = &cobra.Command{
	Use:   "list <issue-id>",
	Short: "List comments on an issue",
	Long:  `List all comments on a specific issue.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		issueID, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Check cache
		cacheKey := fmt.Sprintf("comments-%s", issueID)
		var comments interface{}
		if !noCacheFlag {
			if found, err := cacheInstance.Get(cacheKey, &comments); err == nil && found {
				return formatter.Output(comments)
			}
		}

		// Fetch from API
		comments, err = apiClient.ListIssueComments(getContext(), issueID)
		if err != nil {
			return fmt.Errorf("failed to list comments: %w", err)
		}

		// Cache result
		if !noCacheFlag {
			cacheInstance.Set(cacheKey, comments)
		}

		return formatter.Output(comments)
	},
}

// commentAddCmd represents the comment add command
var commentAddCmd = &cobra.Command{
	Use:   "add <issue-id>",
	Short: "Add a comment to an issue",
	Long: `Add a comment to an issue.

Examples:
  lirt comment add ENG-123 --body "This looks good"
  lirt comment add ENG-123 --body-file comment.md`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Resolve issue ID
		issueID, err := apiClient.ResolveIssueID(getContext(), args[0])
		if err != nil {
			return err
		}

		// Get body from flag or file
		body := commentBodyFlag
		if commentFileFlag != "" {
			content, err := os.ReadFile(commentFileFlag)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", commentFileFlag, err)
			}
			body = string(content)
		}

		if body == "" {
			return fmt.Errorf("comment body is required (use --body or --body-file)")
		}

		// Create comment
		input := &client.CreateCommentInput{
			IssueID: &issueID,
			Body:    body,
		}

		comment, err := apiClient.CreateComment(getContext(), input)
		if err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Added comment to %s\n", args[0])
		}

		return formatter.Output(comment)
	},
}

// commentEditCmd represents the comment edit command
var commentEditCmd = &cobra.Command{
	Use:   "edit <comment-id>",
	Short: "Edit a comment",
	Long: `Edit an existing comment.

Examples:
  lirt comment edit <id> --body "Updated text"
  lirt comment edit <id> --body-file comment.md`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		commentID := args[0]

		// Get body from flag or file
		body := commentBodyFlag
		if commentFileFlag != "" {
			content, err := os.ReadFile(commentFileFlag)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", commentFileFlag, err)
			}
			body = string(content)
		}

		if body == "" {
			return fmt.Errorf("comment body is required (use --body or --body-file)")
		}

		// Update comment
		input := &client.UpdateCommentInput{
			Body: body,
		}

		if err := apiClient.UpdateComment(getContext(), commentID, input); err != nil {
			return fmt.Errorf("failed to update comment: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Updated comment %s\n", commentID)
		}

		return nil
	},
}

// commentDeleteCmd represents the comment delete command
var commentDeleteCmd = &cobra.Command{
	Use:   "delete <comment-id>",
	Short: "Delete a comment",
	Long:  `Permanently delete a comment.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		commentID := args[0]

		// Confirm
		if !quietFlag {
			fmt.Printf("Delete comment %s? This cannot be undone. (y/N): ", commentID)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Delete comment
		if err := apiClient.DeleteComment(getContext(), commentID); err != nil {
			return fmt.Errorf("failed to delete comment: %w", err)
		}

		if !quietFlag {
			fmt.Printf("✓ Deleted comment %s\n", commentID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)

	// Add subcommands
	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentAddCmd)
	commentCmd.AddCommand(commentEditCmd)
	commentCmd.AddCommand(commentDeleteCmd)

	// Flags for comment add
	commentAddCmd.Flags().StringVar(&commentBodyFlag, "body", "", "Comment body text")
	commentAddCmd.Flags().StringVar(&commentFileFlag, "body-file", "", "File containing comment body (markdown)")

	// Flags for comment edit
	commentEditCmd.Flags().StringVar(&commentBodyFlag, "body", "", "Comment body text")
	commentEditCmd.Flags().StringVar(&commentFileFlag, "body-file", "", "File containing comment body (markdown)")
}
