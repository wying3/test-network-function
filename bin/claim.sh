#!/bin/bash

set -e
set -o pipefail

[ $# -lt 1 ] && printf 'usage: %s filename\n\n' "$(basename "$0")" >&2 && exit 2

SCHEMA="$(dirname "$0")"/../share/claim.json
FILENAME="$1"

python3 -mrsk_wax.jsonschema.tool "$SCHEMA" "$FILENAME"
