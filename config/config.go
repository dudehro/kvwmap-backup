package config

import (
//    "fmt"
    "log"
    "encoding/json"
    "io/ioutil"
//    "kvwmap-backup/config"
)

func NewConfig(Path string) *Backup {
    p := Backup{BackupPath: Path}
    return &p
}

func AddNetwork2Backup(backup *Backup, name string, subnet string) (*Network){
    if backup == nil {
        log.Fatal()
    }
    network := Network{Name: name, Subnet: subnet}
    backup.Networks = append(backup.Networks, &network)
    return &network
}

func AddNetwork2Container(container *Service, networkname string) {
    if container == nil {
        log.Fatal()
    }
    container.Networks = append(container.Networks, networkname)
}

func AddContainer2Backup(backup *Backup, image string, name string, network string) (*Service) {
    if backup == nil {
        log.Fatal()
    }
    service := Service{Image: image, Name: name, Networks: []string{network} }
    backup.Services = append(backup.Services, &service)
    return &service
}

func AddMount(backup *Backup, container *Service, source string, destination string) (*Mount) {
    if backup == nil || container == nil {
        log.Fatal()
    }
    mount := Mount{MountSource: source, MountDestination: destination, Services: []string{container.Name}}
    backup.Mounts = append(backup.Mounts, &mount)
    return &mount
}

func AddContainer2Mount(container *Service, mount *Mount) {
    if container == nil || mount == nil {
        log.Fatal()
    }
    mount.Services = append(mount.Services, container.Name)
}

func AddPostgres2Backup(backup *Backup, container *Service, dbName string, dbUser string, dbHost string) (*PgDump) {
    if container == nil {
        log.Fatal()
    }
    pg := PgDump{DbName: dbName, DbUser: dbUser, DbHost: dbHost, Services: []string{container.Name}}
    backup.PgDumps = append(backup.PgDumps, &pg)
    return &pg
}

func AddSchema2PgDump(pgdump *PgDump, schema string) {
    pgdump.Schemas = append(pgdump.Schemas, schema)
}

func AddPgDumpAll2Backup(backup *Backup, container *Service, dbName string, dbUser string, dbHost string, parameters []string) (*PgDumpall) {
    if backup == nil || container == nil {
        log.Fatal()
    }
    pgdumpall := PgDumpall{DbHost: dbHost, DbName: dbName, DbUser: dbUser, Parameters: parameters, Services: []string{container.Name}}
    backup.PgDumpalls = append(backup.PgDumpalls, &pgdumpall)
    return &pgdumpall
}

func AddMysql2Backup(backup *Backup, container *Service, dbUser string, dbPassword string, databases []string, services []string, parameters []string) (*Mysql) {
    if backup == nil || container == nil {
        log.Fatal()
    }
    mysql := Mysql{DbUser: dbUser, DbPassword: dbPassword, Databases: databases, Services: services, Parameters: parameters}
    backup.Mysqls = append(backup.Mysqls, &mysql)
    return &mysql
}

func AddMysqlDB2Container(mysql *Mysql, container *Service, database string) {
    if mysql == nil || container == nil {
        log.Fatal()
    }
    mysql.Databases = append(mysql.Databases, database)
}

func WriteFile(location string, backup *Backup) {
    file, _ := json.MarshalIndent(backup, "", " ")
	_ = ioutil.WriteFile(location, file, 0644)
}

func IsMountInBackup(backup *Backup, source string) (bool) {
    for _, mount := range backup.Mounts {
        if mount.MountSource == source {
            return true
        }
    }
    return false
}
