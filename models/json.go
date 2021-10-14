// Code generated by schema-generate. DO NOT EDIT.

package structs

//type Backup struct {
//	Backup	GdiBackup
//	Meta	map[string]string
//}

// GdiBackup 
type GdiBackup struct {
  BackupFolder string `json:"backup_folder,omitempty"`
  BackupPath string `json:"backup_path,omitempty"`
  Beschreibung string `json:"beschreibung,omitempty"`
  CronInterval string `json:"cron_interval,omitempty"`
  DeleteAfterNDays int `json:"delete_after_n_days,omitempty"`
  DifferentialBackupDuration int `json:"differential_backup_duration,omitempty"`
  Id int `json:"id,omitempty"`
  IntervalParameter1 string `json:"interval_parameter_1,omitempty"`
  IntervalParameter2 string `json:"interval_parameter_2,omitempty"`
  IntervalType string `json:"interval_type,omitempty"`
  MysqlDump []*MysqlDumpItem `json:"mysql_dump,omitempty"`
  Name string `json:"name,omitempty"`
  PgDump []*PgDumpItem `json:"pg_dump,omitempty"`
  PgDumpall []*PgDumpallItem `json:"pg_dumpall,omitempty"`
  Rsync []*RsyncItem `json:"rsync,omitempty"`
  Tar []*TarItem `json:"tar,omitempty"`
}

// MysqlDumpItem 
type MysqlDumpItem struct {
  ConnectionId int `json:"connection_id,omitempty"`
  ContainerId string `json:"container_id,omitempty"`
  DbName string `json:"db_name,omitempty"`
  DockerNetwork string `json:"docker_network,omitempty"`
  MysqlDumpParameter string `json:"mysql_dump_parameter,omitempty"`
  TargetName string `json:"target_name,omitempty"`
}

// PgDumpItem 
type PgDumpItem struct {
  ConnectionId int `json:"connection_id,omitempty"`
  ContainerId string `json:"container_id,omitempty"`
  DbName string `json:"db_name,omitempty"`
  DbUser string `json:"db_user,omitempty"`
  DockerNetwork string `json:"docker_network,omitempty"`
  PgDumpColumnInserts bool `json:"pg_dump_column_inserts,omitempty"`
  PgDumpInExcludeSchemas string `json:"pg_dump_in_exclude_schemas,omitempty"`
  PgDumpInExcludeTables string `json:"pg_dump_in_exclude_tables,omitempty"`
  PgDumpInserts bool `json:"pg_dump_inserts,omitempty"`
  PgDumpParameter string `json:"pg_dump_parameter,omitempty"`
  PgDumpSchemas []string `json:"pg_dump_schemas,omitempty"`
  PgDumpTables []string `json:"pg_dump_tables,omitempty"`
  TargetName string `json:"target_name,omitempty"`
}

// PgDumpallItem 
type PgDumpallItem struct {
  ConnectionId int `json:"connection_id,omitempty"`
  ContainerId string `json:"container_id,omitempty"`
  DbName string `json:"db_name,omitempty"`
  DbUser string `json:"db_user,omitempty"`
  DockerNetwork string `json:"docker_network,omitempty"`
  PgDumpallParameter string `json:"pg_dumpall_parameter,omitempty"`
  TargetName string `json:"target_name,omitempty"`
}

// RsyncItem 
type RsyncItem struct {
  Destination string `json:"destination,omitempty"`
  Parameter string `json:"parameter,omitempty"`
  Source string `json:"source,omitempty"`
}

// TarItem 
type TarItem struct {
  Exclude string `json:"exclude,omitempty"`
  Source string `json:"source,omitempty"`
  TarCompress bool `json:"tar_compress,omitempty"`
  TargetName string `json:"target_name,omitempty"`
}
