package crud

import (
	"github.com/kvwmap-backup/models"
	"github.com/kvwmap-backup/configuration"
	"encoding/json"
	"io/ioutil"
	"log"
//        "fmt"
)

var dir string

func Load_backup_config(file string) (*structs.GdiBackup, error) {
	dir = config.GetConfigValFor(config.KeyBackupConfigDir)
	log.Printf("loading backup-config %s", dir+file)
	content, err := ioutil.ReadFile(dir+file)
	var payload structs.GdiBackup
	if err != nil {
		log.Print("Fehler beim öffnen: ", err)
	} else {
		err = json.Unmarshal(content, &payload)
		if err != nil {
			log.Print("Fehler beim unmarshalen: ", err)
		}
	}
	return &payload, err
}

func Write_json(j *structs.GdiBackup, file string) structs.Request {
	dir = config.GetConfigValFor(config.KeyBackupConfigDir)
	data, err := json.MarshalIndent(j, "", " ")
//        fmt.Printf("Write_json output: \n %s\n", string(data))
	s := structs.Request{}
	s.Success=true
	if err != nil {
		s.Success = false
		s.Errors = append(s.Errors, err.Error())
	}
	log.Printf("writing config file to %s", dir+file)
	err = ioutil.WriteFile(dir+file, data, 0644)
	if err != nil {
		s.Success = false
		s.Errors = append(s.Errors, err.Error())
	}

	return s

}
