package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// Config represents lirt configuration
type Config struct {
	Profile   string
	APIKey    string
	Team      string
	Format    string
	CacheTTL  string
	PageSize  int
	Workspace string // Display-only, set by auth login
}

// GetConfigDir returns the lirt config directory
func GetConfigDir() string {
	if dir := os.Getenv("LIRT_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".config/lirt"
	}
	return filepath.Join(home, ".config", "lirt")
}

// GetCredentialsFile returns the path to the credentials file
func GetCredentialsFile() string {
	if file := os.Getenv("LIRT_CREDENTIALS_FILE"); file != "" {
		return file
	}
	return filepath.Join(GetConfigDir(), "credentials")
}

// GetConfigFile returns the path to the config file
func GetConfigFile() string {
	if file := os.Getenv("LIRT_CONFIG_FILE"); file != "" {
		return file
	}
	return filepath.Join(GetConfigDir(), "config")
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	dir := GetConfigDir()
	return os.MkdirAll(dir, 0700)
}

// LoadConfig loads configuration for the given profile
func LoadConfig(profile string) (*Config, error) {
	cfg := &Config{
		Profile:  profile,
		Format:   "table",
		CacheTTL: "5m",
		PageSize: 50,
	}

	// Load config file
	configFile := GetConfigFile()
	if _, err := os.Stat(configFile); err == nil {
		iniFile, err := ini.Load(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}

		section := "default"
		if profile != "default" {
			section = "profile " + profile
		}

		if iniFile.HasSection(section) {
			sec := iniFile.Section(section)
			if sec.HasKey("workspace") {
				cfg.Workspace = sec.Key("workspace").String()
			}
			if sec.HasKey("team") {
				cfg.Team = sec.Key("team").String()
			}
			if sec.HasKey("format") {
				cfg.Format = sec.Key("format").String()
			}
			if sec.HasKey("cache_ttl") {
				cfg.CacheTTL = sec.Key("cache_ttl").String()
			}
			if sec.HasKey("page_size") {
				cfg.PageSize, _ = sec.Key("page_size").Int()
			}
		}
	}

	// Load API key from credentials
	apiKey, err := LoadAPIKey(profile)
	if err == nil {
		cfg.APIKey = apiKey
	}

	return cfg, nil
}

// LoadAPIKey loads the API key for the given profile
// Resolution order: LIRT_API_KEY, --api-key flag (handled by caller), credentials file, LINEAR_API_KEY
func LoadAPIKey(profile string) (string, error) {
	// Check LIRT_API_KEY env var (highest priority)
	if key := os.Getenv("LIRT_API_KEY"); key != "" {
		return key, nil
	}

	// Load from credentials file
	credFile := GetCredentialsFile()
	if _, err := os.Stat(credFile); err == nil {
		iniFile, err := ini.Load(credFile)
		if err != nil {
			return "", fmt.Errorf("failed to load credentials file: %w", err)
		}

		if iniFile.HasSection(profile) {
			sec := iniFile.Section(profile)
			if sec.HasKey("api_key") {
				return sec.Key("api_key").String(), nil
			}
		}
	}

	// Check LINEAR_API_KEY env var (fallback)
	if key := os.Getenv("LINEAR_API_KEY"); key != "" {
		return key, nil
	}

	return "", fmt.Errorf("no API key found for profile %q", profile)
}

// SaveAPIKey saves an API key to the credentials file
func SaveAPIKey(profile, apiKey string) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	credFile := GetCredentialsFile()

	var iniFile *ini.File
	if _, err := os.Stat(credFile); err == nil {
		iniFile, err = ini.Load(credFile)
		if err != nil {
			return fmt.Errorf("failed to load credentials file: %w", err)
		}
	} else {
		iniFile = ini.Empty()
	}

	section, err := iniFile.NewSection(profile)
	if err != nil {
		section = iniFile.Section(profile)
	}

	section.Key("api_key").SetValue(apiKey)

	if err := iniFile.SaveTo(credFile); err != nil {
		return fmt.Errorf("failed to save credentials file: %w", err)
	}

	// Set secure permissions
	return os.Chmod(credFile, 0600)
}

// SaveConfigValue saves a config value to the config file
func SaveConfigValue(profile, key, value string) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configFile := GetConfigFile()

	var iniFile *ini.File
	if _, err := os.Stat(configFile); err == nil {
		iniFile, err = ini.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
	} else {
		iniFile = ini.Empty()
	}

	sectionName := "default"
	if profile != "default" {
		sectionName = "profile " + profile
	}

	section, err := iniFile.NewSection(sectionName)
	if err != nil {
		section = iniFile.Section(sectionName)
	}

	section.Key(key).SetValue(value)

	if err := iniFile.SaveTo(configFile); err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	// Set readable permissions for config
	return os.Chmod(configFile, 0644)
}

// DeleteProfile removes a profile from both credentials and config files
func DeleteProfile(profile string) error {
	// Remove from credentials
	credFile := GetCredentialsFile()
	if _, err := os.Stat(credFile); err == nil {
		iniFile, err := ini.Load(credFile)
		if err != nil {
			return fmt.Errorf("failed to load credentials file: %w", err)
		}
		iniFile.DeleteSection(profile)
		if err := iniFile.SaveTo(credFile); err != nil {
			return fmt.Errorf("failed to save credentials file: %w", err)
		}
	}

	// Remove from config
	configFile := GetConfigFile()
	if _, err := os.Stat(configFile); err == nil {
		iniFile, err := ini.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
		sectionName := "default"
		if profile != "default" {
			sectionName = "profile " + profile
		}
		iniFile.DeleteSection(sectionName)
		if err := iniFile.SaveTo(configFile); err != nil {
			return fmt.Errorf("failed to save config file: %w", err)
		}
	}

	return nil
}

// ListProfiles returns all configured profiles
func ListProfiles() (map[string]string, error) {
	profiles := make(map[string]string)

	configFile := GetConfigFile()
	if _, err := os.Stat(configFile); err == nil {
		iniFile, err := ini.Load(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}

		for _, section := range iniFile.Sections() {
			name := section.Name()
			if name == ini.DefaultSection || name == "" {
				continue
			}

			profile := name
			if len(name) > 8 && name[:8] == "profile " {
				profile = name[8:]
			}

			workspace := ""
			if section.HasKey("workspace") {
				workspace = section.Key("workspace").String()
			}
			profiles[profile] = workspace
		}

		// Check default section
		if iniFile.HasSection("default") {
			sec := iniFile.Section("default")
			workspace := ""
			if sec.HasKey("workspace") {
				workspace = sec.Key("workspace").String()
			}
			profiles["default"] = workspace
		}
	}

	return profiles, nil
}

// GetProfile returns the profile name to use based on flags and env vars
func GetProfile(profileFlag string) string {
	if profileFlag != "" {
		return profileFlag
	}
	if profile := os.Getenv("LIRT_PROFILE"); profile != "" {
		return profile
	}
	return "default"
}
