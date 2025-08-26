#!/bin/bash

#!/usr/bin/env bash
set -euo pipefail

set -a
. ./.env
set +a

archive="${1:?Usage: $0 <archive-name>}"

path="${PGBACKRESTREPO#/}"
IFS='/' read -r -a parts <<< "$path"
strip=$(( ${#parts[@]} - 1 ))

borg extract --strip-components "$strip" "${BORGREPO}::${archive}" "$path"
