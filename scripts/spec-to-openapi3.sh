#!/usr/bin/env bash
# spec-to-openapi3.sh - fetch an OpenAPI 2.0 spec and upgrade it to 3.0.
#
# Usage: spec-to-openapi3.sh <input-url-or-path> <output-path>
#
# Example:
#   scripts/spec-to-openapi3.sh \
#     https://api.apis.guru/v2/specs/azure.com/web-WebApps/2018-11-01/swagger.json \
#     /tmp/azure-functions-openapi3.json
set -euo pipefail

if [ $# -ne 2 ]; then
  echo "Usage: $0 <input-url-or-path> <output-path>" >&2
  exit 2
fi

INPUT="$1"
OUTPUT="$2"

if ! command -v swagger2openapi >/dev/null 2>&1; then
  echo "ERROR: swagger2openapi not installed. Run scripts/install.sh first." >&2
  exit 3
fi

# swagger2openapi accepts URL or local path as positional arg.
echo "Converting $INPUT -> $OUTPUT (OpenAPI 2.0 -> 3.0)..." >&2
swagger2openapi "$INPUT" -o "$OUTPUT"
echo "Wrote $OUTPUT" >&2
