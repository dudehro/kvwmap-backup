package config

import (
	"log"
	"os"
	"encoding/csv"
)

const KeyBackupConfigDir 	= "backupconfig-directory"
const KeyHTTPPort		= "httpport"
const KeyConfigFile		= "kvwmapBackupConfig"

var config map[string]string


func printConfig() {
	for k, v := range config {
		log.Printf("%s=%s", k, v)
	}
}

func LoadConfig(f string) {

	config = make(map[string]string)

	if f == "" {
		f = GetConfigValFor(KeyConfigFile)
		log.Printf("No config-file specified, trying environment KVWMAPBACKUPCONFIG: %s", f)
	}


	if len(f) > 0 {
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

		for _, line := range data {
			for j, _ := range line {
				if j == 1 {
					config[line[0]] = line[1]
				}
			}
		}
		printConfig()
	} else {
		log.Printf("No configfile specified or found, using os-environment")
	}
}

func GetConfigValFor(s string) string {
	var ret string
	val, ok := config[s]
	if ok {
		ret = val
	} else {
		ret = os.Getenv(s)
		config[s] = ret
	}
	log.Printf("GetConfigValFor(%s) returns %s", s, ret)
	return ret
}
