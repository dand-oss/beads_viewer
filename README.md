# Beads Viewer (bv)

A polished, high-performance TUI for managing and exploring [Beads](https://github.com/steveyegge/beads) issue trackers.

## Features

### 🖥️ Slick Dashboard
*   **Adaptive Split-View**: Automatically transitions to a master-detail dashboard on wide screens (>100 cols).
*   **Ultra-Wide Mode**: Shows extended metadata (Age, Updated At, Comments) on large displays (>180 cols).
*   **Live Stats**: Persistent stats bar showing Open, Ready, Blocked, and Closed counts.

### 🎨 Rich Visualization
*   **Markdown Rendering**: Beautiful rendering of issue content with syntax highlighting (via `glamour`).
*   **Dracula Theme**: A vibrant, high-contrast color scheme.
*   **Dependency Graph**: Visualizes blockers (⛔) and related (🔗) issues.

### ⚡ Workflow
*   **Instant Filtering**: 
    *   `o`: **Open** only
    *   `r`: **Ready** work (Open & Unblocked)
    *   `c`: **Closed**
    *   `a`: **All**
*   **Search**: Use the list's fuzzy search (type to filter).
*   **Export**: Generate comprehensive Markdown reports with `--export-md`.

### 🛠️ Robust & Reliable
*   **Self-Updating**: Automatically checks for new releases.
*   **Resilient Loader**: Handles large or partially malformed JSONL databases gracefully.

## Installation

### One-line Install
```bash
curl -fsSL https://raw.githubusercontent.com/Dicklesworthstone/beads_viewer/main/install.sh | bash
```

### From Source
```bash
go install github.com/Dicklesworthstone/beads_viewer/cmd/bv@latest
```

## Usage

Navigate to any project initialized with `bd init` and run:

```bash
bv
```

To export a report:
```bash
bv --export-md report.md
```

### Keybindings

| Key | Context | Action |
| :--- | :--- | :--- |
| `Tab` | Split View | Switch focus between List and Details |
| `j` / `k` | Global | Navigate list or scroll details |
| `Enter` | List | Open details (Mobile) or Focus details (Split) |
| `o` / `r` / `c` / `a` | Global | Filter by status |
| `q` | Global | Quit |

## CI/CD

This project uses GitHub Actions for:
*   **Tests**: Runs full unit and integration suite on every push.
*   **Releases**: Automatically builds and attaches optimized binaries for Linux, macOS, and Windows to every release tag.

## License

MIT
