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

## Phase 3: Advisor System
- [ ] Complete metric advisor mappings
- [ ] Complete tool advisor mappings
- [ ] Complete stage advisor mappings
- [ ] Score-based consultation frequency
- [ ] Mode-aware advisor selection (AGENT/ASK/MANUAL)

## Phase 4: MCP Protocol Implementation
- [ ] Implement JSON-RPC 2.0 handler
- [ ] Register tools:
  - [ ] consult_advisor
  - [ ] get_wisdom
  - [ ] get_daily_briefing
  - [ ] get_consultation_log
  - [ ] export_for_podcast
- [ ] Register resources:
  - [ ] wisdom://sources
  - [ ] wisdom://advisors
  - [ ] wisdom://advisor/{id}
  - [ ] wisdom://consultations/{days}
- [ ] Handle stdio transport
- [ ] Error handling and logging

## Phase 5: Consultation Logging
- [ ] JSONL log file format
- [ ] Consultation tracking
- [ ] Log retrieval and filtering
- [ ] Date-based log rotation

## Phase 6: Daily Random Source Selection
- [x] Date-seeded random selection
- [x] Consistent daily source
- [x] Random source rotation

## Phase 7: Optional Features
- [ ] Sefaria API integration (Hebrew texts)
- [ ] Voice/TTS support (edge-tts/pyttsx3 equivalent)
- [ ] Podcast export formatting

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

## Future: Cursor Extension (Very Low Priority)
- [ ] Research Cursor Extension Architecture
- [ ] Phase 1: Extension Foundation (MVP)
- [ ] Phase 2: Command Palette Integration
- [ ] Phase 3: Sidebar Panel Implementation
- [ ] Phase 4: Notifications & Polish

**Note**: Extension tasks are documented in `docs/CURSOR_EXTENSION.md`. These are optional enhancements for better UX. MCP server works standalone.
