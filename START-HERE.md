# START HERE 👋

## Super Quick Start (3 steps!)

### 1. Setup Config
```bash
cp configs/auto-detect.yaml configs/my-bot.yaml
nano configs/my-bot.yaml
```

Change ONLY these two lines:
```yaml
username: "YOUR_USERNAME"  # ← Your Margonem username
password: "YOUR_PASSWORD"  # ← Your Margonem password
```

Save and close.

### 2. Login & Walk to Hunting Spot

- Open Margonem in your browser
- Login manually
- Walk your character to where you want to farm mobs

### 3. Run the Bot

```bash
./bin/margonem-bot --config configs/my-bot.yaml
```

**That's it!** The bot will:
- ✅ Auto-detect your location
- ✅ Find all nearby mobs
- ✅ Start attacking them automatically

## Stop the Bot

Press `Ctrl+C` in the terminal.

## First Time?

Run with visible browser to watch it work:

Edit `configs/my-bot.yaml`:
```yaml
runtime:
  headless: false  # ← Change true to false
```

You'll see the browser window and what the bot is doing!

## Need More Help?

- **Quick Guide**: [AUTO-DETECT-GUIDE.md](AUTO-DETECT-GUIDE.md)
- **Full Docs**: [README.md](README.md)
- **What Changed**: [CHANGES.md](CHANGES.md)

## Common Questions

**Q: Do I need to configure waypoints/coordinates?**  
A: No! Auto-detect mode does this for you.

**Q: Do I need to list mob names?**  
A: No! Bot attacks all mobs it finds.

**Q: What if I die?**  
A: Bot respawns you. Walk back to your spot manually, bot continues.

**Q: Can I use potions?**  
A: Press the potion keys manually (1, 2, etc). Auto-potions were removed.

**Q: How do I farm different spots?**  
A: Just walk to a new spot and restart the bot!

## That's All!

No complicated setup. No coordinates. No hassle.

**Just walk somewhere → run bot → profit!** 💰

---

Having issues? Check [AUTO-DETECT-GUIDE.md](AUTO-DETECT-GUIDE.md) for troubleshooting.
