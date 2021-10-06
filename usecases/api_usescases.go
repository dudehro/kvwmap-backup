package api_usecases

import (
	"github.com/kvwmap-backup/repository"
	"github.com/kvwmap-backup/models"
	"io/ioutil"
//	"fmt"
	"log"
	"strings"
)

func List_backup_configs() []string {
	root := "./"
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

func Load_backup_config(file string) *json_struct.GdiBackup {
	return crud.Load_backup_config(file)
}
