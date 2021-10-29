package main

import (
	"github.com/kvwmap-backup/interface"
	"github.com/kvwmap-backup/configuration"
)

func main(){
	config.LoadConfig("")
	delivery.URIHandler()
}
