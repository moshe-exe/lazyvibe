# Textual TUI Patterns Reference

## Layout System

### Container Types

```python
from textual.containers import Horizontal, Vertical, Grid, ScrollableContainer

# Horizontal: children side by side
with Horizontal():
    yield LeftPanel()
    yield RightPanel()

# Vertical: children stacked
with Vertical():
    yield Header()
    yield Content()
    yield Footer()

# Grid: CSS grid layout
with Grid():
    yield Cell1()
    yield Cell2()
```

### Sizing in TCSS

```css
/* Fixed size */
#panel { width: 40; height: 10; }

/* Percentage */
#panel { width: 50%; }

/* Fraction (flexible) */
#left { width: 1fr; }
#right { width: 2fr; }  /* Takes 2x space */

/* Auto (content-based) */
#header { height: auto; }

/* Min/Max constraints */
#panel { min-width: 30; max-width: 80; }
```

## Widget Patterns

### Basic Widget Structure

```python
from textual.widget import Widget
from textual.reactive import reactive

class MyWidget(Widget):
    """Widget with reactive data."""

    data = reactive([])  # Auto-triggers refresh on change

    def compose(self) -> ComposeResult:
        """Define child widgets."""
        yield Static("Content", id="content")

    def watch_data(self, new_data: list) -> None:
        """Called when data changes."""
        self.query_one("#content").update(str(new_data))

    def update_data(self, new_data: list) -> None:
        """Public method to update widget."""
        self.data = new_data
```

### DataTable Widget

```python
from textual.widgets import DataTable

class ProjectsTable(Widget):
    def compose(self) -> ComposeResult:
        table = DataTable()
        table.add_columns("Name", "Count", "Last Active")
        table.cursor_type = "row"
        yield table

    def update_data(self, projects: list[ProjectSummary]) -> None:
        table = self.query_one(DataTable)
        table.clear()
        for p in projects:
            table.add_row(p.name, str(p.count), p.last_active)
```

### Sparkline Widget

```python
from textual_plotext import PlotextPlot

class ActivitySparkline(Widget):
    def compose(self) -> ComposeResult:
        yield PlotextPlot()

    def update_data(self, activity: list[DailyActivity]) -> None:
        plot = self.query_one(PlotextPlot)
        plot.plt.clear_data()
        values = [a.message_count for a in activity]
        plot.plt.bar(range(len(values)), values)
        plot.refresh()
```

## Styling Patterns

### Border Styles

```css
/* Box border types */
#panel { border: solid blue; }
#panel { border: double green; }
#panel { border: round cyan; }
#panel { border: heavy white; }

/* Border titles */
#panel {
    border: solid blue;
    border-title-align: center;
}
```

### Color Palette (One Dark Theme)

```css
/* Background colors */
$bg-dark: #282c34;
$bg-medium: #21252b;
$bg-light: #2c313c;

/* Text colors */
$text-primary: #abb2bf;
$text-muted: #5c6370;

/* Accent colors */
$blue: #61afef;
$green: #98c379;
$yellow: #e5c07b;
$red: #e06c75;
$purple: #c678dd;
$cyan: #56b6c2;
```

### Responsive Layouts

```css
/* Hide element on small screens */
@media (width < 100) {
    #sidebar { display: none; }
}

/* Adjust proportions */
@media (width >= 120) {
    #left { width: 1fr; }
    #right { width: 2fr; }
}

@media (width < 120) {
    #left { width: 1fr; }
    #right { width: 1fr; }
}
```

## Common Issues & Solutions

### Issue: Content Truncated

```css
/* Enable scrolling */
#content {
    overflow-y: auto;
    overflow-x: hidden;
}
```

### Issue: Widget Not Updating

```python
# Force refresh after data change
self.refresh()

# Or use reactive properties
data = reactive([])  # Auto-refreshes
```

### Issue: Layout Not Responding to Size

```css
/* Use fr units instead of fixed */
#panel { width: 1fr; }  /* Good */
#panel { width: 50; }   /* Fixed - won't adapt */
```

### Issue: Borders Overlapping

```css
/* Add margin between panels */
#panel {
    margin: 1;
}

/* Or use container padding */
#container {
    padding: 1;
}
```

## Debugging

### Print Debug Info

```python
def on_mount(self) -> None:
    self.log(f"Widget mounted: {self.size}")

def on_resize(self) -> None:
    self.log(f"New size: {self.size}")
```

### Inspect Widget Tree

```python
# In app
def action_debug(self) -> None:
    self.log(self.tree)
```

### Check Computed Styles

```python
widget = self.query_one("#panel")
self.log(f"Width: {widget.styles.width}")
self.log(f"Height: {widget.styles.height}")
```
