# Quick Start Guide

Get your Margonem bot up and running in 5 minutes!

## Step 1: Prerequisites

Ensure you have:
- Go 1.21+ installed
- Chrome or Chromium browser
- A Margonem account (for testing purposes)

## Step 2: Build the Bot

```bash
make build
# Or: go build -o bin/margonem-bot ./cmd/bot
```

## Step 3: Configure

1. Copy the example configuration:
```bash
cp configs/config.example.yaml configs/config.yaml
```

2. Edit `configs/config.yaml` with your settings:
```bash
nano configs/config.yaml  # or use your preferred editor
```

**Minimum required changes:**
- `account.username`: Your Margonem username
- `account.password`: Your Margonem password
- `account.startUrl`: The Margonem game URL
- `profile.huntingGround`: Coordinates for your hunting area
- `profile.waypoints`: Path from spawn to hunting ground

## Step 4: Test Run (Visible Browser)

For first-time setup, run with visible browser to see what's happening:

1. Set in `configs/config.yaml`:
```yaml
runtime:
  headless: false
  debug: true
```

2. Run the bot:
```bash
./bin/margonem-bot
```

3. Watch the browser window - the bot should:
   - Navigate to Margonem
   - Login (you may need to help with CAPTCHA)
   - Wait for game to load
   - Navigate to hunting ground
   - Start hunting mobs

4. Press `Ctrl+C` to stop

## Step 5: Customize for Your Hunting Ground

### Find Your Coordinates

1. In the game, press F12 to open developer console
2. Type: `window.hero.x` and `window.hero.y` to see your position
3. Write down coordinates of:
   - Town spawn point
   - Waypoints along the route
   - Center of hunting ground

### Update Waypoints

```yaml
profile:
  waypoints:
    - mapId: "town-1"
      x: 600
      y: 520
      description: "Town gate"
      action: "walk"
    - mapId: "meadow"
      x: 420
      y: 380
      description: "Hunting ground center"
      action: "walk"
```

### Set Target Mobs

```yaml
combat:
  targetPriority: ["Wolf", "Boar", "Fox"]  # Your preferred mobs
  minLevel: 1    # Minimum mob level
  maxLevel: 10   # Maximum mob level
```

## Step 6: Production Run

Once everything works:

1. Enable headless mode:
```yaml
runtime:
  headless: true
  debug: false
```

2. Run the bot:
```bash
./bin/margonem-bot
```

The bot will now run in the background!

## Troubleshooting

### Bot doesn't login
- Make sure username/password are correct
- Check if Margonem has CAPTCHA
- Try running with `headless: false` to see what's happening

### Character doesn't move
- Verify coordinates are correct
- Check browser console (F12) for JavaScript errors
- Make sure waypoints are on the correct map

### Bot can't find mobs
- Check `combat.targetPriority` mob names match exactly
- Use `minLevel` and `maxLevel` to filter mobs
- Increase `combat.maxEngageDistance`

### Character dies repeatedly
- Lower `combat.minLevel` and `combat.maxLevel`
- Increase `combat.hpThreshold` (retreat sooner)
- Set `potions.hpBelow` higher (use potions earlier)
- Choose easier hunting ground

## Monitoring

Watch the logs for:
- `Phase: HUNT` - Bot is working normally
- `New target acquired` - Found a mob
- `Using HP potion` - Health management
- `Phase: DEAD` - Character died (will auto-respawn)
- `Reconnected successfully` - Recovered from disconnect

## Advanced Tips

### Multiple Profiles

Create multiple config files for different hunting grounds:
```bash
cp configs/config.yaml configs/wolves.yaml
cp configs/config.yaml configs/bears.yaml
```

Run with specific config:
```bash
./bin/margonem-bot --config configs/wolves.yaml
```

### Human-like Behavior

Adjust randomness to be more/less human-like:
```yaml
behavior:
  minDelayMs: 800      # Faster actions
  maxDelayMs: 2500     # More variation
  pathJitter: 10       # More random movement
  idleBreakEvery: 25   # More frequent breaks
```

### Performance Tuning

For lower-level mobs (faster farming):
```yaml
combat:
  maxEngageDistance: 300  # Chase further
  retargetOnDeath: true   # Immediately find next target
behavior:
  minDelayMs: 500         # Faster actions
  maxDelayMs: 1000
```

## Safety Reminders

- ⚠️ Use only on accounts you own
- ⚠️ Respect the game's Terms of Service
- ⚠️ Test on low-level characters first
- ⚠️ Monitor the bot periodically
- ⚠️ Use reasonable behavior settings (don't make it obvious)

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Experiment with different hunting grounds
- Fine-tune behavior settings for your playstyle
- Set up logging and monitoring

Happy hunting! 🏹
