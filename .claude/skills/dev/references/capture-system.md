# Capture System Reference

The lazyvibe provides two capture modes for development and debugging. Understanding their capabilities and limitations is crucial for effective development.

## Capture Modes Overview

| Mode | Command | Output | Best For |
|------|---------|--------|----------|
| `--dump` | `python -m lazyvibe --dump` | JSON | Data validation, debugging data flow |
| `--capture WxH` | `python -m lazyvibe --capture 120x40` | ASCII text | Layout analysis, CI testing |

## Standard Test Sizes

| Size | Use Case |
|------|----------|
| 80x24 | Minimum supported terminal |
| 120x40 | Standard development size |
| 160x50 | Large terminal |

## Implementation Details

### Text Capture (`--capture`)

The text capture uses a custom `capture.py` module that bypasses Textual's test mode limitations:

```python
# Located at: src/lazyvibe/capture.py
```

**How it works:**

1. Creates the app with `app.run_test(size=(width, height))`
2. Calls `app._update_all_widgets()` to populate data
3. Renders Static widgets directly via Rich Console (bypassing Textual's cache)
4. Composites widget content onto a text grid
5. Overlays structural elements (borders) from compositor

**Why this approach is needed:**

Textual's test mode has a cache issue where `Static._render_cache` gets set with `Size(width=0, height=0)` before layout runs. The cache never invalidates because test mode doesn't trigger normal render cycles.

### Rendering Pipeline

```
widget.render() → Rich Console → export_text() → grid placement
                                                        ↓
compositor.render_full_update() → borders/structure → grid overlay
```

## Known Limitations

### 1. RecentSessions (ListView) Shows Empty

**Status:** Known issue, partially addressed

**Cause:** `SessionListItem` widgets have their `compose()` method, but in test mode, the children aren't fully mounted. The `capture.py` module has special handling for this:

```python
# Special handling for RecentSessions ListView
for sessions_widget in app.query(RecentSessions):
    items = list(sessions_widget.query(SessionListItem))[:10]
    for item in items:
        for static_child in item.query(Static):
            lines = render_widget_content(static_child, content_width)
            # ... render to grid
```

**Current behavior:** Sessions may appear empty or partial in text capture but display correctly in live app.

**Workaround:** Run the app interactively with `make run` or use `--dump` to verify session data.

### 2. Box-Drawing Characters

Box-drawing characters (┌┐└┘│─├┤┬┴┼) are preserved from the compositor output to maintain panel borders.

### 3. Sparkline/Plots

Sparkline renders as ASCII characters in capture mode.

## Quick Visibility Commands

### See current state at standard size
```bash
python -m lazyvibe --capture 120x40
```

### See raw data
```bash
python -m lazyvibe --dump | head -80
```

### Full capture suite (using Makefile)
```bash
make capture-all
```

### Individual captures
```bash
make capture          # 120x40
make capture-small    # 80x24
make capture-large    # 160x50
```

## Debugging Capture Issues

### Widget not appearing in capture?

1. Check if widget is a Static subclass (should work)
2. Check if widget uses ListView/ListItem (known limitation)
3. Verify widget has `can_focus = True` and an ID
4. Run `--dump` to confirm data is present

### Content truncated?

1. Check terminal size - try larger dimensions
2. Check widget width constraints in `monitor.tcss`
3. Verify `overflow` settings

### Borders broken?

1. Box-drawing chars require UTF-8 terminal
2. Check pipe redirection isn't mangling output

## Files

| File | Purpose |
|------|---------|
| `src/lazyvibe/capture.py` | Custom capture implementation |
| `src/lazyvibe/__main__.py` | CLI entry point with capture modes |
| `Makefile` | Convenient capture targets |
| `.claude/skills/dev/scripts/capture-all.sh` | Capture at multiple sizes |
