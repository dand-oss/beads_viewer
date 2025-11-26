# Agents Guide to Beads Viewer (bv)

`bv` is a high-performance TUI for the [Beads](https://github.com/steveyegge/beads) issue tracker.

## Features

- **Split View Dashboard**: On wide terminals (>100 cols), shows a list on the left and rich details on the right.
- **Markdown Rendering**: Renders issue descriptions, notes, and comments with syntax highlighting.
- **Live Filtering**: Filter by Open (`o`), Closed (`c`), Ready (`r`), or All (`a`).
- **Dependency Graph**: Visualizes blockers and dependencies.

## Navigation

### Global
- `q` / `Ctrl+C`: Quit
- `Tab`: Switch focus between List and Details pane (Split View only)

### List View
- `j` / `↓`: Next issue
- `k` / `↑`: Previous issue
- `Enter`: Open details (Mobile view) or Focus details (Split view)
- `o`: Filter Open
- `c`: Filter Closed
- `r`: Filter Ready (Open + Unblocked)
- `a`: Show All

### Details View
- `j` / `k` / Arrows: Scroll content
- `Esc`: Back to list (Mobile view)

## Installation

```bash
./install.sh
```

## Development

Built with Go + Charmbracelet (Bubble Tea, Lipgloss, Glamour).
Follows `GOLANG_BEST_PRACTICES.md`.


---

### ast-grep vs ripgrep (quick guidance)

**Use `ast-grep` when structure matters.** It parses code and matches AST nodes, so results ignore comments/strings, understand syntax, and can **safely rewrite** code.

* Refactors/codemods: rename APIs, change import forms, rewrite call sites or variable kinds.
* Policy checks: enforce patterns across a repo (`scan` with rules + `test`).
* Editor/automation: LSP mode; `--json` output for tooling.

**Use `ripgrep` when text is enough.** It’s the fastest way to grep literals/regex across files.

* Recon: find strings, TODOs, log lines, config values, or non-code assets.
* Pre-filter: narrow candidate files before a precise pass.

**Rule of thumb**

* Need correctness over speed, or you’ll **apply changes** → start with `ast-grep`.
* Need raw speed or you’re just **hunting text** → start with `rg`.
* Often combine: `rg` to shortlist files, then `ast-grep` to match/modify with precision.

**Snippets**

Find structured code (ignores comments/strings):

```bash
ast-grep run -l TypeScript -p 'import $X from "$P"'
```

Codemod (only real `var` declarations become `let`):

```bash
ast-grep run -l JavaScript -p 'var $A = $B' -r 'let $A = $B' -U
```

Quick textual hunt:

```bash
rg -n 'console\.log\(' -t js
```

Combine speed + precision:

```bash
rg -l -t ts 'useQuery\(' | xargs ast-grep run -l TypeScript -p 'useQuery($A)' -r 'useSuspenseQuery($A)' -U
```

**Mental model**

* Unit of match: `ast-grep` = node; `rg` = line.
* False positives: `ast-grep` low; `rg` depends on your regex.
* Rewrites: `ast-grep` first-class; `rg` requires ad-hoc sed/awk and risks collateral edits.

---

## UBS Quick Reference for AI Agents

UBS stands for "Ultimate Bug Scanner": **The AI Coding Agent's Secret Weapon: Flagging Likely Bugs for Fixing Early On**

**Install:**

```bash
curl -sSL https://raw.githubusercontent.com/Dicklesworthstone/ultimate_bug_scanner/main/install.sh | bash
```

**Golden Rule:** `ubs <changed-files>` before every commit. Exit 0 = safe. Exit >0 = fix & re-run.

**Commands:**

```bash
ubs file.ts file2.ts                    # Specific files (< 1s) — USE THIS
ubs $(git diff --name-only --cached)    # Staged files — before commit
ubs --only=js,ts src/                   # Language filter (3-5x faster)
ubs --ci --fail-on-warning .            # CI mode — before PR
ubs --help                              # Full command reference
ubs sessions --entries 1                # Tail the latest install session log
ubs .                                   # Whole project (ignores things like .next, node_modules automatically)
```

**Output Format:**

```text
⚠️  Category (N errors)
    file.ts:42:5 – Issue description
    💡 Suggested fix
Exit code: 1
```

Parse: `file:line:col` → location | 💡 → how to fix | Exit 0/1 → pass/fail

**Fix Workflow:**

1. Read finding → category + fix suggestion.
2. Navigate `file:line:col` → view context.
3. Verify real issue (not false positive).
4. Fix root cause (not symptom).
5. Re-run `ubs <file>` → exit 0.
6. Commit.

**Speed Critical:** Scope to changed files. `ubs src/file.ts` (< 1s) vs `ubs .` (30s). Never full scan for small edits.

**Bug Severity:**

* **Critical** (always fix): null/undefined safety, injection vulnerabilities, race conditions, resource leaks.
* **Important** (production): type narrowing, error handling, performance landmines.
* **Contextual** (judgment): TODO/FIXME, excessive console logs.

**Anti-Patterns:**

* ❌ Ignore findings → ✅ Investigate each.
* ❌ Full scan per edit → ✅ Scope to changed files.
* ❌ Fix symptom only → ✅ Fix root cause.
