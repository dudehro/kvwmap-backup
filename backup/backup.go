package backup

import (
	"fmt"
	"kvwmap-backup/config"
	"kvwmap-backup/docker"
	"kvwmap-backup/fileops"
	"kvwmap-backup/logging"
	"path/filepath"
	"strings"
	"time"
)

func StartBackup(configFile, dirFile string, loglevel string) {
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

		netzwerk_tarfile := filepath.Join(backupPath, "network."+netzwerk.Name+".tar")
		netzwerk_source_dir := fileops.GetTemplateDir(dirFile, fileops.TemplateNetwork, netzwerk.Name, "", "")

		fileops.TarFile(filepath.Join(netzwerk_source_dir, "docker-compose.yml"), "", netzwerk_tarfile)
		fileops.TarFile(filepath.Join(netzwerk_source_dir, "env"), "", netzwerk_tarfile)
		fileops.TarFile(filepath.Join(netzwerk_source_dir, "env2"), "", netzwerk_tarfile)

		log.Debug("Services sichern")
		for _, service := range backup.Services {
			service_tarfile := filepath.Join(backupPath, netzwerk.Name+"."+strings.TrimPrefix(service.Name, "/")+".tar")
			log.Debug(service_tarfile)

			//docker-compose.yamls
			for _, configfile := range docker.GetContainerConfigFiles(service.Name) {
				fmt.Printf("TarFile(%s, %s)\n", configfile, service_tarfile)
				//				fileops.TarFile(configfile, "", service_tarfile)
			}

			log.Debug("Mounts sichern")
			for _, mount := range backup.Mounts {
				if mount.Service == service.Name {
					fmt.Printf("TarDir(%s, %s)", mount.MountSource, service_tarfile)
					//					fileops.TarDir(mount.MountSource, service_tarbfile)
				}
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
