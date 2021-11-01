package api_usecases

import (
	"github.com/kvwmap-backup/repository"
	"github.com/kvwmap-backup/models"
	"github.com/kvwmap-backup/configuration"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
)

func List_backup_configs() []string {
	root := config.GetConfigValFor(config.KeyBackupConfigDir)
	var files []string
	dir_content, err := ioutil.ReadDir(root)
        if err != nil {
                log.Fatal(err)
        }
        for _, item := range dir_content {
		if strings.ToUpper(item.Name())[len(item.Name())-4:] == "JSON" {
			files = append(files, item.Name())
		}
        }
	return files
}

func Load_backup_config(file string) (*structs.GdiBackup, error) {
	return crud.Load_backup_config(file)
}

func Save_Backup_Config(data *structs.GdiBackup, file string) structs.Request {
	s := structs.Request{}
	write_result := crud.Write_json(data, file)
	if write_result.Success {
		s.Success = true
	} else {
		s.Errors = write_result.Errors
		s.Success = false
	}
	return s
}

func SaveItem(item interface{}, file string, id string) structs.Request {
	data, err := Load_backup_config(file)
	s := structs.Request{}
	if err != nil {
		s.Errors = append(s.Errors, err.Error())
		s.Success = false
		return s
	}

	var id_int int
	if len(id) > 0 {
		id_int, err = strconv.Atoi(id)
		if err != nil {
			s.Errors = append(s.Errors, err.Error())
			s.Success = false
			return s
		}
	} else {
		id_int = -1
	}

	error_ndf := "Angegebenes Element existiert nicht!"
	switch item.(type) {
		case structs.MysqlDumpItem:
			new_item := item.(structs.MysqlDumpItem)
			if id_int > -1 {					//update
				if len(data.MysqlDump) >= id_int {
					data.MysqlDump[id_int] = &new_item
				} else {
					s.Errors = append(s.Errors, error_ndf)
				}
			} else {						//insert
				mysqlDumps := data.MysqlDump
				mysqlDumps = append(mysqlDumps, &new_item)
				data.MysqlDump = mysqlDumps
			}

		case structs.TarItem:
			new_item := item.(structs.TarItem)
			if id_int > -1 {
				if len(data.Tar) >= id_int {
					data.Tar[id_int] = &new_item
				} else {
					s.Errors = append(s.Errors, error_ndf)
				}
			} else {
				tarItems := data.Tar
				tarItems = append(tarItems, &new_item)
				data.Tar = tarItems
			}
		case structs.PgDumpItem:
			new_item := item.(structs.PgDumpItem)
			if id_int > -1 {
				if len(data.PgDump) >= id_int {
					data.PgDump[id_int] = &new_item
				} else {
					s.Errors = append(s.Errors, error_ndf)
				}
			} else {
				pgDumps := data.PgDump
				pgDumps = append(pgDumps, &new_item)
				data.PgDump = pgDumps
			}
		case structs.PgDumpallItem:
			new_item := item.(structs.PgDumpallItem)
			if id_int > -1 {
				if len(data.PgDumpall) >= id_int {
					data.PgDumpall[id_int] = &new_item
				} else {
					s.Errors = append(s.Errors, error_ndf)
				}
			} else {
				log.Println("save new item")
				pgDumps := data.PgDumpall
				pgDumps = append(pgDumps, &new_item)
				data.PgDumpall = pgDumps
                        }
		default:
			log.Println("Keine Ahnung?")
	}

	req := Save_Backup_Config(data, file)
	if req.Success {
		s.Success = true
	} else {
		s.Errors = append(s.Errors, req.Errors[0])
		s.Success = false
	}

	return s
}

func DeleteItem(file string, id string, t interface{}) structs.Request {
	log.Printf("DeleteItem(%s, %s, %T)", file, id, t)
	data, err := Load_backup_config(file)
	s := structs.Request{}
	if err != nil {
		s.Errors = append(s.Errors, err.Error())
		s.Success = false
		return s
	}
        id_int, err := strconv.Atoi(id)
        if err != nil {
                s.Errors = append(s.Errors, err.Error())
                s.Success = false
                return s
        }

	switch t.(type) {
		case structs.TarItem:
			if id_int <= len(data.Tar) {
				copy(data.Tar[id_int:], data.Tar[id_int+1:])
				data.Tar[len(data.Tar)-1] = nil
				data.Tar = data.Tar[:len(data.Tar)-1]
				s.Success = true
			} else {
				s.Errors = append(s.Errors, "Es existiert kein Eintrag mit ID="+id)
				s.Success = false
			}
		case structs.MysqlDumpItem:
			if id_int <= len(data.MysqlDump) {
				copy(data.MysqlDump[id_int:], data.MysqlDump[id_int+1:])
				data.MysqlDump[len(data.MysqlDump)-1] = nil
				data.MysqlDump = data.MysqlDump[:len(data.MysqlDump)-1]
				s.Success = true
			} else {
				s.Errors = append(s.Errors, "Es existiert kein Eintrag mit ID="+id)
				s.Success = false
			}
		case structs.PgDumpItem:
			if id_int <= len(data.PgDump) {
				copy(data.PgDump[id_int:], data.PgDump[id_int+1:])
				data.PgDump[len(data.PgDump)-1] = nil
				data.PgDump = data.PgDump[:len(data.PgDump)-1]
				s.Success = true
			} else {
				s.Errors = append(s.Errors, "Es existiert kein Eintrag mit ID="+id)
				s.Success = false
			}
		case structs.PgDumpallItem:
			if id_int <= len(data.PgDumpall) {
				copy(data.PgDumpall[id_int:], data.PgDumpall[id_int+1:])
				data.PgDumpall[len(data.PgDumpall)-1] = nil
				data.PgDumpall = data.PgDumpall[:len(data.PgDumpall)-1]
				s.Success = true
			} else {
				s.Errors = append(s.Errors, "Es existiert kein Eintrag mit ID="+id)
				s.Success = false
			}

	}

	req := Save_Backup_Config(data, file)
	if req.Success {
		s.Success = true
	} else {
		s.Errors = append(s.Errors, req.Errors[0])
		s.Success = false
	}

	return s
}
