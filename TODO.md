# devwisdom-go TODO

## Phase 1: Core Structure ✅
- [x] Project structure
- [x] Basic types (Quote, Source, Consultation)
- [x] Engine skeleton
- [x] Config management
- [x] Advisor registry skeleton
- [x] MCP server structure

## Phase 2: Wisdom Data Porting
- [x] Port all 21+ wisdom sources from Python
  - [x] pistis_sophia
  - [x] stoic
  - [x] tao
  - [x] art_of_war
  - [x] bible
  - [x] confucius
  - [x] bofh
  - [x] tao_of_programming
  - [x] murphy
  - [x] shakespeare
  - [x] kybalion
  - [x] gracian
  - [x] enochian
  - [x] Hebrew sources (rebbe, tzaddik, chacham) - Local sources complete
  - [ ] Hebrew sources (pirkei_avot, proverbs, ecclesiastes, psalms) - Phase 7 (Sefaria API)
  - [x] random source selector

## Phase 3: Advisor System ✅
- [x] Complete metric advisor mappings
- [x] Complete tool advisor mappings
- [x] Complete stage advisor mappings
- [x] Score-based consultation frequency
- [x] Mode-aware advisor selection (AGENT/ASK/MANUAL)

## Phase 4: MCP Protocol Implementation ✅
- [x] Implement JSON-RPC 2.0 handler
- [x] Register tools:
  - [x] consult_advisor
  - [x] get_wisdom
  - [x] get_daily_briefing
  - [x] get_consultation_log (stub - Phase 5)
- [x] Register resources:
  - [x] wisdom://sources
  - [x] wisdom://advisors
  - [x] wisdom://advisor/{id}
  - [x] wisdom://consultations/{days} (stub - Phase 5)
- [x] Handle stdio transport
- [x] Error handling and logging

## Phase 5: Consultation Logging ✅
- [x] JSONL log file format
- [x] Consultation tracking
- [x] Log retrieval and filtering
- [ ] Date-based log rotation (optional - Phase 5.4, T-7)

## Phase 6: Daily Random Source Selection
- [x] Date-seeded random selection
- [x] Consistent daily source
- [x] Random source rotation

## Phase 7: Optional Features
- [ ] Sefaria API integration (Hebrew texts)

## Phase 8: Testing
- [ ] Unit tests for wisdom engine
- [ ] Unit tests for advisors
- [ ] Integration tests for MCP server
- [ ] Test with Cursor MCP client

## Phase 9: Documentation
- [ ] API documentation
- [ ] Usage examples
- [ ] Migration guide from Python version
- [ ] Performance benchmarks

## Phase 10: Polish
- [ ] Error messages
- [ ] Logging improvements
- [ ] Performance optimization
- [ ] Cross-compilation (Windows, Linux, macOS)

## Future Goals

### Cursor Extension (Future Enhancement)
A Cursor/VS Code extension is documented in `docs/CURSOR_EXTENSION.md` as a potential future enhancement. This would provide visual UI, status bar integration, and command palette access. **This is not currently being developed** - the MCP server works standalone and provides all core functionality through Cursor's built-in MCP integration.
