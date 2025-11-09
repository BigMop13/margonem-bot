# Auto-Detect Mode - Super Simple Setup! 🚀

No manual configuration needed! Just walk to your hunting spot and run the bot.

## Quick Start (3 Steps!)

### 1. Edit Configuration

```bash
cp configs/auto-detect.yaml configs/my-bot.yaml
nano configs/my-bot.yaml  # or your editor
```

**Only change these 2 lines:**
```yaml
account:
  username: "YOUR_USERNAME"  # ← Your Margonem username
  password: "YOUR_PASSWORD"  # ← Your Margonem password
```

**That's it!** Everything else is already configured.

### 2. Login to Margonem

Open your browser and login to Margonem normally. Walk your character to where you want to hunt.

### 3. Run the Bot

```bash
./bin/margonem-bot --config configs/my-bot.yaml
```

**The bot will:**
- ✅ Detect your current location automatically
- ✅ Scan for all nearby mobs
- ✅ Start hunting everything it finds
- ✅ No waypoints needed!
- ✅ No mob names needed!
- ✅ No coordinates needed!

## What Happens

```
INFO Starting Margonem Bot...
INFO Phase: LOGIN
INFO Phase: WAIT_GAME_READY
INFO Game engine is ready!
INFO AUTO-DETECT MODE: Bot will hunt at current location
INFO Current position detected mapId="meadow" x=523.5 y=412.8
INFO Hunting ground auto-detected, starting combat!
INFO Phase: HUNT
INFO New target acquired mob="Wilk" level=5
INFO Attacking target="Wilk"
```

## How It Works

1. **You login manually** (or bot uses saved session)
2. **You walk to hunting spot** with your character
3. **Bot detects** your position and map
4. **Bot hunts** within ~200 pixels radius
5. **Bot attacks** any mob it finds nearby

## What If I Die?

In auto-detect mode:
1. Bot respawns you automatically
2. **You walk back** to hunting ground manually
3. Bot detects new location and continues
4. OR just press `Ctrl+C` to stop

## Configuration Options

You can adjust these in your config:

```yaml
combat:
  hpThreshold: 20          # Retreat when HP below 20%
  maxEngageDistance: 300   # How far to chase mobs
  minLevel: 0              # Attack mobs level 0+
  maxLevel: 999            # Attack mobs up to level 999

behavior:
  minDelayMs: 800          # Fast: 500, Normal: 800, Slow: 1200
  maxDelayMs: 1800         # Fast: 1000, Normal: 1800, Slow: 3000
  pathJitter: 8            # Random movement offset
```

## Tips

### Start in Visible Mode

First time, run with visible browser to watch:
```yaml
runtime:
  headless: false   # Show browser
  debug: true       # Verbose logs
```

### Then Go Headless

Once it works, hide the browser:
```yaml
runtime:
  headless: true    # Hide browser
  debug: false      # Less logging
```

### Multiple Hunting Spots

Create different configs for different spots:
```bash
cp configs/my-bot.yaml configs/wolves.yaml
cp configs/my-bot.yaml configs/bears.yaml
```

Run specific one:
```bash
./bin/margonem-bot --config configs/wolves.yaml
```

### Save Login Session

The bot saves your login with `userDataDir`:
```yaml
runtime:
  userDataDir: "./.chrome-profile"
```

You'll only need to login once!

## Comparison: Auto-Detect vs Manual Mode

| Feature | Auto-Detect | Manual Mode |
|---------|-------------|-------------|
| Setup Time | **30 seconds** | 10+ minutes |
| Waypoints | **Not needed** | Must configure |
| Coordinates | **Auto-detected** | Must find manually |
| Mob Names | **Auto-finds all** | Must list names |
| Death Recovery | Manual walk back | Auto-returns |
| Best For | **Quick farming** | AFK overnight |

## Troubleshooting

### "Bot doesn't attack"
- Make sure you're in an area with mobs
- Check `maxEngageDistance` (increase to 500)
- Verify `minLevel`/`maxLevel` includes your mobs

### "Bot attacks too strong mobs"
- Set `maxLevel: 10` to limit mob level
- Increase `hpThreshold: 30` to retreat sooner

### "Character moves weird"
- Reduce `pathJitter: 3` for less randomness
- Increase delays: `minDelayMs: 1200`

### "Need to login every time"
- Set `userDataDir: "./.chrome-profile"`
- Run once with visible browser to login
- Future runs will reuse session

## Advanced: Mix with Manual Config

You can use auto-detect for location, but still configure combat:

```yaml
runtime:
  autoDetectMode: true      # Auto-detect location

combat:
  targetPriority: ["Boss"]  # But only attack bosses
  minLevel: 50              # Level 50+
  maxLevel: 60              # Up to level 60
```

## Stop the Bot

Press `Ctrl+C` in the terminal.

The bot will stop gracefully.

## Example Session

```bash
# 1. Edit config
nano configs/my-bot.yaml

# 2. Login to game manually, walk to hunting spot

# 3. Run bot
./bin/margonem-bot --config configs/my-bot.yaml

# Watch it work!
# Press Ctrl+C to stop
```

---

**That's it!** No complicated setup, no coordinates, no hassle. Just walk somewhere and run the bot! 🎮

Need help? Check the main [README.md](README.md) for more details.
