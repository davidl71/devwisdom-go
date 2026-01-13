# Protobuf Decision for devwisdom-go

**Date:** 2026-01-13  
**Status:** Decision Made - **Not Recommended**  
**Related Analysis:** [PROTOBUF_ANALYSIS.md](PROTOBUF_ANALYSIS.md)

---

## Quick Summary

**Question:** Should devwisdom-go use protobuf for data serialization?  
**Answer:** **No** - JSON is sufficient for current scale and usage patterns.

---

## Key Findings

### Current State
- **Data volume:** 1.36 KB total (2 consultations)
- **Performance:** 17.48 ns operations (exceeds goals by 57,000x)
- **Communication:** JSON-RPC 2.0 (MCP standard) + Direct Go library calls
- **Storage:** JSONL files (consultation logs)

### Protobuf Impact
- **Benefit:** ~2-5% improvement (minimal)
- **Cost:** +20-30% code complexity, +protobuf dependency, +codegen step
- **ROI:** **Negative** - complexity outweighs benefit

### Communication Patterns
1. **MCP Communication:** JSON-RPC 2.0 (protobuf would break compatibility)
2. **Direct Library Calls:** Go structs (no serialization needed)
3. **Consultation Logs:** JSONL files (tiny data volume)

---

## Decision Rationale

### Why NOT Protobuf?

1. **MCP Requires JSON** - MCP specification uses JSON-RPC 2.0, protobuf would break compatibility
2. **Tiny Data Volume** - 1.36 KB total doesn't justify binary serialization
3. **Already Fast** - 17ns operations are already optimal
4. **Complexity Cost** - Protobuf adds schema management, codegen, and maintenance overhead
5. **No Serialization Needed** - Direct Go library calls don't need serialization

### When Protobuf Would Help

Only if:
- >10,000 consultations/day (currently ~2 total)
- >100 consultations/second (currently very low)
- >10 KB per consultation (currently ~680 bytes)
- Network transfer needed (currently local files)
- Database storage with size constraints (currently JSONL files)

**None of these conditions are met.**

---

## Migration File Status

**Current:** `migrations/002_add_protobuf_support.sql` exists (copied from exarp-go)  
**Status:** Not used, not needed  
**Recommendation:** Can be removed or kept for future reference

---

## Related Documentation

- **[PROTOBUF_ANALYSIS.md](PROTOBUF_ANALYSIS.md)** - Detailed analysis with metrics
- **[exarp-go Protobuf Analysis](../../exarp-go/docs/PROTOBUF_ANALYSIS.md)** - Analysis for exarp-go (high-volume use case)
- **[PERFORMANCE.md](PERFORMANCE.md)** - Current performance benchmarks

---

## Conclusion

**Protobuf is not recommended for devwisdom-go.** The current JSON approach is appropriate for the project's scale, and protobuf would add complexity without meaningful benefit.

**Priority:** Low - Focus development effort on features that provide more value.
