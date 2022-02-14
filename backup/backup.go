package backup

import (
	"fmt"
	"kvwmap-backup/config"
	"kvwmap-backup/fileops"
	"kvwmap-backup/logging"
	"path/filepath"
	"time"
)

func StartBackup(configFile string, loglevel string) {
	logging.Printf("Starte Backup mit Datei %s \n", configFile)

	backup, isNewBackup := config.GetConfig(configFile)

	if isNewBackup {
		logging.Fatal("Backup-Config nicht gefunden!")
	}

	//BackupPath und Folder prüfen/anlegen
	isDir := fileops.IsDir(backup.BackupPath)
	if !isDir {
		logging.Fatal("BackupPath ist kein Verzeichnis!")
	}

	backupPath := filepath.Join(backup.BackupPath, backup.BackupFolder)
	fileops.Mkdir(backupPath)

	//Logging
	if len(loglevel) == 0 {
		loglevel = "error,warning,info"
	}
	logfilename := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()) + ".log"
	log := logging.InitLog(filepath.Join(backupPath, logfilename), loglevel)

	logging.Printf("weitere Logausgabe erfolgt in %s", filepath.Join(backupPath, logfilename))
	log.Info(fmt.Sprintf("Sicherungspfad: %s", backupPath))

	//Backup
	log.Debug("Netwerke sichern")
	for _, netzwerk := range backup.Networks {

		netzwerk_tarball := "network." + netzwerk.Name + " .tar"
        log.Debug(netzwerk_tarball)

		fileops.TarFile( filepath.Join(config.GetNetworkPath(netzwerk.Name), "env"), netzwerk_tarball)

		log.Debug("Services sichern")
		for _, service := range backup.Services {
			service_tarball := netzwerk.Name + "." + service.Name + ".tar"

			log.Debug(service_tarball)

			log.Debug("Mounts sichern")
			for _, mount := range backup.Mounts {
				fmt.Println(mount.Service)
			}

			log.Debug("mysql")
			for _, mysql := range backup.Mysqls {
				fmt.Println(mysql.Services)
			}

			log.Debug("pg_dump")
			for _, pgdump := range backup.PgDumps {
				fmt.Println(pgdump.Services)
			}

			log.Debug("pg_dumpall")
			for _, pgdumpall := range backup.PgDumpalls {
				fmt.Println(pgdumpall.Services)
			}
		}
	}
}
