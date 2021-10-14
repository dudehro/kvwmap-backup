package crud

import (
	"github.com/kvwmap-backup/models"
	"encoding/json"
	"io/ioutil"
	"log"
//        "fmt"
)

func Load_backup_config(file string) (*structs.GdiBackup, error) {
	content, err := ioutil.ReadFile("./backup-config/"+file)
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
	data, err := json.MarshalIndent(j, "", " ")
//        fmt.Printf("Write_json output: \n %s\n", string(data))
	s := structs.Request{}
	s.Success=true
	if err != nil {
		s.Success = false
		s.Errors = append(s.Errors, err.Error())
	}
	err = ioutil.WriteFile("./backup-config/"+file, data, 0644)
	if err != nil {
		s.Success = false
		s.Errors = append(s.Errors, err.Error())
	}
	
	return s

}
