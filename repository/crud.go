package crud

import (
	"github.com/kvwmap-backup/models"
	"github.com/kvwmap-backup/configuration"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

var dir string

func Load_backup_config(file string) (*structs.GdiBackup, error) {
	location := filepath.Join( config.GetConfigValFor(config.KeyBackupConfigDir), file)

	log.Printf("loading backup-config %s", location )
	content, err := ioutil.ReadFile(location)
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
	location := filepath.Join( config.GetConfigValFor(config.KeyBackupConfigDir), file)
	data, err := json.MarshalIndent(j, "", " ")
//        fmt.Printf("Write_json output: \n %s\n", string(data))
	s := structs.Request{}
	s.Success=true
	if err != nil {
		s.Success = false
		s.Errors = append(s.Errors, err.Error())
	}
	log.Printf("writing config file to %s", location)
	err = ioutil.WriteFile(location, data, 0644)
	if err != nil {
		s.Success = false
		s.Errors = append(s.Errors, err.Error())
	}

	return s

}
