# devwisdom-go TODO

## Phase 1: Core Structure ✅
- [x] Project structure
- [x] Basic types (Quote, Source, Consultation)
- [x] Engine skeleton
- [x] Config management
- [x] Advisor registry skeleton
- [x] MCP server structure

## Phase 2: Wisdom Data Porting
- [ ] Port all 21+ wisdom sources from Python
  - [ ] pistis_sophia
  - [ ] stoic
  - [ ] tao
  - [ ] art_of_war
  - [ ] bible
  - [ ] confucius
  - [ ] bofh
  - [ ] tao_of_programming
  - [ ] murphy
  - [ ] shakespeare
  - [ ] kybalion
  - [ ] gracian
  - [ ] enochian
  - [ ] Hebrew sources (rebbe, tzaddik, chacham, pirkei_avot, proverbs, ecclesiastes, psalms)
  - [ ] random source selector

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
- [ ] Date-seeded random selection
- [ ] Consistent daily source
- [ ] Random source rotation

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
