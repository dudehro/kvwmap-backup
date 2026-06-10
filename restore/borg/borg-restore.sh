#!/usr/bin/env bash
#
# borg-restore.sh — letztes Archiv eines Borg-Repos per SSH auf diesen Server holen
#
# Läuft auf dem ZIELSERVER (Pull-Prinzip: borg extract schreibt nur lokal).
#
# Modi:
#   ./borg-restore.sh list              # Archive im Repo anzeigen
#   ./borg-restore.sh contents          # Inhalt des letzten Archivs listen
#   ./borg-restore.sh dry-run           # Extraktion simulieren
#   ./borg-restore.sh extract           # letztes Archiv komplett nach $TARGET
#   ./borg-restore.sh extract etc/ var/lib/postgresql/
#                                       # nur bestimmte Pfade (ohne führenden /)
# Setup:
# 1. separaten SSH-Key erstellen: ssh-keygen -t ed25519 -f /root/.ssh/borg_restore -N ""
# 2. erstellten public-key auf REPO_HOST einrichten
# 3. Script-Konfiguration anpassen: REPO_PORT, REPO_PATH
# 4. wenn Repo verschlüssel ist, PASSFILE anlegen und im Script eintragen
#

set -euo pipefail

# ── Konfiguration ────────────────────────────────────────────────
REPO_HOST="storage.gdi-service.de"
REPO_PORT="22"
REPO_USER="borgbackup"
REPO_PATH="/storage/borgrepos/<hostname>"
REPO="ssh://${REPO_USER}@${REPO_HOST}:${REPO_PORT}${REPO_PATH}"

TARGET="/restore"                          # Staging-Verzeichnis, NICHT /
SSH_KEY="/root/.ssh/borg_restore"
PASSFILE=""          # chmod 600, eine Zeile
LOGFILE="/var/log/borg-restore-$(date +%Y%m%d-%H%M%S).log"

export BORG_RSH="ssh -i ${SSH_KEY} -o BatchMode=yes -o ConnectTimeout=10"
[[ -r "$PASSFILE" ]] && export BORG_PASSCOMMAND="cat ${PASSFILE}"
# Repo-Fingerprint ist auf dem neuen Server unbekannt; bewusst akzeptieren:
export BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK=yes
export BORG_RELOCATED_REPO_ACCESS_IS_OK=yes

# ── Logging mit Zeitstempel ──────────────────────────────────────
log() { echo "[$(date +%Y-%m-%d\ %H:%M:%S)] $*"; }

MODE="${1:-}"; shift || true

# ── Preflight ────────────────────────────────────────────────────
command -v borg >/dev/null || { log "FEHLER: borg nicht installiert (apt install borgbackup)"; exit 1; }

# ── Letztes Archiv ermitteln (dient zugleich als Erreichbarkeits- und Passphrase-Check, ohne teures borg info) ──
LATEST="$(borg list --last 1 --format '{archive}' "$REPO")" \
  || { log "FEHLER: Repo nicht erreichbar oder Passphrase falsch"; exit 1; }
[[ -n "$LATEST" ]] || { log "FEHLER: kein Archiv im Repo gefunden"; exit 1; }
log "==> Letztes Archiv: ${LATEST}"


case "$MODE" in
  list)
    borg list "$REPO"
    ;;
  contents)
    borg list "${REPO}::${LATEST}" | less
    ;;
  dry-run)
    mkdir -p "$TARGET"; cd "$TARGET"
    log "==> Simuliere Extraktion (dry-run)"
    borg extract --dry-run --list "${REPO}::${LATEST}" "$@"
    ;;
  extract)
    mkdir -p "$TARGET"
    # Sicherung gegen versehentliches Extrahieren nach /
    [[ "$(realpath "$TARGET")" == "/" ]] && { log "FEHLER: TARGET ist / — nutz das Staging-Verzeichnis."; exit 1; }
    cd "$TARGET"
    log "==> Extrahiere nach ${TARGET} (Log: ${LOGFILE})"
    # --numeric-ids: UIDs/GIDs numerisch erhalten, kein Name-Mapping
    #   (borg < 1.2: stattdessen --numeric-owner)
    borg extract --numeric-ids --list --progress \
      "${REPO}::${LATEST}" "$@" 2>&1 | tee "$LOGFILE"
    log "==> Inhalt liegt unter ${TARGET}/ (Pfade ohne führenden /)."
    ;;
  *)
    echo "Usage: $0 {list|contents|dry-run|extract} [pfad ...]"
    exit 1
    ;;
esac

log "==> Script fertig (Modus: ${MODE})."
