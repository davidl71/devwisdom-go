# Shared Go / Ansible / Make Patterns — Export to Shared Repository

*Analysis: exarp-go vs devwisdom-go. What can be extracted into a shared repo for reuse.*

---

## Summary

| Area | exarp-go | devwisdom-go | Shared export? |
|------|----------|--------------|----------------|
| **Go rules** | `.cursor/rules/go-development.mdc` (large) | Same + count sync from exarp | **Yes** — canonical rule file |
| **Makefile** | Config, sanity-check, many targets | Simpler; build/test/fmt/lint/clean | **Yes** — include fragment or template |
| **Ansible** | `ansible/` with roles: common, golang, linters, python | None | **Yes** — reusable roles (any Go project) |
| **Cursor MCP / rules** | `.cursor/mcp.json` + `mcp-configuration.mdc` in each project | Same servers documented in both; paths differ | **Yes** — single source of truth, optionally Ansible-managed |

---

## 1. Go development rules (Cursor / AI)

**What’s shared today:** Both use a `go-development.mdc` with the same core (style, errors, interfaces, testing). exarp-go adds: count sync, Makefile section, graph/Gonum, Todo2/SQLite, Python/uv. devwisdom-go has count sync + Makefile pointer.

**Export idea:** A single **canonical** `go-development.mdc` in a shared repo that:

- Keeps: principles, style, error handling, interfaces, testing, **count synchronization**, Makefile reminder.
- Marks or omits project-specific blocks: Gonum (optional), Todo2/SQLite (exarp-only), Python/uv (exarp-only).

**Shared repo layout (example):**

```
shared-go-patterns/   # or davidl71/cursor-rules-go
├── README.md
├── go-development.mdc    # canonical rule (copy into .cursor/rules/)
└── go-development-full.mdc   # optional: exarp-go full version (Gonum, DB, uv)
```

**Usage:** Projects copy `go-development.mdc` into `.cursor/rules/` (or symlink). Optionally keep project-specific addenda in a second file.

---

## 2. Makefile patterns

**Common subset both need:**

- Variables: `PROJECT_NAME`, `BINARY_NAME`, `VERSION`, `GO`
- Targets: `build`, `test`, `fmt`, `lint`, `clean`, `install`, `help`
- Optional: `test-coverage`, `test-html`, `lint-fix`

**exarp-go extras:** `config` + `.make.config`, `sanity-check`, `dev`/`dev-watch`, Python/uv targets, sprint/migrate. These are project-specific.

**Export idea:** Provide a **Makefile fragment** that projects include, plus a **minimal standalone** Makefile for simple Go projects.

**Shared repo layout (example):**

```
shared-go-make/
├── README.md
├── make/
│   ├── go-common.mk      # Variables + build/test/fmt/lint/clean/help
│   └── go-config.mk      # Optional: config target + HAVE_* (from exarp-go)
└── templates/
    └── Makefile.go-minimal   # Standalone minimal Makefile (like devwisdom-go)
```

**Usage:**

- **Include pattern:** In project `Makefile`: set `PROJECT_NAME`/`BINARY_NAME`, then `include path/to/shared-go-make/make/go-common.mk`. Add project-specific targets (e.g. `build-cli`, `watchdog`) in the same Makefile.
- **Copy pattern:** Copy `Makefile.go-minimal` and adjust `PROJECT_NAME`/`BINARY_NAME`/paths.

**Suggested content of `go-common.mk`:**

- `PROJECT_NAME`, `BINARY_NAME`, `VERSION`, `GO` detection.
- Targets: `build`, `test`, `fmt`, `lint`, `lint-fix`, `clean`, `install`, `help` (with `##` comments for help).
- No config, no sanity-check, no Python — keep those in project Makefiles or in an optional `go-config.mk`.

---

## 3. Ansible patterns

**exarp-go today:** `ansible/` with:

- **common:** Base packages (curl, git, make, gcc, protobuf), SQLite, optional gh, project user/dirs; Debian/RedHat/macOS.
- **golang:** Install Go from go.dev, GOROOT/GOPATH, protoc-gen-go; uses vars `go_version`, `go_install_path`, `go_gopath`, `go_bin_path`, `go_user`.
- **linters:** golangci-lint, govulncheck, shellcheck, shfmt, gomarklint, markdownlint-cli, cspell; all driven by `linters` list.
- **python:** (exarp-go also has Python; devwisdom-go doesn’t need it for Ansible.)

Roles get **no** `defaults/`; all vars come from `group_vars/all.yml` (project_name, project_path, go_version, linters list, etc.). So the roles are already generic.

**Export idea:** Put **common**, **golang**, and **linters** in a shared Ansible repo. Each project keeps its own playbooks and `group_vars` (project name, paths, go version, which linters to install).

**Shared repo layout (example):**

```
shared-ansible-dev/   # or davidl71/ansible-go-dev
├── README.md
├── requirements.yml   # collections only; roles live in this repo
└── roles/
    ├── common/        # Base system, SQLite, optional gh, project dirs
    │   ├── tasks/
    │   └── templates/
    ├── golang/        # Go install, env, protoc-gen-go
    │   └── tasks/
    └── linters/       # golangci-lint, govulncheck, shellcheck, etc.
        └── tasks/
```

**Project usage (e.g. devwisdom-go):**

- **Option A — Git submodule / subtree:** Add `shared-ansible-dev` under `ansible/shared-roles` or similar; playbooks reference `roles/common`, `roles/golang`, `roles/linters` from there.
- **Option B — Ansible Galaxy:** Publish roles as Galaxy roles (e.g. `davidl71.common_dev`, `davidl71.golang`, `davidl71.linters`). Project `requirements.yml`:

  ```yaml
  roles:
    - name: davidl71.common_dev
      src: https://github.com/davidl71/shared-ansible-dev
      version: main
      scm: git
    - name: davidl71.golang
      src: ...
    - name: davidl71.linters
      src: ...
  ```

- **Option C — Single repo clone:** Clone shared repo next to project; in project playbook set `roles_path` or reference roles by path.

**Project-side playbook (minimal):**

```yaml
# devwisdom-go/ansible/playbooks/development.yml
- hosts: development
  become: true
  vars_files:
    - ../inventories/development/group_vars/all.yml
  roles:
    - common
    - golang
    - linters  # when: install_linters | default(false)
```

**Project-side vars (`group_vars/all.yml`):**

- `project_name: "devwisdom-go"`
- `project_path: "{{ lookup('env', 'HOME') }}/Projects/devwisdom-go"`
- `go_version: "1.24.0"`
- `go_install_path: "/usr/local/go"`
- `go_user`, `go_gopath`, `go_bin_path` (same as exarp-go)
- `install_linters: true`, `linters: [golangci-lint, govulncheck, shellcheck, ...]`

So: **shared repo = roles only.** Playbooks and group_vars stay per-project.

### 3.1 Missing / recommended Ansible dependencies

Compared to what exarp-go and devwisdom-go actually use at build/test time, these are **missing or should be explicit** in the shared Ansible roles:

| Dependency | Role | Reason |
|------------|------|--------|
| **ca-certificates** | common | Required for HTTPS (`get_url`, curl, Go module proxy). Minimal Debian/RedHat images often omit it. Add: Debian `ca-certificates`, RedHat `ca-certificates`. macOS usually has it via Xcode. |
| **pkg-config** (macOS) | common | common installs pkg-config on Debian/RedHat but **not** on macOS. Some CGO/FFI builds need it. Add: `brew install pkg-config` when `ansible_system == "Darwin"`. |
| **Node.js for npm-based linters** | linters (or doc) | markdownlint-cli and cspell use `npm install -g`. Node is installed in the **python** role. If a project runs only common + golang + linters (no python), `npm` is missing. **Options:** (1) Document that the python role (or a node role) must run before linters when using markdownlint-cli/cspell; (2) Add an optional “install Node when markdownlint-cli or cspell in linters list” task to the linters role. |
| **goimports** (optional) | golang | exarp-go Makefile references gofmt/goimports; some projects run `goimports` for formatting. golang role installs protoc-gen-go but not goimports. **Optional:** add task `go install golang.org/x/tools/cmd/goimports@latest` (e.g. when `install_goimports \| default(false)`). |
| **protoc-gen-go-grpc** (optional) | golang | exarp-go only uses protoc-gen-go (no gRPC services in .proto). Projects that generate gRPC code need `protoc-gen-go-grpc`. **Optional:** add when `use_grpc \| default(false)`. |

**Python dev tools (pytest, black, ruff):** exarp-go lists them in `pyproject.toml` and installs via `uv sync` in the project. They do **not** need to be installed globally by Ansible; the Python role only needs to provide Python + uv so that `uv sync` works.

**Summary:** Add **ca-certificates** and **pkg-config (macOS)** to common; document or optionally satisfy **Node for npm linters**; optionally add **goimports** and **protoc-gen-go-grpc** in golang when variables are set.

---

## 4. Suggested shared repository layout

One repo per “kind” of shared content keeps concerns separate and lets each project adopt only what it needs.

| Repo | Purpose | Consumed by |
|------|---------|-------------|
| **shared-go-patterns** (or **cursor-rules-go**) | Canonical `go-development.mdc` (+ optional full variant) | exarp-go, devwisdom-go, future Go projects |
| **shared-go-make** | `go-common.mk` (+ optional `go-config.mk`), minimal Makefile template | exarp-go, devwisdom-go, future Go projects |
| **shared-ansible-dev** (or **ansible-go-dev**) | Ansible roles: common, golang, linters | exarp-go (refactor to use shared), devwisdom-go (add ansible), future Go projects |

Alternative: **Single repo** (e.g. `davidl71/go-project-shared`) with subdirs:

- `cursor-rules/` → go-development.mdc
- `make/` → go-common.mk, templates
- `ansible/roles/` → common, golang, linters

Then each project can submodule or copy the subdirs it needs.

---

## 5. Migration order

1. **Go rules:** Create shared repo (or subdir), add canonical `go-development.mdc`. devwisdom-go and exarp-go copy or link it; remove duplicated content and keep project-specific addenda in a separate rule file if needed.
2. **Make:** Add `shared-go-make` with `go-common.mk` and optional `go-config.mk`. devwisdom-go: optionally switch to `include go-common.mk` and trim duplication. exarp-go: optionally factor common targets into `go-common.mk` and keep config/sanity/sprint in main Makefile.
3. **Ansible:** Add `shared-ansible-dev` with roles common, golang, linters. exarp-go: point playbook to shared roles (submodule or Galaxy). devwisdom-go: add minimal `ansible/` (inventory + playbook + group_vars) that use shared roles.

---

## 6. What to keep project-specific

- **exarp-go:** Todo2/SQLite, Gonum, Python/uv, sprint/migrate, sanity-check, CGO/Apple FM logic — stay in exarp-go.
- **devwisdom-go:** Single binary + CLI, watchdog, no Ansible today — add Ansible only if you want dev environment automation; then use shared roles.
- **Cursor rules:** Project-specific rules (e.g. Todo2, openmemory, agent-locking) stay in each project; the **Go** development rule and **generic MCP configuration** (format, troubleshooting) are canonical in the shared repo. **MCP server list** can be centralized (user-level mcp.json or Ansible template) so projects don’t duplicate the same servers; see §8.

---

## 7. Architecture differences

The plan **does** rely on architecture-aware content in Ansible and **does not** yet spell out how Make/Go build handle OS and CPU arch. Below makes that explicit and keeps shared vs project responsibilities clear.

### 7.1 What “architecture” means here

- **OS:** Linux, Darwin (macOS), Windows.
- **CPU arch:** amd64 (x86_64), arm64 (aarch64).
- **Context:** Build host (native `go build`), cross-compilation (e.g. `GOOS=linux GOARCH=arm64`), and target hosts for Ansible (where dev tools run).

### 7.2 Ansible — already architecture-aware

Shared roles **must** stay OS- and arch-aware:

- **common:** Uses `ansible_os_family` (Debian, RedHat) and `ansible_system` (Darwin, Linux) for packages (apt vs yum vs Homebrew), SQLite, gh, C compiler checks. Some tasks are `when: ansible_system == "Darwin"` or `when: ansible_system == "Linux"`.
- **golang:** Go tarball URL uses `ansible_system` and `ansible_architecture` (e.g. `go1.24.0.darwin-arm64.tar.gz`, `go1.24.0.linux-amd64.tar.gz`). Go expects arch names **amd64** and **arm64**; Ansible reports **x86_64** and **aarch64** on Linux. The shared role must **map** `x86_64` → `amd64` and `aarch64` → `arm64` (e.g. via a `go_arch` variable) so the download URL is correct on all supported hosts.
- **linters:** Uses `ansible_os_family` and `ansible_system` for package installs (apt/yum/Homebrew) and for placing binaries (e.g. `go_bin_path`).

**When exporting:** Keep all `when:` conditions and URL construction; add a short **supported matrix** to the shared repo README (e.g. Debian/Ubuntu amd64 & arm64, RHEL/CentOS amd64, Darwin amd64 & arm64). Document any mapping (e.g. `aarch64` → `arm64` for Go downloads).

### 7.3 Makefile — native vs cross-build and CGO

- **Native build:** `go build` on the current host. Arch = host OS/arch.
- **Cross-build:** `GOOS`/`GOARCH` (e.g. devwisdom-go’s `build-linux`, `build-darwin` with amd64 and arm64). Used for releases and CI.
- **CGO:** exarp-go enables CGO on Mac Silicon (Darwin + arm64) for Apple Foundation Models; elsewhere it uses `CGO_ENABLED=0`. That’s **build-host** dependent.

**Shared vs project:**

- **Shared `go-common.mk`:** Assume **native build only** (single OS/arch, the host). Variables: `GO`, `VERSION`, `BINARY_NAME`. No `GOOS`/`GOARCH` in the default `build` target. That keeps the fragment simple and works for local dev on all supported OS/arch.
- **Project or optional fragment:** Cross-build targets (`build-linux`, `build-darwin`, `build-windows`, multi-arch) stay in the **project Makefile** or in an optional **`go-cross.mk`** (if you want to share a standard cross-build pattern). Same for **CGO-by-platform** logic (e.g. “on Darwin arm64 use CGO, else CGO_ENABLED=0”) — keep in project or in an optional `go-config.mk` so the shared base stays portable.

So: **architecture differences are accounted for** by (1) Ansible roles that already branch on OS/arch and must keep doing so in the shared repo, with a documented matrix and any arch-name mapping; (2) Make shared fragment that stays native-only, with cross-build and CGO handled in project or optional fragments.

### 7.4 CI and multi-arch

If CI builds or tests on several OS/arch (e.g. GitHub Actions matrix: linux-amd64, darwin-arm64, windows-amd64), that’s **project-level** (or org-level) workflow. Shared content doesn’t need to define it; the plan stays valid as long as each job runs `make build` (native) or the project’s own cross-build targets. Optional: shared repo can document “typical CI matrix” as a suggestion.

---

## 8. Cursor rules and MCP config — centralize to avoid duplication

**Problem:** The same MCP servers (devwisdom, exarp-go, tractatus_thinking, context7) are configured and documented in multiple places: exarp-go has a full `.cursor/mcp.json` and both projects have an `mcp-configuration.mdc` that lists servers with different names/paths. That causes duplicate config and drift when you add or change a server.

**What to move to a single source (optionally Ansible):**

### 8.1 MCP server list — single source of truth

- **Today:** exarp-go `.cursor/mcp.json` has 4 servers with absolute paths (`/Users/davidl/Projects/...`); devwisdom-go `.cursor/mcp.json` is empty. Each project’s `mcp-configuration.mdc` documents a slightly different list (different server names, paths, or “current” section).
- **Goal:** One list of MCP servers (name, command, args, env, description). Paths that depend on machine or project root are **templated** (e.g. `{{ projects_base }}/devwisdom-go/devwisdom`, `{{ projects_base }}/exarp-go/run_server.sh`).

**Options:**

| Approach | Where the list lives | How Cursor gets config |
|----------|----------------------|-------------------------|
| **A. User-level only** | Ansible vars (or YAML/JSON in shared repo). Ansible templates `mcp.json` to **~/.cursor/mcp.json**. | All workspaces use the same global MCP config. Cursor substitutes `{{PROJECT_ROOT}}` per workspace. |
| **B. Shared template in repo** | Shared repo holds `mcp.json.j2` (or `mcp.json` with placeholders). Script or Ansible runs once per machine (or per clone) to write `~/.cursor/mcp.json` or each project’s `.cursor/mcp.json` with paths filled in. | One template, many targets; no duplicate server definitions. |
| **C. Per-project override only** | User-level has the full set; each project’s `.cursor/mcp.json` only adds **project-specific** servers (e.g. “advisor” pointing at this repo’s binary) or is empty. | Shared servers live in ~/.cursor/mcp.json; projects only add local overrides. |

**Recommended:** **A** or **C**. Use **A** if you want Ansible to own the full MCP list and deploy to ~/.cursor/mcp.json so every workspace sees the same servers. Use **C** if you want user-level to hold the common set and projects to only add/override what’s specific (e.g. path to this project’s MCP binary).

### 8.2 What Ansible can do

- **Role or tasks:** e.g. `cursor_config` or a “developer dotfiles” role.
- **Vars:** `cursor_projects_base: "/Users/{{ ansible_user_id }}/Projects"` (or from env), plus a list of MCP servers, e.g.:

  ```yaml
  cursor_mcp_servers:
    - name: devwisdom
      command: "{{ cursor_projects_base }}/devwisdom-go/run_devwisdom.sh"
      args: []
      env: { PROJECT_ROOT: "{{PROJECT_ROOT}}" }
      description: "DevWisdom Go MCP Server - Wisdom, advisors, guidance"
    - name: exarp-go
      command: "{{ cursor_projects_base }}/exarp-go/run_server.sh"
      ...
    - name: tractatus_thinking
      command: npx
      args: ["-y", "tractatus_thinking"]
      ...
    - name: context7
      command: npx
      args: ["-y", "@upstash/context7-mcp"]
      ...
  ```

- **Task:** Template a `mcp.json.j2` that loops over `cursor_mcp_servers` and write to **~/.cursor/mcp.json** (or to a shared path that projects symlink). Use `ansible_user_id` and `cursor_projects_base` so paths are correct per user/machine.
- **Optional:** Same vars can drive a “Current MCP Servers” section in a shared rule file so docs and config stay in sync.

### 8.3 Cursor rules to move to shared (or Ansible-deployed)

| Rule / content | Today | Move to shared? |
|-----------------|--------|------------------|
| **mcp-configuration.mdc** | In both projects; “Current MCP Servers” and format/troubleshooting duplicated. | **Yes.** Put **generic** content (server format, env vars, transport, troubleshooting, best practices) in a **shared** rule (e.g. in shared-go-patterns or cursor-rules-go). “Current MCP Servers” either (1) generated from same vars that generate mcp.json, or (2) replaced by “See ~/.cursor/mcp.json; standard set is devwisdom, exarp-go, tractatus_thinking, context7.” |
| **go-development.mdc** | Already planned for shared canonical version. | As in §1. |
| **Project-only rules** | e.g. agent-locking, agentic-ci, session-prime (exarp-go); exarp-go-patterns, openmemory refs (devwisdom-go). | **No.** Keep in each project; they describe project-specific behavior. |

So: **move MCP server list + generic mcp-configuration content** to a single source (Ansible vars + template, or shared repo); **optionally** have Ansible deploy `mcp.json` and shared rules to user Cursor dir or to each project’s `.cursor/` so you don’t duplicate MCP servers or the same rule text across projects.

### 8.4 What stays in each project

- **.cursor/mcp.json:** Either **empty** (rely on user-level) or **minimal override** (e.g. only this project’s server with a path relative to repo or `{{PROJECT_ROOT}}`).
- **.cursor/rules/:** Project-specific rules only (agent-locking, exarp-go-patterns, etc.). Shared rules (go-development, mcp-configuration generic) come from shared repo or Ansible; project copies/symlinks or Ansible deploys into `.cursor/rules/` or into user-level.

---

## 9. File reference

| Source | Path |
|--------|------|
| exarp-go Go rule | `exarp-go/.cursor/rules/go-development.mdc` |
| exarp-go Makefile | `exarp-go/Makefile` |
| exarp-go Ansible | `exarp-go/ansible/` (roles: common, golang, linters, python; playbooks; group_vars) |
| exarp-go MCP config | `exarp-go/.cursor/mcp.json` |
| exarp-go MCP rule | `exarp-go/.cursor/rules/mcp-configuration.mdc` |
| devwisdom-go Makefile | `devwisdom-go/Makefile` |
| devwisdom-go Go rule | `devwisdom-go/.cursor/rules/go-development.mdc` |
| devwisdom-go MCP config | `devwisdom-go/.cursor/mcp.json` (currently empty) |
| devwisdom-go MCP rule | `devwisdom-go/.cursor/rules/mcp-configuration.mdc` |

This doc can live in both exarp-go and devwisdom-go (e.g. `docs/SHARED_GO_ANSIBLE_MAKE_PATTERNS.md`) and be updated when shared repos are created or layout changes.
