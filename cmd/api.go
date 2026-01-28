package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiInputFlag string
	apiVarsFlag  []string
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api <query>",
	Short: "Execute raw GraphQL queries",
	Long: `Execute raw GraphQL queries against the Linear API.

This is an escape hatch for operations not covered by built-in commands.
Always outputs JSON.

Examples:
  # Simple query
  lirt api 'query { viewer { name } }'

  # Query from file
  lirt api --input query.graphql

  # Query with variables
  lirt api 'query($id: String!) { issue(id: $id) { title } }' -f id=abc123`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, err := getClient()
		if err != nil {
			return err
		}

		// Get query from arg or file
		var query string
		if apiInputFlag != "" {
			content, err := os.ReadFile(apiInputFlag)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", apiInputFlag, err)
			}
			query = string(content)
		} else if len(args) > 0 {
			query = args[0]
		} else {
			return fmt.Errorf("query is required (provide as argument or use --input)")
		}

		// Parse variables
		variables := make(map[string]interface{})
		for _, v := range apiVarsFlag {
			parts := splitOnce(v, "=")
			if len(parts) != 2 {
				return fmt.Errorf("invalid variable format: %s (expected key=value)", v)
			}
			variables[parts[0]] = parts[1]
		}

		// Execute raw query
		var result interface{}
		if err := apiClient.Query(getContext(), query, variables); err != nil {
			return fmt.Errorf("query failed: %w", err)
		}

		// Always output as JSON
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

// Helper to split string on first occurrence of separator
func splitOnce(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			return []string{s[:i], s[i+len(sep):]}
		}
	}
	return []string{s}
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Flags
	apiCmd.Flags().StringVar(&apiInputFlag, "input", "", "Read query from file")
	apiCmd.Flags().StringSliceVarP(&apiVarsFlag, "var", "f", []string{}, "Query variables (key=value)")
}
