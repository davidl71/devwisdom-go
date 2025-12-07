# devwisdom-go

**Wisdom Module Extraction (Go Proof of Concept)**

A standalone Go MCP server providing wisdom quotes, trusted advisors, and inspirational guidance for developers. Extracted from the exarp project as a proof of concept for using compiled languages (Go) for exarp modules.

## 🎯 Project Status

**Phase 1**: ✅ Complete (Core Structure)  
**Current Phase**: Phase 2 (Wisdom Data Porting)  
**Language**: Go 1.21+  
**Type**: MCP Server / Developer Tools

## 📋 Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd devwisdom-go

# Build
make build

# Run
make run

# Test
make test
```

## 🏗️ Project Structure

```
devwisdom-go/
├── cmd/server/          # MCP server entry point
├── internal/
│   ├── wisdom/         # Wisdom engine (quotes, sources, advisors)
│   ├── mcp/            # MCP protocol handler
│   └── config/         # Configuration management
├── docs/               # Documentation
├── Makefile           # Build commands
└── go.mod             # Go dependencies
```

## 📊 Planning & Status

**Todo2 Tasks**: 37 tasks across 9 phases (tracked in agentic-tools MCP)  
**Project ID**: `039bb05a-6f78-492b-88b5-28fdfa3ebce7`

See `PROJECT_GOALS.md` for detailed phase breakdown and `PRD.md` for full requirements.

## 🚀 Phases

1. ✅ **Phase 1**: Core Structure (Complete)
2. 🔄 **Phase 2**: Wisdom Data Porting (21+ sources)
3. ⏳ **Phase 3**: Advisor System
4. ⏳ **Phase 4**: MCP Protocol Implementation
5. ⏳ **Phase 5**: Consultation Logging
6. ⏳ **Phase 6**: Daily Random Source Selection
7. ⏳ **Phase 7**: Optional Features (Sefaria, TTS)
8. ⏳ **Phase 8**: Testing
9. ⏳ **Phase 9**: Documentation
10. ⏳ **Phase 10**: Polish & Deployment

## 📚 Documentation

- **PROJECT_GOALS.md** - Strategic phases and goals
- **PRD.md** - Product Requirements Document (129 user stories)
- **TODO.md** - Task breakdown by phase
- **EXARP_PLANNING_COMPLETE.md** - Planning analysis summary

## 🔗 Related

- **Source**: Python wisdom module in `exarp` project
- **MCP Spec**: https://modelcontextprotocol.io/
- **Go Docs**: https://go.dev/doc/effective_go

## 📝 License

[Add your license here]

## 👤 Author

Extracted from exarp project as compiled language PoC.
