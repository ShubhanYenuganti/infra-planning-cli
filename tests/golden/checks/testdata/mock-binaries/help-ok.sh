#!/usr/bin/env bash
# Mock CLI that handles --help correctly and unknown commands with exit 2.
if [[ "$1" == "--help" ]]; then
  echo "Usage: mock [command]"
  exit 0
fi
echo "unknown command" >&2
exit 2
