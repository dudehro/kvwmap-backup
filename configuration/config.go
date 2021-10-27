package config

import (
	"log"
	"os"
	"encoding/csv"
)

const KeyBackupConfigDir 	= "backupconfig-directory"
const KeyHTTPPort		= "httpport"

var config map[string]string

//func SetConfigFile(f string) {
//	configfile = f
//}

func printConfig() {
	for k, v := range config {
		log.Printf("%s=%s", k, v)
	}
}



func LoadConfig(f string) {
	log.Printf("Loading config from %s", f)
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = '='
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	config = make(map[string]string)
	for _, line := range data {
		for j, _ := range line {
			if j == 1 {
				config[line[0]] = line[1]
			}
		}
	}
	printConfig()
}

func GetConfigValFor(s string) string {
	var ret string
	val, ok := config[s]
	if ok {
		ret = val
	} else {
		ret = ""
	}
	log.Printf("GetConfigValFor(%s) returns %s", s, ret)
	return ret
}
