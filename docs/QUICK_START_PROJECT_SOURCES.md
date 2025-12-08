# Quick Start: Adding Project Sources

**The easiest way to add custom wisdom sources to your project**

---

## 🚀 30-Second Setup

### Step 1: Create `.wisdom/sources.json` in your project root

```bash
mkdir -p .wisdom
```

### Step 2: Add your source

Create `.wisdom/sources.json`:

```json
{
  "version": "1.0",
  "sources": {
    "my_project": {
      "id": "my_project",
      "name": "My Project Wisdom",
      "icon": "🚀",
      "quotes": {
        "chaos": [
          {
            "quote": "When everything breaks, remember: you built this.",
            "source": "Project Wisdom",
            "encouragement": "You can fix it."
          }
        ],
        "treasury": [
          {
            "quote": "Everything is working perfectly!",
            "source": "Project Wisdom",
            "encouragement": "Enjoy the moment."
          }
        ]
      }
    }
  }
}
```

### Step 3: Done! 🎉

The system automatically detects and loads sources from `.wisdom/sources.json` in your project root.

---

## 📍 Where to Put It

**Recommended**: `.wisdom/sources.json` in project root

```
your-project/
├── .git/
├── go.mod
├── .wisdom/              ← Create this
│   └── sources.json      ← Your sources here
└── ...
```

**Why `.wisdom/`?**
- ✅ Automatically detected as project root marker
- ✅ High priority (overrides global sources)
- ✅ Easy to version control or ignore
- ✅ Clean project structure

---

## 🎯 Using Your Source

```go
engine := wisdom.NewEngine()
engine.Initialize()

// Use your custom source
quote, _ := engine.GetWisdom(75.0, "my_project")
fmt.Println(quote.Quote)
```

---

## 📝 Full Example

See `examples/sources.json` for a complete example with all aeon levels.

---

## 🔄 Hot Reload

After editing `.wisdom/sources.json`:

```go
engine.ReloadSources() // Pick up changes without restart
```

---

## 💡 Tips

- **Version Control**: Add `.wisdom/sources.json` to git for team sharing
- **Private Sources**: Add to `.gitignore` for personal notes
- **Multiple Files**: Split sources across multiple JSON files
- **Override Global**: Your project sources override global sources with same ID

---

## 📚 More Info

- **Full Guide**: `docs/ADDING_PROJECT_SOURCES.md`
- **Configuration**: `docs/CONFIGURABLE_SOURCES.md`
- **Example Code**: `examples/add_project_source.go`

