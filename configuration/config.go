package config

import (
	"log"
        "flag"
)

const KeyBackupConfigDir 	= "backupconfigdir"
const KeyHTTPPort		= "port"

var config map[string]string

func InitConfig(){
	config = make(map[string]string)
	cliFlagPort             := flag.String(KeyHTTPPort, "8500", "port for server to listen on")
	cliFlagBackupConfigDir  := flag.String(KeyBackupConfigDir, "./backup-config", "location of backup-configs")
	flag.Parse()

	SetVal(KeyHTTPPort, *cliFlagPort)
	SetVal(KeyBackupConfigDir, *cliFlagBackupConfigDir)
	log.Println("printing config:")
	printConfig()
}

func printConfig() {
	for k, v := range config {
		log.Printf("%s=%s", k, v)
	}
}

func GetConfigValFor(s string) string {
	var ret string
	val, ok := config[s]
	if ok {
		ret = val
	} else {
		log.Printf("config %s not found", s)
		ret = ""
	}
	log.Printf("GetConfigValFor(%s) returns %s", s, ret)
	return ret
}

func SetVal(k string, v string) {
	config[k] = v
}
