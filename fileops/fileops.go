package fileops

import (
	log "kvwmap-backup/logging"
	"os"
)

func Mkdir(dir string) bool {
	var err error
	if !DirExists(dir) {
		err := os.Mkdir(dir, 0640)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return err == nil
}

func DirExists(dir string) bool {
	_, err := os.Stat(dir)
	return os.IsNotExist(err)
}

func IsDir(path string) bool {
	info, _ := os.Stat(path)
	return info.IsDir()
}

func IsFile(file string) bool {
	return true
}
