#!/usr/bin/env bash
set -euo pipefail

while read -r volume; do
    # Log nur auf STDERR
    echo "Speichere Docker-Volume $volume" >&2

    docker run --rm \
        -v "$volume":/input:ro \
        debian:stable \
        bash -c '
            set -e
            tar -C /input -cf - .
        ' \
        | borg create --stdin-name "dockervolume.tar" "::dockervolume-$volume.{now}" -
done
