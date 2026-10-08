#!/usr/bin/env bash
# borg-growth-report.sh
#
# Wrapper um borg-growth-diff.py: durchsucht ein Basisverzeichnis nach Borg-Repos
# und wertet je Repo die Archive einer Serie (Default: files.*) aus.
#
# Verglichen werden die letzten N Archive paarweise aufeinanderfolgend
# (N-1 Vergleiche). Die Summe dieser Deltas entspricht dem realen Repo-Wachstum;
# ein einzelner Diff aeltestes<->neuestes wuerde rotierende Dateien (Logs, Dumps)
# unterschaetzen, weil deren taeglich neue Chunks wegen der Retention liegen bleiben.
#
# Nutzung: borg-growth-report.sh [OPTIONEN] <basisverzeichnis>
#
# Nicht zugaengliche Repos (verschluesselt ohne Passphrase, fehlende Rechte,
# verschoben) werden uebersprungen und im Report mit Grund gemeldet -- das ist
# weder Fehler noch Warnung.
#
# Exit-Codes:
#   0  Vergleiche gelaufen (uebersprungene Repos aendern daran nichts)
#   1  Nutzungs-/Umgebungsfehler
#   2  Repos gefunden, aber kein einziger Vergleich moeglich
#   3  ein Vergleich in einem zugaenglichen Repo ist fehlgeschlagen

set -euo pipefail

SCRIPTDIR="$(cd "$(dirname "$0")" && pwd)"
GROWTH_DIFF="$SCRIPTDIR/borg-growth-diff.py"

# Cron/nicht-interaktiv: nicht auf eine Bestaetigung fuer unbekannte
# unverschluesselte Repos warten.
export BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK="${BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK:-yes}"

GLOB="files.*"
LAST=7
DEPTH=2
TOP=20
SORT="growth"
OUTDIR=""
MAXDEPTH=3
SUMMARY_ONLY=0

TAB=$'\t'

usage() {
    cat <<'EOF'
Nutzung: borg-growth-report.sh [OPTIONEN] <basisverzeichnis>

  --glob GLOB     Archiv-Glob je Repo (Default: 'files.*')
  --last N        Nur die letzten N Archive je Repo (Default: 7, min. 2)
  --depth N       Aggregations-Verzeichnistiefe (Default: 2)
  --top N         Nur die N groessten Gruppen je Vergleich (0 = alle, Default: 20)
  --sort KEY      growth | added | removed | churn (Default: growth)
  --outdir DIR    JSON-Export je Vergleich dort ablegen
  --max-depth N   Suchtiefe der Repo-Suche (Default: 3)
  --summary-only  Einzeltabellen unterdruecken, nur die Schlussuebersicht
  -h, --help      Diese Hilfe

Der Report geht nach stdout, Fortschritt und Meldungen nach stderr.

Umgebung: BORG_LOCK_WAIT setzen, falls parallel gesichert wird.
          BORG_PASSPHRASE fuer verschluesselte Repos (sonst uebersprungen).
EOF
}

fail() {
    echo "$*" >&2
    exit 1
}

need_value() {
    # $1 = Optionsname, $2 = Anzahl verbleibender Argumente
    [ "$2" -ge 2 ] || fail "Option $1 erwartet einen Wert."
}

is_uint() {
    case "$1" in
        ''|*[!0-9]*) return 1 ;;
        *) return 0 ;;
    esac
}

BASEDIR=""
while [ "$#" -gt 0 ]; do
    case "$1" in
        --glob)        need_value "$1" "$#"; GLOB="$2"; shift 2 ;;
        --last)        need_value "$1" "$#"; LAST="$2"; shift 2 ;;
        --depth)       need_value "$1" "$#"; DEPTH="$2"; shift 2 ;;
        --top)         need_value "$1" "$#"; TOP="$2"; shift 2 ;;
        --sort)        need_value "$1" "$#"; SORT="$2"; shift 2 ;;
        --outdir)      need_value "$1" "$#"; OUTDIR="$2"; shift 2 ;;
        --max-depth)   need_value "$1" "$#"; MAXDEPTH="$2"; shift 2 ;;
        --summary-only) SUMMARY_ONLY=1; shift ;;
        -h|--help)     usage; exit 0 ;;
        --)            shift; break ;;
        -*)            fail "Unbekannte Option: $1 (--help fuer Hilfe)" ;;
        *)
            [ -z "$BASEDIR" ] || fail "Nur ein Basisverzeichnis erlaubt (bereits: $BASEDIR)."
            BASEDIR="$1"; shift ;;
    esac
done
# Reste nach '--'
for arg in "$@"; do
    [ -z "$BASEDIR" ] || fail "Nur ein Basisverzeichnis erlaubt (bereits: $BASEDIR)."
    BASEDIR="$arg"
done

[ -n "$BASEDIR" ] || { usage >&2; exit 1; }
BASEDIR="${BASEDIR%/}"
[ -d "$BASEDIR" ] || fail "Basisverzeichnis existiert nicht: $BASEDIR"

is_uint "$LAST" && [ "$LAST" -ge 2 ] || fail "--last muss eine Zahl >= 2 sein."
is_uint "$DEPTH" && [ "$DEPTH" -ge 1 ] || fail "--depth muss eine Zahl >= 1 sein."
is_uint "$TOP" || fail "--top muss eine Zahl >= 0 sein."
is_uint "$MAXDEPTH" && [ "$MAXDEPTH" -ge 1 ] || fail "--max-depth muss eine Zahl >= 1 sein."
case "$SORT" in
    growth|added|removed|churn) ;;
    *) fail "--sort muss growth, added, removed oder churn sein." ;;
esac

[ -x "$GROWTH_DIFF" ] || fail "borg-growth-diff.py nicht gefunden/ausfuehrbar: $GROWTH_DIFF"
command -v borg >/dev/null 2>&1 || fail "borg nicht gefunden (PATH?)."
command -v jq   >/dev/null 2>&1 || fail "jq nicht gefunden (PATH?)."

if [ -n "$OUTDIR" ]; then
    mkdir -p "$OUTDIR" || fail "Kann --outdir nicht anlegen: $OUTDIR"
fi

TMPDIR_RUN="$(mktemp -d)"
trap 'rm -rf "$TMPDIR_RUN"' EXIT

SUMMARY="$TMPDIR_RUN/summary.tsv"
: > "$SUMMARY"
SKIPPED="$TMPDIR_RUN/skipped.tsv"
: > "$SKIPPED"

skipped=0
errors=0
compares=0

# Grund fuer einen uebersprungenen Repo-Zugriff aus dem borg-stderr ableiten.
# Fallback: erste nicht-leere Zeile der borg-Ausgabe.
skip_reason() {
    local errfile="$1"
    if grep -qE 'can not acquire a passphrase|passphrase supplied' "$errfile" 2>/dev/null; then
        echo "verschluesselt, Passphrase nicht verfuegbar"
    elif grep -qE 'Repository access aborted|was previously located at' "$errfile" 2>/dev/null; then
        echo "Repo verschoben, Zugriff nicht bestaetigt (BORG_RELOCATED_REPO_ACCESS_IS_OK)"
    # Bewusst eng gefasst: ein blosses 'lock' wuerde auch auf Repo-Pfade
    # matchen, die borg in der sys.argv-Zeile jeder Fehlermeldung ausgibt.
    elif grep -qE 'Permission denied|acquire the lock|lock\.exclusive|LockTimeout|Failed to create/acquire' "$errfile" 2>/dev/null; then
        echo "Lock oder Dateirechte"
    else
        grep -m1 . "$errfile" 2>/dev/null || echo "borg-Zugriff fehlgeschlagen"
    fi
}

note_skip() {
    # $1 = Repo-Label, $2 = Grund
    printf '%s\t%s\n' "$1" "$2" >> "$SKIPPED"
    skipped=$((skipped + 1))
}

# --------------------------------------------------------------------------- #
# Repos finden
# --------------------------------------------------------------------------- #
# Ein Repo wird am data/-Verzeichnis erkannt, nicht am Inhalt von config: config
# kann einem anderen Nutzer mit Modus 600 gehoeren und waere dann nicht lesbar,
# das Repo bliebe unsichtbar. `-prune -printf '%h\n'` liefert das Elternver-
# zeichnis von data/ -- also das Repo -- ohne je in data/ abzusteigen.
# lock.exclusive wird mitgeprunt: ein stehengebliebener Lock eines anderen
# Nutzers ist nicht lesbar und liesse find sonst mit rc 1 enden.
#
# Bewusst keine Pipeline: mit `set -o pipefail` wuerde ein Teilfehler von find
# (z.B. ein einzelnes unlesbares Verzeichnis) den ganzen Lauf abbrechen.
found_raw="$TMPDIR_RUN/found.txt"
find_err="$TMPDIR_RUN/find.err"
repos_raw="$TMPDIR_RUN/repos.txt"

find "$BASEDIR" -mindepth 1 -maxdepth "$((MAXDEPTH + 1))" \
     \( -type d -name lock.exclusive -prune \) -o \
     \( -type d -name data -prune -printf '%h\n' \) \
     > "$found_raw" 2>"$find_err" || true
sort -u "$found_raw" > "$repos_raw"

if [ -s "$find_err" ]; then
    echo "Hinweis: die Repo-Suche konnte nicht alles lesen:" >&2
    sed 's/^/  /' "$find_err" >&2
fi

repos=()
while IFS= read -r dir; do
    # Existenzpruefung, kein Lesezugriff.
    [ -f "$dir/config" ] || continue
    repos+=("$dir")
done < "$repos_raw"

if [ "${#repos[@]}" -eq 0 ]; then
    fail "Keine Borg-Repos unter $BASEDIR gefunden (Suchtiefe $MAXDEPTH)."
fi

echo "Gefundene Borg-Repos: ${#repos[@]}   Serie: $GLOB   letzte Archive: $LAST"
echo

# --------------------------------------------------------------------------- #
# Je Repo: Archive listen, Paare vergleichen
# --------------------------------------------------------------------------- #
repo_index=0
for repo in "${repos[@]}"; do
    repo_index=$((repo_index + 1))
    # Label = Pfad relativ zum Basisverzeichnis, '/' -> '_'. Bleibt kurz genug fuer
    # die Tabelle, ist aber eindeutig -- sonst kollidieren gleichnamige Repos
    # verschiedener Hosts in der Uebersicht und in den --outdir-Dateinamen.
    reponame="${repo#"$BASEDIR"}"
    reponame="${reponame#/}"
    reponame="${reponame//\//_}"
    [ -n "$reponame" ] || reponame="$(basename "$repo")"
    fortschritt="[$repo_index/${#repos[@]}] $reponame"

    # stdin auf /dev/null: verschluesselte oder verschobene Repos fragen sonst
    # interaktiv nach (Passphrase bzw. [yN]) und blockieren den Lauf.
    listing="$TMPDIR_RUN/list.json"
    listerr="$TMPDIR_RUN/list.err"
    if ! borg list --json "$repo" </dev/null > "$listing" 2>"$listerr"; then
        # Nicht zugaenglich ist kein Fehler und keine Warnung, nur eine Meldung.
        reason="$(skip_reason "$listerr")"
        echo "$fortschritt: uebersprungen ($reason)" >&2
        note_skip "$reponame" "$reason"
        continue
    fi

    # Namen chronologisch; Glob-Filterung client-seitig, weil der borg-Flagname
    # zwischen 1.2 (--glob-archives) und 1.4 (--match-archives) wechselt.
    names=()
    while IFS= read -r name; do
        [ -n "$name" ] || continue
        # shellcheck disable=SC2254  # GLOB soll als Muster wirken
        case "$name" in
            $GLOB) names+=("$name") ;;
        esac
    done < <(jq -r '.archives | sort_by(.time) | .[].name' "$listing")

    count="${#names[@]}"
    if [ "$count" -lt 2 ]; then
        reason="nur $count Archiv(e) passend zu '$GLOB'"
        echo "$fortschritt: uebersprungen ($reason)" >&2
        note_skip "$reponame" "$reason"
        continue
    fi

    # Auf die letzten LAST Archive kuerzen.
    if [ "$count" -gt "$LAST" ]; then
        names=("${names[@]: -$LAST}")
        count="$LAST"
    fi

    repo_growth=0
    repo_churn=0
    repo_pairs=0

    for ((i = 1; i < count; i++)); do
        alt="${names[i-1]}"
        neu="${names[i]}"
        json="$TMPDIR_RUN/pair.json"

        args=(--repo "$repo" --depth "$DEPTH" --top "$TOP" --sort "$SORT" --json "$json")
        if [ "$SUMMARY_ONLY" -eq 1 ]; then
            args+=(--no-table)
        fi

        if [ "$SUMMARY_ONLY" -eq 0 ]; then
            echo "==================================================================="
            echo "Repo: $repo   $alt -> $neu"
            echo "==================================================================="
        fi

        # stderr zwischenpuffern: die Erfolgsmeldung des JSON-Exports ist hier
        # nur Rauschen (die Tempdatei interessiert niemanden), echte Fehler und
        # borg-Warnungen sollen aber durchkommen.
        # stdin auf /dev/null vererbt sich auf das von Python gestartete borg diff.
        differr="$TMPDIR_RUN/diff.err"
        t0=$SECONDS
        if ! "$GROWTH_DIFF" "${args[@]}" "$alt" "$neu" </dev/null 2>"$differr"; then
            grep -v '^JSON geschrieben: ' "$differr" >&2 || true
            # borg list war erfolgreich, das Repo ist also zugaenglich -- ein
            # Scheitern hier ist ein echter Fehler.
            echo "$fortschritt: FEHLER beim Vergleich $alt -> $neu" >&2
            errors=$((errors + 1))
            continue
        fi
        grep -v '^JSON geschrieben: ' "$differr" >&2 || true

        growth="$(jq -r '.totals.growth_bytes' "$json")"
        churn="$(jq -r '.totals.churn_bytes' "$json")"

        # Fortschritt auf stderr: ein Lauf ueber viele Repos dauert Stunden und
        # liefert mit --summary-only sonst bis zum Ende keine Rueckmeldung.
        printf '%s: %s -> %s  %ss  %s\n' "$fortschritt" "$alt" "$neu" \
            "$((SECONDS - t0))" "$(LC_ALL=C numfmt --to=iec-i --suffix=B -- "$growth")" >&2

        printf '%s\t%s\t%s\t%s\t%s\n' "$reponame" "$alt" "$neu" "$growth" "$churn" >> "$SUMMARY"

        repo_growth=$((repo_growth + growth))
        repo_churn=$((repo_churn + churn))
        repo_pairs=$((repo_pairs + 1))
        compares=$((compares + 1))

        if [ -n "$OUTDIR" ]; then
            cp "$json" "$OUTDIR/${reponame}.${alt}__${neu}.json"
        fi
        if [ "$SUMMARY_ONLY" -eq 0 ]; then
            echo
        fi
    done

    if [ "$repo_pairs" -gt 0 ]; then
        # Platzhalter '-' statt leerem Feld: Tab ist IFS-Whitespace, aufeinander-
        # folgende Tabs wuerden beim spaeteren `read` zu einem Trenner kollabieren.
        if [ "$repo_pairs" -eq 1 ]; then wort="Vergleich"; else wort="Vergleiche"; fi
        printf '%s\tSUMME (%s %s)\t-\t%s\t%s\n' \
            "$reponame" "$repo_pairs" "$wort" "$repo_growth" "$repo_churn" >> "$SUMMARY"
    fi
done

# --------------------------------------------------------------------------- #
# Schlussuebersicht
# --------------------------------------------------------------------------- #
if [ "$compares" -gt 0 ]; then
    echo
    echo "Uebersicht (Wachstum je Vergleich)"
    {
        printf 'Repo\tvon\tbis\tWachstum\tChurn\n'
        while IFS="$TAB" read -r r alt neu growth churn; do
            printf '%s\t%s\t%s\t%s\t%s\n' \
                "$r" "$alt" "$neu" \
                "$(LC_ALL=C numfmt --to=iec-i --suffix=B -- "$growth")" \
                "$(LC_ALL=C numfmt --to=iec-i --suffix=B -- "$churn")"
        done < "$SUMMARY"
    } | column -t -s "$TAB"
    echo
    echo "Hinweis: Byte-Werte sind logische Dateigroessen/-deltas aus \`borg diff\`,"
    echo "nicht der reale dedup/komprimierte Plattenzuwachs im Repo."
else
    echo "Keine Vergleiche durchgefuehrt."
fi

# Uebersprungene Repos gehoeren in den Report (stdout), nicht in den Fehlerkanal:
# sie sind weder Fehler noch Warnung, aber man will wissen, was fehlt.
if [ "$skipped" -gt 0 ]; then
    echo
    echo "Uebersprungen ($skipped von ${#repos[@]} Repos)"
    {
        printf 'Repo\tGrund\n'
        cat "$SKIPPED"
    } | column -t -s "$TAB"
fi

if [ "$errors" -gt 0 ]; then
    echo "Abgeschlossen: $compares Vergleiche, $skipped uebersprungen, $errors Fehler." >&2
    exit 3
fi
if [ "$compares" -eq 0 ]; then
    echo "Abgeschlossen: kein einziger Vergleich moeglich ($skipped Repos uebersprungen)." >&2
    exit 2
fi
echo "Abgeschlossen: $compares Vergleiche, $skipped uebersprungen." >&2
exit 0
