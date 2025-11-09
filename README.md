# Margonem Bot

An automated bot for the browser MMORPG Margonem, written in Go.

## ⚠️ Disclaimer

**IMPORTANT:** This bot is provided for educational and research purposes only. Use of automated bots may violate the game's Terms of Service and could result in account suspension or ban. Only use this software:
- On accounts where you have explicit permission
- In compliance with applicable laws and regulations
- In controlled testing environments

The authors are not responsible for any consequences resulting from the use of this software.

## Features

- **🎯 Auto-Detect Mode**: Walk to any location and bot automatically detects your position and hunts there (NEW!)
- **Automated Combat**: Intelligently selects and attacks mobs based on configurable priorities
- **Smart Navigation**: Waypoint-based navigation to hunting grounds with pathfinding
- **Death Recovery**: Automatic respawn (manual walk-back in auto-detect mode)
- **Human-like Behavior**: 
  - Random delays between actions (1-2 seconds configurable)
  - Non-linear movement paths with bezier curves
  - Random position jitter to avoid predictable patterns
  - Periodic idle breaks
- **Stuck Detection**: Detects when character is stuck and attempts recovery
- **Disconnect Handling**: Automatic reconnection with exponential backoff
- **Configurable**: YAML-based configuration for easy customization
- **Robust Error Handling**: Comprehensive logging and screenshot capture on errors

## Architecture

The bot is built with a modular architecture:

```
margonem-bot/
├── cmd/bot/          # Main entry point with state machine
├── internal/
│   ├── browser/      # Chrome browser automation (chromedp)
│   ├── game/         # Game client with JavaScript bridge
│   ├── combat/       # Combat engine and target selection
│   ├── navigation/   # Waypoint navigation and pathfinding
│   ├── behavior/     # Randomization and human-like patterns
│   └── config/       # Configuration management
└── configs/          # YAML configuration files
```

## Prerequisites

- **Go 1.21+**: [Install Go](https://go.dev/doc/install)
- **Google Chrome or Chromium**: The bot uses Chrome for browser automation
- **macOS, Linux, or Windows**: Cross-platform support

## Installation

1. Clone the repository:
```bash
git clone https://github.com/kamilkurek/margonem-bot.git
cd margonem-bot
```

2. Install dependencies:
```bash
go mod download
```

3. Build:
```bash
make build
```

4. **EASY MODE**: Use auto-detect configuration:
```bash
cp configs/auto-detect.yaml configs/my-bot.yaml
# Edit configs/my-bot.yaml - only change username/password!
```

**OR Manual Mode**: Create full configuration:
```bash
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with all settings
```

## Configuration

The bot is configured via YAML files. See `configs/config.example.yaml` for a complete example with comments.

### Key Configuration Sections

#### Account
```yaml
account:
  username: "your_username"
  password: "your_password"
  server: "berufs"
  startUrl: "https://www.margonem.pl/"
```

#### Hunting Profile
```yaml
profile:
  name: "meadow-farm"
  huntingGround:
    mapId: "meadow"
    centerX: 500
    centerY: 500
    radius: 200
  waypoints:
    - mapId: "town"
      x: 600
      y: 520
      description: "Exit gate"
      action: "portal"
```

#### Combat
```yaml
combat:
  hpThreshold: 30          # Retreat below this HP%
  targetPriority: ["Wolf", "Boar", "Fox"]
  maxEngageDistance: 250
  minLevel: 1
  maxLevel: 50
```

#### Behavior (Anti-Detection)
```yaml
behavior:
  minDelayMs: 1000         # Min delay between actions
  maxDelayMs: 2000         # Max delay between actions
  pathJitter: 6            # Random position offset (pixels)
  idleBreakEvery: 30       # Take break every N actions
  idleBreakDuration: 6     # Break duration (seconds)
```

## Usage

### 🚀 Quick Start (Auto-Detect Mode)

**The easiest way to use the bot:**

1. Edit `configs/auto-detect.yaml` - change username/password only
2. Login to Margonem and walk to your hunting spot
3. Run: `./bin/margonem-bot --config configs/auto-detect.yaml`
4. Bot auto-detects location and starts hunting!

**See [AUTO-DETECT-GUIDE.md](AUTO-DETECT-GUIDE.md) for detailed instructions.**

### Build the Bot

```bash
make build
# or: go build -o bin/margonem-bot ./cmd/bot
```

### Run the Bot

```bash
# Run with default config
./bin/margonem-bot

# Run with custom config
./bin/margonem-bot --config /path/to/config.yaml

# Run in debug mode (more verbose logging)
# Edit config.yaml and set runtime.debug: true

# Show version
./bin/margonem-bot --version
```

### Command-Line Flags

- `--config <path>`: Path to configuration file (default: `configs/config.yaml`)
- `--version`: Show version and exit

## How It Works

### State Machine

The bot operates as a state machine with the following phases:

1. **LOGIN**: Navigates to game and logs in
2. **WAIT_GAME_READY**: Waits for game engine to load
3. **NAVIGATE**: Travels to configured hunting ground
4. **HUNT**: Main combat loop
   - Polls game state every second
   - Selects targets based on priority
   - Engages mobs with random delays
   - Uses potions when HP/MP low
   - Patrols when no mobs found
5. **DEAD**: Detected death
6. **RECOVER**: Respawns and returns to hunting ground
7. **DISCONNECTED**: Handles disconnection and reconnects

### Browser Automation

The bot uses `chromedp` to control a Chrome browser instance. It:
- Injects JavaScript to interact with the game engine
- Reads game state (hero position, HP, mobs, etc.)
- Simulates mouse clicks and keyboard input
- Can run headless or with visible browser for debugging

### JavaScript Bridge

The bot probes multiple possible JavaScript namespaces to find game objects:
- `window.hero`, `window.Hero`, `window.g.hero`
- `window.npcs`, `window.NPC`, `window.g.npcs`
- `window.Engine`, `window.map`

This defensive approach makes the bot more resilient to game updates.

### Human-like Behavior

To avoid detection, the bot implements:
- **Random Delays**: 1-2 seconds (configurable) between actions
- **Bezier Curves**: Movement follows curved paths instead of straight lines
- **Position Jitter**: Random offsets when clicking coordinates
- **Idle Breaks**: Periodic pauses to simulate human breaks
- **Variable Patterns**: Randomized patrol routes within hunting grounds

## Troubleshooting

### Bot doesn't login

- Check that `account.username` and `account.password` are correct
- Run with `runtime.headless: false` to see the browser
- Margonem may have changed their login page - you may need to update the login selectors in `cmd/bot/main.go`

### Game state not detected

- The JavaScript bridge may need updating if Margonem changed their client
- Check browser console (F12) for JavaScript errors
- Run with `runtime.debug: true` for detailed logging

### Character appears stuck

- The bot has built-in stuck detection
- Check that waypoints are correct for your map
- Ensure `pathJitter` is not too large

### Screenshots

Screenshots are automatically saved to `runtime.screenshotDir` (default: `./screenshots`) when errors occur.

## Development

### Project Structure

- `cmd/bot/main.go`: Entry point and main state machine
- `internal/browser/`: Browser automation wrapper
- `internal/game/client.go`: JavaScript bridge to game
- `internal/game/state.go`: Game state manager with thread safety
- `internal/combat/`: Target selection and combat logic
- `internal/navigation/`: Waypoint-based navigation
- `internal/behavior/`: Randomization and delays
- `internal/config/`: Configuration loading and validation

### Adding New Features

1. **New Combat Strategy**: Modify `internal/combat/targeting.go`
2. **New Movement Pattern**: Update `internal/navigation/waypoints.go`
3. **New Game State**: Add fields to `internal/game/state.go`
4. **New Behavior**: Add functions to `internal/behavior/`

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/combat
```

## Known Limitations

- Login flow is simplified and may need customization for different Margonem login pages
- Portal/door interactions may need manual selector configuration
- No support for trading, crafting, or other non-combat activities
- Pathfinding is basic - no collision detection with walls/obstacles

## Future Enhancements

- [ ] More sophisticated pathfinding with A* algorithm and collision maps
- [ ] Support for multiple hunting profiles with auto-switching
- [ ] Loot collection and inventory management
- [ ] Party/group support
- [ ] Webhook notifications for deaths, disconnects, level-ups
- [ ] Web dashboard for monitoring bot status
- [ ] Machine learning for improved target selection

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## License

This project is provided as-is for educational purposes. Use at your own risk.

## Credits

Built with:
- [chromedp](https://github.com/chromedp/chromedp) - Chrome browser automation
- [logrus](https://github.com/sirupsen/logrus) - Structured logging
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML parsing

---

**Remember**: Always respect the game's Terms of Service and use automation responsibly!
