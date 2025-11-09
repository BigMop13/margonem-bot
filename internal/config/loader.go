package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.SetDefaults()

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// validate checks the configuration for errors
func validate(cfg *Config) error {
	if cfg.Account.Username == "" {
		return fmt.Errorf("account.username is required")
	}
	if cfg.Account.Password == "" {
		return fmt.Errorf("account.password is required")
	}
	if cfg.Account.StartURL == "" {
		return fmt.Errorf("account.startUrl is required")
	}

	// Profile validation only if not in auto-detect mode
	if !cfg.Runtime.AutoDetectMode {
		if cfg.Profile.Name == "" {
			return fmt.Errorf("profile.name is required (or enable autoDetectMode)")
		}
		if cfg.Profile.HuntingGround.MapID == "" {
			return fmt.Errorf("profile.huntingGround.mapId is required (or enable autoDetectMode)")
		}
		if cfg.Profile.HuntingGround.Radius <= 0 {
			return fmt.Errorf("profile.huntingGround.radius must be positive (or enable autoDetectMode)")
		}
	}

	if cfg.Combat.HPThreshold < 0 || cfg.Combat.HPThreshold > 100 {
		return fmt.Errorf("combat.hpThreshold must be between 0 and 100")
	}
	if cfg.Combat.MPThreshold < 0 || cfg.Combat.MPThreshold > 100 {
		return fmt.Errorf("combat.mpThreshold must be between 0 and 100")
	}


	if cfg.Behavior.MinDelayMs < 0 {
		return fmt.Errorf("behavior.minDelayMs must be non-negative")
	}
	if cfg.Behavior.MaxDelayMs < cfg.Behavior.MinDelayMs {
		return fmt.Errorf("behavior.maxDelayMs must be >= minDelayMs")
	}

	return nil
}
