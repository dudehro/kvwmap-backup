package api_usecases

import (
	"github.com/kvwmap-backup/repository"
	"github.com/kvwmap-backup/models"
	"io/ioutil"
//	"fmt"
	"log"
	"strings"
	"strconv"
)

func List_backup_configs() []string {
	root := "./backup-config/"
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



func Save_TarItem(item structs.TarItem, file string) structs.Request {
	data, err := Load_backup_config(file)
	s := structs.Request{}
	if err != nil {
		s.Errors = append(s.Errors, err.Error())
		s.Success = false
		return s
	}

	found := false
	for i, taritem := range data.Tar {
		if taritem.Source == item.Source {
			found = true
			data.Tar[i] = &item	//update
		}
	}
	if !found {				//insert
		tararray := data.Tar		//refactor
		tararray = append(tararray, &item)
		data.Tar = tararray
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

func Delete_TarItem(file string, id string) structs.Request {
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
	if id_int <= len(data.Tar) {

		copy(data.Tar[id_int:], data.Tar[id_int+1:])
		data.Tar[len(data.Tar)-1] = nil
		data.Tar = data.Tar[:len(data.Tar)-1]

		s.Success = true
	} else {
		s.Errors = append(s.Errors, "Es existiert kein Eintrag mit ID="+id)
		s.Success = false
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
