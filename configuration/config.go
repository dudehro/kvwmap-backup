package config

import (

)



var configfile string

func SetConfigFile(f string) {
	configfile = f
}

func load_config() {
	
}

func GetConfigValFor(s string) string {
	return configfile
}
