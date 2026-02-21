# lazyvibe

A lazygit-inspired terminal dashboard for monitoring Claude Code sessions and activity.

[![Go 1.21+](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

```
┌─────────────────────────────────┬──────────────────────────────────────────────┐
│ Stats 1 [All Time]              │ Projects 3 [Activity ↓]              j/k u/i │
│                                 │                                              │
│   Projects  20                  │   Project            Sessions  Msgs    Active│
│   Sessions  674                 │ ▶ lazyvibe               2    21    12h ago│
│   Messages  14,372              │   mini-moshe              94  1782    20h ago│
│   Tools     42,914              │   electron-app            35   960     1d ago│
│   Tokens    ~26.5M              │   raymosh                  3    27     5d ago│
├─────────────────────────────────┼──────────────────────────────────────────────┤
│ Activity 2 [Messages]           │ Sessions 4 [Time ↓]                   j/k u/i│
│      Mo Tu We Th Fr Sa Su       │                                              │
│ Feb  ▓▓ ▓▓ ░░ ▓▓ ▓▓ ░░ ░░       │   12h: Migrating Python TUI to Go [main]    │
│      ▓▓ ▓▓ ▓▓ ▓▓ ░░ ░░ ░░       │     lazyvibe | 2 msgs | 1m                │
│ Jan  ░░ ▓▓ ▓▓ ▓▓ ▓▓ ░░ ░░       │   20h: Fixed Makefile venv [main]           │
│      ▓▓ ▓▓ ▓▓ ▓▓ ▓▓ ░░ ░░       │     lazyvibe | 19 msgs | 7m               │
│      Less ░░ ▒▒ ▓▓ ██ More      │   1d: CLI dev workflow capture               │
│                                 │     mini-moshe | 35 msgs | 3h                │
└─────────────────────────────────┴──────────────────────────────────────────────┘
 q quit  r refresh  t time  1-4 panels  h/l navigate  j/k scroll  ? help
```

## Why lazyvibe?

**Think btop meets lazygit for Claude Code.** Monitor your AI coding sessions with the same keyboard-driven workflow you love from lazygit:

- **Panel navigation** with `1-4` keys and `h/j/k/l`
- **Vim-style scrolling** through lists
- **Sort and filter** with `s` and `/`
- **Modal details** with `Enter`
- **Context-aware footer** showing available actions

## Features

| Feature | Description |
|---------|-------------|
| **Session Tracking** | Recent sessions with summaries, message counts, durations, git branches |
| **Project Overview** | All projects ranked by activity with session/message stats |
| **Activity Heatmap** | GitHub-style contribution graph (messages, sessions, tools, tokens) |
| **VM Monitoring** | Claude Desktop VM status with CPU/memory bars |
| **Time Filtering** | Filter by Today, This Week, This Month, or All Time |
| **Theming** | 5 themes: One Dark, Dracula, Nord, Gruvbox, Catppuccin |
| **Single Binary** | No dependencies, instant startup |

## Installation

### Download Binary (macOS)

```bash
# Apple Silicon
curl -L https://github.com/moshe-exe/lazyvibe/releases/latest/download/lazyvibe-darwin-arm64 -o lazyvibe
chmod +x lazyvibe && sudo mv lazyvibe /usr/local/bin/

# Intel
curl -L https://github.com/moshe-exe/lazyvibe/releases/latest/download/lazyvibe-darwin-amd64 -o lazyvibe
chmod +x lazyvibe && sudo mv lazyvibe /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/moshe-exe/lazyvibe.git
cd lazyvibe
make build
./lazyvibe
```

## Key Bindings

### Navigation (lazygit-style)

| Key | Action |
|-----|--------|
| `1` `2` `3` `4` | Jump to panel (Stats, Activity, Projects, Sessions) |
| `h` / `l` | Move focus left/right |
| `Tab` | Next panel |
| `j` / `k` | Scroll up/down in lists |
| `u` / `i` | Page up/down (5 items) |

### Actions

| Key | Action |
|-----|--------|
| `Enter` | Open session detail modal |
| `y` | Copy session ID to clipboard |
| `s` / `S` | Cycle sort field / toggle direction |
| `/` | Filter current list |
| `Esc` | Clear filter / close modal |

### Global

| Key | Action |
|-----|--------|
| `q` | Quit |
| `r` | Force refresh |
| `p` | Pause/resume auto-refresh |
| `t` | Cycle time range |
| `T` | Cycle theme |
| `m` | Cycle heatmap metric (Activity panel) |
| `?` | Toggle help |

## Panels

```
┌─────────────┬─────────────┐
│ 1 Stats     │ 3 Projects  │  Stats: Aggregate metrics
├─────────────┼─────────────┤  Activity: Heatmap visualization
│ 2 Activity  │ 4 Sessions  │  Projects: Sortable project table
└─────────────┴─────────────┘  Sessions: Recent session list
```

## Data Sources

lazyvibe reads from Claude Code's local data (macOS):

| Path | Data |
|------|------|
| `~/.claude/projects/*/sessions-index.json` | Session metadata |
| `~/.claude/stats-cache.json` | Daily activity stats |
| VM process | Claude Desktop CPU/memory |

## CLI Options

```bash
lazyvibe              # Run dashboard
lazyvibe --dump       # Dump raw JSON data
lazyvibe --capture 120x40  # Capture ASCII at terminal size
```

## Development

```bash
make run            # Build and run
make capture        # Capture ASCII snapshot (120x40)
make dump           # View raw JSON data
make build-all      # Cross-compile all platforms
```

See [CLAUDE.md](CLAUDE.md) for development workflow.

## Requirements

- **macOS** (Linux planned)
- **Claude Code** installed with session history
- **Go 1.21+** (build from source only)

## License

[MIT](LICENSE)
