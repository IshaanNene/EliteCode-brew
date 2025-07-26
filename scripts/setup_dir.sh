#!/usr/bin/env bash
set -euo pipefail

OWNER="IshaanNene"
REPO="AlgoRank"
BRANCH="main"
FOLDER_PATH="${1:-.}"         # e.g., Problems1
EXTENSION_FILTER="${2:-}"     # e.g., cpp
TARGET_DIR="Problem"

mkdir -p "$TARGET_DIR"

RAW_CODE="https://raw.githubusercontent.com/$OWNER/$REPO/$BRANCH/$FOLDER_PATH/starter_code.$EXTENSION_FILTER"
echo "Downloading starter code → $TARGET_DIR/starter_code.$EXTENSION_FILTER"
curl -sSf "$RAW_CODE" -o "$TARGET_DIR/starter_code.$EXTENSION_FILTER" || {
    echo "Starter code not found at $RAW_CODE"
    exit 1
}