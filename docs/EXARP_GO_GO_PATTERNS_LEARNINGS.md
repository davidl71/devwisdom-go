# exarp-go Go Patterns → devwisdom-go Learnings

*Checked: exarp-go `.cursor/rules` and internal structure vs devwisdom-go.*

## Summary

exarp-go’s **go-development.mdc** has grown with patterns that are worth reusing in devwisdom-go where they apply. Below: what’s new in exarp-go, what we adopted, and what we skipped.

---

## 1. Count synchronization (adopt)

**exarp-go:** A **“Count Synchronization (CRITICAL)”** section: when registering tools/prompts/resources, keep these in sync:

- Count in implementation
- Comments that mention the count
- Test assertions (e.g. `want 17` → `want 18`)
- Expected lists in tests

**Relevance for devwisdom-go:** We already assert counts in tests (`len(tools)`, `len(sources)`, cache size, advisor counts in `internal/mcp/sdk_adapter_test.go`, `server_test.go`, wisdom tests). Adding/removing a tool or resource can easily desync tests.

**Action:** The same “count sync” guidance and checklist from exarp-go were added to devwisdom-go’s **go-development.mdc** so we follow it when changing tools/resources/prompts.

---

## 2. Makefile usage (reference only)

**exarp-go:** Long, project-specific list of Makefile targets (build, dev, test, lint, config, sprint automation, migrate, uv/python, etc.) and “never run `go build` / `go test` directly” rules.

**devwisdom-go:** Makefile is smaller and different (no config, no sanity-check, no sprint/migrate). **.cursorrules** already has the “always use Makefile targets” rule and the list of devwisdom-go targets.

**Action:** No duplication of the full exarp-go list. We added a short pointer in **go-development.mdc** to use Makefile targets and to see .cursorrules for the full list.

---

## 3. Graph algorithms / Gonum (skip for now)

**exarp-go:** Detailed guidance on **gonum.org/v1/gonum/graph**: TaskGraph wrapper, cycle detection, critical path, topo sort, when to use vs “when NOT to use Gonum.”

**devwisdom-go:** No task dependency graph or similar; wisdom/MCP server only.

**Action:** Not added to devwisdom-go. If we later add dependency/graph logic, we can copy this section from exarp-go’s go-development.mdc.

---

## 4. Todo2 / SQLite (exarp-go–specific)

**exarp-go:** Todo2 SQLite as primary storage, `LoadTodo2Tasks`/`SaveTodo2Tasks`, direct DB ops, migration tool, when to use DB vs JSON.

**devwisdom-go:** No Todo2/SQLite; different domain.

**Action:** Not added. Keep as exarp-go–specific reference only.

---

## 5. Python bridge / uv (exarp-go–specific)

**exarp-go:** “PREFER `uv run python`”, “Check for uv first”, Python command guidelines.

**devwisdom-go:** Go-only MCP server; no Python bridge.

**Action:** Not added.

---

## 6. Agent locking (exarp-go–specific)

**exarp-go:** `.cursor/rules/agent-locking.mdc` for agent IDs and task locking in parallel execution.

**devwisdom-go:** Single-process MCP server; no multi-agent task locking.

**Action:** Not added. Revisit if we add parallel agents/task locking.

---

## 7. Other exarp-go rules

- **agentic-ci.mdc, session-prime.mdc, llm-tools.mdc, mcp-configuration.mdc** – project/process specific; no Go pattern changes pulled into devwisdom-go.

---

## Files updated in devwisdom-go

| File | Change |
|------|--------|
| `.cursor/rules/go-development.mdc` | Added “Count Synchronization” section; added one-sentence Makefile pointer to .cursorrules. |
| `docs/EXARP_GO_GO_PATTERNS_LEARNINGS.md` | New: this learnings doc. |

---

## Reference: exarp-go locations

- **Go patterns:** `exarp-go/.cursor/rules/go-development.mdc`
- **Exarp patterns we already mirror:** `devwisdom-go/.cursor/rules/exarp-go-patterns.mdc` and `docs/EXARP_GO_LESSONS.md`
