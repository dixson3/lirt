package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dixson3/lirt/internal/cache"
	"github.com/dixson3/lirt/internal/client"
	"github.com/dixson3/lirt/internal/config"
	"github.com/dixson3/lirt/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Global flags
	profileFlag  string
	apiKeyFlag   string
	teamFlag     string
	formatFlag   string
	jsonFlag     string
	jqFlag       string
	noCacheFlag  bool
	quietFlag    bool
	verboseFlag  bool

	// Version is injected at build time
	Version = "dev"

	// Shared context
	cfg    *config.Config
	apiClient *client.Client
	cacheInstance *cache.Cache
	formatter *output.Formatter
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "lirt",
	Short: "Linear CLI tool",
	Long: `lirt is a command-line interface for Linear, following gh CLI semantics.

It provides fast, scriptable access to Linear's API with support for multiple workspaces via named profiles.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize configuration
		profile := config.GetProfile(profileFlag)

		var err error
		cfg, err = config.LoadConfig(profile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override with flags
		if apiKeyFlag != "" {
			cfg.APIKey = apiKeyFlag
		}
		if teamFlag != "" {
			cfg.Team = teamFlag
		}
		if formatFlag != "" {
			cfg.Format = formatFlag
		}

		// Override with environment variables
		if team := os.Getenv("LIRT_TEAM"); team != "" && teamFlag == "" {
			cfg.Team = team
		}
		if format := os.Getenv("LIRT_FORMAT"); format != "" && formatFlag == "" {
			cfg.Format = format
		}

		// Parse cache TTL
		cacheTTL := 5 * time.Minute
		if cfg.CacheTTL != "" {
			if duration, err := time.ParseDuration(cfg.CacheTTL); err == nil {
				cacheTTL = duration
			}
		}

		// Initialize cache
		cacheInstance = cache.New(profile, cacheTTL)

		// Initialize formatter (auto-detect if piped)
		format := output.Format(cfg.Format)
		if !isTerminal() && formatFlag == "" {
			format = output.FormatJSON
		}
		formatter = output.New(format, os.Stdout)

		return nil
	},
}

// Execute runs the root command
func Execute(version string) error {
	Version = version
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	// Global persistent flags
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "profile", "P", "", "Named profile to use")
	rootCmd.PersistentFlags().StringVar(&apiKeyFlag, "api-key", "", "Override API key for this invocation")
	rootCmd.PersistentFlags().StringVarP(&teamFlag, "team", "t", "", "Team key context (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "", "Output format: table, json, csv, plain")
	rootCmd.PersistentFlags().StringVar(&jsonFlag, "json", "", "Output specific fields as JSON (comma-separated)")
	rootCmd.PersistentFlags().StringVar(&jqFlag, "jq", "", "Apply jq expression to JSON output")
	rootCmd.PersistentFlags().BoolVar(&noCacheFlag, "no-cache", false, "Bypass cached data")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Debug output")

	// Bind flags to viper
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("team", rootCmd.PersistentFlags().Lookup("team"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	viper.BindPFlag("no-cache", rootCmd.PersistentFlags().Lookup("no-cache"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// getClient returns an authenticated Linear API client
func getClient() (*client.Client, error) {
	if apiClient != nil {
		return apiClient, nil
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("not authenticated - run 'lirt auth login' to set up credentials")
	}

	var err error
	apiClient, err = client.New(cfg.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return apiClient, nil
}

// getContext returns a context for API calls
func getContext() context.Context {
	return context.Background()
}

// ExitCode constants
const (
	ExitSuccess         = 0
	ExitError           = 1
	ExitUsageError      = 2
	ExitAuthError       = 3
	ExitNotFound        = 4
)
