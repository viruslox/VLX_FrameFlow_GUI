package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	ScriptPath string
}

func LoadConfig() *Config {
	cfg := &Config{}

	// Default path from issue description
	defaultSettingsPath := "/opt/VLX_FrameFlow/config/FrameFlow_user.settings"

	// Check if default setting file exists
	if _, err := os.Stat(defaultSettingsPath); err == nil {
		cfg.ScriptPath = "/opt/VLX_FrameFlow"
		return cfg
	}

	// Fallback to VLX_PATH environment variable
	if path := os.Getenv("VLX_PATH"); path != "" {
		cfg.ScriptPath = path
		return cfg
	}

	// Fallback to local scripts folder for testing
	pwd, err := os.Getwd()
	if err == nil {
		cfg.ScriptPath = filepath.Join(pwd, "scripts")
	} else {
		cfg.ScriptPath = "./scripts"
	}

	return cfg
}
