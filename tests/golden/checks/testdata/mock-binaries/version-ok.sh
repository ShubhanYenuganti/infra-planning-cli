#!/usr/bin/env bash
[[ "$1" == "--version" ]] && { echo "v0.1.0"; exit 0; }
[[ "$1" == "--help" ]] && { echo usage; exit 0; }
exit 2
