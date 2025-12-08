# Python Wisdom Sources Analysis

**Date**: 2025-01-26  
**Purpose**: Understanding Python wisdom module structure for Go porting

---

## 📁 File Structure

```
project_management_automation/tools/wisdom/
├── __init__.py          # Public API exports
├── sources.py           # Main wisdom sources (17 local sources)
├── pistis_sophia.py     # Separate Pistis Sophia module
├── advisors.py           # Advisor system (metric/tool/stage mappings)
├── sefaria.py           # Sefaria API integration (Hebrew texts)
└── voice.py             # TTS/voice synthesis (optional)
```

---

## 📊 Data Structure

### Quote Format
```python
{
    "quote": "The quote text",
    "source": "Source attribution (chapter, book, etc.)",
    "encouragement": "Encouraging message for developers"
}
```

### Source Structure
```python
{
    "name": "Source Name",
    "icon": "📜",
    "chaos": [Quote, Quote, ...],           # 0-30% health
    "lower_aeons": [Quote, Quote, ...],     # 31-50% health
    "middle_aeons": [Quote, Quote, ...],    # 51-70% health
    "upper_aeons": [Quote, Quote, ...],    # 71-85% health
    "treasury": [Quote, Quote, ...],        # 86-100% health
    "language": "hebrew"  # Optional: for Hebrew sources
}
```

### Aeon Level Mapping
```python
def get_aeon_level(health_score: float) -> str:
    if health_score <= 30:    return "chaos"
    elif health_score <= 50:  return "lower_aeons"
    elif health_score <= 70:  return "middle_aeons"
    elif health_score <= 85:  return "upper_aeons"
    else:                     return "treasury"
```

---

## 📚 Sources Breakdown

### Local Sources (17 total, in `sources.py`)

1. **bofh** - Bastard Operator From Hell (tech humor) 😈
   - 5 aeon levels × ~3 quotes each = ~15 quotes

2. **tao** - Tao Te Ching (Lao Tzu) ☯️
   - 5 aeon levels × ~3 quotes each = ~15 quotes

3. **art_of_war** - Sun Tzu ⚔️
   - 5 aeon levels × ~3 quotes each = ~15 quotes

4. **stoic** - Marcus Aurelius & Epictetus 🏛️
   - 5 aeon levels × ~3 quotes each = ~15 quotes

5. **bible** - Proverbs & Ecclesiastes (KJV) 📖
   - 5 aeon levels × ~3 quotes each = ~15 quotes

6. **tao_of_programming** - Geoffrey James 💻
   - 5 aeon levels × ~3 quotes each = ~15 quotes

7. **murphy** - Murphy's Laws 🎲
   - 5 aeon levels × ~3 quotes each = ~15 quotes

8. **shakespeare** - The Bard 🎭
   - 5 aeon levels × ~3 quotes each = ~15 quotes

9. **confucius** - The Analects 🎓
   - 5 aeon levels × ~3 quotes each = ~15 quotes

10. **kybalion** - Hermetic Philosophy ⚗️
    - 5 aeon levels × ~3 quotes each = ~15 quotes

11. **gracian** - Art of Worldly Wisdom 🎭
    - 5 aeon levels × ~3 quotes each = ~15 quotes

12. **enochian** - John Dee's Mystical Calls 🔮
    - 5 aeon levels × ~3 quotes each = ~15 quotes

13. **rebbe** - Chassidic Wisdom (Hebrew) 🕎
    - 5 aeon levels × ~3 quotes each = ~15 quotes
    - Language: Hebrew (bilingual)

14. **tzaddik** - Path of Righteousness (Hebrew) ✡️
    - 5 aeon levels × ~3 quotes each = ~15 quotes
    - Language: Hebrew (bilingual)

15. **chacham** - Sage Wisdom (Hebrew) 📜
    - 5 aeon levels × ~3 quotes each = ~15 quotes
    - Language: Hebrew (bilingual)

### Separate Module

16. **pistis_sophia** - Gnostic Text (in `pistis_sophia.py`) 📜
    - 5 aeon levels × ~3 quotes each = ~15 quotes
    - Has additional metadata: `chapter`, `context`, `aeon_number`

### Sefaria API Sources (4 total, in `sefaria.py`)

17. **pirkei_avot** - Ethics of the Fathers (Hebrew) 🕎
18. **proverbs** - Mishlei/Proverbs (Hebrew) 📜
19. **ecclesiastes** - Kohelet/Ecclesiastes (Hebrew) 🌅
20. **psalms** - Tehillim/Psalms (Hebrew) 🎵

### Special Source

21. **random** - Daily random source selector 🎲
    - Uses date-seeded random selection
    - Picks from all available sources

---

## 🔑 Key Functions

### Core Functions (`sources.py`)

```python
get_wisdom(health_score, source=None, seed_date=True) -> dict
    # Main function to get wisdom quote
    # Returns: {quote, source, encouragement, wisdom_source, wisdom_icon, aeon_level, health_score}

get_aeon_level(health_score: float) -> str
    # Maps health score (0-100) to aeon level

get_random_source(seed_date=True) -> str
    # Returns random source ID (date-seeded for consistency)

list_available_sources() -> list[dict]
    # Returns all available sources with metadata

list_hebrew_sources() -> list[dict]
    # Returns only Hebrew sources
```

### Pistis Sophia (`pistis_sophia.py`)

```python
get_daily_wisdom(health_score, seed_date=True) -> dict
    # Returns Pistis Sophia quote with additional metadata
    # Returns: {quote, chapter, context, encouragement, aeon_level, aeon_number, health_score, source}
```

### Configuration

```python
load_config() -> dict
    # Loads from .exarp_wisdom_config or environment variables

save_config(config: dict) -> None
    # Saves configuration to file
```

---

## 📝 Quote Selection Logic

1. **Determine Aeon Level**: Based on health score (0-100)
2. **Select Source**: 
   - If `source="random"`: Use date-seeded random selection
   - If `source="pistis_sophia"`: Use separate module
   - If Sefaria source: Use API integration
   - Otherwise: Use local source from `WISDOM_SOURCES`
3. **Select Quote**: 
   - Get quotes for aeon level
   - Use date-seeded random if `seed_date=True`
   - Return first quote if no seed
4. **Format Response**: Return structured dict with quote data

---

## 🌐 Hebrew Support

### Hebrew Sources
- **Local**: `rebbe`, `tzaddik`, `chacham` (bilingual by default)
- **Sefaria API**: `pirkei_avot`, `proverbs`, `ecclesiastes`, `psalms`

### Hebrew Modes
- `EXARP_WISDOM_HEBREW=1` - Bilingual (Hebrew + English)
- `EXARP_WISDOM_HEBREW_ONLY=1` - Hebrew only

### Hebrew Quote Format
```python
{
    "quote": "Hebrew text (עברית)",
    "source": "Source in Hebrew/English",
    "encouragement": "English encouragement",
    "language": "hebrew",
    "bilingual": True  # For local sources
}
```

---

## 🎯 Porting Strategy

### Phase 1: Data Structures ✅
- ✅ Quote struct (matches Python dict)
- ✅ Source struct (matches Python dict)
- ✅ Aeon level mapping function

### Phase 2: Local Sources (Priority Order)

**Easy (Simple structure, no special handling):**
1. `stoic` - Simple quotes, no special features
2. `tao` - Simple quotes, well-structured
3. `bofh` - Tech humor, straightforward
4. `tao_of_programming` - Tech philosophy
5. `murphy` - Pragmatic laws

**Medium (More quotes, standard structure):**
6. `art_of_war` - Strategy quotes
7. `bible` - Biblical wisdom
8. `shakespeare` - Literary quotes
9. `confucius` - Ethical wisdom
10. `kybalion` - Hermetic philosophy
11. `gracian` - Pragmatic maxims
12. `enochian` - Mystical calls

**Complex (Special handling needed):**
13. `pistis_sophia` - Separate module, extra metadata
14. `rebbe` - Hebrew, bilingual
15. `tzaddik` - Hebrew, bilingual
16. `chacham` - Hebrew, bilingual

**API-Dependent (Phase 7 - Optional):**
17. `pirkei_avot` - Sefaria API
18. `proverbs` - Sefaria API
19. `ecclesiastes` - Sefaria API
20. `psalms` - Sefaria API

**System:**
21. `random` - Daily random selector (depends on other sources)

---

## 📋 Go Implementation Notes

### Data Structure Mapping

**Python → Go:**
```python
# Python
{
    "quote": "...",
    "source": "...",
    "encouragement": "..."
}
```

```go
// Go
type Quote struct {
    Quote        string `json:"quote"`
    Source       string `json:"source"`
    Encouragement string `json:"encouragement"`
    WisdomSource string `json:"wisdom_source,omitempty"`
    WisdomIcon   string `json:"wisdom_icon,omitempty"`
}
```

### Source Storage

**Option 1: Embedded in Go code** (Current approach)
- Pros: Single binary, fast access
- Cons: Large binary size, harder to update

**Option 2: JSON files**
- Pros: Easy to update, smaller binary
- Cons: File I/O, deployment complexity

**Option 3: Embedded JSON** (Recommended)
- Use `embed` package to embed JSON files
- Best of both worlds: single binary + easy updates

### Random Selection

**Date-Seeded Random:**
```go
// Python: random.seed(int(today) + hash(source))
// Go equivalent:
today := time.Now().Format("20060102")
seed := hash(today + source)
rand.Seed(seed)
```

---

## 🔍 Key Differences: Python vs Go

| Feature | Python | Go |
|---------|--------|-----|
| **Data Storage** | Dict in code | Struct + embedded JSON |
| **Random** | `random.seed()` | `rand.Seed()` |
| **Date Format** | `"%Y%m%d"` | `"20060102"` |
| **Config** | `.exarp_wisdom_config` | Same (JSON) |
| **Hebrew** | String with Unicode | Same (Go supports Unicode) |
| **API Calls** | `requests` library | `net/http` or `resty` |

---

## ✅ Next Steps

1. **Create Go data structures** matching Python format
2. **Port first 5 sources** (stoic, tao, bofh, tao_of_programming, murphy)
3. **Test quote retrieval** with different health scores
4. **Implement random selector** with date seeding
5. **Port remaining local sources**
6. **Handle Pistis Sophia** separately (extra metadata)
7. **Port Hebrew sources** (bilingual support)
8. **Sefaria API** (Phase 7 - optional)

---

## 📊 Statistics

- **Total Sources**: 21
- **Local Sources**: 17
- **API Sources**: 4 (Sefaria)
- **Total Quotes**: ~255 (17 sources × 5 levels × ~3 quotes)
- **Hebrew Sources**: 7 (3 local + 4 API)
- **Average Quotes per Source**: ~15

---

**Analysis Complete** ✅  
**Ready for Porting** ✅

