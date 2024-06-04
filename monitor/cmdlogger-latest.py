#!/usr/bin/python3

#return codes
# 0 keine Fehler
# 1 nicht implementiert
# 2 Warnungen
# 3 Fehler

import os
import glob
import sys
import json
from datetime import datetime, timedelta

def get_latest_logfile(befehlsdatei):
    # Extrahiere den Dateinamen ohne Erweiterung
    normalized_name = befehlsdatei.replace('/', '_').replace('.','_')
    befehlsdatei_name = os.path.splitext(normalized_name)[0]
    globString = f'/var/log/cmdlogger/{befehlsdatei_name}_*.json'

    # Suche nach Logdateien im angegebenen Format
    log_files = glob.glob(globString)

    # Sortiere die Logdateien nach Änderungsdatum (aufsteigend)
    log_files.sort(key=os.path.getmtime, reverse=True)

    # Gib den Pfad der neuesten Logdatei zurück (falls vorhanden)
    if log_files:
        return log_files[0]
    else:
        return None

def find_exec_errors(logfile):
    exec_errors = []
    exec_warnings = []

    with open(logfile, 'r') as f:
        for line in f:
            log_entry = json.loads(line.strip())
            if log_entry['message'] == 'exec':
                if log_entry['cmd'].startswith('borg'):
                    if log_entry['exitcode'] == 1:
                        exec_warnings.append(log_entry)
                    elif log_entry['exitcode'] == 2:
                        exec_errors.append(log_entry)
                elif log_entry['exitcode'] > 0:
                    exec_errors.append(log_entry)
    return exec_errors, exec_warnings

if __name__ == "__main__":
    befehlsdatei = sys.argv[1]  # Annahme: Die Befehlsdatei wird als Argument übergeben
    latest_logfile = get_latest_logfile(befehlsdatei)

    if latest_logfile:
        print(f"Aktuellste Logdatei: {latest_logfile}")
        # Überprüfe, ob die Datei älter als 3 Tage ist
        file_creation_time = datetime.fromtimestamp(os.path.getctime(latest_logfile))
        if datetime.now() - file_creation_time > timedelta(days=3):
            print('Feher: Logfile älter als 3 Tage')
            sys.exit(3)
        #Fehler aus Logdatei auswerden
        exec_errors, exec_warnings = find_exec_errors(latest_logfile)
        if exec_errors:
            print('Fehler in Logdatei')
            sys.exit(3)
        elif exec_warnings:
            print('Warnungen in Logdatei')
            sys.exit(2)
        else:
            print('keine Fehler')
            sys.exit(0)
    else:
        print('Fehler: kein Logfile gefunden')
        sys.exit(3)
