#!/bin/bash
# Compare two captures side by side
# Usage: ./compare.sh <before-dir> <after-dir> [size]

set -e

BEFORE="${1:?Usage: compare.sh <before-dir> <after-dir> [size]}"
AFTER="${2:?Usage: compare.sh <before-dir> <after-dir> [size]}"
SIZE="${3:-120x40}"

BEFORE_FILE="$BEFORE/capture-$SIZE.txt"
AFTER_FILE="$AFTER/capture-$SIZE.txt"

if [[ ! -f "$BEFORE_FILE" ]]; then
    echo "Error: $BEFORE_FILE not found"
    exit 1
fi

if [[ ! -f "$AFTER_FILE" ]]; then
    echo "Error: $AFTER_FILE not found"
    exit 1
fi

echo "Comparing captures at $SIZE"
echo "Before: $BEFORE_FILE"
echo "After:  $AFTER_FILE"
echo ""
echo "=== DIFF ==="
diff --side-by-side --suppress-common-lines "$BEFORE_FILE" "$AFTER_FILE" || true
echo ""
echo "=== STATS ==="
echo "Before: $(wc -l < "$BEFORE_FILE") lines, $(wc -c < "$BEFORE_FILE") bytes"
echo "After:  $(wc -l < "$AFTER_FILE") lines, $(wc -c < "$AFTER_FILE") bytes"
