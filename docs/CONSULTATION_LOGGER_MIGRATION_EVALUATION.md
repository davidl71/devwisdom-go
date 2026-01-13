# ConsultationLogger Migration Evaluation

**Date:** 2026-01-13  
**Task:** T-1768325421830  
**Status:** 📊 Evaluation Complete

---

## Executive Summary

**Recommendation:** ❌ **DO NOT MIGRATE** ConsultationLogger to mcp-go-core

**Reason:** ConsultationLogger is **highly project-specific** and tightly coupled to devwisdom-go's domain model. The date-based rotation pattern could potentially be extracted, but the overall implementation is too domain-specific for a generic library.

---

## Analysis

### ConsultationLogger Characteristics

**Location:** `internal/logging/consultation_log.go`  
**Size:** ~265 lines  
**Dependencies:**
- `internal/wisdom.Consultation` struct (domain-specific)
- JSONL file format (project-specific)
- Date-based rotation with specific naming pattern

### Key Features

1. **Date-Based Log Rotation**
   - Rotates log files when date changes
   - Naming pattern: `consultations-YYYY-MM-DD.jsonl`
   - Automatic rotation on date boundary

2. **JSONL Format**
   - Writes consultations as JSON lines
   - One consultation per line
   - Thread-safe writes with mutex

3. **Log Reading with Date Filtering**
   - Reads from current and rotated files
   - Filters by date range (last N days)
   - Handles timestamp parsing and validation

4. **Thread Safety**
   - Mutex-protected writes
   - Safe for concurrent access

---

## Migration Feasibility

### ✅ Could Be Extracted (Generic Patterns)

1. **Date-Based Rotation Pattern**
   - Generic file rotation by date
   - Could be: `pkg/mcp/logging/rotating_file.go`
   - Pattern: `RotatingFileLogger` with configurable rotation strategy

2. **JSONL Writer Pattern**
   - Generic JSONL file writer
   - Could be: `pkg/mcp/logging/jsonl_writer.go`
   - Pattern: `JSONLWriter` with generic type parameter

### ❌ Should Stay in devwisdom-go (Project-Specific)

1. **Consultation Domain Model**
   - Tightly coupled to `wisdom.Consultation` struct
   - Domain-specific fields (Advisor, Quote, Score, etc.)
   - Not reusable across projects

2. **Project-Specific Naming**
   - File naming: `consultations.jsonl`
   - Directory: `.devwisdom`
   - Project-specific conventions

3. **Business Logic**
   - Consultation-specific filtering logic
   - Timestamp parsing for consultations
   - Domain-specific error handling

---

## Comparison with mcp-go-core Logging

### mcp-go-core/pkg/mcp/logging
- **Purpose:** Structured application logging (stderr)
- **Format:** Text-based log lines with levels
- **Use Case:** General application logging, request tracing
- **Scope:** Generic, reusable across all MCP servers

### devwisdom-go/internal/logging/consultation_log.go
- **Purpose:** Domain-specific consultation logging (JSONL files)
- **Format:** JSON lines (one consultation per line)
- **Use Case:** Wisdom consultation history tracking
- **Scope:** Project-specific, tightly coupled to wisdom domain

**Conclusion:** They serve **different purposes** and are **complementary**, not redundant.

---

## Alternative Approaches

### Option 1: Keep as-is (Recommended) ✅
- **Pros:**
  - No migration effort
  - Maintains project-specific functionality
  - Clear separation of concerns
- **Cons:**
  - Code stays in devwisdom-go
- **Effort:** 0 hours

### Option 2: Extract Generic Patterns
- **Pros:**
  - Reusable date-based rotation
  - Reusable JSONL writer
- **Cons:**
  - Significant refactoring effort
  - ConsultationLogger still project-specific
  - Limited benefit (pattern extraction only)
- **Effort:** 4-6 hours
- **Value:** Low (pattern extraction, not full migration)

### Option 3: Full Migration with Generics
- **Pros:**
  - Fully generic implementation
- **Cons:**
  - Requires Go generics (Go 1.18+)
  - Complex refactoring
  - ConsultationLogger still needs wrapper
- **Effort:** 8-10 hours
- **Value:** Low (over-engineering for single use case)

---

## Recommendation

### ✅ **Keep ConsultationLogger in devwisdom-go**

**Rationale:**
1. **Domain-Specific:** Tightly coupled to `wisdom.Consultation` domain model
2. **Project-Specific:** File naming, directory structure, business logic all specific to devwisdom-go
3. **Low Reusability:** Unlikely to be useful for other MCP servers
4. **Complementary:** Works alongside mcp-go-core logging (different purposes)
5. **Maintenance:** Easier to maintain in project where it's used

### 📝 **Optional: Extract Generic Patterns (Future)**

If date-based rotation or JSONL writing patterns become needed in other projects:
- Extract `RotatingFileLogger` pattern to mcp-go-core
- Extract `JSONLWriter` pattern to mcp-go-core
- Keep ConsultationLogger in devwisdom-go but use extracted patterns internally

**Priority:** Low (only if pattern is needed elsewhere)

---

## Current State

**Status:** ✅ **Appropriately Located**

- ConsultationLogger is in `internal/logging/` (project-specific)
- Uses mcp-go-core logging for application logs (complementary)
- Clear separation: domain logging vs application logging
- No migration needed

---

## Files Using ConsultationLogger

1. `internal/mcp/sdk_adapter.go` - Creates and uses logger
2. `internal/mcp/server.go` - Creates and uses logger (legacy)
3. `internal/mcp/handlers.go` - Uses logger for consultation logging
4. `internal/logging/consultation_log_test.go` - Comprehensive tests

**All usage is appropriate and project-specific.**

---

## Conclusion

**Decision:** ❌ **Do not migrate ConsultationLogger**

**Action Items:**
- ✅ Keep ConsultationLogger in `internal/logging/`
- ✅ Continue using mcp-go-core logging for application logs
- ✅ Document that ConsultationLogger is intentionally project-specific
- ⏸️ Consider pattern extraction only if needed in other projects

**Migration Status:** ✅ **Complete** (evaluation shows migration not needed)

---

**Last Updated:** 2026-01-13  
**Next Review:** Only if pattern extraction becomes needed
