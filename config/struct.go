// Code generated by schema-generate. DO NOT EDIT.

package config

// Root
type Root *Backup

// Backup 
type Backup struct {
  BackupPath string `json:"backup_path,omitempty"`
  Networks []*Network `json:"networks,omitempty"`
}

// Mysql 
type Mysql struct {
  DbPassword string `json:"db_password,omitempty"`
  DbUser string `json:"db_user,omitempty"`
  Mysqldump []*MysqldumpItems `json:"mysqldump,omitempty"`
}

// MysqldumpItems 
type MysqldumpItems struct {
  DbName string `json:"db_name,omitempty"`
  MysqlDumpParameter string `json:"mysql_dump_parameter,omitempty"`
}

// Network 
type Network struct {
  Name string `json:"name,omitempty"`
  SaveNetwork bool `json:"save_network,omitempty"`
  Services []*Service `json:"services,omitempty"`
  Subnet string `json:"subnet,omitempty"`
}

// Pgdump 
type Pgdump struct {
  IncludeListedSchemas bool `json:"include_listed_schemas,omitempty"`
  PgDumpParameter string `json:"pg_dump_parameter,omitempty"`
  Schemas []string `json:"schemas,omitempty"`
}

// PgdumpallItems 
type PgdumpallItems struct {
  PgDumpallParameter string `json:"pg_dumpall_parameter,omitempty"`
}

// Postgres 
type Postgres struct {
  DbName string `json:"db_name,omitempty"`
  DbUser string `json:"db_user,omitempty"`
  Host string `json:"host,omitempty"`
  Pgdump *Pgdump `json:"pgdump,omitempty"`
  Pgdumpall []*PgdumpallItems `json:"pgdumpall,omitempty"`
}

// Service 
type Service struct {
  Image string `json:"image,omitempty"`
  Mysql *Mysql `json:"mysql,omitempty"`
  Name string `json:"name,omitempty"`
  Postgres *Postgres `json:"postgres,omitempty"`
  SaveService bool `json:"save_service,omitempty"`
  Tar *Tar `json:"tar,omitempty"`
}

// Tar 
type Tar struct {
  DiffBackupDays string `json:"diff_backup_days,omitempty"`
  Directories []*Taritem `json:"directories,omitempty"`
}

// Taritem 
type Taritem struct {
  ExcludeDirs string `json:"exclude_dirs,omitempty"`
  MountDestination string `json:"mount_destination,omitempty"`
  MountSource string `json:"mount_source,omitempty"`
  SaveData bool `json:"save_data,omitempty"`
  SavedByService string `json:"saved_by_service,omitempty"`
}