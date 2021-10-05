package crud

import (
	"github.com/kvwmap-backup/models"
	"encoding/json"
	"io/ioutil"
	"log"
)

func Load_json(){
	content, err := ioutil.ReadFile("./sicherung.json")
	if err != nil {
		log.Fatal("Fehler beim öffnen: ",err)
	}

	var payload json_struct.GdiBackup
	err = json.Unmarshal(content, &payload)
       	if err != nil {
		log.Fatal("Fehler beim unmarshalen: ", err)
	}


	log.Printf("Erfolg")
}

func Write_json(j *json_struct.GdiBackup){

}
