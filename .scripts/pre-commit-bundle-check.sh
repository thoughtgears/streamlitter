#!/bin/bash

VERSION_FILE=".version"
BIN_DIR="bin"

if [ ! -f "$VERSION_FILE" ]; then
  echo "Version file not found!"
  exit 1
fi

EXPECTED_VERSION=$(cat "$VERSION_FILE")


for bin in "$BIN_DIR"/*; do
  # Get the basename of the file (without the directory)
  BASENAME=$(basename "$bin")

  # Check if the file is suffixed with the expected version and git SHA
  if [[ "$BASENAME" != *"$EXPECTED_VERSION" ]]; then
    echo "Version mismatch in $bin: expected suffix '$EXPECTED_VERSION'"
    exit 1
  fi
done

echo "All binaries have the correct version and SHA."
exit 0
