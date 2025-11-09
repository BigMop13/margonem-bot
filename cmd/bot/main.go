package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kamilkurek/margonem-bot/internal/behavior"
	"github.com/kamilkurek/margonem-bot/internal/browser"
	"github.com/kamilkurek/margonem-bot/internal/combat"
	"github.com/kamilkurek/margonem-bot/internal/config"
	"github.com/kamilkurek/margonem-bot/internal/game"
	"github.com/kamilkurek/margonem-bot/internal/navigation"
	"github.com/sirupsen/logrus"
)

var (
	version    = "dev"
	configPath = flag.String("config", "configs/config.yaml", "Path to config file")
	showVersion = flag.Bool("version", false, "Show version")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("Margonem Bot %s\n", version)
		return
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Setup logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting Margonem Bot...")
	log.WithField("version", version).Info("Bot version")

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	if cfg.Runtime.Debug {
		log.SetLevel(logrus.DebugLevel)
		log.Debug("Debug mode enabled")
	}

	log.WithField("profile", cfg.Profile.Name).Info("Loaded profile")

	// Create screenshot directory
	if err := os.MkdirAll(cfg.Runtime.ScreenshotDir, 0755); err != nil {
		log.WithError(err).Warn("Failed to create screenshot directory")
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Shutdown signal received")
		cancel()
	}()

	// Run the bot
	if err := run(ctx, cfg, log); err != nil {
		log.WithError(err).Fatal("Bot execution failed")
	}

	log.Info("Bot stopped successfully")
}

func run(ctx context.Context, cfg *config.Config, log *logrus.Logger) error {
	// Initialize browser
	browserCtrl, err := browser.New(
		cfg.Runtime.Headless,
		cfg.Runtime.UserDataDir,
		cfg.Runtime.ViewportWidth,
		cfg.Runtime.ViewportHeight,
		cfg.Runtime.Debug,
		log,
	)
	if err != nil {
		return fmt.Errorf("failed to create browser controller: %w", err)
	}

	if err := browserCtrl.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %w", err)
	}
	defer browserCtrl.Stop()

	// Initialize components
	gameClient := game.NewClient(browserCtrl, log)
	stateMgr := game.NewStateManager()
	combatEngine := combat.NewEngine(gameClient, cfg, log)
	navigator := navigation.NewNavigator(gameClient, cfg, log)

	// State machine
	phase := game.PhaseLogin

	// Login
	log.Info("Phase: LOGIN")
	if err := performLogin(browserCtrl, cfg, log); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Wait for game ready
	log.Info("Phase: WAIT_GAME_READY")
	phase = game.PhaseWaitGameReady
	stateMgr.SetPhase(phase)

	if err := gameClient.EnsureReady(); err != nil {
		return fmt.Errorf("game not ready: %w", err)
	}

	// Start state polling
	go pollGameState(ctx, gameClient, stateMgr, log)

	// Give it a moment to gather initial state
	time.Sleep(2 * time.Second)

	// Auto-detect mode: just detect current location, skip navigation
	if cfg.Runtime.AutoDetectMode {
		log.Info("AUTO-DETECT MODE: Bot will hunt at current location")
		hero := stateMgr.GetHero()
		log.WithFields(logrus.Fields{
			"mapId": hero.MapID,
			"x":     hero.X,
			"y":     hero.Y,
		}).Info("Current position detected")
		
		// Auto-populate hunting ground config based on current position
		cfg.Profile.HuntingGround.MapID = hero.MapID
		cfg.Profile.HuntingGround.CenterX = hero.X
		cfg.Profile.HuntingGround.CenterY = hero.Y
		cfg.Profile.HuntingGround.Radius = 200 // Default 200 pixel radius
		
		log.Info("Hunting ground auto-detected, starting combat!")
	} else {
		// Manual mode: Navigate to configured hunting ground
		log.Info("Phase: NAVIGATE")
		phase = game.PhaseNavigate
		stateMgr.SetPhase(phase)

		if err := navigator.GoToHuntingGround(stateMgr); err != nil {
			return fmt.Errorf("navigation failed: %w", err)
		}
	}

	// Main bot loop
	log.Info("Phase: HUNT")
	phase = game.PhaseHunt
	stateMgr.SetPhase(phase)

	ticker := time.NewTicker(2 * time.Second) // Combat tick every 2 seconds
	defer ticker.Stop()

	lastPatrol := time.Now()
	patrolInterval := 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			log.Info("Context cancelled, stopping bot")
			return nil

		case <-ticker.C:
			hero := stateMgr.GetHero()

			// Check if dead
			if hero.Dead {
				log.Warn("Phase: DEAD")
				
				if cfg.Runtime.AutoDetectMode {
					// In auto-detect mode, just wait and respawn
					log.Warn("Character died! Waiting for respawn...")
					time.Sleep(5 * time.Second)
					
					// Try to respawn
					if err := gameClient.Respawn(); err != nil {
						log.WithError(err).Warn("Respawn failed")
					}
					
					log.Info("Respawned - return to your hunting ground manually!")
					log.Info("Press Ctrl+C to stop, or wait here...")
					time.Sleep(30 * time.Second) // Wait for manual return
					
					// Update hunting ground to new location after respawn
					hero = stateMgr.GetHero()
					cfg.Profile.HuntingGround.MapID = hero.MapID
					cfg.Profile.HuntingGround.CenterX = hero.X
					cfg.Profile.HuntingGround.CenterY = hero.Y
				} else {
					phase = game.PhaseDead
					stateMgr.SetPhase(phase)

					if err := handleDeath(gameClient, navigator, stateMgr, log); err != nil {
						log.WithError(err).Error("Failed to handle death")
					}
				}

				log.Info("Phase: HUNT")
				phase = game.PhaseHunt
				stateMgr.SetPhase(phase)
				combatEngine.Reset()
				continue
			}

			// Check connection
			if connected, err := gameClient.IsConnected(); err != nil || !connected {
				log.Warn("Phase: DISCONNECTED")
				phase = game.PhaseDisconnected
				stateMgr.SetPhase(phase)

				if err := handleDisconnect(browserCtrl, gameClient, cfg, log); err != nil {
					return fmt.Errorf("reconnection failed: %w", err)
				}

				log.Info("Phase: HUNT")
				phase = game.PhaseHunt
				stateMgr.SetPhase(phase)
				continue
			}

			// Check if stuck
			if stateMgr.IsStuck(10, 30*time.Second) {
				log.Warn("Character appears stuck, attempting recovery")
				if err := navigator.PatrolArea(stateMgr); err != nil {
					log.WithError(err).Warn("Failed to recover from stuck state")
				}
			}

			// Periodic patrol to find mobs
			if time.Since(lastPatrol) > patrolInterval {
				mobs := stateMgr.GetMobs()
				if len(mobs) == 0 {
					log.Debug("No mobs nearby, patrolling")
					if err := navigator.PatrolArea(stateMgr); err != nil {
						log.WithError(err).Warn("Patrol failed")
					}
				}
				lastPatrol = time.Now()
			}

			// Combat tick
			if err := combatEngine.Tick(stateMgr); err != nil {
				log.WithError(err).Warn("Combat tick failed")
			}

			// Check for idle breaks
			actionCount := stateMgr.IncrementAction()
			if behavior.ShouldTakeBreak(actionCount, cfg.Behavior.IdleBreakEvery) {
				log.Info("Taking idle break")
				time.Sleep(time.Duration(cfg.Behavior.IdleBreakDuration) * time.Second)
			}
		}
	}
}

// performLogin logs into the game
func performLogin(browserCtrl *browser.Controller, cfg *config.Config, log *logrus.Logger) error {
	log.Info("Navigating to game...")

	if err := browserCtrl.Navigate(cfg.Account.StartURL); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}

	if err := browserCtrl.WaitReady(); err != nil {
		return fmt.Errorf("page not ready: %w", err)
	}

	// Wait a bit for page to fully load
	time.Sleep(3 * time.Second)

	// Try to find login form
	// This is a simplified version - actual Margonem login might be different
	log.Info("Attempting login...")

	// Check if already logged in by looking for game canvas
	script := `
	(function() {
		return window.hero !== undefined && window.hero !== null;
	})()
	`

	var alreadyLoggedIn bool
	if err := browserCtrl.Eval(script, &alreadyLoggedIn); err == nil && alreadyLoggedIn {
		log.Info("Already logged in (using saved session)")
		return nil
	}

	// Perform login
	// NOTE: This is a placeholder - actual Margonem login flow will need to be adapted
	loginScript := fmt.Sprintf(`
	(function() {
		try {
			// Try to fill login form
			let userField = document.querySelector('input[name="login"], input[name="username"], #login, #username');
			let passField = document.querySelector('input[name="password"], input[type="password"], #password');
			let loginBtn = document.querySelector('button[type="submit"], input[type="submit"], .login-button, #login-button');

			if (userField && passField) {
				userField.value = '%s';
				passField.value = '%s';
				if (loginBtn) {
					loginBtn.click();
					return true;
				}
			}
			return false;
		} catch(e) {
			console.error('Login error:', e);
			return false;
		}
	})()
	`, cfg.Account.Username, cfg.Account.Password)

	var loginSuccess bool
	if err := browserCtrl.Eval(loginScript, &loginSuccess); err != nil {
		log.Warn("Automated login failed, manual intervention may be required")
	}

	// Wait for login to complete
	time.Sleep(5 * time.Second)

	// Select server if needed
	if cfg.Account.Server != "" {
		log.WithField("server", cfg.Account.Server).Info("Selecting server...")
		// Server selection logic would go here
		time.Sleep(2 * time.Second)
	}

	log.Info("Login complete")
	return nil
}

// pollGameState continuously updates game state
func pollGameState(ctx context.Context, gameClient *game.Client, stateMgr *game.StateManager, log *logrus.Logger) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Get hero state
			hero, err := gameClient.GetHeroState()
			if err != nil {
				log.WithError(err).Debug("Failed to get hero state")
				continue
			}
			stateMgr.UpdateHero(hero)

			// Get mobs
			mobs, err := gameClient.GetMobs()
			if err != nil {
				log.WithError(err).Debug("Failed to get mobs")
				continue
			}
			stateMgr.UpdateMobs(mobs)

			// Check connection
			connected, err := gameClient.IsConnected()
			if err != nil {
				connected = false
			}
			stateMgr.UpdateConnection(connected)
		}
	}
}

// handleDeath handles character death and respawn
func handleDeath(gameClient *game.Client, navigator *navigation.Navigator, stateMgr *game.StateManager, log *logrus.Logger) error {
	log.Info("Handling death...")

	// Wait a moment
	time.Sleep(2 * time.Second)

	// Respawn
	if err := gameClient.Respawn(); err != nil {
		return fmt.Errorf("respawn failed: %w", err)
	}

	log.Info("Respawned successfully")

	// Wait for respawn to complete
	time.Sleep(3 * time.Second)

	// Return to hunting ground
	log.Info("Phase: RECOVER")
	stateMgr.SetPhase(game.PhaseRecover)

	if err := navigator.ReturnFromDeath(stateMgr); err != nil {
		return fmt.Errorf("failed to return to hunting ground: %w", err)
	}

	log.Info("Recovery complete")
	return nil
}

// handleDisconnect handles disconnection and reconnection
func handleDisconnect(browserCtrl *browser.Controller, gameClient *game.Client, cfg *config.Config, log *logrus.Logger) error {
	log.Warn("Handling disconnection...")

	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		log.WithField("attempt", attempt+1).Info("Attempting to reconnect...")

		// Wait with exponential backoff
		backoff := behavior.Backoff(5*time.Second, 2.0, 60*time.Second, attempt)
		time.Sleep(backoff)

		// Try to reload page
		if err := browserCtrl.Navigate(cfg.Account.StartURL); err != nil {
			log.WithError(err).Warn("Failed to navigate")
			continue
		}

		// Wait for game to be ready
		if err := gameClient.EnsureReady(); err != nil {
			log.WithError(err).Warn("Game not ready after reload")
			continue
		}

		log.Info("Reconnected successfully")
		return nil
	}

	return fmt.Errorf("failed to reconnect after %d attempts", maxRetries)
}
