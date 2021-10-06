package crud

import (
	"github.com/kvwmap-backup/models"
	"encoding/json"
	"io/ioutil"
	"log"
)

func Load_backup_config(file string) *json_struct.GdiBackup {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Fehler beim öffnen: ", err)
	}

	var payload json_struct.GdiBackup
	err = json.Unmarshal(content, &payload)
       	if err != nil {
		log.Fatal("Fehler beim unmarshalen: ", err)
	}
//	log.Printf("Erfolg")
	return &payload
}

func Write_json(j *json_struct.GdiBackup){
}
