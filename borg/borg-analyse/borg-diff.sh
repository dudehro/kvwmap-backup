#!/bin/sh
# borg-diff.sh
#
# borg diff --json-lines zwischen zwei Archiven.
#
# Emuliert --content-only (erst ab borg 1.4 verfuegbar, hier borg 1.2.3):
# Reine Metadaten-Aenderungen (mode, owner, mtime/ctime, ...) werden
# herausgefiltert, es bleiben nur inhaltliche Aenderungen uebrig.
# Sortierung nach geaenderter Datenmenge, gross nach klein.
#
# Nutzung: borg-diff.sh <archiv1> <archiv2> [PFAD ...]
# BORG_REPO muss in der Umgebung gesetzt sein.

set -eu

if [ -z "${BORG_REPO:-}" ]; then
    echo "BORG_REPO nicht gesetzt. Abbruch!" >&2
    exit 1
fi

if [ "$#" -lt 2 ]; then
    echo "Nutzung: $0 <archiv1> <archiv2> [PFAD ...]" >&2
    exit 1
fi

ARCHIVE1="$1"
ARCHIVE2="$2"
shift 2

# Nur bei interaktivem Terminal paginieren; per Cron/Pipe unveraendert durchreichen.
pager="cat"
if [ -t 1 ]; then
    pager="${PAGER:-less -R}"
fi

borg diff --json-lines "::${ARCHIVE1}" "${ARCHIVE2}" "$@" \
    | jq -sr '
        # Bytes in passende Groessenordnung umrechnen (1.2 MiB, 512 KiB, ...).
        def human:
            if . == 0 then "0 B"
            else
                (log / (1024 | log) | floor) as $e
                | (. / pow(1024; $e)) as $n
                | "\($n * 10 | round / 10) \(["B","KiB","MiB","GiB","TiB","PiB"][$e])"
            end;

        # Aenderungstypen, die als inhaltlich gelten (--content-only-Emulation).
        ["added", "removed", "modified",
         "added link", "removed link", "changed link",
         "added directory", "removed directory"] as $content

        | map(
            # nur inhaltliche Aenderungen je Pfad behalten
            (.changes |= map(select(.type as $t | $content | index($t))))
            # Pfade ohne verbleibende inhaltliche Aenderung verwerfen
            | select(.changes | length > 0)
            # geaenderte Byte-Menge als Sortierschluessel anhaengen
            | . + {_bytes: (.changes
                | map((.size // 0) + (.added // 0) + (.removed // 0))
                | add)}
          )
        # gross nach klein
        | sort_by(-._bytes)
        # je Datei eine Zeile: type, size, pfad (tab-getrennt fuer column)
        | .[]
        | [ (.changes | map(.type) | join(",")), (._bytes | human), .path ]
        | @tsv
    ' \
    | column -t -s "$(printf '\t')" \
    | $pager
