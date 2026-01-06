# CLI Usage Examples

This guide provides practical examples for using the `devwisdom` command-line interface.

## Basic Commands

### Get a Random Wisdom Quote

```bash
devwisdom quote
```

**Output:**
```
📜 "The path of wisdom is found in daily practice."
   — Pistis Sophia
   
💡 Encouragement: Every step forward is progress.
```

### Get Quote from Specific Source

```bash
devwisdom quote --source stoic
```

**Output:**
```
🏛️ "The impediment to action advances action. What stands in the way becomes the way."
   — Stoic Philosophy
   
💡 Encouragement: Obstacles are opportunities in disguise.
```

### Get Quote with Score Context

The score determines which aeon level quotes are selected:

```bash
# Low score (chaos level)
devwisdom quote --source art_of_war --score 25

# Medium score (middle_aeons level)
devwisdom quote --source art_of_war --score 60

# High score (treasury level)
devwisdom quote --source art_of_war --score 90
```

## Advisor Consultations

### Consult Advisor for a Metric

```bash
devwisdom consult --metric security --score 40
```

**Output:**
```
🔒 Security Advisor Consultation

📊 Score: 40.0 (Lower Aeons)
🎯 Mode: Building
📅 Frequency: milestones

📜 Quote: "Security is not a destination, but a journey."
   — Pistis Sophia

💡 Rationale: Your security score indicates foundational work is needed.
   Focus on establishing core security practices and regular reviews.

🌱 Encouragement: Every security improvement builds a stronger foundation.
```

### Consult Advisor for a Tool

```bash
devwisdom consult --tool project_scorecard --score 75
```

**Output:**
```
📊 Project Scorecard Advisor Consultation

📊 Score: 75.0 (Upper Aeons)
🎯 Mode: Maturing
📅 Frequency: milestones

📜 Quote: "Excellence is not a skill, it's an attitude."
   — Stoic Philosophy

💡 Rationale: Your project shows strong maturity. Continue refining
   processes and maintaining quality standards.

🌱 Encouragement: You're building something meaningful. Keep going!
```

### Consult Advisor for a Workflow Stage

```bash
devwisdom consult --stage daily_checkin --score 65
```

## Listing Resources

### List All Wisdom Sources

```bash
devwisdom sources
```

**Output:**
```
Available Wisdom Sources:

📜 pistis_sophia - Pistis Sophia
🏛️ stoic - Stoic Philosophy
☯️ tao - Tao Te Ching
⚔️ art_of_war - The Art of War
📖 bible - Biblical Wisdom
... (21+ sources)
```

### List All Advisors

```bash
devwisdom advisors
```

**Output:**
```
Available Advisors:

Metric Advisors (14):
  - security
  - testing
  - documentation
  ...

Tool Advisors (12):
  - project_scorecard
  - check_documentation_health
  ...

Stage Advisors (10):
  - daily_checkin
  - sprint_planning
  ...
```

## Daily Briefing

### Get Today's Briefing

```bash
devwisdom briefing
```

**Output:**
```
📅 Daily Wisdom Briefing - 2026-01-06

🎯 Project Health: 75.0 (Upper Aeons)

📜 Today's Wisdom Source: stoic

💬 Quote: "The impediment to action advances action."
   — Stoic Philosophy

🌱 Encouragement: What stands in the way becomes the way.

📊 Available Sources: 21
👥 Available Advisors: 36
```

### Get Briefing for Last N Days

```bash
devwisdom briefing --days 7
```

## JSON Output Format

All commands support `--json` flag for machine-readable output:

```bash
devwisdom quote --json
```

**Output:**
```json
{
  "quote": "The path of wisdom is found in daily practice.",
  "source": "pistis_sophia",
  "encouragement": "Every step forward is progress.",
  "wisdom_source": "Pistis Sophia",
  "wisdom_icon": "📜"
}
```

```bash
devwisdom consult --metric security --score 40 --json
```

**Output:**
```json
{
  "timestamp": "2026-01-06T20:00:00Z",
  "consultation_type": "advisor",
  "advisor": "pistis_sophia",
  "advisor_icon": "📜",
  "advisor_name": "Pistis Sophia",
  "rationale": "Your security score indicates foundational work is needed.",
  "metric": "security",
  "score_at_time": 40.0,
  "consultation_mode": "building",
  "mode_icon": "🌱",
  "mode_frequency": "milestones",
  "quote": "Security is not a destination, but a journey.",
  "quote_source": "pistis_sophia",
  "encouragement": "Every security improvement builds a stronger foundation."
}
```

## Error Handling

### Invalid Source

```bash
devwisdom quote --source invalid_source
```

**Output:**
```
Error: unknown source: invalid_source
Available sources: pistis_sophia, stoic, tao, ...
```

### Invalid Score

```bash
devwisdom quote --source stoic --score invalid
```

**Output:**
```
Error: invalid score: invalid (must be a number)
```

### Missing Required Parameters

```bash
devwisdom consult --metric security
```

**Output:**
```
Error: score is required for metric consultations
Usage: devwisdom consult --metric <metric> --score <score>
```

## Advanced Usage

### Combining Commands

```bash
# Get quote and then consult advisor
devwisdom quote --source stoic --score 50
devwisdom consult --metric testing --score 50
```

### Using in Scripts

```bash
#!/bin/bash
# Get daily wisdom for automation
QUOTE=$(devwisdom quote --json)
echo "Daily Wisdom: $(echo $QUOTE | jq -r '.quote')"
```

### Environment Variables

Some configuration can be set via environment variables:

```bash
export EXARP_WISDOM_SOURCE=stoic
devwisdom quote  # Will use stoic as default source
```

## Tips

1. **Score Context Matters**: Different scores select different aeon level quotes
2. **JSON Output**: Use `--json` for integration with other tools
3. **Source Selection**: Use `random` to get a date-seeded random source
4. **Advisor Selection**: Advisors are automatically selected based on metric/tool/stage
5. **Daily Consistency**: Using `random` source gives the same source for the entire day

## See Also

- [MCP Integration Examples](./mcp_integration.md)
- [Main README](../README.md)
- [Project Goals](../PROJECT_GOALS.md)

