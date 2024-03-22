#!/usr/bin/python3
import subprocess
import logging
import sys
import json
import os
from datetime import datetime

# Aktuelles Datum als YYYY-MM-DD-Format
current_date = datetime.now().strftime('%Y-%m-%d')

# Funktion zum Extrahieren des Dateinamens ohne Erweiterung
def get_file_name_without_extension(file_path):
#    return file_path.split('.')[0]
#    file_name = os.path.basename(file_path)
    normalized_name = file_path.replace('/', '_').replace('.','_')
    filename = os.path.splitext(normalized_name)[0]
    return filename


# Logging-Konfiguration
logfilepath = f'/var/log/cmdlogger/{get_file_name_without_extension(sys.argv[1])}_{current_date}.json'
print(f'logfile: ', logfilepath)
log_handler = logging.FileHandler(filename=logfilepath, encoding='utf-8')
log_handler.setLevel(logging.INFO)
logger = logging.getLogger()
logger.setLevel(logging.INFO)
logger.addHandler(log_handler)

# Eigener JSON-Formatter für das Logging
class JSONFormatter(logging.Formatter):
    def format(self, record):
        log_data = {
            'timestamp': datetime.fromtimestamp(record.created).isoformat(),
            'level': record.levelname,
            'message': record.getMessage(),
            'file': record.filename,
            'line': record.lineno,
            'function': record.funcName,
            'cmd': getattr(record, 'cmd', None),
            'stdout': getattr(record, 'stdout', None),
            'stderr': getattr(record, 'stderr', None),
            'exitcode': getattr(record, 'exitcode', None)
        }
        return json.dumps(log_data)

log_formatter = JSONFormatter()
log_handler.setFormatter(log_formatter)

# Handler für die Konsole hinzufügen
#console_handler = logging.StreamHandler(sys.stdout)
#console_handler.setFormatter(log_formatter)
#logger.addHandler(console_handler)

# Funktion zum Ausführen der Befehle aus einer Datei
def ausfuehren_befehle_aus_datei(datei):
    try:
        with open(datei, 'r') as f:
            for line in f:
                befehl = line.strip()  # Entferne Leerzeichen und Zeilenumbrüche
                #print(befehl)
                try:
                    result = subprocess.run(befehl, shell=True, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
                    logger.info("exec", extra={'cmd': befehl, 'stdout': result.stdout.strip(), 'stderr': result.stderr.strip(), 'exitcode': result.returncode})
                except subprocess.CalledProcessError as e:
                    logger.error("exec", extra={'cmd': befehl, 'stdout': result.stdout.strip(), 'stderr': result.stderr.strip(), 'exitcode': e.returncode})
    except FileNotFoundError:
        logger.error("FileNotFoundError")
    except Exception as e:
        logger.error("Execetpion", extra={'stderr': {e}, 'exitcode': e.returncode })

if __name__ == "__main__":
    if not sys.argv[1]:
        print("Verwendung: python script.py <Dateiname>")
        sys.exit(1)  # Rückgabecode für falsche Verwendung

    ausfuehren_befehle_aus_datei(sys.argv[1])
