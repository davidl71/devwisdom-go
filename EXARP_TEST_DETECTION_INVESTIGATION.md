# EXARP Test Detection Investigation

**Date**: 2026-01-06  
**Status**: ✅ Investigation Complete  
**Issue**: EXARP scorecard reported 0% test coverage despite 14 Go test files existing

---

## 🔍 Problem Statement

EXARP project scorecard tool reported:
- **0 test files**
- **0 test lines**  
- **0% test coverage**

However, the project actually has:
- **14 Go test files** (`*_test.go`)
- **672 lines of test code**
- **75.7% actual test coverage** (verified via `go test`)

---

## 📊 Actual Test Coverage

### Test Files Inventory

**internal/wisdom/** (6 test files):
- `engine_test.go`
- `sources_test.go`
- `sources_config_test.go`
- `sources_cache_test.go`
- `sources_api_test.go`
- `advisors_test.go`
- `types_test.go`

**internal/cli/** (6 test files):
- `app_test.go`
- `quote_test.go`
- `consult_test.go`
- `sources_test.go`
- `advisors_test.go`
- `briefing_test.go`

**internal/logging/** (1 test file):
- `consultation_log_test.go`

### Verified Coverage Metrics

From `TESTING_SUMMARY.md` (2025-01-26):
- **Wisdom package**: 75.7% coverage
- **All tests pass**: ✅
- **Test patterns**: Table-driven, HTTP test servers, temporary files, time-based tests

---

## 🔬 Root Cause Analysis

### EXARP Tool Design

**EXARP** is part of the `project-management-automation` Python-based tool suite. The scorecard tool appears to be designed primarily for **Python projects** and looks for Python test patterns:

- ✅ `test_*.py` files
- ✅ `*_test.py` files

**Go test patterns are NOT detected:**
- ❌ `*_test.go` files (Go standard)

### Evidence

1. **EXARP metrics output**:
   ```json
   "testing": {
     "test_files": 0,
     "test_lines": 0,
     "test_ratio": 0.0
   }
   ```

2. **No Go-specific configuration found**:
   - No `.exarp/config.json` or similar
   - No test file pattern configuration
   - No language-specific settings

3. **EXARP is Python-focused**:
   - Part of Python project (`project-management-automation`)
   - Uses Python patterns for detection
   - No evidence of multi-language support

---

## 💡 Solutions & Recommendations

### Short-Term (Immediate)

1. **Document Actual Coverage Manually**:
   - Add note to scorecard: "Actual Go test coverage: 75.7% (14 test files, 672 lines)"
   - Update scorecard recommendations to reflect real testing status

2. **Use Go Native Tools**:
   ```bash
   go test ./... -coverprofile=coverage.out
   go tool cover -func=coverage.out
   ```

### Medium-Term

1. **Check EXARP Configuration Options**:
   - Investigate if EXARP supports custom test file patterns
   - Look for language-specific configuration
   - Check if `.exarp/config.json` can override defaults

2. **Create Workaround Script**:
   - Generate Go coverage report
   - Parse and inject into EXARP scorecard output
   - Maintain separate Go coverage tracking

### Long-Term

1. **Request Go Language Support**:
   - File feature request with EXARP maintainers
   - Request `*_test.go` pattern detection
   - Request Go test coverage integration

2. **Alternative Tools**:
   - Use Go-specific coverage tools (gocov, goveralls)
   - Integrate with CI/CD for automated coverage tracking
   - Consider Go-native project health tools

---

## 📝 Current Status

**Actual Test Coverage**: 75.7% ✅  
**EXARP Reported Coverage**: 0% ❌  
**Gap**: EXARP detection limitation, not actual test deficiency

**Recommendation**: 
- Continue using Go native tools (`go test -cover`) for accurate metrics
- Document EXARP limitation in project notes
- Manually update scorecard with real coverage numbers when needed

---

## 🔗 References

- [Go Test File Detection Issues](https://community.sonarsource.com/t/golang-test-files-are-not-analyzed-after-global-exclusions-pattern-change/107849)
- [Go Test Runner Documentation](https://pkg.go.dev/github.com/exercism/go-test-runner)
- `TESTING_SUMMARY.md` - Actual test coverage documentation
- `devwisdom-go-scorecard.md` - EXARP generated scorecard

---

**Investigation Complete**: 2026-01-06  
**Next Steps**: Document limitation, use Go native tools for accurate metrics

