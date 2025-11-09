package navigation

import (
	"fmt"
	"time"

	"github.com/kamilkurek/margonem-bot/internal/behavior"
	"github.com/kamilkurek/margonem-bot/internal/config"
	"github.com/kamilkurek/margonem-bot/internal/game"
	"github.com/sirupsen/logrus"
)

// Navigator handles waypoint-based navigation
type Navigator struct {
	gameClient *game.Client
	cfg        *config.Config
	log        *logrus.Logger
}

// NewNavigator creates a new navigator
func NewNavigator(gameClient *game.Client, cfg *config.Config, log *logrus.Logger) *Navigator {
	return &Navigator{
		gameClient: gameClient,
		cfg:        cfg,
		log:        log,
	}
}

// GoToHuntingGround navigates to the configured hunting ground
func (n *Navigator) GoToHuntingGround(stateMgr *game.StateManager) error {
	n.log.Info("Navigating to hunting ground...")
	
	hero := stateMgr.GetHero()
	target := n.cfg.Profile.HuntingGround
	
	// Check if already at hunting ground
	if hero.MapID == target.MapID {
		dist := behavior.Distance(
			behavior.Point{X: hero.X, Y: hero.Y},
			behavior.Point{X: target.CenterX, Y: target.CenterY},
		)
		
		if dist <= target.Radius {
			n.log.Info("Already at hunting ground")
			return nil
		}
	}
	
	// Follow waypoints
	for i, wp := range n.cfg.Profile.Waypoints {
		n.log.WithFields(logrus.Fields{
			"waypoint": i + 1,
			"total":    len(n.cfg.Profile.Waypoints),
			"desc":     wp.Description,
		}).Info("Following waypoint")
		
		if err := n.followWaypoint(wp, stateMgr); err != nil {
			return fmt.Errorf("failed to follow waypoint %d: %w", i, err)
		}
		
		// Random delay between waypoints
		behavior.SleepRange(
			n.cfg.Behavior.GetMinDelay(),
			n.cfg.Behavior.GetMaxDelay(),
		)
	}
	
	n.log.Info("Arrived at hunting ground")
	return nil
}

// followWaypoint navigates to a single waypoint
func (n *Navigator) followWaypoint(wp config.Waypoint, stateMgr *game.StateManager) error {
	hero := stateMgr.GetHero()
	
	// Check if we need to change maps first
	if hero.MapID != wp.MapID && wp.Action == "portal" {
		n.log.WithField("map", wp.MapID).Debug("Using portal to change map")
		
		// Click portal if selector provided
		if wp.Selector != "" {
			// Try to click the portal selector
			// This is a simplified version - might need more robust handling
			n.log.WithField("selector", wp.Selector).Debug("Clicking portal")
		}
		
		// Wait for map change
		timeout := 10 * time.Second
		start := time.Now()
		
		for time.Since(start) < timeout {
			hero = stateMgr.GetHero()
			if hero.MapID == wp.MapID {
				n.log.Debug("Map changed successfully")
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		
		if hero.MapID != wp.MapID {
			return fmt.Errorf("failed to change map after portal")
		}
	}
	
	// Move to waypoint coordinates with jitter
	targetPos := behavior.Point{X: wp.X, Y: wp.Y}
	jitteredPos := behavior.AddJitter(targetPos, n.cfg.Behavior.PathJitter)
	
	n.log.WithFields(logrus.Fields{
		"x": jitteredPos.X,
		"y": jitteredPos.Y,
	}).Debug("Moving to waypoint")
	
	if err := n.gameClient.MoveTo(jitteredPos.X, jitteredPos.Y); err != nil {
		return fmt.Errorf("failed to move to waypoint: %w", err)
	}
	
	// Wait for movement to complete
	time.Sleep(2 * time.Second)
	
	return nil
}

// PatrolArea moves around the hunting ground
func (n *Navigator) PatrolArea(stateMgr *game.StateManager) error {
	hero := stateMgr.GetHero()
	huntingGround := n.cfg.Profile.HuntingGround
	
	// Generate a random position within the hunting ground
	center := behavior.Point{X: huntingGround.CenterX, Y: huntingGround.CenterY}
	offset := behavior.RandomOffset(huntingGround.Radius * 0.8) // Stay within 80% of radius
	
	targetPos := behavior.Point{
		X: center.X + offset.X,
		Y: center.Y + offset.Y,
	}
	
	n.log.WithFields(logrus.Fields{
		"x": targetPos.X,
		"y": targetPos.Y,
	}).Debug("Patrolling to position")
	
	// Check if we're too far from current position
	currentPos := behavior.Point{X: hero.X, Y: hero.Y}
	dist := behavior.Distance(currentPos, targetPos)
	
	if dist > 150 { // Don't make huge jumps
		// Move in steps
		path := behavior.GeneratePath(currentPos, targetPos, 50)
		
		for _, point := range path {
			if err := n.gameClient.MoveTo(point.X, point.Y); err != nil {
				return fmt.Errorf("failed to patrol: %w", err)
			}
			
			behavior.SleepRange(
				n.cfg.Behavior.GetMinDelay()/2,
				n.cfg.Behavior.GetMaxDelay()/2,
			)
		}
	} else {
		// Move directly
		if err := n.gameClient.MoveTo(targetPos.X, targetPos.Y); err != nil {
			return fmt.Errorf("failed to patrol: %w", err)
		}
	}
	
	return nil
}

// ReturnFromDeath navigates from respawn point to hunting ground
func (n *Navigator) ReturnFromDeath(stateMgr *game.StateManager) error {
	n.log.Info("Returning from death to hunting ground...")
	
	// Wait a bit after respawn
	behavior.LongPause()
	
	// Use same logic as GoToHuntingGround
	return n.GoToHuntingGround(stateMgr)
}
