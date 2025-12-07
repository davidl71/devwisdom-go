# devwisdom-go

> Developer wisdom and trusted advisors - Compiled in Go for performance

**MCP Server** providing inspirational quotes, trusted advisors, and wisdom guidance for development workflows.

## Features

- 🎯 **21+ Wisdom Sources** - BOFH, Stoic, Tao, Art of War, Hebrew advisors, and more
- 📊 **Trusted Advisor System** - Metric/tool/stage → advisor mapping
- 📝 **Consultation Logging** - Track advisor consultations for analysis
- 🎙️ **Podcast Export** - Export consultations for audio/video generation
- 🌐 **Hebrew Support** - Sefaria API integration (optional)
- 🔊 **Voice/TTS** - Text-to-speech for podcast generation (optional)

## Quick Start

```bash
# Build
go build -o devwisdom ./cmd/server

# Run as MCP server (stdio)
./devwisdom

# Or install globally
go install ./cmd/server
devwisdom
```

## MCP Client Configuration

Add to your MCP client config (Cursor, Claude Desktop):

```json
{
  "mcpServers": {
    "devwisdom": {
      "command": "/path/to/devwisdom-go/devwisdom",
      "args": [],
      "description": "Developer wisdom and trusted advisors"
    }
  }
}
```

## Usage

### Available Tools

- `consult_advisor` - Consult trusted advisor for metric/tool/stage
- `get_wisdom` - Get wisdom quote by source and score
- `get_daily_briefing` - Daily advisor briefing based on scores
- `get_consultation_log` - Retrieve consultation history
- `export_for_podcast` - Export consultations for podcast generation

### Available Resources

- `wisdom://sources` - List all wisdom sources
- `wisdom://advisors` - List all advisors
- `wisdom://advisor/{id}` - Get advisor details
- `wisdom://consultations/{days}` - Get recent consultations

## Why Go?

- ✅ **Fast compilation** - Quick iteration
- ✅ **Excellent JSON** - Built-in `encoding/json`
- ✅ **Single binary** - Easy deployment
- ✅ **Proven MCP framework** - Foxy Contexts
- ✅ **Zero dependencies** - All stdlib for core features

## Wisdom Sources

**Classical**: pistis_sophia, stoic, tao, art_of_war, bible, confucius  
**Tech**: bofh, tao_of_programming, murphy  
**Creative**: shakespeare, kybalion, gracian  
**Hebrew** (עברית): rebbe, tzaddik, chacham, pirkei_avot, proverbs, ecclesiastes, psalms  
**Mystical**: enochian  
**Random**: Daily random source selection

## License

MIT License

## Credits

- Wisdom texts from [sacred-texts.com](https://sacred-texts.com/) (public domain)
- Hebrew texts from [Sefaria.org](https://sefaria.org/) (open API)
- Extracted from [exarp](https://github.com/davidl71/project-management-automation)
