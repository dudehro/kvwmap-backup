#!/usr/bin/env bash
set -euo pipefail

while read -r volume; do
    echo "Speichere Docker-Volume $volume nach Borg"

    docker run --rm \
        -v "$volume":/input:ro \
        debian:stable \
        bash -c '
            set -e
            tar -C /input -cf - .
        ' \
        | borg create --stdin-name "dockervolume.tar" "::dockervolume-$volume.{now}" -
done
