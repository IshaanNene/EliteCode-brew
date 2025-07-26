#!/usr/bin/env bash
set -euo pipefail

OWNER="IshaanNene"
REPO="AlgoRank"
BRANCH="main"
FOLDER_PATH="${1:-.}"         # e.g., Problems1
EXTENSION_FILTER="${2:-}"     # e.g., cpp
TARGET_DIR="Problem"

PROBLEM_ID="problem${FOLDER_PATH//[!0-9]/}"
RAW_JSON="https://raw.githubusercontent.com/$OWNER/$REPO/$BRANCH/$FOLDER_PATH/${PROBLEM_ID}_testcases.json"
echo "Downloading testcases → $TARGET_DIR/testcases.json"
curl -sSf "$RAW_JSON" -o "$TARGET_DIR/testcases.json" || {
    echo "Testcases not found at $RAW_JSON"
    exit 1
}

RAW_DOCKER="https://raw.githubusercontent.com/$OWNER/$REPO/$BRANCH/dockerfiles/${EXTENSION_FILTER}.Dockerfile"
echo "Downloading Dockerfile → $TARGET_DIR/Dockerfile"
curl -sSf "$RAW_DOCKER" -o "$TARGET_DIR/Dockerfile" || {
    echo "Dockerfile not found at $RAW_DOCKER"
    exit 1
}
