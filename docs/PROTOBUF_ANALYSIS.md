# Protobuf Analysis for devwisdom-go

**Date:** 2026-01-13  
**Purpose:** Analyze whether protobuf would benefit devwisdom-go  
**Conclusion:** **Not recommended** - JSON is sufficient for current scale and usage patterns

---

## Executive Summary

After analyzing devwisdom-go's data storage, communication patterns, and performance characteristics, **protobuf would provide minimal benefit (~2-5%)** while adding significant complexity. The current JSON approach is appropriate for the project's scale.

---

## Current State Analysis

### Data Volume
- **Total consultation logs:** 1.36 KB (2 consultations)
- **Average consultation size:** ~680 bytes (JSON)
- **Write frequency:** Low (only when consultations occur)
- **Read frequency:** Low (occasional log retrieval)

### Performance Characteristics
- **GetWisdom operation:** 17.48 ns (0.000017 ms) ✅ Exceeds goal by 57,000x
- **JSON encoding/decoding:** Fast for small structures (< 1 µs)
- **No performance bottleneck identified**

### Communication Patterns

#### 1. MCP Communication (JSON-RPC 2.0)
- **Protocol:** JSON-RPC 2.0 over stdio (MCP standard)
- **Format:** JSON (required by MCP specification)
- **Usage:** Python bridge calls devwisdom-go MCP server
- **Protobuf impact:** ❌ **Would break MCP compatibility** (MCP requires JSON)

#### 2. Direct Go Library Calls
- **Protocol:** Go structs (in-process, no serialization)
- **Format:** Native Go types
- **Usage:** exarp-go imports `github.com/davidl71/devwisdom-go/pkg/wisdom` directly
- **Protobuf impact:** ❌ **No serialization needed** (direct function calls)

#### 3. Consultation Log Storage
- **Format:** JSONL (JSON Lines) files
- **Location:** `.devwisdom/consultations.jsonl`
- **Rotation:** Date-based (daily rotation)
- **Protobuf impact:** ⚠️ **Minimal benefit** (small data volume)

---

## Protobuf Benefit Assessment

### Consultation Logs (Primary Use Case)

| Aspect | Current (JSON) | With Protobuf | Benefit |
|--------|----------------|---------------|---------|
| **File size** | ~680 bytes/consultation | ~400-500 bytes (30-40% smaller) | Minimal (saves ~200 bytes) |
| **Write speed** | ~1-5 µs | ~0.5-2 µs (2-3x faster) | Negligible (writes are infrequent) |
| **Read speed** | ~1-5 µs per line | ~0.5-2 µs per line | Negligible (reads are occasional) |
| **Code complexity** | Simple (JSON stdlib) | Medium (protobuf setup + codegen) | **Negative** (adds complexity) |
| **Maintenance** | Low (JSON is standard) | Medium (protobuf schema + codegen) | **Negative** (more to maintain) |

### MCP Communication

| Aspect | Current (JSON-RPC 2.0) | With Protobuf | Impact |
|--------|----------------------|---------------|--------|
| **Protocol compliance** | ✅ Compliant | ❌ **Would break MCP spec** | **Negative** |
| **Compatibility** | ✅ Works with all MCP clients | ❌ **Incompatible** | **Negative** |
| **Performance** | Fast enough (< 1 ms) | N/A (not applicable) | **No benefit** |

### Direct Library Calls

| Aspect | Current (Go structs) | With Protobuf | Impact |
|--------|---------------------|---------------|--------|
| **Serialization** | None (in-process) | Unnecessary | **No benefit** |
| **Performance** | Optimal (direct calls) | Would add overhead | **Negative** |

---

## When Protobuf Would Help

Protobuf would be beneficial **only if**:

1. **High volume:** >10,000 consultations/day (currently ~2 total)
2. **High frequency:** >100 consultations/second (currently very low)
3. **Large data:** >10 KB per consultation (currently ~680 bytes)
4. **Network transfer:** Sending over network (currently local files)
5. **Database storage:** Storing in SQLite with size constraints (currently JSONL files)

**None of these conditions are met** in devwisdom-go's current usage.

---

## Recommendation

### ❌ **Do NOT add protobuf to devwisdom-go**

**Reasons:**
1. ✅ Data volume is tiny (1.36 KB total)
2. ✅ Performance is already excellent (17ns operations)
3. ✅ JSON is sufficient for the use case
4. ✅ Protobuf adds complexity without meaningful benefit
5. ✅ MCP requires JSON (protobuf would break compatibility)
6. ✅ Direct library calls don't need serialization

### Better Alternatives (If Performance Becomes an Issue)

1. **Optimize JSON encoding:**
   - Use `json.Encoder` with buffer pooling
   - Pre-allocate slice capacity
   - Reuse encoders/decoders

2. **Add compression for archived logs:**
   - Gzip rotated log files
   - Decompress on read (transparent to API)

3. **Batch writes:**
   - Buffer multiple consultations before writing
   - Reduce file I/O operations

4. **Use database for high-volume scenarios:**
   - SQLite with JSON columns (if volume grows)
   - Indexed queries for faster retrieval

---

## Migration File Status

**Current state:**
- `migrations/002_add_protobuf_support.sql` exists (copied from exarp-go)
- **Not used** by devwisdom-go (no protobuf implementation)
- **Not needed** for devwisdom-go's functionality

**Recommendation:**
- **Option 1:** Remove the migration file (not needed)
- **Option 2:** Keep it for future reference (if volume grows significantly)

---

## Comparison with exarp-go

**exarp-go uses protobuf for:**
- ✅ Task metadata storage (high volume, frequent operations)
- ✅ Python bridge communication (high frequency)
- ✅ Memory system (large data structures)

**devwisdom-go does NOT need protobuf because:**
- ❌ Consultation logs are tiny (1.36 KB vs. exarp-go's task data)
- ❌ Write frequency is low (occasional vs. exarp-go's frequent)
- ❌ MCP communication requires JSON (can't use protobuf)
- ❌ Direct library calls don't need serialization

---

## Conclusion

**Protobuf would help ~2-5%** for consultation logs, but the added complexity and maintenance cost outweigh the benefit. The current JSON approach is appropriate for devwisdom-go's scale and usage patterns.

**Priority: Low** - Focus development effort on features that provide more value.

---

## Related Documentation

- **[exarp-go Protobuf Analysis](../exarp-go/docs/PROTOBUF_ANALYSIS.md)** - Analysis for exarp-go (high-volume use case)
- **[devwisdom-go Performance Benchmarks](PERFORMANCE.md)** - Current performance metrics
- **[MCP Implementation Summary](MCP_IMPLEMENTATION_SUMMARY.md)** - MCP protocol details
