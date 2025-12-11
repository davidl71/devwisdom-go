# Cursor Extension for devwisdom-go

## Overview

This document outlines the architecture and implementation plan for a Cursor/VS Code extension to complement the devwisdom-go MCP server. The extension would provide visual UI, better discoverability, and enhanced user experience for accessing wisdom quotes and advisor consultations.

## Current State

### MCP Server (Implemented)
- ✅ Core MCP server with 5 tools
- ✅ Wisdom sources (16/21 local sources)
- ✅ Advisor system foundation
- ✅ Resources for tools and consultations
- ✅ Accessible via Cursor's MCP integration

### Limitations
- No visual UI or persistent display
- Requires AI chat interaction to access
- No status bar indicators
- No command palette shortcuts
- No sidebar panels

## Extension Value Proposition

### User Benefits
1. **Visual Feedback**: See project health and wisdom quotes at a glance
2. **Discoverability**: Easy access via command palette and status bar
3. **Persistent Awareness**: Daily briefings visible without chat interaction
4. **Better Integration**: Native Cursor UI components
5. **Notifications**: Proactive advisor recommendations

### Technical Benefits
1. **TypeScript/Node.js**: Leverages existing ecosystem
2. **VS Code API**: Rich UI components available
3. **MCP Integration**: Can call MCP tools directly
4. **Extension Host**: Isolated from main process

## Architecture

### Extension Structure

```
devwisdom-extension/
├── package.json              # Extension manifest
├── tsconfig.json             # TypeScript config
├── src/
│   ├── extension.ts         # Main entry point
│   ├── mcpClient.ts         # MCP server communication
│   ├── statusBar.ts         # Status bar integration
│   ├── commands.ts           # Command palette handlers
│   ├── sidebar/
│   │   ├── briefingView.ts  # Daily briefing panel
│   │   ├── sourcesView.ts   # Sources browser
│   │   └── consultationsView.ts # Consultation history
│   └── notifications.ts     # Notification system
├── media/
│   └── icons/               # Extension icons
└── README.md
```

### MCP Communication

The extension will communicate with the devwisdom-go MCP server using:

1. **Direct MCP Protocol**: Use MCP client library to call tools
2. **Stdio Transport**: Connect to devwisdom binary via stdio
3. **Resource Access**: Read wisdom:// resources for data

### Key Components

#### 1. Status Bar Integration
- Display current project health score
- Show daily wisdom quote (truncated)
- Click to open full quote
- Color-coded health indicator

#### 2. Command Palette
- `DevWisdom: Get Daily Briefing`
- `DevWisdom: Consult Advisor`
- `DevWisdom: Get Wisdom Quote`
- `DevWisdom: View Sources`
- `DevWisdom: View Consultation Log`

#### 3. Sidebar Panel
- **Daily Briefing View**: Shows today's wisdom quotes and guidance
- **Sources Browser**: Browse all available wisdom sources
- **Consultation History**: View past advisor consultations
- **Project Health**: Visual health score dashboard

#### 4. Notifications
- Daily wisdom quote on startup (optional)
- Advisor recommendations for low scores
- Project health alerts
- Source change notifications

## Implementation Phases

### Phase 1: Foundation (MVP)
**Priority**: Very Low  
**Estimated Time**: 8-12 hours

**Goals:**
- Basic extension structure
- MCP client integration
- Status bar with health score
- Single command: "Get Wisdom Quote"

**Deliverables:**
- Extension package.json
- Basic TypeScript setup
- MCP stdio client
- Status bar component
- One working command

### Phase 2: Command Palette
**Priority**: Very Low  
**Estimated Time**: 4-6 hours

**Goals:**
- All command palette commands
- Command handlers for each tool
- Output formatting

**Deliverables:**
- 5 command handlers
- Output views (webview or output channel)
- Error handling

### Phase 3: Sidebar Panel
**Priority**: Very Low  
**Estimated Time**: 8-10 hours

**Goals:**
- Webview-based sidebar
- Daily briefing view
- Sources browser
- Consultation history

**Deliverables:**
- Sidebar panel implementation
- Three view components
- Data refresh mechanism

### Phase 4: Notifications & Polish
**Priority**: Very Low  
**Estimated Time**: 4-6 hours

**Goals:**
- Notification system
- Settings/configuration
- Icon assets
- Documentation

**Deliverables:**
- Notification handlers
- Settings UI
- Extension README
- Marketplace assets

## Technical Decisions

### MCP Client Library
**Options:**
1. **@modelcontextprotocol/sdk** (Official MCP SDK)
   - Pros: Official, well-maintained, TypeScript
   - Cons: May be overkill for simple use case
2. **Custom stdio client**
   - Pros: Lightweight, full control
   - Cons: More implementation work

**Recommendation**: Start with official SDK, fallback to custom if needed

### UI Framework
**Options:**
1. **VS Code Webview API** (Native)
   - Pros: Integrated, no external deps
   - Cons: Limited styling options
2. **React + Webview**
   - Pros: Rich UI, component reuse
   - Cons: Larger bundle, more complexity

**Recommendation**: Start with native Webview API, consider React if UI becomes complex

### Data Refresh
- **Polling**: Check for updates every N minutes
- **Event-driven**: Listen to file changes (sources.json)
- **Manual**: User-triggered refresh

**Recommendation**: Hybrid - polling for health scores, event-driven for file changes

## Dependencies

### Required
- `@types/vscode`: VS Code API types
- `@modelcontextprotocol/sdk`: MCP client (if using SDK)
- `typescript`: TypeScript compiler

### Optional
- `react`: If using React for UI
- `react-dom`: React rendering
- `@vscode/webview-ui-toolkit`: VS Code UI components

## Configuration

### Extension Settings

```json
{
  "devwisdom.showStatusBar": true,
  "devwisdom.statusBarQuote": true,
  "devwisdom.dailyBriefingOnStartup": false,
  "devwisdom.healthScoreThreshold": 70,
  "devwisdom.refreshInterval": 300,
  "devwisdom.mcpServerPath": "./devwisdom"
}
```

## User Stories

### US-EXT-1: Status Bar Health Indicator
*As a* developer,  
*I want* to see my project health score in the status bar,  
*So that* I'm aware of project status at a glance.

### US-EXT-2: Daily Briefing Command
*As a* developer,  
*I want* to run a command to get my daily briefing,  
*So that* I can start my day with wisdom and guidance.

### US-EXT-3: Wisdom Quote Notification
*As a* developer,  
*I want* to receive a daily wisdom quote notification,  
*So that* I'm reminded of project health and guidance.

### US-EXT-4: Sources Browser
*As a* developer,  
*I want* to browse available wisdom sources,  
*So that* I can explore different philosophical perspectives.

### US-EXT-5: Consultation History
*As a* developer,  
*I want* to view my consultation history,  
*So that* I can track advisor recommendations over time.

## Future Enhancements

### Potential Features
1. **Project Health Dashboard**: Visual charts and metrics
2. **Wisdom Quote Widget**: Floating quote display
3. **Advisor Chat Integration**: Direct advisor access in chat
4. **Custom Source Editor**: Add/edit sources from UI
5. **Export/Share**: Export consultations as markdown/PDF
6. **Themes**: Customizable quote display themes
7. **Keyboard Shortcuts**: Quick access to common actions

### Integration Ideas
1. **Git Integration**: Show wisdom on commit/PR
2. **Task Integration**: Link consultations to Todo2 tasks
3. **Code Review**: Advisor suggestions in review comments
4. **Sprint Planning**: Wisdom quotes in planning sessions

## Research Needed

### Before Implementation
1. ✅ VS Code Extension API documentation review
2. ✅ MCP SDK usage patterns
3. ⏳ Cursor-specific extension requirements
4. ⏳ Performance considerations (stdio overhead)
5. ⏳ Extension marketplace requirements

### Technical Questions
1. How does Cursor handle MCP extensions?
2. Can extension call MCP tools directly?
3. What's the best way to refresh data?
4. How to handle MCP server restarts?
5. Extension activation strategies?

## Success Criteria

### MVP Success
- Extension installs and activates
- Status bar shows health score
- One command works end-to-end
- MCP communication functional

### Full Success
- All commands implemented
- Sidebar panel functional
- Notifications working
- User feedback positive
- Extension published to marketplace

## Risks & Mitigations

### Risk: MCP Communication Complexity
**Mitigation**: Start with simple stdio client, add SDK if needed

### Risk: Performance Overhead
**Mitigation**: Lazy loading, efficient polling, caching

### Risk: Maintenance Burden
**Mitigation**: Keep scope minimal, focus on core features

### Risk: User Adoption
**Mitigation**: Gather feedback early, iterate based on usage

## Related Documentation

- [VS Code Extension API](https://code.visualstudio.com/api)
- [MCP Specification](https://modelcontextprotocol.io/)
- [MCP SDK Documentation](https://github.com/modelcontextprotocol/typescript-sdk)
- [Cursor Extension Guidelines](https://cursor.sh/docs) (if available)

## Notes

- This extension is **optional** - MCP server works standalone
- Priority is **very low** - focus on core MCP functionality first
- Can be implemented incrementally
- Consider community contributions if interest exists

---

**Status**: Documented for future implementation  
**Priority**: Very Low  
**Last Updated**: 2025-12-09

