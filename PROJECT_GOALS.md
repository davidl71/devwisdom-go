# devwisdom-go Project Goals

**Project**: devwisdom-go - Wisdom Module Extraction (Go Proof of Concept)  
**Type**: MCP Server / Developer Tools  
**Language**: Go  
**Status**: In Progress  
**Updated**: 2025-01-26

---

## Vision

Extract the wisdom module from exarp into a standalone, compiled Go MCP server. This serves as a proof of concept for using compiled languages (Go) for exarp modules, demonstrating performance, deployment simplicity, and language diversity.

---

## Strategic Phases

### Phase 1: Core Structure ✅ **COMPLETE**
**Timeline**: Foundation  
**Priority**: Critical

**Goals**:
- ✅ Go project structure
- ✅ Basic types and interfaces
- ✅ Engine skeleton
- ✅ Config management
- ✅ Advisor registry skeleton
- ✅ MCP server structure

**Keywords**: go, project-structure, types, engine, config, mcp, foundation

---

### Phase 2: Wisdom Data Porting
**Timeline**: Data Migration  
**Priority**: High

**Goals**:
- Port 21+ wisdom sources from Python to Go
- Implement quote storage and retrieval
- Maintain Python API compatibility
- Preserve all wisdom quotes and metadata

**Sources to Port**:
- pistis_sophia, stoic, tao, art_of_war, bible, confucius
- bofh, tao_of_programming, murphy
- shakespeare, kybalion, gracian, enochian
- Hebrew sources: rebbe, tzaddik, chacham, pirkei_avot, proverbs, ecclesiastes, psalms
- random source selector

**Keywords**: porting, data-migration, wisdom-sources, quotes, python-to-go

---

### Phase 3: Advisor System
**Timeline**: Logic Implementation  
**Priority**: High

**Goals**:
- Complete metric → advisor mappings
- Complete tool → advisor mappings
- Complete stage → advisor mappings
- Implement score-based consultation frequency
- Mode-aware advisor selection (AGENT/ASK/MANUAL)

**Keywords**: advisors, mappings, consultation, mode-aware, logic

---

### Phase 4: MCP Protocol Implementation
**Timeline**: MCP Integration  
**Priority**: Critical

**Goals**:
- Implement JSON-RPC 2.0 handler
- Register 5 tools: consult_advisor, get_wisdom, get_daily_briefing, get_consultation_log, export_for_podcast
- Register 4 resources: wisdom://sources, wisdom://advisors, wisdom://advisor/{id}, wisdom://consultations/{days}
- Handle stdio transport
- Error handling and logging

**Keywords**: mcp, json-rpc, tools, resources, stdio, transport, protocol

---

### Phase 5: Consultation Logging
**Timeline**: Persistence  
**Priority**: Medium

**Goals**:
- JSONL log file format
- Consultation tracking
- Log retrieval and filtering
- Date-based log rotation

**Keywords**: logging, persistence, jsonl, consultation-tracking

---

### Phase 6: Daily Random Source Selection
**Timeline**: Enhancement  
**Priority**: Medium

**Goals**:
- Date-seeded random selection
- Consistent daily source (same date = same source)
- Random source rotation

**Keywords**: random, daily-selection, date-seeding, consistency

---

### Phase 7: Optional Features
**Timeline**: Enhancement  
**Priority**: Low

**Goals**:
- Sefaria API integration (Hebrew texts)
- Voice/TTS support (Go equivalent of edge-tts)
- Podcast export formatting

**Keywords**: optional, sefaria, voice, tts, podcast, hebrew

---

### Phase 8: Testing
**Timeline**: Quality  
**Priority**: High

**Goals**:
- Unit tests for wisdom engine
- Unit tests for advisors
- Integration tests for MCP server
- Test with Cursor MCP client

**Keywords**: testing, unit-tests, integration-tests, cursor, quality

---

### Phase 9: Documentation
**Timeline**: Documentation  
**Priority**: Medium

**Goals**:
- API documentation (godoc)
- Usage examples
- Migration guide from Python version
- Performance benchmarks (Go vs Python)

**Keywords**: documentation, godoc, examples, migration-guide, benchmarks

---

### Phase 10: Polish & Deployment
**Timeline**: Finalization  
**Priority**: Medium

**Goals**:
- Error message improvements
- Logging improvements
- Performance optimization
- Cross-compilation (Windows, Linux, macOS)
- CI/CD setup

**Keywords**: polish, deployment, cross-compilation, ci-cd, optimization

---

## Success Criteria

### Minimum Viable Product (MVP)
- ✅ Phase 1: Core structure complete
- [ ] Phase 2: At least 10 wisdom sources ported
- [ ] Phase 4: MCP server functional with stdio transport
- [ ] Phase 8: Basic tests passing

### Full Release
- [ ] All 21+ wisdom sources ported
- [ ] Complete advisor system
- [ ] All MCP tools and resources implemented
- [ ] Comprehensive test coverage
- [ ] Documentation complete
- [ ] Performance benchmarks show improvement over Python

---

## Technical Goals

### Performance
- **Startup time**: < 50ms (vs Python ~200ms)
- **Response time**: < 10ms per tool call
- **Binary size**: < 10MB (single file)

### Compatibility
- **API compatibility**: Match Python version API
- **MCP compatibility**: Full MCP specification compliance
- **Platform support**: Linux, macOS, Windows

### Code Quality
- **Test coverage**: > 80%
- **Linting**: golangci-lint clean
- **Documentation**: All public APIs documented

---

## Dependencies

### Required
- Go 1.21+
- Standard library only (core features)

### Optional
- MCP framework (Foxy Contexts or official SDK) - for Phase 4
- HTTP client - for Sefaria API (Phase 7)

---

## Risks & Mitigations

### Risk: MCP Framework Maturity
- **Risk**: Foxy Contexts may not be production-ready
- **Mitigation**: Fallback to official MCP SDK or direct JSON-RPC implementation

### Risk: Data Porting Complexity
- **Risk**: 21+ sources with complex quote structures
- **Mitigation**: Automated conversion script from Python to Go

### Risk: Performance Not Better
- **Risk**: Go version not significantly faster
- **Mitigation**: Benchmark early, optimize hot paths

---

## References

- **Source**: Python wisdom module in `project_management_automation/tools/wisdom/`
- **MCP Spec**: https://modelcontextprotocol.io/
- **Go Best Practices**: https://go.dev/doc/effective_go

---

**Last Updated**: 2025-01-26
