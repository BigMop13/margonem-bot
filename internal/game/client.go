package game

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kamilkurek/margonem-bot/internal/browser"
	"github.com/sirupsen/logrus"
)

// Client handles game-specific operations
type Client struct {
	browser *browser.Controller
	log     *logrus.Logger
}

// NewClient creates a new game client
func NewClient(browser *browser.Controller, log *logrus.Logger) *Client {
	return &Client{
		browser: browser,
		log:     log,
	}
}

// EnsureReady waits for the game engine to be ready
func (c *Client) EnsureReady() error {
	c.log.Info("Waiting for game engine to be ready...")
	
	script := `
	(function() {
		return (window.hero !== undefined && window.hero !== null) || 
		       (window.Hero !== undefined && window.Hero !== null) ||
		       (window.g && window.g.hero !== undefined) ||
		       (window.Engine !== undefined);
	})()
	`
	
	// Wait up to 30 seconds for game to be ready
	timeout := 30 * time.Second
	start := time.Now()
	
	for time.Since(start) < timeout {
		var ready bool
		if err := c.browser.Eval(script, &ready); err == nil && ready {
			c.log.Info("Game engine is ready!")
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	
	return fmt.Errorf("game engine not ready after %v", timeout)
}

// GetHeroState retrieves the current hero state from the game
func (c *Client) GetHeroState() (*HeroState, error) {
	script := `
	(function() {
		try {
			let hero = window.hero || window.Hero || (window.g && window.g.hero);
			if (!hero) return null;
			
			return {
				x: hero.x || hero.posX || 0,
				y: hero.y || hero.posY || 0,
				mapId: (hero.map && hero.map.id) || hero.mapId || "",
				hp: hero.hp || hero.HP || 0,
				hpMax: hero.maxhp || hero.maxHP || hero.hpMax || 100,
				mp: hero.mp || hero.MP || 0,
				mpMax: hero.maxmp || hero.maxMP || hero.mpMax || 100,
				level: hero.lvl || hero.level || 1,
				exp: hero.exp || 0,
				inCombat: hero.inCombat || hero.incombat || false,
				dead: hero.dead || hero.isDead || hero.hp <= 0
			};
		} catch(e) {
			console.error("Error getting hero state:", e);
			return null;
		}
	})()
	`
	
	var result map[string]interface{}
	if err := c.browser.Eval(script, &result); err != nil {
		return nil, fmt.Errorf("failed to get hero state: %w", err)
	}
	
	if result == nil {
		return nil, fmt.Errorf("hero object not found")
	}
	
	state := &HeroState{
		X:        getFloat(result, "x"),
		Y:        getFloat(result, "y"),
		MapID:    getString(result, "mapId"),
		HP:       getInt(result, "hp"),
		HPMax:    getInt(result, "hpMax"),
		MP:       getInt(result, "mp"),
		MPMax:    getInt(result, "mpMax"),
		Level:    getInt(result, "level"),
		Exp:      getInt(result, "exp"),
		InCombat: getBool(result, "inCombat"),
		Dead:     getBool(result, "dead"),
	}
	
	return state, nil
}

// GetMobs retrieves nearby mobs from the game
func (c *Client) GetMobs() ([]*Mob, error) {
	script := `
	(function() {
		try {
			let npcList = window.npcs || window.NPC || (window.g && window.g.npcs) || [];
			let mobs = [];
			
			for (let id in npcList) {
				let npc = npcList[id];
				if (!npc || npc.type !== 1) continue; // type 1 = monster
				
				mobs.push({
					id: id,
					name: npc.nick || npc.name || "",
					level: npc.lvl || npc.level || 1,
					x: npc.x || npc.posX || 0,
					y: npc.y || npc.posY || 0,
					hp: npc.hp || 0,
					hpMax: npc.maxhp || npc.hpMax || 100,
					alive: !npc.dead && npc.hp > 0,
					attackable: !npc.dead && npc.hp > 0
				});
			}
			
			return mobs;
		} catch(e) {
			console.error("Error getting mobs:", e);
			return [];
		}
	})()
	`
	
	var result []map[string]interface{}
	if err := c.browser.Eval(script, &result); err != nil {
		return nil, fmt.Errorf("failed to get mobs: %w", err)
	}
	
	mobs := make([]*Mob, 0, len(result))
	for _, m := range result {
		mob := &Mob{
			ID:         getString(m, "id"),
			Name:       getString(m, "name"),
			Level:      getInt(m, "level"),
			X:          getFloat(m, "x"),
			Y:          getFloat(m, "y"),
			HP:         getInt(m, "hp"),
			HPMax:      getInt(m, "hpMax"),
			Alive:      getBool(m, "alive"),
			Attackable: getBool(m, "attackable"),
		}
		mobs = append(mobs, mob)
	}
	
	return mobs, nil
}

// MoveTo moves the hero to specific coordinates
func (c *Client) MoveTo(x, y float64) error {
	c.log.WithFields(logrus.Fields{
		"x": x,
		"y": y,
	}).Debug("Moving hero...")
	
	// Try to use game API first
	script := fmt.Sprintf(`
	(function() {
		try {
			if (window.hero && window.hero.moveTo) {
				window.hero.moveTo(%f, %f);
				return true;
			}
			if (window.Engine && window.Engine.moveHero) {
				window.Engine.moveHero(%f, %f);
				return true;
			}
			return false;
		} catch(e) {
			console.error("Error moving:", e);
			return false;
		}
	})()
	`, x, y, x, y)
	
	var success bool
	if err := c.browser.Eval(script, &success); err != nil {
		return fmt.Errorf("failed to move: %w", err)
	}
	
	if !success {
		// Fallback: click on map at coordinates
		return c.browser.ClickAt(x, y)
	}
	
	return nil
}

// AttackMob attacks a specific mob
func (c *Client) AttackMob(mobID string) error {
	c.log.WithField("mobId", mobID).Debug("Attacking mob...")
	
	script := fmt.Sprintf(`
	(function() {
		try {
			let npcList = window.npcs || window.NPC || (window.g && window.g.npcs);
			let target = npcList['%s'];
			
			if (!target) return false;
			
			if (window.hero && window.hero.attack) {
				window.hero.attack(target);
				return true;
			}
			if (window.Engine && window.Engine.attack) {
				window.Engine.attack(target);
				return true;
			}
			
			// Try clicking on mob
			if (target.x && target.y) {
				// This will be handled by click fallback
				return false;
			}
			
			return false;
		} catch(e) {
			console.error("Error attacking:", e);
			return false;
		}
	})()
	`, mobID)
	
	var success bool
	if err := c.browser.Eval(script, &success); err != nil {
		return fmt.Errorf("failed to attack: %w", err)
	}
	
	if !success {
		return fmt.Errorf("could not attack mob")
	}
	
	return nil
}

// UsePotion uses a potion via hotkey
func (c *Client) UsePotion(key string) error {
	c.log.WithField("key", key).Debug("Using potion...")
	return c.browser.PressKey(key)
}

// Respawn clicks the respawn button when dead
func (c *Client) Respawn() error {
	c.log.Info("Attempting to respawn...")
	
	// Try to find and click respawn button
	script := `
	(function() {
		try {
			// Look for common respawn button selectors
			let btn = document.querySelector('.respawn-button') ||
			          document.querySelector('#respawn') ||
			          document.querySelector('[data-action="respawn"]');
			          
			if (btn) {
				btn.click();
				return true;
			}
			
			// Try game API
			if (window.hero && window.hero.respawn) {
				window.hero.respawn();
				return true;
			}
			
			return false;
		} catch(e) {
			console.error("Error respawning:", e);
			return false;
		}
	})()
	`
	
	var success bool
	if err := c.browser.Eval(script, &success); err != nil {
		return fmt.Errorf("failed to respawn: %w", err)
	}
	
	if !success {
		return fmt.Errorf("could not find respawn button")
	}
	
	time.Sleep(2 * time.Second) // Wait for respawn
	return nil
}

// IsConnected checks if the game is still connected
func (c *Client) IsConnected() (bool, error) {
	script := `
	(function() {
		try {
			// Check if socket/connection is alive
			if (window.connection && window.connection.connected !== undefined) {
				return window.connection.connected;
			}
			if (window.socket && window.socket.connected !== undefined) {
				return window.socket.connected;
			}
			if (window.ws && window.ws.readyState !== undefined) {
				return window.ws.readyState === 1; // WebSocket.OPEN
			}
			// Assume connected if hero exists
			return window.hero !== undefined && window.hero !== null;
		} catch(e) {
			return false;
		}
	})()
	`
	
	var connected bool
	if err := c.browser.Eval(script, &connected); err != nil {
		return false, err
	}
	
	return connected, nil
}

// Helper functions to safely extract values from map
func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		case int64:
			return float64(val)
		}
	}
	return 0
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		case int64:
			return int(val)
		}
	}
	return 0
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// Dump for debugging
func (c *Client) DumpGameState() (string, error) {
	hero, err := c.GetHeroState()
	if err != nil {
		return "", err
	}
	
	mobs, err := c.GetMobs()
	if err != nil {
		return "", err
	}
	
	data := map[string]interface{}{
		"hero": hero,
		"mobs": mobs,
	}
	
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(b), nil
}
