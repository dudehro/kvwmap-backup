# borg-analyse

Werkzeuge, um Unterschiede zwischen Borg-Archiven und das Wachstum von Borg-Repos auszuwerten. Ausgelegt für borg 1.2.

| Skript | Zweck |
|---|---|
| `borg-diff.sh` | Dateiliste der inhaltlichen Änderungen zwischen zwei Archiven, nach Größe sortiert |
| `borg-growth-diff.py` | Ein Archiv-Paar eines Repos, Änderungen nach Verzeichnis aggregiert |
| `borg-growth-report.sh` | Alle Repos unter einem Basisverzeichnis, je Repo die letzten N Archive einer Serie |

Abhängigkeiten: `borg`, `jq`, `column`, `numfmt`, `python3` ≥ 3.7.

Die Byte-Werte sind logische Dateigrößen aus `borg diff`, **nicht** der deduplizierte und komprimierte Platzbedarf im Repo. Sie zeigen, *welche* Pfade sich ändern, und messen nicht exakt das Plattenwachstum.

---

## borg-diff.sh

```
borg-diff.sh <archiv1> <archiv2> [PFAD ...]
```

| Parameter | Bedeutung |
|---|---|
| `archiv1`, `archiv2` | Archivnamen (ohne `::`), älteres zuerst |
| `PFAD ...` | optional: nur diese Pfade im Archiv vergleichen |
| `BORG_REPO` | **Pflicht**, Repo kommt nur aus der Umgebung |

Reine Metadaten-Änderungen (mode, owner, mtime …) werden herausgefiltert; das entspricht `--content-only` aus borg 1.4. Ausgabe: Änderungstyp, Bytes, Pfad. Am Terminal läuft sie durch `$PAGER` (Default `less -R`), in Pipes unverändert.

```bash
export BORG_REPO=/storage/borgrepos/gdi-server1
borg list --short "$BORG_REPO" | tail -2          # Archivnamen ermitteln
borg-diff.sh files.2026-10-06T22:00:04 files.2026-10-07T22:00:03
borg-diff.sh files.2026-10-06T22:00:04 files.2026-10-07T22:00:03 home/gisadmin/www | head -50
```

---

## borg-growth-diff.py

```
borg-growth-diff.py [OPTIONEN] [archiv_alt archiv_neu]
```

| Option | Default | Bedeutung |
|---|---|---|
| `archiv_alt archiv_neu` | – | beide oder keins angeben. Ohne Archive: die zwei zeitlich jüngsten im Repo (**über alle Serien hinweg**, siehe Hinweise) |
| `--repo REPO` | `BORG_REPO`, sonst `repopath` aus `--config` | Repo-Pfad |
| `--config FILE` | `/etc/backup/borg.conf` | bash-sourcebare Config mit `repopath` |
| `--depth N` | 2 | Verzeichnistiefe für die Aggregation |
| `--top N` | 20 | nur die N größten Gruppen anzeigen, `0` = alle |
| `--sort KEY` | `growth` | `growth`, `churn`, `added`, `removed` (siehe unten) |
| `--json FILE` | – | vollständige Aggregation als JSON (`-` = stdout) |
| `--csv FILE` | – | vollständige Aggregation als CSV (`-` = stdout) |
| `--no-table` | – | keine Tabelle, nur Export |
| `--diff-file FILE` | – | `borg diff --json-lines`-Ausgabe aus Datei/stdin lesen statt borg aufzurufen |

Kennzahlen je Gruppe:
- `growth` = Bytes neuer Dateien + hinzugekommene Bytes geänderter Dateien
- `churn` = `growth` + entfernte Bytes (gelöschte und geänderte Dateien)
- `added` / `removed` = nur Bytes **neuer** bzw. **gelöschter** Dateien (ohne Änderungen an bestehenden Dateien)

Exit-Codes: `0` OK, `1` Nutzungsfehler / kein Repo / zu wenige Archive, `2` borg-Aufruf fehlgeschlagen.

```bash
R=/storage/borgrepos/gdi-server1

# ein bestimmtes Paar, Aggregation auf 3 Ebenen, alle Gruppen
borg-growth-diff.py --repo $R --depth 3 --top 0 files.2026-10-06T22:00:04 files.2026-10-07T22:00:03

# nur Export, keine Tabelle
borg-growth-diff.py --repo $R --no-table --csv growth.csv files.2026-10-06T22:00:04 files.2026-10-07T22:00:03

# Diff einmal holen, dann offline mit verschiedenen Tiefen auswerten
borg diff --json-lines $R::files.2026-10-06T22:00:04 files.2026-10-07T22:00:03 > diff.jsonl
borg-growth-diff.py --diff-file diff.jsonl --depth 4
```

---

## borg-growth-report.sh

```
borg-growth-report.sh [OPTIONEN] <basisverzeichnis>
```

Sucht alle Borg-Repos unter dem Basisverzeichnis, nimmt je Repo die letzten N Archive der Serie und vergleicht jeweils **aufeinanderfolgende** Paare (N−1 Diffs je Repo). Am Ende kommt eine Übersicht mit Wachstum/Churn je Vergleich und der Summe je Repo.

| Option | Default | Bedeutung |
|---|---|---|
| `--glob GLOB` | `files.*` | Archiv-Serie (Shell-Muster auf den Archivnamen) |
| `--last N` | 7 | letzte N Archive je Repo, mindestens 2 |
| `--depth N` | 2 | wie bei `borg-growth-diff.py` |
| `--top N` | 20 | wie bei `borg-growth-diff.py` |
| `--sort KEY` | `growth` | wie bei `borg-growth-diff.py` |
| `--outdir DIR` | – | JSON je Vergleich als `<repo>.<alt>__<neu>.json` ablegen |
| `--max-depth N` | 3 | Suchtiefe der Repo-Suche unter dem Basisverzeichnis |
| `--summary-only` | – | nur die Schlussübersicht, keine Einzeltabellen |

Umgebung: `BORG_LOCK_WAIT` (wenn parallel gesichert wird), `BORG_PASSPHRASE` (für verschlüsselte Repos, sonst werden sie übersprungen).

Nicht zugängliche Repos und Repos mit weniger als 2 passenden Archiven werden übersprungen und am Ende mit Grund gelistet.

Exit-Codes: `0` Vergleiche gelaufen, `1` Nutzungsfehler / kein Repo gefunden, `2` kein einziger Vergleich möglich, `3` mindestens ein Vergleich fehlgeschlagen.

Laufzeit: ein Diff dauert je nach Repo 1–100 s. Mit `--last 7` über ~46 Repos sind das mehrere Stunden. Deshalb außerhalb des Sicherungsfensters (22:00–00:00) laufen lassen und stdout umlenken; der Fortschritt geht nach stderr.

```bash
B=/storage/borgrepos

# Überblick über alle Repos, letzte Woche, nur Schlussübersicht
BORG_LOCK_WAIT=600 borg-growth-report.sh --summary-only $B > growth-report.txt

# andere Serie: Docker-Volumes, nur die letzten 3 Archive
borg-growth-report.sh --glob 'dockervolume-*' --last 3 --summary-only $B

# ausführlich mit JSON-Export je Vergleich
borg-growth-report.sh --depth 3 --top 10 --outdir ./growth-json $B > growth-report.txt
```

---

## Hinweise

- **Archiv-Serien:** Ein Repo kann mehrere unabhängige Serien enthalten (`files.*`, `dockervolume-<volume>.*`, `postgres.<container>.<db>.*`). Der Auto-Modus von `borg-growth-diff.py` (ohne Archivnamen) vergleicht die zwei jüngsten Archive *irgendeiner* Serie. Bei gemischten Repos deshalb Archive explizit angeben oder `borg-growth-report.sh --glob` verwenden.
- **stdin-Archive** (`dockervolume-*`, Dumps) enthalten nur eine einzige Datei. Die Aggregation nach Verzeichnissen sagt dort nichts aus, nur die Summe ist brauchbar.
- **Summe statt ältestes↔neuestes:** Nur die Summe der aufeinanderfolgenden Deltas spiegelt das reale Wachstum. Rotierende Dateien (Logs, Dumps) würden ein einzelnes Diff zwischen ältestem und neuestem Archiv unterschätzen.
