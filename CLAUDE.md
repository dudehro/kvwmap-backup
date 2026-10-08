# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Zweck des Repositories

Sammlung von Bash- und Python-Skripten zur Sicherung und Wiederherstellung der kvwmap-GIS-Infrastruktur. kvwmap läuft in Docker-Containern auf Linux-Hosts. Die Skripte werden per Cron über den Job-Runner `cmdlogger.py` aufgerufen (Standard auf praktisch allen Servern). `runjob.py` ist deprecated.

## Verzeichnisstruktur und Verantwortlichkeiten

```
cmdlogger.py             Job-Runner (Standard): Befehlsdatei zeilenweise ausführen → JSON-Log in /var/log/cmdlogger/
runjob.py, job_funcs.py  DEPRECATED Job-Runner: JSON-Jobdefinition → joblog.json im datierten workdir
bash/                    Ältere, monolithische Backup-Lösung (JSON-getrieben, kvwmap-spezifisch)
                         + rm-alte-tagessicherung.sh (gehört zu runjob.py, deprecated)
borg/                    Borg-Backup: Docker-Volumes sichern, Repository bereinigen
borg/borg-analyse/       Analyse von Archiv-Unterschieden und Repo-Wachstum
postgresql/              PostgreSQL-Dumps aus Docker-Containern; Legacy-Restore-Skripte
mariadb/                 MariaDB/MySQL-Dumps aus Docker-Containern (mehrere Varianten)
restore/borg/            Letztes Archiv eines entfernten Borg-Repos per SSH lokal wiederherstellen
restore/postgres/        Nummerierter Restore-Workflow: Borg → pgbackrest → PostgreSQL
monitor/                 Statusauswertung der Job-Logs (für Monitoring-Integration)
custom-scripts/          Host-spezifische Skripte, von git ignoriert, nicht von Ansible erfasst
```

`backup-config/` ist ebenfalls von git ignoriert (lokale Konfigurationen).

## Zwei Backup-Generationen

**Legacy (`bash/backup.sh`):** Alles-in-einem-Skript. JSON-Konfiguration. Eng an kvwmap-Serverstruktur unter `/home/gisadmin/` gekoppelt. Unterstützt tar (mit differenzieller Sicherung), pg\_dump, pg\_dumpall, mysqldump, rsync in einer Pipeline.

**Modern (alle anderen Verzeichnisse):** Einzelne, zweckgebundene Skripte. Die Dump-Skripte sourcen eine einfache Konfigurationsdatei als erstes Argument (`source ${1}`). Orchestrierung über `cmdlogger.py`.

## Job-Runner

### `cmdlogger.py <befehlsdatei>` — Standard
Führt jede Zeile der Befehlsdatei einzeln mit `/bin/bash -o pipefail -c` aus; Zeilen mit `#` am Anfang werden ignoriert. Ein Fehler bricht die übrigen Zeilen **nicht** ab. Log als JSON-Lines nach `/var/log/cmdlogger/<pfad-normalisiert>_<YYYY-MM-DD>.json` (`/` und `.` im Pfad der Befehlsdatei werden zu `_`). Das Verzeichnis `/var/log/cmdlogger/` muss existieren.

### `runjob.py <jobs.json> <jobname>` — DEPRECATED
Nur noch auf einzelnen Altservern im Einsatz. Nicht für neue Jobs verwenden und nicht weiterentwickeln; neue Funktionen gehören in den `cmdlogger.py`-Weg. Ebenfalls nur für diesen Weg relevant: `monitor/latest-status.sh`, `monitor/dsplog.sh`, `bash/rm-alte-tagessicherung.sh`, `postgresql/bootstrap-restore.sh` (alle arbeiten mit `joblog.json` bzw. `WORKDIR`).

Startet den Job `<jobname>` aus der JSON-Definition (relativer Pfad wird relativ zum Skriptverzeichnis aufgelöst) und arbeitet eine Job-Queue ab. Aufbau:
```json
{
  "workdir": "/pfad/zu/backups/$today$",
  "jobs": [
    {"name": "...", "command": ["/pfad/script.sh", "arg"],
     "next-job": "...", "start-job-on-success": "...", "start-job-on-error": "...",
     "exit-queue-on-error": true}
  ]
}
```
- `$today$` wird in `workdir` und `command` durch `YYYY-MM-DD` ersetzt; `workdir` wird angelegt.
- `command` ist eine Argumentliste (kein Shell-String). Jobs bekommen `WORKDIR` als Umgebungsvariable.
- Ergebnis je Job (Start/Ende, Exit-Code, stdout/stderr, args) in `<workdir>/joblog.json`. Achtung: der Startzeit-Schlüssel heißt `startime` (Tippfehler, wird so ausgewertet).
- Exit-Code = Anzahl fehlgeschlagener Jobs (1 auch bei Abbruch durch Exception).
- Jobnamen, die mit `borg` beginnen, werden vom Monitoring gesondert behandelt (siehe unten).

## Konfigurationskonventionen

### Moderne Dump-Skripte (postgresql/, mariadb/)
Alle akzeptieren einen Pfad zu einer bash-sourcebaren Konfigurationsdatei als `$1`. Erwartete Variablen:

- PostgreSQL (`dump-postgresql.sh`): `CONTAINER`, `DBUSER`, `TARGET`, `DATABASES[]`
- MariaDB (`dump-mariadb.sh`, `dump-mysql.sh`, `dump-mariadb-all-databases.sh`): `CONTAINER`, `CREDENTIALS_FILE`, `TARGET`, `DATABASES[]`
- MariaDB via Zwischencontainer (`dump-maria-via.sh`): zusätzlich `VIA_CONTAINER`
- MySQL/MariaDB-Credentials werden aus `credentials.php` gelesen (kvwmap-Konvention): `grep MYSQL_USER ... | cut -d "'" -f 4`

`postgresql/bootstrap-restore.sh` (deprecated, nur mit `runjob.py`): rechnet `WORKDIR` auf den Container-Pfad unter `/dumps/` um und startet `restore_backup_database.sh` im Container `kvwmap_prod_pgsql` (host-spezifisch, LKMSE).

### Legacy (`bash/backup.sh`)
JSON-Konfiguration. Relevante Felder:
```json
{
  "backup_path": "...",
  "backup_folder": "$(date +%F)",
  "delete_after_n_days": 10,
  "delete_diff_on_dow": "1",
  "tar": [{"source": "...", "target_name": "...", "exclude": "..."}],
  "pg_dump": [{"db_name": "...", "container_id": "...", "db_user": "...", "target_name": "...", "docker_network": "..."}],
  "pg_dumpall": [...],
  "mysql_dump": [...],
  "rsync": [{"source": "...", "destination": "...", "parameter": "..."}]
}
```
Sourcet global `/home/gisadmin/kvwmap-server/config/config`.

### Docker-Volumes sichern (`borg/backup-docker-volumes.sh`)
Liest Volume-Namen zeilenweise von **stdin**, tart jedes Volume read-only über einen `debian:stable`-Container und streamt es per `borg create --stdin-name dockervolume.tar` ins Repo. Repo kommt aus `BORG_REPO` (Umgebung). Archivname: `dockervolume-<volume>.{now}`.

### Borg-Prune (`borg/borg-prune.sh`)
Bash-sourcebar, Standard: `/etc/backup/borg.conf`. Variablen: `repopath`, `keeplastdays`, `keepweekly`, `keepmonthly`.

### Borg-Diff (`borg/borg-analyse/borg-diff.sh`)
`borg-diff.sh <archiv1> <archiv2> [PFAD ...]`, `BORG_REPO` muss gesetzt sein. Emuliert `--content-only` und Sortierung nach Größe für borg 1.2 (beides erst ab 1.4): filtert reine Metadaten-Änderungen per `jq` heraus, sortiert nach geänderter Datenmenge absteigend. Paginiert nur bei interaktivem Terminal. POSIX-`sh`, nicht bash.

### Borg-Wachstumsanalyse (`borg/borg-analyse/borg-growth-diff.py`, `borg/borg-analyse/borg-growth-report.sh`)
`borg-growth-diff.py` vergleicht **ein** Archiv-Paar **eines** Repos und aggregiert die Änderungen nach Verzeichnis-Präfix. Repo aus `--repo`, sonst `BORG_REPO`, sonst `repopath` aus der Borg-Config. Nur Python-Standardbibliothek. Die ausgewiesenen Bytes sind logische Größen aus `borg diff`, nicht die deduplizierten/komprimierten Chunk-Bytes — Proxy für *welche* Pfade churnen, keine exakte Messung des Plattenzuwachses.

`borg-growth-report.sh <basisverzeichnis>` ist der Wrapper darüber: sucht Borg-Repos unterhalb des Basisverzeichnisses (Suchtiefe `--max-depth`, Default 3), nimmt je Repo die letzten `--last` Archive der Serie `--glob` (Default `files.*`) und vergleicht alle **aufeinanderfolgenden** Paare. Am Ende eine Übersichtstabelle mit Wachstum/Churn je Vergleich und Summe je Repo; `--outdir` legt zusätzlich den JSON-Export je Vergleich ab. Weitere Optionen: `--depth`, `--top`, `--sort`, `--summary-only`.

Wichtig:
- Die Serien-Filterung ist nötig, weil ein Repo mehrere unabhängige Archiv-Serien enthalten kann (z. B. `dockervolume-<volume>.{now}`, `postgres.<container>.<db>.{now}`). Auf gdi-borgstorageBig haben Repos bis zu 8 Serien und bis zu 8121 Archive. Der Auto-Modus von `borg-growth-diff.py` nimmt sonst die zwei zeitlich jüngsten Archive — quer über Serien hinweg ergibt das Unsinn.
- Aufeinanderfolgende Paare statt ältestes↔neuestes: nur die Summe der Einzeldeltas entspricht dem realen Repo-Wachstum, da rotierende Dateien täglich neue Chunks belegen, die wegen der Retention liegen bleiben.
- Glob-Filterung erfolgt client-seitig im Wrapper, nicht per borg-Flag — der Flagname wechselt zwischen borg 1.2 (`--glob-archives`) und 1.4 (`--match-archives`).
- Repo-Erkennung über das `data/`-Verzeichnis, nicht über den Inhalt von `config`: die `config` kann einem anderen Nutzer mit Modus 600 gehören und wäre dann nicht lesbar. `lock.exclusive` wird beim Suchen ausgespart, weil ein stehengebliebener Lock eines anderen Nutzers `find` sonst mit rc 1 beendet.
- Alle borg-Aufrufe bekommen stdin von `/dev/null`. Verschlüsselte Repos (Passphrase-Abfrage) und verschobene Repos (`[yN]`-Rückfrage) würden den Lauf sonst blockieren.

Nicht zugängliche Repos (Passphrase fehlt, Dateirechte/Lock, verschoben) und Repos ohne passende Archive werden **übersprungen** und im Report mit Grund gelistet — das ist weder Fehler noch Warnung.

Exit-Codes:
- `0` Vergleiche gelaufen (übersprungene Repos ändern daran nichts)
- `1` Nutzungsfehler / kein Repo gefunden
- `2` Repos gefunden, aber kein einziger Vergleich möglich
- `3` ein Vergleich in einem zugänglichen Repo ist fehlgeschlagen

Laufzeit (gemessen auf gdi-borgstorageBig): ein `borg diff` dauert 1 s bis 98 s je nach Repo-Größe. Der Default `--last 7` bedeutet bei 46 auswertbaren Repos rund 276 Diffs, also mehrere Stunden. Außerhalb des Sicherungsfensters (22:00–00:00) laufen lassen und `BORG_LOCK_WAIT` setzen. Fortschritt geht nach stderr, der Report nach stdout.

### Monitor für runjob.py (`monitor/latest-status.sh`, deprecated)
JSON-Konfiguration, Standard: `/etc/backup/jobs.json` (dieselbe Datei wie für `runjob.py`). Das Elternverzeichnis von `workdir` enthält datierte Unterordner (`YYYY-MM-DD/joblog.json`); ausgewertet wird der jüngste der letzten drei Tage.

`bash/rm-alte-tagessicherung.sh [jobs.json]` löscht in demselben Elternverzeichnis alle Unterordner, die älter als 10 Tage sind.

### Restore (`restore/postgres/`)
Benötigt `.env`-Datei im Arbeitsverzeichnis mit: `PGBACKRESTREPO` (Pfad des pgbackrest-Repos), `BORGREPO` (Borg-Repository-Pfad).

## Borg-Restore von entferntem Repo (`restore/borg/borg-restore.sh`)

Läuft auf dem **Zielserver** (Pull-Prinzip) und holt das jeweils **letzte** Archiv per SSH. Konfiguration steht direkt im Skript (`REPO_HOST`, `REPO_PORT`, `REPO_USER`, `REPO_PATH`, `SSH_KEY`, optional `PASSFILE`); Setup-Schritte im Skriptkopf.

```bash
./borg-restore.sh list                     # Archive im Repo
./borg-restore.sh contents [pfad ...]      # Inhalt des letzten Archivs + Symlink-Bericht
./borg-restore.sh dry-run [pfad ...]       # Extraktion simulieren
./borg-restore.sh extract [pfad ...]       # nach $TARGET (Default /restore) + Symlink-Bericht
```
- Pfade ohne führenden `/` angeben. Extrahieren nach `/` wird verweigert.
- Gesamte Ausgabe wird mit Zeitstempel zusätzlich nach `/var/log/borg-restore-<datum>.log` gespiegelt.
- Symlink-Bericht: borg sichert nur den Link, nicht das Ziel. Jeder Symlink in der Auswahl wird gegen das gesamte Archiv geprüft (`ok` / `ausserhalb-auswahl` / `FEHLT`). Benötigt ein vollständiges Archiv-Listing — bei großen Archiven langsam.
- Muss mit borg 1.x funktionieren (`{source}` statt `{linktarget}`, Hinweis zu `--numeric-owner` für borg < 1.2).

## PostgreSQL-Restore-Workflow (nummerierte Skripte)

Schritte in `restore/postgres/` der Reihe nach ausführen:

```bash
# 1. pgbackrest-Repo aus Borg-Archiv extrahieren
./10-restore-from-borg.sh <archiv-name>

# 2. Backup-Integrität prüfen (Docker: pkorduan/postgis:15-3.3)
./11-verify-backups.sh

# 3. Verfügbare Backup-Sets auflisten
./20-list-backups.sh

# 4. Wiederherstellen (bringt dcm-Services hoch/runter)
./30-restore.sh [set-name]

# Notfall: WAL zurücksetzen
./pg_resetwal.sh
```

`30-restore.sh` nutzt `dcm` (Docker Compose Wrapper) und erwartet `./data/` und `./backup/` relativ zum Arbeitsverzeichnis.

## Docker-Dump-Muster

Alle Dump-Skripte folgen diesem Muster:
1. Dump in Containervolume erstellen (`docker exec ... pg_dump/mysqldump ... > /var/lib/postgresql/data/...`)
2. Volume-Pfad auf dem Host ermitteln: `docker inspect --format "{{json .Mounts}}" $CONTAINER | jq -r '.[]|select(.Destination=="...").Source'`
3. Dump-Datei auf den Host verschieben: `mv $DUMPDIR/*.dump $TARGET`

Dump-Dateinamen folgen der Konvention `{CONTAINER}.{DB}.dump`.

## Monitoring-Integration

Für Monitoring-Systeme (z. B. Checkmk), einheitliche Exit-Codes:
- `0`: alles OK
- `2`: Warnung (borg Exit-Code 1 = nicht-fatale Fehler)
- `3`: Fehler (auch: kein aktuelles Log gefunden)

- `monitor/cmdlogger-latest.py <befehlsdatei>` — **Standard**. Wertet das jüngste Log von `cmdlogger.py` aus; Log älter als 3 Tage gilt als Fehler. borg-Sonderbehandlung greift bei Befehlen, die mit `borg` beginnen.
- `monitor/latest-status.sh [jobs.json]` — deprecated, wertet `joblog.json` von `runjob.py` aus. borg-Sonderbehandlung greift bei Jobnamen, die mit `borg` beginnen.
- `monitor/dsplog.sh <joblog.json>` — deprecated, gibt Job-Laufzeiten und Fehler aus `joblog.json` menschenlesbar aus.

## Wichtige Abhängigkeiten

- `jq` — wird in allen Skripten intensiv verwendet
- `docker` — alle Datenbankoperationen laufen in Containern
- `borg` — für Borg-Backups und pgbackrest-Archivierung; Skripte müssen mit borg 1.2 laufen (Unterschiede zu 1.4 beachten)
- `pgbackrest` — für PITR-fähige PostgreSQL-Sicherung (läuft im Docker-Image `pkorduan/postgis:15-3.3`)
- `rsync` — für Dateiübertragungen im Legacy-Skript
- `python3` (≥ 3.7 wegen `capture_output`/f-Strings, nur Standardbibliothek) — Job-Runner, Monitoring, Wachstumsanalyse
