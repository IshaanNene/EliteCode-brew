#!/bin/bash

OWNER="IshaanNene"
REPO="AlgoRank"
BRANCH="main"
FOLDER_PATH="Solutions/Problem1"
TARGET_DIR="Problem1"

mkdir -p "$TARGET_DIR"

echo "Fetching file list for $FOLDER_PATH ..."
FILES=$(curl -s "https://api.github.com/repos/$OWNER/$REPO/contents/$FOLDER_PATH?ref=$BRANCH" | grep '"name"' | cut -d '"' -f 4)

for FILE in $FILES; do
    RAW_URL="https://raw.githubusercontent.com/$OWNER/$REPO/$BRANCH/$FOLDER_PATH/$FILE"
    echo "Downloading $FILE ..."
    curl -s "$RAW_URL" -o "$TARGET_DIR/$FILE"
done

echo "✅ Download complete. Files saved in $TARGET_DIR/"
