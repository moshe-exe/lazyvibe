# CLAUDE.md

Development guide for lazyvibe - a lazygit-inspired terminal dashboard for monitoring Claude Code sessions.

## Quick Reference

```bash
make build     # Build binary
make run       # Build and run
./lazyvibe   # Run dashboard
```

## Architecture

```
lazyvibe/
├── cmd/lazyvibe/main.go    # Entry point, CLI flags
├── internal/
│   ├── config/               # Themes (5 built-in) and settings
│   ├── data/                 # Data layer
│   │   ├── models.go         # Session, Project, Stats structs
│   │   ├── sessions.go       # Parse ~/.claude/projects/*/sessions-index.json
│   │   ├── stats.go          # Parse ~/.claude/stats-cache.json
│   │   ├── vm.go             # Claude Desktop VM process monitoring
│   │   └── manager.go        # Caching data aggregator
│   ├── ui/                   # Bubbletea TUI
│   │   ├── app.go            # Main model, panel navigation
│   │   ├── styles.go         # Lipgloss theme-aware styles
│   │   ├── header.go         # Top bar: title, VM status, CPU/mem bars
│   │   ├── stats.go          # Panel 1: aggregate metrics
│   │   ├── activity.go       # Panel 2: GitHub-style heatmap
│   │   ├── projects.go       # Panel 3: sortable project table
│   │   ├── sessions.go       # Panel 4: recent sessions list
│   │   ├── detail.go         # Session detail modal
│   │   └── help.go           # Help modal
│   └── util/time.go          # Relative time formatting
├── go.mod
└── Makefile
```

## UI Layout

```
┌─────────────────────────────────────────────────────────────────────────┐
│ Header: lazyvibe │ VM: Running (PID) │ CPU [████░░] │ MEM [██░░░░]    │
├───────────────────────────────────┬─────────────────────────────────────┤
│ Stats (35% width, 35% height)     │ Projects (65% width, 35% height)   │
│ Panel 1 - Read-only metrics       │ Panel 3 - Sortable table           │
├───────────────────────────────────┼─────────────────────────────────────┤
│ Activity (35% width, 65% height)  │ Sessions (65% width, 65% height)   │
│ Panel 2 - Heatmap with m cycling  │ Panel 4 - Scrollable list          │
├───────────────────────────────────┴─────────────────────────────────────┤
│ Footer: Context-aware keybindings                        │ Theme name   │
└─────────────────────────────────────────────────────────────────────────┘
```

## Navigation Model (lazygit-style)

### Panel Grid

```
Stats(1) ←→ Projects(3)
   ↕           ↕
Activity(2) ←→ Sessions(4)
```

- `1-4`: Direct panel jump
- `h/l`: Horizontal movement
- `Tab`: Circular next

### List Navigation (Projects, Sessions)

- `j/k`: Single item
- `u/i`: Page (5 items)
- `/`: Filter mode
- `s/S`: Sort field/direction

## Key Bindings

| Key | Global | Lists | Activity |
|-----|--------|-------|----------|
| `q` | Quit/Close | | |
| `r` | Refresh | | |
| `p` | Pause | | |
| `t` | Time range | | |
| `T` | Theme | | |
| `?` | Help | | |
| `j/k` | | Scroll | |
| `u/i` | | Page | |
| `s/S` | | Sort | |
| `/` | | Filter | |
| `y` | | Copy ID | |
| `Enter` | | Detail | |
| `m` | | | Metric |

## Dev Workflow

Use `/dev` to run an iterative capture-analyze-fix loop.

### Capture Modes

```bash
./lazyvibe --dump          # Raw JSON data
./lazyvibe --capture 120x40  # ASCII at specific size
```

### Make Targets

```bash
make capture         # 120x40 (standard)
make capture-small   # 80x24 (minimum)
make capture-large   # 160x50 (large)
make dump            # JSON output
```

### Standard Sizes

| Size | Use Case |
|------|----------|
| 80x24 | Minimum supported |
| 120x40 | Standard dev |
| 160x50 | Large terminal |

## Data Sources

| Location | Contents | Refresh |
|----------|----------|---------|
| `~/.claude/projects/*/sessions-index.json` | Session metadata | 10s |
| `~/.claude/stats-cache.json` | Daily stats | 10s |
| VM process (`pgrep`) | CPU/memory | 2s |

## Themes

Cycle with `T`: One Dark → Dracula → Nord → Gruvbox → Catppuccin

## Cross-Platform Build

```bash
make build-all   # darwin-arm64, darwin-amd64, linux-amd64
```
