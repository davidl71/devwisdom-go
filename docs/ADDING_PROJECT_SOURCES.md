# Adding Project-Specific Wisdom Sources

**Guide**: How to add custom wisdom sources to your project

---

## Quick Start

### Option 1: Create `.wisdom/sources.json` in Your Project Root

```bash
# In your project root
mkdir -p .wisdom
cat > .wisdom/sources.json << 'EOF'
{
  "version": "1.0",
  "sources": {
    "my_project": {
      "id": "my_project",
      "name": "My Project Wisdom",
      "icon": "🚀",
      "description": "Custom wisdom for my project",
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
EOF
```

### Option 2: Use the Go API

```go
package main

import (
    "github.com/davidl71/devwisdom-go/internal/wisdom"
)

func main() {
    engine := wisdom.NewEngine()
    engine.Initialize()

    // Create a new source
    config := &wisdom.SourceConfig{
        ID:   "my_project",
        Name: "My Project Wisdom",
        Icon: "🚀",
        Quotes: map[string][]wisdom.Quote{
            "chaos": {
                {
                    Quote:        "When everything breaks, remember: you built this.",
                    Source:       "Project Wisdom",
                    Encouragement: "You can fix it.",
                },
            },
            "treasury": {
                {
                    Quote:        "Everything is working perfectly!",
                    Source:       "Project Wisdom",
                    Encouragement: "Enjoy the moment.",
                },
            },
        },
    }

    // Save to project directory
    loader := engine.GetLoader() // If exposed, or use NewSourceLoader()
    if err := loader.SaveProjectSource(config); err != nil {
        log.Fatal(err)
    }

    // Reload to pick up new source
    engine.ReloadSources()
}
```

---

## How It Works

### Project Root Detection

The system automatically detects your project root by looking for:

- `.git` directory
- `.todo2` directory
- `go.mod` file
- `package.json` file
- `CMakeLists.txt` file
- `Makefile` file
- `.wisdom` directory (if it exists, that's your project root)

### Source Loading Priority

Sources are loaded in this order (later sources override earlier ones):

1. **Embedded sources** (compiled into binary)
2. **Explicit paths** (via `WithConfigPaths()`)
3. **Project-specific** (`.wisdom/sources.json` in project root) ⭐ **YOUR PROJECT**
4. **Project root** (`sources.json` in project root)
5. **Current directory** (`sources.json` in CWD)
6. **Home directory** (`~/.wisdom/sources.json`)
7. **XDG config** (`~/.config/wisdom/sources.json`)

**Project-specific sources have high priority** and will override global sources with the same ID.

---

## Project-Specific Location

### Recommended: `.wisdom/sources.json`

Create a `.wisdom` directory in your project root:

```
your-project/
├── .git/
├── go.mod
├── .wisdom/              ← Create this
│   └── sources.json      ← Your project sources
└── ...
```

**Benefits:**
- ✅ Version controlled (add to `.gitignore` if you want private sources)
- ✅ Project-specific (doesn't affect other projects)
- ✅ High priority (overrides global sources)
- ✅ Easy to share with team

### Alternative: `sources.json` in Project Root

You can also place `sources.json` directly in your project root:

```
your-project/
├── .git/
├── go.mod
├── sources.json          ← Your project sources
└── ...
```

---

## Example: Team Project Sources

Create shared wisdom sources for your team:

```json
{
  "version": "1.0",
  "author": "My Team",
  "sources": {
    "team_motto": {
      "id": "team_motto",
      "name": "Team Motto",
      "icon": "👥",
      "description": "Our team's shared wisdom",
      "quotes": {
        "chaos": [
          {
            "quote": "We've been here before. We'll get through this.",
            "source": "Team Wisdom",
            "encouragement": "Remember past victories."
          }
        ],
        "middle_aeons": [
          {
            "quote": "Code review is a gift, not a criticism.",
            "source": "Team Wisdom",
            "encouragement": "We're all learning together."
          }
        ],
        "treasury": [
          {
            "quote": "We ship quality code, on time, together.",
            "source": "Team Wisdom",
            "encouragement": "This is who we are."
          }
        ]
      }
    }
  }
}
```

---

## Example: Personal Project Sources

Add your own personal wisdom:

```json
{
  "version": "1.0",
  "sources": {
    "my_notes": {
      "id": "my_notes",
      "name": "My Development Notes",
      "icon": "📝",
      "quotes": {
        "chaos": [
          {
            "quote": "Remember: The last time this happened, the fix was...",
            "source": "My Notes",
            "encouragement": "Check your notes."
          }
        ]
      }
    }
  }
}
```

---

## Using Your Custom Source

Once you've added your source, use it like any other:

```go
engine := wisdom.NewEngine()
engine.Initialize()

// Use your custom source
quote, err := engine.GetWisdom(75.0, "my_project")
if err != nil {
    log.Fatal(err)
}

fmt.Println(quote.Quote)
```

Or set it as default in your config:

```json
{
  "source": "my_project",
  "disabled": false
}
```

---

## Tips

### 1. Version Control

**Public sources** (safe to share):
```bash
# Add to git
git add .wisdom/sources.json
```

**Private sources** (personal notes):
```bash
# Add to .gitignore
echo ".wisdom/sources.json" >> .gitignore
```

### 2. Team Collaboration

Create a shared source file that everyone can contribute to:

```bash
# Everyone adds their favorite quotes
git pull
# Edit .wisdom/sources.json
git commit -m "Add team wisdom quote"
git push
```

### 3. Multiple Source Files

You can split sources across multiple files:

```go
loader := wisdom.NewSourceLoader().
    WithConfigPaths(
        ".wisdom/team_sources.json",
        ".wisdom/personal_sources.json",
        ".wisdom/project_sources.json",
    )
```

### 4. Hot Reloading

Reload sources without restarting:

```go
// After editing .wisdom/sources.json
engine.ReloadSources()
```

---

## Full Example

Complete example with all aeon levels:

```json
{
  "version": "1.0",
  "sources": {
    "my_project": {
      "id": "my_project",
      "name": "My Project Wisdom",
      "icon": "🚀",
      "description": "Custom wisdom for my project",
      "quotes": {
        "chaos": [
          {
            "quote": "When everything breaks, remember: you built this.",
            "source": "Project Wisdom",
            "encouragement": "You can fix it."
          }
        ],
        "lower_aeons": [
          {
            "quote": "Progress is progress, even if it's slow.",
            "source": "Project Wisdom",
            "encouragement": "Keep moving forward."
          }
        ],
        "middle_aeons": [
          {
            "quote": "We're getting there, one commit at a time.",
            "source": "Project Wisdom",
            "encouragement": "Consistency wins."
          }
        ],
        "upper_aeons": [
          {
            "quote": "The architecture is solid, the tests are passing.",
            "source": "Project Wisdom",
            "encouragement": "You've built something good."
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

---

## Troubleshooting

### Source Not Found

1. **Check file location**: Ensure `.wisdom/sources.json` is in project root
2. **Check JSON syntax**: Validate JSON with `jq` or online validator
3. **Reload sources**: Call `engine.ReloadSources()`
4. **Check priority**: Your source might be overridden by a later file

### Source Overridden

If your source is being overridden:
- Check for duplicate IDs in other config files
- Later files override earlier ones
- Project sources have higher priority than global sources

### Project Root Not Detected

If project root isn't detected:
- Create `.wisdom` directory manually
- Or use `WithProjectRoot()` to set it explicitly:

```go
loader := wisdom.NewSourceLoader().
    WithProjectRoot("/path/to/your/project")
```

---

## See Also

- `CONFIGURABLE_SOURCES.md` - Full configuration documentation
- `examples/sources.json` - Example source configuration
- `PYTHON_SOURCES_ANALYSIS.md` - Reference for quote structure

