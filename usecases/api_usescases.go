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

func Load_Taritem(file string, id string) (*structs.TarItem, structs.Request) {
	data, err := Load_backup_config(file)
	s := structs.Request{}
	if err != nil {
		s.Errors = append(s.Errors, err.Error())
		s.Success = false
		return nil, s
	}
        id_int, err := strconv.Atoi(id)
        if err != nil {
                s.Errors = append(s.Errors, err.Error())
		s.Success = false
		return nil, s
        }
	taritem := &structs.TarItem{}
	if len(data.Tar) >= id_int {
		taritem = data.Tar[id_int]
		s.Success = true
	}
	return taritem, s
}

func SaveItem(item interface{}, file string) structs.Request {
	data, err := Load_backup_config(file)
	s := structs.Request{}
	if err != nil {
		s.Errors = append(s.Errors, err.Error())
		s.Success = false
		return s
	}

	found := false
	switch item.(type){
		case structs.MysqlDumpItem:
			new_item := item.(structs.MysqlDumpItem)
			for i, d := range data.MysqlDump {
				if d.ContainerId == new_item.ContainerId &&
				   d.DbName == new_item.DbName {
					found = true
					data.MysqlDump[i] = &new_item     //update
				}
			}
			if !found {                             //insert
				mysqlDumps := data.MysqlDump
				mysqlDumps = append(mysqlDumps, &new_item)
				data.MysqlDump = mysqlDumps
			}

		case structs.TarItem:
			new_item := item.(structs.TarItem)
			for i, taritem := range data.Tar {
				if taritem.Source == new_item.Source {
					found = true
					data.Tar[i] = &new_item	//update
				}
			}
			if !found {				//insert
				tararray := data.Tar
				tararray = append(tararray, &new_item)
				data.Tar = tararray
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
