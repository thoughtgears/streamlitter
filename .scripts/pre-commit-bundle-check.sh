#!/bin/bash

VERSION_FILE=".version"
BIN_DIR="bin"

# Check if the version file exists
if [ ! -f "$VERSION_FILE" ]; then
  echo "Version file not found!"
  exit 1
fi

# Read the expected version from the version file
EXPECTED_VERSION=$(cat "$VERSION_FILE")

# Get the current Git short SHA
GIT_SHA=$(git rev-parse --short HEAD)

# Check if git command was successful
if [ $? -ne 0 ]; then
  echo "Failed to get current Git SHA"
  exit 1
fi

# Iterate over each file in the bin directory
for bin in "$BIN_DIR"/*; do
  # Get the basename of the file (without the directory)
  BASENAME=$(basename "$bin")

  # Check if the file is suffixed with the expected version and git SHA
  if [[ "$BASENAME" != *"$EXPECTED_VERSION-$GIT_SHA" ]]; then
    echo "Version mismatch in $bin: expected suffix '$EXPECTED_VERSION-$GIT_SHA'"
    exit 1
  fi
done

echo "All binaries have the correct version and SHA."
exit 0
