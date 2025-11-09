package config

import "time"

// Config represents the complete bot configuration
type Config struct {
	Account  AccountConfig  `yaml:"account"`
	Profile  ProfileConfig  `yaml:"profile,omitempty"`
	Combat   CombatConfig   `yaml:"combat"`
	Behavior BehaviorConfig `yaml:"behavior"`
	Runtime  RuntimeConfig  `yaml:"runtime"`
}

// AccountConfig holds account credentials and connection info
type AccountConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	StartURL string `yaml:"startUrl"`
}

// ProfileConfig defines the hunting profile
type ProfileConfig struct {
	Name          string         `yaml:"name"`
	HuntingGround HuntingGround  `yaml:"huntingGround"`
	Waypoints     []Waypoint     `yaml:"waypoints"`
	TownRespawn   RespawnPoint   `yaml:"townRespawn"`
}

// HuntingGround defines the area to hunt in
type HuntingGround struct {
	MapID  string  `yaml:"mapId"`
	CenterX float64 `yaml:"centerX"`
	CenterY float64 `yaml:"centerY"`
	Radius  float64 `yaml:"radius"`
}

// Waypoint represents a navigation point
type Waypoint struct {
	MapID       string  `yaml:"mapId"`
	X           float64 `yaml:"x"`
	Y           float64 `yaml:"y"`
	Description string  `yaml:"description"`
	Action      string  `yaml:"action"` // "walk", "portal", "door"
	Selector    string  `yaml:"selector,omitempty"`
}

// RespawnPoint defines where the character respawns
type RespawnPoint struct {
	MapID string  `yaml:"mapId"`
	X     float64 `yaml:"x"`
	Y     float64 `yaml:"y"`
}

// CombatConfig defines combat behavior
type CombatConfig struct {
	HPThreshold       int      `yaml:"hpThreshold"`       // % HP to retreat
	MPThreshold       int      `yaml:"mpThreshold"`       // % MP to consider
	TargetPriority    []string `yaml:"targetPriority"`    // mob names by priority
	RetargetOnDeath   bool     `yaml:"retargetOnDeath"`   // find new target immediately
	MaxEngageDistance float64  `yaml:"maxEngageDistance"` // max distance to chase
	MinLevel          int      `yaml:"minLevel"`          // min mob level
	MaxLevel          int      `yaml:"maxLevel"`          // max mob level
}

// BehaviorConfig defines human-like behavior patterns
type BehaviorConfig struct {
	MinDelayMs       int     `yaml:"minDelayMs"`       // min delay between actions
	MaxDelayMs       int     `yaml:"maxDelayMs"`       // max delay between actions
	MouseSpeedRange  int     `yaml:"mouseSpeedRange"`  // variance in mouse speed (ms)
	PathJitter       float64 `yaml:"pathJitter"`       // random offset for waypoints (pixels)
	IdleBreakEvery   int     `yaml:"idleBreakEvery"`   // idle break every N actions
	IdleBreakDuration int    `yaml:"idleBreakDuration"` // idle break duration (seconds)
}


// RuntimeConfig defines runtime behavior
type RuntimeConfig struct {
	Headless       bool   `yaml:"headless"`
	Debug          bool   `yaml:"debug"`
	UserDataDir    string `yaml:"userDataDir,omitempty"`
	ViewportWidth  int    `yaml:"viewportWidth"`
	ViewportHeight int    `yaml:"viewportHeight"`
	ScreenshotDir  string `yaml:"screenshotDir"`
	AutoDetectMode bool   `yaml:"autoDetectMode"` // Auto-detect location and mobs
}

// GetMinDelay returns minimum delay as duration
func (b *BehaviorConfig) GetMinDelay() time.Duration {
	return time.Duration(b.MinDelayMs) * time.Millisecond
}

// GetMaxDelay returns maximum delay as duration
func (b *BehaviorConfig) GetMaxDelay() time.Duration {
	return time.Duration(b.MaxDelayMs) * time.Millisecond
}


// SetDefaults applies default values to the configuration
func (c *Config) SetDefaults() {
	if c.Runtime.ViewportWidth == 0 {
		c.Runtime.ViewportWidth = 1280
	}
	if c.Runtime.ViewportHeight == 0 {
		c.Runtime.ViewportHeight = 720
	}
	if c.Runtime.ScreenshotDir == "" {
		c.Runtime.ScreenshotDir = "./screenshots"
	}
	if c.Behavior.MinDelayMs == 0 {
		c.Behavior.MinDelayMs = 1000
	}
	if c.Behavior.MaxDelayMs == 0 {
		c.Behavior.MaxDelayMs = 2000
	}
	if c.Behavior.MouseSpeedRange == 0 {
		c.Behavior.MouseSpeedRange = 200
	}
	if c.Behavior.PathJitter == 0 {
		c.Behavior.PathJitter = 5.0
	}
	if c.Combat.MaxEngageDistance == 0 {
		c.Combat.MaxEngageDistance = 200.0
	}
	if c.Combat.HPThreshold == 0 {
		c.Combat.HPThreshold = 30
	}
}
