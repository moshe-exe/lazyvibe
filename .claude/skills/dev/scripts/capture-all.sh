#!/bin/bash
# Capture the lazyclaude at multiple terminal sizes
# Usage: ./capture-all.sh [output-dir]
#
# Creates timestamped capture directory with:
#   - ASCII text captures at multiple sizes
#   - JSON data dump

set -e

OUTPUT_DIR="${1:-.dev-captures}"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
CAPTURE_DIR="$OUTPUT_DIR/$TIMESTAMP"

mkdir -p "$CAPTURE_DIR"

echo "Capturing lazyclaude at $(date)"
echo "Output directory: $CAPTURE_DIR"
echo ""

# Activate venv if not already active
if [[ -z "$VIRTUAL_ENV" ]]; then
    if [[ -f ".venv/bin/activate" ]]; then
        source .venv/bin/activate
    fi
fi

# Sizes to capture
SIZES=("80x24" "120x40" "160x50")

# Capture text at each size
for size in "${SIZES[@]}"; do
    echo "Capturing at $size..."
    if python -m lazyclaude --capture "$size" > "$CAPTURE_DIR/capture-$size.txt" 2>&1; then
        echo "  ✓ capture-$size.txt"
    else
        echo "  ✗ Failed to capture at $size"
    fi
done

# Capture JSON data
echo "Capturing JSON data..."
if python -m lazyclaude --dump > "$CAPTURE_DIR/data.json" 2>&1; then
    echo "  ✓ data.json"
else
    echo "  ✗ Failed to dump data"
fi

echo ""
echo "Done! Captures saved to: $CAPTURE_DIR"
echo ""
echo "Files:"
ls -la "$CAPTURE_DIR"

# Create symlink to latest
ln -sf "$TIMESTAMP" "$OUTPUT_DIR/latest"
echo ""
echo "Symlink created: $OUTPUT_DIR/latest -> $TIMESTAMP"
