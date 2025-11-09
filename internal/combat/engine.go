package combat

import (
	"fmt"

	"github.com/kamilkurek/margonem-bot/internal/behavior"
	"github.com/kamilkurek/margonem-bot/internal/config"
	"github.com/kamilkurek/margonem-bot/internal/game"
	"github.com/sirupsen/logrus"
)

// Engine manages combat operations
type Engine struct {
	gameClient    *game.Client
	cfg           *config.Config
	log           *logrus.Logger
	currentTarget *game.Mob
}

// NewEngine creates a new combat engine
func NewEngine(gameClient *game.Client, cfg *config.Config, log *logrus.Logger) *Engine {
	return &Engine{
		gameClient: gameClient,
		cfg:        cfg,
		log:        log,
	}
}

// Tick performs one combat cycle
func (e *Engine) Tick(stateMgr *game.StateManager) error {
	hero := stateMgr.GetHero()
	
	// Check if HP is critically low
	if hero.HPPercent() < e.cfg.Combat.HPThreshold {
		e.log.Warn("HP critically low, retreating")
		e.currentTarget = nil
		return e.retreat(&hero)
	}
	
	// Get available mobs
	mobs := stateMgr.GetMobs()
	
	// If we have a current target, check if it's still valid
	if e.currentTarget != nil {
		valid := false
		for _, m := range mobs {
			if m.ID == e.currentTarget.ID && m.Alive && m.Attackable {
				valid = true
				e.currentTarget = m // Update with fresh data
				break
			}
		}
		
		if !valid {
			e.log.Debug("Current target no longer valid")
			e.currentTarget = nil
			
			if e.cfg.Combat.RetargetOnDeath {
				// Immediately find a new target
				e.currentTarget = SelectTarget(&hero, mobs, &e.cfg.Combat)
			}
		}
	}
	
	// If no target, find one
	if e.currentTarget == nil {
		e.currentTarget = SelectTarget(&hero, mobs, &e.cfg.Combat)
		
		if e.currentTarget == nil {
			// No targets available
			return nil
		}
		
		e.log.WithFields(logrus.Fields{
			"mob":   e.currentTarget.Name,
			"level": e.currentTarget.Level,
			"id":    e.currentTarget.ID,
		}).Info("New target acquired")
	}
	
	// Engage the target
	return e.engage(&hero, e.currentTarget)
}

// engage attacks the target
func (e *Engine) engage(hero *game.HeroState, target *game.Mob) error {
	// Calculate distance to target
	dist := distance(hero.X, hero.Y, target.X, target.Y)
	
	e.log.WithFields(logrus.Fields{
		"target":   target.Name,
		"distance": dist,
	}).Debug("Engaging target")
	
	// If too far, move closer first
	if dist > 50 { // Assuming attack range is ~50 units
		e.log.Debug("Moving closer to target")
		
		// Add jitter to movement
		targetPos := behavior.Point{X: target.X, Y: target.Y}
		jitteredPos := behavior.AddJitter(targetPos, e.cfg.Behavior.PathJitter)
		
		if err := e.gameClient.MoveTo(jitteredPos.X, jitteredPos.Y); err != nil {
			return fmt.Errorf("failed to move to target: %w", err)
		}
		
		// Random delay before attacking
		behavior.SleepRange(
			e.cfg.Behavior.GetMinDelay(),
			e.cfg.Behavior.GetMaxDelay(),
		)
	}
	
	// Attack the mob
	e.log.WithField("target", target.Name).Debug("Attacking")
	
	if err := e.gameClient.AttackMob(target.ID); err != nil {
		return fmt.Errorf("failed to attack: %w", err)
	}
	
	// Random delay after attack
	behavior.SleepRange(
		e.cfg.Behavior.GetMinDelay(),
		e.cfg.Behavior.GetMaxDelay(),
	)
	
	return nil
}


// retreat moves the hero away from danger
func (e *Engine) retreat(hero *game.HeroState) error {
	e.log.Warn("Retreating from combat")
	
	// Move to a random nearby position away from current location
	offset := behavior.RandomOffset(100)
	newX := hero.X + offset.X
	newY := hero.Y + offset.Y
	
	if err := e.gameClient.MoveTo(newX, newY); err != nil {
		return fmt.Errorf("failed to retreat: %w", err)
	}
	
	// Wait a bit
	behavior.LongPause()
	
	return nil
}

// Reset resets the combat state
func (e *Engine) Reset() {
	e.currentTarget = nil
}
