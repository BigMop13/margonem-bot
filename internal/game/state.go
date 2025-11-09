package game

import (
	"sync"
	"time"
)

// BotPhase represents the current bot state
type BotPhase int

const (
	PhaseStartup BotPhase = iota
	PhaseLogin
	PhaseWaitGameReady
	PhaseNavigate
	PhaseHunt
	PhaseDead
	PhaseRecover
	PhaseDisconnected
	PhaseShutdown
)

func (p BotPhase) String() string {
	switch p {
	case PhaseStartup:
		return "STARTUP"
	case PhaseLogin:
		return "LOGIN"
	case PhaseWaitGameReady:
		return "WAIT_GAME_READY"
	case PhaseNavigate:
		return "NAVIGATE"
	case PhaseHunt:
		return "HUNT"
	case PhaseDead:
		return "DEAD"
	case PhaseRecover:
		return "RECOVER"
	case PhaseDisconnected:
		return "DISCONNECTED"
	case PhaseShutdown:
		return "SHUTDOWN"
	default:
		return "UNKNOWN"
	}
}

// HeroState represents the character's current state
type HeroState struct {
	X        float64
	Y        float64
	MapID    string
	HP       int
	HPMax    int
	MP       int
	MPMax    int
	Level    int
	Exp      int
	InCombat bool
	Dead     bool
	LastUpdate time.Time
}

// HPPercent returns HP as a percentage
func (h *HeroState) HPPercent() int {
	if h.HPMax == 0 {
		return 0
	}
	return (h.HP * 100) / h.HPMax
}

// MPPercent returns MP as a percentage
func (h *HeroState) MPPercent() int {
	if h.MPMax == 0 {
		return 0
	}
	return (h.MP * 100) / h.MPMax
}

// Mob represents a monster in the game
type Mob struct {
	ID         string
	Name       string
	Level      int
	X          float64
	Y          float64
	HP         int
	HPMax      int
	Alive      bool
	Attackable bool
}

// ConnectionState represents connection status
type ConnectionState struct {
	Connected  bool
	LastCheck  time.Time
	Retries    int
}

// StateManager manages game state with thread safety
type StateManager struct {
	mu              sync.RWMutex
	hero            *HeroState
	mobs            []*Mob
	connection      ConnectionState
	phase           BotPhase
	positionHistory []PositionRecord
	actionCount     int
}

// PositionRecord tracks position for stuck detection
type PositionRecord struct {
	X         float64
	Y         float64
	Timestamp time.Time
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		hero:            &HeroState{},
		mobs:            make([]*Mob, 0),
		connection:      ConnectionState{Connected: true},
		phase:           PhaseStartup,
		positionHistory: make([]PositionRecord, 0, 10),
	}
}

// UpdateHero updates hero state
func (sm *StateManager) UpdateHero(hero *HeroState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	hero.LastUpdate = time.Now()
	sm.hero = hero
	
	// Track position for stuck detection
	sm.positionHistory = append(sm.positionHistory, PositionRecord{
		X:         hero.X,
		Y:         hero.Y,
		Timestamp: time.Now(),
	})
	
	// Keep only last 10 positions
	if len(sm.positionHistory) > 10 {
		sm.positionHistory = sm.positionHistory[1:]
	}
}

// GetHero returns a copy of the hero state
func (sm *StateManager) GetHero() HeroState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return *sm.hero
}

// UpdateMobs updates mob list
func (sm *StateManager) UpdateMobs(mobs []*Mob) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.mobs = mobs
}

// GetMobs returns a copy of the mob list
func (sm *StateManager) GetMobs() []*Mob {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	result := make([]*Mob, len(sm.mobs))
	copy(result, sm.mobs)
	return result
}

// SetPhase sets the current bot phase
func (sm *StateManager) SetPhase(phase BotPhase) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.phase = phase
}

// GetPhase returns the current bot phase
func (sm *StateManager) GetPhase() BotPhase {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.phase
}

// IncrementAction increments action counter
func (sm *StateManager) IncrementAction() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.actionCount++
	return sm.actionCount
}

// GetActionCount returns action count
func (sm *StateManager) GetActionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.actionCount
}

// IsStuck checks if character is stuck based on position history
func (sm *StateManager) IsStuck(threshold float64, duration time.Duration) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	if len(sm.positionHistory) < 2 {
		return false
	}
	
	// Check if we've been in the same general area for too long
	recent := sm.positionHistory[len(sm.positionHistory)-1]
	for i := len(sm.positionHistory) - 2; i >= 0; i-- {
		pos := sm.positionHistory[i]
		if time.Since(pos.Timestamp) > duration {
			break
		}
		
		dx := recent.X - pos.X
		dy := recent.Y - pos.Y
		dist := dx*dx + dy*dy
		
		if dist < threshold*threshold {
			return true
		}
	}
	
	return false
}

// UpdateConnection updates connection state
func (sm *StateManager) UpdateConnection(connected bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sm.connection.Connected = connected
	sm.connection.LastCheck = time.Now()
	
	if !connected {
		sm.connection.Retries++
	} else {
		sm.connection.Retries = 0
	}
}

// IsConnected returns connection status
func (sm *StateManager) IsConnected() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.connection.Connected
}

// GetRetries returns reconnection attempt count
func (sm *StateManager) GetRetries() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.connection.Retries
}
