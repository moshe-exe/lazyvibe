#!/bin/bash
# Quick visibility check - runs capture at standard size and shows output
# Usage: ./quick-check.sh [size]
#
# This script is optimized for fast feedback during development.
# It captures and immediately displays the output.

set -e

SIZE="${1:-120x40}"

echo "=== lazyclaude @ $SIZE ==="
echo ""

# Activate venv if not already active
if [[ -z "$VIRTUAL_ENV" ]]; then
    if [[ -f ".venv/bin/activate" ]]; then
        source .venv/bin/activate
    fi
fi

# Run capture and display
python -m lazyclaude --capture "$SIZE"

echo ""
echo "=== Capture complete ==="
