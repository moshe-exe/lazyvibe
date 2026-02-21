# Dev Workflow Command

Run an iterative development loop on the lazyvibe TUI.

## Instructions

When this command is invoked, perform the following development workflow:

### 1. Capture Current State

Use the Makefile for quick visibility:

```bash
# Standard capture (120x40) - ALWAYS run this first
make capture

# If data issues suspected
make dump | head -80
```

Or capture at multiple sizes:
```bash
make capture-small   # 80x24
make capture         # 120x40
make capture-large   # 160x50
```

### 2. Analyze the Captured Output

Review the text output for:

**Layout Issues:**
- Borders rendering correctly (box-drawing characters ┌┐└┘│─)
- Columns properly aligned
- Proper spacing between panels
- Content not truncated unexpectedly

**Data Display:**
- All expected fields visible (StatsPanel shows totals, ActivitySparkline shows trend)
- Numbers formatted correctly
- Timestamps showing relative time
- Empty states handled

**Known Limitations to Ignore:**
- RecentSessions (ListView) may appear empty in capture - this is a test mode limitation
- Sparkline bars may not render fully in text mode

### 3. If Data Issues Suspected

Check raw JSON data:
```bash
make dump | head -80
```

Compare against what's displayed to identify:
- Missing fields not being rendered
- Incorrect data transformations
- Parsing errors

### 4. Make Targeted Improvements

Based on analysis, edit the appropriate file:

| Issue Type | File to Edit |
|------------|--------------|
| Layout/styling | `src/lazyvibe/styles/monitor.tcss` |
| StatsPanel | `src/lazyvibe/widgets/stats_panel.py` |
| ActivitySparkline | `src/lazyvibe/widgets/activity.py` |
| ProjectsTable | `src/lazyvibe/widgets/projects.py` |
| RecentSessions | `src/lazyvibe/widgets/sessions.py` |
| Header | `src/lazyvibe/widgets/header.py` |
| Data fetching | `src/lazyvibe/data/*.py` |
| App structure | `src/lazyvibe/app.py` |
| Capture system | `src/lazyvibe/capture.py` |

### 5. Verify Changes

After making changes, re-capture:
```bash
make capture
```

Then test at minimum size to ensure no regression:
```bash
make capture-small
```

### 6. Report

Summarize:
- What was found in the capture analysis
- What changes were made
- What was verified

## User Arguments

$ARGUMENTS

If the user provides arguments like "fix the header" or "check layout at 80x24", focus on that specific aspect.

## Reference Files

- `references/capture-system.md` - How the capture system works, known limitations
- `references/textual-patterns.md` - Textual TUI patterns and common fixes
