// Code generated by schema-generate. DO NOT EDIT.

package config

// Root
type Root *Backup

// Backup
type Backup struct {
	BackupFolder string       `json:"backup_folder,omitempty"`
	BackupPath   string       `json:"backup_path,omitempty"`
	Mounts       []*Mount     `json:"mounts,omitempty"`
	Mysqls       []*Mysql     `json:"mysqls,omitempty"`
	Networks     []*Network   `json:"networks,omitempty"`
	PgDumpalls   []*PgDumpall `json:"pg_dumpalls,omitempty"`
	PgDumps      []*PgDump    `json:"pg_dumps,omitempty"`
	Services     []*Service   `json:"services,omitempty"`
}

// Mount
type Mount struct {
	ExcludeDirs      string `json:"exclude_dirs,omitempty"`
	MountDestination string `json:"mount_destination,omitempty"`
	MountSource      string `json:"mount_source,omitempty"`
	Service          string `json:"service,omitempty"`
}

// Mysql
type Mysql struct {
	Databases  []string `json:"databases,omitempty"`
	DbPassword string   `json:"db_password,omitempty"`
	DbUser     string   `json:"db_user,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
	Services   []string `json:"services,omitempty"`
}

// Network
type Network struct {
	Name   string `json:"name,omitempty"`
	Subnet string `json:"subnet,omitempty"`
}

// PgDump
type PgDump struct {
	DbHost     string   `json:"db_host,omitempty"`
	DbName     string   `json:"db_name,omitempty"`
	DbUser     string   `json:"db_user,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
	Schemas    []string `json:"schemas,omitempty"`
	Services   []string `json:"services,omitempty"`
}

// PgDumpall
type PgDumpall struct {
	DbHost     string   `json:"db_host,omitempty"`
	DbName     string   `json:"db_name,omitempty"`
	DbUser     string   `json:"db_user,omitempty"`
	Parameters []string `json:"parameters,omitempty"`
	Services   []string `json:"services,omitempty"`
}

// Service
type Service struct {
	Image    string   `json:"image,omitempty"`
	Name     string   `json:"name,omitempty"`
	Networks []string `json:"networks,omitempty"`
}
