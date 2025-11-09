package combat

import (
	"math"

	"github.com/kamilkurek/margonem-bot/internal/config"
	"github.com/kamilkurek/margonem-bot/internal/game"
)

// TargetScore represents a mob with its score
type TargetScore struct {
	Mob      *game.Mob
	Score    float64
	Distance float64
}

// SelectTarget finds the best mob to attack based on configuration
func SelectTarget(hero *game.HeroState, mobs []*game.Mob, cfg *config.CombatConfig) *game.Mob {
	if len(mobs) == 0 {
		return nil
	}
	
	candidates := make([]TargetScore, 0)
	
	for _, mob := range mobs {
		// Filter by basic criteria
		if !mob.Alive || !mob.Attackable {
			continue
		}
		
		// Filter by level range
		if cfg.MinLevel > 0 && mob.Level < cfg.MinLevel {
			continue
		}
		if cfg.MaxLevel > 0 && mob.Level > cfg.MaxLevel {
			continue
		}
		
		// Calculate distance
		dist := distance(hero.X, hero.Y, mob.X, mob.Y)
		
		// Filter by max engage distance
		if cfg.MaxEngageDistance > 0 && dist > cfg.MaxEngageDistance {
			continue
		}
		
		// Calculate score
		score := scoreMob(mob, dist, cfg.TargetPriority)
		
		candidates = append(candidates, TargetScore{
			Mob:      mob,
			Score:    score,
			Distance: dist,
		})
	}
	
	if len(candidates) == 0 {
		return nil
	}
	
	// Find highest score
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.Score > best.Score {
			best = c
		}
	}
	
	return best.Mob
}

// scoreMob calculates a score for a mob based on priority and distance
func scoreMob(mob *game.Mob, distance float64, priority []string) float64 {
	score := 100.0
	
	// Priority bonus (higher priority = higher score)
	for i, name := range priority {
		if mob.Name == name {
			score += float64(len(priority)-i) * 50.0
			break
		}
	}
	
	// Distance penalty (closer = better)
	// Normalize distance to 0-100 range assuming max distance is 300
	distPenalty := (distance / 300.0) * 50.0
	score -= distPenalty
	
	// Level consideration (prefer similar level)
	// This is a simple approach - can be made more sophisticated
	
	return score
}

// distance calculates Euclidean distance
func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// FindNearestMob finds the closest attackable mob
func FindNearestMob(hero *game.HeroState, mobs []*game.Mob) *game.Mob {
	if len(mobs) == 0 {
		return nil
	}
	
	var nearest *game.Mob
	minDist := math.MaxFloat64
	
	for _, mob := range mobs {
		if !mob.Alive || !mob.Attackable {
			continue
		}
		
		dist := distance(hero.X, hero.Y, mob.X, mob.Y)
		if dist < minDist {
			minDist = dist
			nearest = mob
		}
	}
	
	return nearest
}
