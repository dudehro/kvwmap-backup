#!/usr/bin/env python3
"""Borg-Wachstumsanalyse: vergleicht zwei Archive und aggregiert die geaenderten
Pfade nach Verzeichnis-Praefix und Typ (added / modified / removed) inkl. Bytes.

Ziel: erkennen, welche Teilbaeume das taegliche Wachstum eines Borg-Repos treiben.

Datenquelle ist `borg diff --json-lines` (borg >= 1.2), das pro Pfad die
Aenderungen maschinenlesbar liefert.

WICHTIGER VORBEHALT
-------------------
`borg diff` gibt *logische* Dateigroessen bzw. Content-Deltas aus, NICHT die
tatsaechlich dedupliziert/komprimiert im Repo belegten Chunk-Bytes. Die
ausgewiesenen Bytes sind ein guter Proxy dafuer, WELCHE Pfade churnen, aber
keine exakte Messung des realen Plattenzuwachses.

Nur Python-Standardbibliothek.
"""

import argparse
import csv
import json
import os
import subprocess
import sys
from collections import defaultdict

DEFAULT_CONFIG = "/etc/backup/borg.conf"

# Exit-Codes
EX_OK = 0
EX_USAGE = 1   # Nutzungs-/Umgebungsfehler
EX_BORG = 2    # borg-Aufruf fehlgeschlagen


def eprint(*args):
    print(*args, file=sys.stderr)


# --------------------------------------------------------------------------- #
# Repo / Archive bestimmen
# --------------------------------------------------------------------------- #

def repopath_from_config(config_path):
    """Sourcet die bash-Config und liefert die Variable `repopath` (Konvention
    wie borg-prune.sh). Gibt None zurueck, wenn Config fehlt oder leer."""
    if not os.path.isfile(config_path):
        return None
    try:
        out = subprocess.run(
            ["bash", "-c", 'source "$1"; printf %s "$repopath"', "_", config_path],
            capture_output=True, text=True, check=True,
        )
    except (OSError, subprocess.CalledProcessError) as exc:
        eprint(f"Warnung: Config {config_path} nicht sourcebar: {exc}")
        return None
    value = out.stdout.strip()
    return value or None


def resolve_repo(args):
    """Repo aus --repo, sonst BORG_REPO, sonst Config `repopath`."""
    if args.repo:
        return args.repo
    env_repo = os.environ.get("BORG_REPO")
    if env_repo:
        return env_repo
    repo = repopath_from_config(args.config)
    if repo:
        return repo
    eprint(
        "Fehler: kein Repo gefunden. --repo angeben, BORG_REPO setzen oder "
        f"`repopath` in {args.config} definieren."
    )
    sys.exit(EX_USAGE)


def latest_two_archives(repo):
    """Ermittelt die zwei juengsten Archive (nach Zeit) via `borg list --json`.
    Liefert (archiv_alt, archiv_neu)."""
    try:
        out = subprocess.run(
            ["borg", "list", "--json", repo],
            capture_output=True, text=True, check=True,
        )
    except FileNotFoundError:
        eprint("Fehler: `borg` nicht gefunden (PATH?).")
        sys.exit(EX_USAGE)
    except subprocess.CalledProcessError as exc:
        eprint(f"Fehler: `borg list` fehlgeschlagen (rc={exc.returncode}):")
        eprint(exc.stderr.strip())
        sys.exit(EX_BORG)

    try:
        archives = json.loads(out.stdout).get("archives", [])
    except json.JSONDecodeError as exc:
        eprint(f"Fehler: `borg list --json` lieferte ungueltiges JSON: {exc}")
        sys.exit(EX_BORG)

    if len(archives) < 2:
        eprint(
            f"Fehler: Repo enthaelt nur {len(archives)} Archiv(e); fuer den "
            "Auto-Vergleich sind zwei noetig. Archive explizit angeben."
        )
        sys.exit(EX_USAGE)

    archives.sort(key=lambda a: a.get("time", ""))
    return archives[-2], archives[-1]


# --------------------------------------------------------------------------- #
# borg diff einlesen
# --------------------------------------------------------------------------- #

def iter_diff_from_borg(repo, archive_old, archive_new):
    """Startet `borg diff --json-lines` und streamt die stdout-Zeilen."""
    cmd = ["borg", "diff", "--json-lines", f"{repo}::{archive_old}", archive_new]
    try:
        proc = subprocess.Popen(
            cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
        )
    except FileNotFoundError:
        eprint("Fehler: `borg` nicht gefunden (PATH?).")
        sys.exit(EX_USAGE)

    for line in proc.stdout:
        line = line.strip()
        if line:
            yield line

    proc.stdout.close()
    rc = proc.wait()
    stderr = proc.stderr.read()
    proc.stderr.close()
    # borg diff: rc 0 = ok, 1 = Warnungen (nicht fatal), >=2 = Fehler.
    if rc >= 2:
        eprint(f"Fehler: `borg diff` fehlgeschlagen (rc={rc}):")
        eprint(stderr.strip())
        sys.exit(EX_BORG)
    if rc == 1 and stderr.strip():
        eprint(f"Warnung von borg diff (rc=1):\n{stderr.strip()}")


def iter_diff_from_file(path):
    """Liest --json-lines aus Datei ('-' = stdin) fuer Tests/Offline-Nutzung."""
    stream = sys.stdin if path == "-" else open(path, encoding="utf-8")
    try:
        for line in stream:
            line = line.strip()
            if line:
                yield line
    finally:
        if stream is not sys.stdin:
            stream.close()


# --------------------------------------------------------------------------- #
# Aggregation
# --------------------------------------------------------------------------- #

class Group:
    """Aggregierte Kennzahlen fuer einen Verzeichnis-Praefix."""

    __slots__ = (
        "added", "modified", "removed",
        "added_bytes", "mod_added", "mod_removed", "removed_bytes",
    )

    def __init__(self):
        self.added = 0
        self.modified = 0
        self.removed = 0
        self.added_bytes = 0
        self.mod_added = 0
        self.mod_removed = 0
        self.removed_bytes = 0

    @property
    def growth(self):
        # Was macht das Repo groesser: neue Dateien + hinzugefuegte Bytes.
        return self.added_bytes + self.mod_added

    @property
    def churn(self):
        return self.added_bytes + self.mod_added + self.mod_removed + self.removed_bytes

    @property
    def net(self):
        return self.growth - (self.mod_removed + self.removed_bytes)

    def as_row(self):
        return {
            "files_added": self.added,
            "files_modified": self.modified,
            "files_removed": self.removed,
            "added_bytes": self.added_bytes,
            "mod_added_bytes": self.mod_added,
            "mod_removed_bytes": self.mod_removed,
            "removed_bytes": self.removed_bytes,
            "growth_bytes": self.growth,
            "churn_bytes": self.churn,
            "net_bytes": self.net,
        }


def group_key(path, depth):
    """Kuerzt einen Pfad auf die ersten `depth` Verzeichniskomponenten."""
    parts = path.strip("/").split("/")
    if len(parts) <= depth:
        # Datei liegt flacher als depth -> Elternverzeichnis als Gruppe.
        key = "/".join(parts[:-1]) if len(parts) > 1 else parts[0]
    else:
        key = "/".join(parts[:depth])
    return key or "/"


def classify_change(change):
    """Wertet einen einzelnen change-Eintrag aus.
    Liefert (typ, added_bytes, removed_bytes) mit typ in
    {added, removed, modified, meta}."""
    ctype = change.get("type", "")
    if ctype == "added":
        return "added", int(change.get("size", 0)), 0
    if ctype == "removed":
        return "removed", 0, int(change.get("size", 0))
    if ctype == "modified":
        return "modified", int(change.get("added", 0)), int(change.get("removed", 0))
    # Alles andere (mode, owner, mtime, ctime, link, directory, ...) ist eine
    # Metadaten-Aenderung ohne nennenswertes Byte-Volumen.
    return "meta", 0, 0


def aggregate(lines, depth):
    """Verarbeitet die JSON-Lines und liefert (groups, totals, path_count)."""
    groups = defaultdict(Group)
    totals = Group()
    path_count = 0

    for line in lines:
        try:
            record = json.loads(line)
        except json.JSONDecodeError as exc:
            eprint(f"Warnung: ungueltige JSON-Zeile uebersprungen: {exc}")
            continue

        path = record.get("path")
        changes = record.get("changes")
        if path is None or not isinstance(changes, list):
            eprint(f"Warnung: unerwartetes Datensatzformat uebersprungen: {line[:120]}")
            continue

        path_count += 1
        grp = groups[group_key(path, depth)]

        # Pro Pfad die Aenderungen zusammenfassen: der "staerkste" Typ bestimmt,
        # ob der Pfad als added/modified/removed gezaehlt wird; Bytes summieren.
        a_bytes = r_bytes = 0
        kinds = set()
        for change in changes:
            if not isinstance(change, dict):
                continue
            kind, ab, rb = classify_change(change)
            kinds.add(kind)
            a_bytes += ab
            r_bytes += rb

        if "added" in kinds:
            category = "added"
        elif "removed" in kinds:
            category = "removed"
        else:  # modified oder nur meta
            category = "modified"

        if category == "added":
            grp.added += 1
            grp.added_bytes += a_bytes
            totals.added += 1
            totals.added_bytes += a_bytes
        elif category == "removed":
            grp.removed += 1
            grp.removed_bytes += r_bytes
            totals.removed += 1
            totals.removed_bytes += r_bytes
        else:
            grp.modified += 1
            grp.mod_added += a_bytes
            grp.mod_removed += r_bytes
            totals.modified += 1
            totals.mod_added += a_bytes
            totals.mod_removed += r_bytes

    return groups, totals, path_count


# --------------------------------------------------------------------------- #
# Ausgabe
# --------------------------------------------------------------------------- #

SORT_KEYS = {
    "growth": lambda g: g.growth,
    "added": lambda g: g.added_bytes,
    "removed": lambda g: g.removed_bytes,
    "churn": lambda g: g.churn,
}


def human(num):
    """Bytes menschenlesbar (SI-Basis 1024)."""
    sign = "-" if num < 0 else ""
    num = abs(num)
    for unit in ("B", "K", "M", "G", "T", "P"):
        if num < 1024 or unit == "P":
            if unit == "B":
                return f"{sign}{num}{unit}"
            return f"{sign}{num:.1f}{unit}"
        num /= 1024.0
    return f"{sign}{num:.1f}P"


def sorted_groups(groups, sort_key):
    keyfn = SORT_KEYS[sort_key]
    return sorted(groups.items(), key=lambda kv: keyfn(kv[1]), reverse=True)


def print_table(groups, totals, meta, sort_key, top):
    print(f"Borg-Wachstumsanalyse")
    print(f"  Repo   : {meta['repo']}")
    print(f"  Alt    : {meta['archive_old']}  ({meta.get('time_old', '?')})")
    print(f"  Neu    : {meta['archive_new']}  ({meta.get('time_new', '?')})")
    print(f"  Tiefe  : {meta['depth']}   Sortierung: {sort_key}   "
          f"geaenderte Pfade: {meta['path_count']}")
    print()

    header = ("Gruppe", "+Dat", "~Dat", "-Dat",
              "+Bytes", "~D+", "~D-", "-Bytes", "Growth")
    rows = sorted_groups(groups, sort_key)
    shown = rows if top <= 0 else rows[:top]

    # Spaltenbreiten fuer die Zahlenspalten dynamisch, Gruppe links.
    name_w = max([len(header[0])] + [len(k) for k, _ in shown] + [len("SUMME")])

    def fmt_row(name, g):
        return (
            f"{name:<{name_w}}  "
            f"{g.added:>5}  {g.modified:>5}  {g.removed:>5}  "
            f"{human(g.added_bytes):>9}  {human(g.mod_added):>9}  "
            f"{human(g.mod_removed):>9}  {human(g.removed_bytes):>9}  "
            f"{human(g.growth):>9}"
        )

    head = (
        f"{header[0]:<{name_w}}  "
        f"{header[1]:>5}  {header[2]:>5}  {header[3]:>5}  "
        f"{header[4]:>9}  {header[5]:>9}  {header[6]:>9}  "
        f"{header[7]:>9}  {header[8]:>9}"
    )
    print(head)
    print("-" * len(head))
    for name, g in shown:
        print(fmt_row(name, g))
    print("-" * len(head))
    print(fmt_row("SUMME", totals))

    if top > 0 and len(rows) > top:
        print(f"\n... {len(rows) - top} weitere Gruppen (siehe --top 0 oder --json/--csv).")

    print()
    print("Hinweis: Byte-Werte sind logische Dateigroessen/-deltas aus `borg diff`,")
    print("nicht der reale dedup/komprimierte Plattenzuwachs im Repo.")


def export_json(groups, totals, meta, sort_key, dest):
    payload = {
        "meta": meta,
        "totals": totals.as_row(),
        "groups": [
            {"group": name, **g.as_row()}
            for name, g in sorted_groups(groups, sort_key)
        ],
    }
    text = json.dumps(payload, indent=2, ensure_ascii=False)
    if dest == "-":
        print(text)
    else:
        with open(dest, "w", encoding="utf-8") as fh:
            fh.write(text + "\n")
        eprint(f"JSON geschrieben: {dest}")


def export_csv(groups, sort_key, dest):
    fieldnames = ["group"] + list(Group().as_row().keys())
    rows = [{"group": name, **g.as_row()} for name, g in sorted_groups(groups, sort_key)]
    if dest == "-":
        writer = csv.DictWriter(sys.stdout, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(rows)
    else:
        with open(dest, "w", encoding="utf-8", newline="") as fh:
            writer = csv.DictWriter(fh, fieldnames=fieldnames)
            writer.writeheader()
            writer.writerows(rows)
        eprint(f"CSV geschrieben: {dest}")


# --------------------------------------------------------------------------- #
# main
# --------------------------------------------------------------------------- #

def parse_args(argv):
    parser = argparse.ArgumentParser(
        description="Vergleicht zwei Borg-Archive und aggregiert die Aenderungen "
                    "nach Verzeichnis-Praefix und Typ (added/modified/removed).",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("archive_old", nargs="?",
                        help="Aelteres Archiv (weglassen fuer Auto: zweitjuengstes).")
    parser.add_argument("archive_new", nargs="?",
                        help="Neueres Archiv (weglassen fuer Auto: juengstes).")
    parser.add_argument("--repo", help="Borg-Repo (Default: BORG_REPO oder Config).")
    parser.add_argument("--config", default=DEFAULT_CONFIG,
                        help=f"bash-sourcebare Config mit `repopath` (Default: {DEFAULT_CONFIG}).")
    parser.add_argument("--depth", type=int, default=2,
                        help="Aggregations-Verzeichnistiefe (Default: 2).")
    parser.add_argument("--top", type=int, default=20,
                        help="Nur die N groessten Gruppen anzeigen (0 = alle; Default: 20).")
    parser.add_argument("--sort", choices=sorted(SORT_KEYS), default="growth",
                        help="Sortierschluessel (Default: growth).")
    parser.add_argument("--json", metavar="FILE",
                        help="Vollstaendige Aggregation als JSON exportieren ('-' = stdout).")
    parser.add_argument("--csv", metavar="FILE",
                        help="Vollstaendige Aggregation als CSV exportieren ('-' = stdout).")
    parser.add_argument("--diff-file", metavar="FILE",
                        help="borg-diff --json-lines aus Datei/stdin ('-') lesen "
                             "statt borg aufzurufen (Test/Offline).")
    parser.add_argument("--no-table", action="store_true",
                        help="Text-Tabelle unterdruecken (nur Export).")
    args = parser.parse_args(argv)

    if args.depth < 1:
        parser.error("--depth muss >= 1 sein.")
    if (args.archive_old is None) != (args.archive_new is None):
        parser.error("Beide Archive angeben oder beide weglassen (Auto-Modus).")
    return args


def main(argv=None):
    args = parse_args(argv if argv is not None else sys.argv[1:])

    meta = {"depth": args.depth}

    if args.diff_file:
        # Offline-/Testmodus: keine borg-Aufrufe.
        meta.update({
            "repo": args.repo or "(diff-file)",
            "archive_old": args.archive_old or "(alt)",
            "archive_new": args.archive_new or "(neu)",
        })
        lines = iter_diff_from_file(args.diff_file)
    else:
        repo = resolve_repo(args)
        if args.archive_old and args.archive_new:
            archive_old, archive_new = args.archive_old, args.archive_new
            time_old = time_new = "?"
        else:
            old, new = latest_two_archives(repo)
            archive_old, time_old = old["name"], old.get("time", "?")
            archive_new, time_new = new["name"], new.get("time", "?")
        meta.update({
            "repo": repo,
            "archive_old": archive_old,
            "archive_new": archive_new,
            "time_old": time_old,
            "time_new": time_new,
        })
        lines = iter_diff_from_borg(repo, archive_old, archive_new)

    groups, totals, path_count = aggregate(lines, args.depth)
    meta["path_count"] = path_count

    if not args.no_table:
        print_table(groups, totals, meta, args.sort, args.top)
    if args.json:
        export_json(groups, totals, meta, args.sort, args.json)
    if args.csv:
        export_csv(groups, args.sort, args.csv)

    return EX_OK


if __name__ == "__main__":
    sys.exit(main())
