package backup

import (
    "kvwmap-backup/config"
    "kvwmap-backup/fileops"
    "kvwmap-backup/logging"
    "path/filepath"
    "fmt"
    "time"
)

const (
    networkDir = "/home/gisadmin/networks"
)

func StartBackup(configFile string, loglevel string) {

    logging.Printf("Starte Backup mit Datei %s \n", configFile)

    backup, newBackup := config.GetConfig(configFile)
    if newBackup {
        logging.Fatal("Backup-Config nicht gefunden!")
    }

    //BackupPath und Folder prüfen/anlegen
    if !fileops.IsDir(backup.BackupPath) {
        logging.Fatal("BackupPath ist kein Verzeichnis!")
    }

    backupPath := filepath.Join(backup.BackupPath, backup.BackupFolder)
    fileops.Mkdir(backupPath)

    if len(loglevel) == 0 {
        loglevel = "error,warning,info"
    }
    logfilename := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()) + ".log"
    log := logging.InitLog(filepath.Join(backupPath, logfilename), loglevel)

    logging.Printf("weitere Logausgabe erfolgt in %s", filepath.Join(backupPath, logfilename) )
    log.Info("")
    log.Info( fmt.Sprintf("Sicherungspfad: %s", backupPath)  )

    log.Debug("Netzwerke sichern")
    for _ , network := range backup.Networks {
        fmt.Println(network.Name)
    }

    log.Debug("Services sichern")
    for _ , service := range backup.Services {
        fmt.Println(service.Name)
    }

    log.Debug("Mounts sichern")
    for _ , mount := range backup.Mounts {
        fmt.Println(mount.Service)
    }

    log.Debug("mysql")
    for _ , mysql := range backup.Mysqls {
        fmt.Println(mysql.Services)
    }

    log.Debug("pg_dump")
    for _ , pgdump := range backup.PgDumps {
        fmt.Println(pgdump.Services)
    }

    log.Debug("pg_dumpall")
    for _ , pgdumpall := range backup.PgDumpalls {
        fmt.Println(pgdumpall.Services)
    }
}
