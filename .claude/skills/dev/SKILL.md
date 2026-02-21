---
name: dev
description: This skill should be used when the user asks to "run dev loop", "iterate on the UI", "check the dashboard", "analyze the TUI", "fix the layout", or wants to make iterative improvements to the lazyvibe terminal app. Provides a structured workflow for capturing, analyzing, and improving the Textual TUI.
version: 0.2.0
---

# lazyvibe Development Workflow

This skill enables iterative development of the lazyvibe TUI application by capturing its state, analyzing output, and making targeted improvements.

## Quick Visibility Commands

Use these to immediately see the app state:

```bash
# Quick capture at standard size (recommended first step)
make capture

# See raw data (for data issues)
make dump

# All sizes at once
make capture-all
```

Or without Make:

```bash
python -m lazyvibe --capture 120x40   # ASCII snapshot
python -m lazyvibe --dump | head -80  # JSON data
```

## Capture Modes

| Mode | Command | Best For |
|------|---------|----------|
| `--capture WxH` | `make capture` | Layout analysis, widget content |
| `--dump` | `make dump` | Data validation, debugging data flow |

## What's Visible in Capture

| Widget | Capture Status | Notes |
|--------|----------------|-------|
| StatusHeader | ✅ Works | VM status, PID, CPU, Memory |
| StatsPanel | ✅ Works | Sessions, Messages, Tool Calls, Projects, Days |
| ActivitySparkline | ✅ Works | ASCII sparkline + Total/Average/Peak stats |
| ProjectsTable | ✅ Works | Full project list with columns |
| RecentSessions | ⚠️ Empty | Test mode limitation - works in live app |

## Known Capture Limitations

1. **RecentSessions ListView** - Shows empty in text capture due to Textual test mode not calling `compose()` on ListItem children. Works correctly in live app.
2. **Sparkline ASCII art** - Renders as basic characters, not full block elements.

See `references/capture-system.md` for technical details on the capture implementation.

## Dev Loop Process

### Step 1: Capture Current State

Always start by capturing the current state at multiple sizes to understand responsive behavior:

```bash
# Quick: Standard development size
make capture

# Or multiple sizes
make capture-small   # 80x24 (minimum)
make capture         # 120x40 (standard)
make capture-large   # 160x50 (large)
```

### Step 2: Analyze the Output

When reviewing captured output, check for:

1. **Layout Issues**
   - Columns properly aligned
   - Borders rendering correctly (box-drawing characters)
   - Content not truncated unexpectedly
   - Proper spacing between panels

2. **Data Display**
   - All expected data fields visible
   - Numbers formatted correctly (commas, units)
   - Timestamps showing relative time properly
   - Empty states handled gracefully

3. **Responsive Behavior**
   - Layout adapts to terminal size
   - Critical content visible at minimum size (80x24)
   - Additional content revealed at larger sizes

### Step 3: Validate Data

If display issues might be data-related, check the raw JSON:

```bash
python -m lazyvibe --dump | head -100
```

Compare JSON data against what's displayed to identify:
- Missing fields not being rendered
- Incorrect data transformations
- Caching issues

### Step 4: Make Targeted Changes

Based on analysis, edit the appropriate file:

| Issue Type | File to Edit |
|------------|--------------|
| Layout/styling | `src/lazyvibe/styles/monitor.tcss` |
| Widget content | `src/lazyvibe/widgets/*.py` |
| Data fetching | `src/lazyvibe/data/*.py` |
| App structure | `src/lazyvibe/app.py` |

### Step 5: Verify Changes

After making changes, re-capture and compare:

```bash
python -m lazyvibe --capture 120x40
```

## Project Structure

```
src/lazyvibe/
├── __main__.py           # Entry point + dev capture modes
├── app.py                # Main Textual App, widget composition
├── data/
│   ├── manager.py        # DataManager with caching
│   ├── models.py         # Dataclasses (VMStatus, SessionEntry, etc.)
│   ├── sessions.py       # Parse sessions-index.json
│   ├── stats.py          # Parse stats-cache.json
│   └── vm.py             # VM process detection
├── widgets/
│   ├── header.py         # StatusHeader (VM status, pause indicator)
│   ├── stats_panel.py    # StatsPanel (aggregate statistics)
│   ├── activity.py       # ActivitySparkline (30-day trend)
│   ├── projects.py       # ProjectsTable (project list)
│   └── sessions.py       # RecentSessions (session list)
└── styles/
    └── monitor.tcss      # Textual CSS styling
```

## Common Patterns

### Adding a New Widget

1. Create widget file in `widgets/`
2. Implement `compose()` for layout, `update_data()` for refresh
3. Add to `app.py` in the `compose()` method
4. Wire data flow in `_update_all_widgets()`
5. Add styling in `monitor.tcss`

### Fixing Layout at Small Sizes

Check `monitor.tcss` for:
- `min-width` / `min-height` constraints
- `fr` units vs fixed sizes
- `overflow` settings for scrollable content

### Debugging Data Flow

1. Add `--dump` to see raw data
2. Check `DataManager` cache TTLs
3. Verify parsing in `sessions.py` / `stats.py`

## Test Sizes Reference

| Size | Use Case |
|------|----------|
| 80x24 | Minimum supported (classic terminal) |
| 100x30 | Small modern terminal |
| 120x40 | Standard development |
| 160x50 | Large/wide terminal |
| 200x60 | Ultra-wide monitor |

## Iterative Workflow Example

```
1. User: "/dev check the layout"
2. Capture at 120x40, analyze output
3. Identify: "Projects table header misaligned"
4. Read widgets/projects.py
5. Fix alignment in table column definitions
6. Re-capture, verify fix
7. Test at 80x24 to ensure no regression
8. Report changes made
```
