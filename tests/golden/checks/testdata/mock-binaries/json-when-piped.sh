#!/usr/bin/env bash
# Mock CLI that emits JSON when stdout is not a TTY.
if [[ "$1" == "list" ]]; then
  if [[ -t 1 ]]; then
    echo "TTY output (table)"
  else
    echo '{"items":[]}'
  fi
  exit 0
fi
[[ "$1" == "--help" ]] && { echo usage; exit 0; }
exit 2
