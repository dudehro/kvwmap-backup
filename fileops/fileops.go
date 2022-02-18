package fileops

import (
	"archive/tar"
	"bufio"
	"bytes"
	"io"
	log "kvwmap-backup/logging"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	//    "errors"
	"fmt"
)

const (
	TemplateNetwork = "network"
	TemplateService = "service"
	TemplateMount   = "mount"
)

type templateData struct {
	Network string
	Service string
	Mount   string
}

func Mkdir(dir string) bool {
	var err error
	if !IsDir(dir) {
		err := os.MkdirAll(dir, 0750)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return err == nil
}

func DirExists(dir string) bool {
	fileinfo, err := os.Stat(dir)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fileinfo.IsDir()
	}
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return false
	} else {
		return info.IsDir()
	}
}

func GetFileInfo(fileDir string) os.FileInfo {
	fileInfo, err := os.Lstat(fileDir)
	if err != nil {
		panic(err)
	}
	return fileInfo
}

/* Adds file to tarball. target can be new or existing tar-archive.
If newDir is set, header.Name will be set to newDir
*/
func TarFile(file, newDir, tarball string) error {

	log.Println("TarFile()")
	log.Printf("file: %s\n", file)
	log.Printf("newDir: %s\n", newDir)
	log.Printf("tarball: %s\n", tarball)

	//FileInfo holen
	info, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var tarfile *os.File
	var appendFile bool
	//var err error

	if !IsFile(file) && !IsDir(file) {
		tarfile, err = os.Create(file)
		appendFile = false
	} else {
		tarfile, err = os.OpenFile(file, os.O_RDWR, os.ModePerm)
		appendFile = true
	}
	if err != nil {
		return err
	}

	if appendFile {
		if _, err = tarfile.Seek(-1024, os.SEEK_END); err != nil {
			log.Println(err)
		}
	}

	tw := tar.NewWriter(tarfile)

	fmt.Println("Name : ", info.Name())
	fmt.Println("Size : ", info.Size())
	fmt.Println("Mode/permission : ", info.Mode())
	fmt.Println("Is a regular file? : ", info.Mode().IsRegular())

	//Verzeichnisse überspringen
	if info.IsDir() {
		return nil
	}

	//nur reguläre Dateien
	if !info.Mode().IsRegular() {
		return nil
	}

	//header erstellen
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		log.Fatal(err)
		return err
	}
	header.Size = info.Size()

	//dir ermitteln
	var dir, base string
	if len(newDir) == 0 {
		dir = filepath.Dir(file)
	} else {
		dir = newDir
	}
	base = filepath.Base(file)

	//Pfad+Dateiname im Header setzen
	if dir != "" {
		header.Name = filepath.Join(dir, base)
		//		fmt.Println(dir)
		//		fmt.Println(base)
		//		fmt.Println(header.Name)
	}

	//header schreiben
	if err := tw.WriteHeader(header); err != nil {
		log.Fatal(err)
		return err
	}

	//Datei öffnen
	fmt.Println("Datei öffnen")
	copyfile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//Datei ins Tar kopieren
	fmt.Println("Datei ins tar kopieren")
	num_bytes, err := io.Copy(tw, copyfile)
	fmt.Printf("Bytes kopiert %d \n", num_bytes)
	fmt.Printf("Size %d \n", info.Size())
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("alles schließen")
	tw.Close()
	copyfile.Close()
	tarfile.Close()
	return err
}

// copypaste from https://golangdocs.com/tar-gzip-in-golang
func TarDir(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

func GetValueForKey(file, key string) string {
	var returnval string
	filehandler, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(filehandler)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		linearray := strings.Split(scanner.Text(), ":")
		pos0 := linearray[0]
		if pos0 == key {
			returnval = linearray[1]
		}
	}
	filehandler.Close()
	return returnval
}

func GetTemplateDir(dirFile, templateType, networkName, serviceName, mountName string) string {
	data := templateData{Network: networkName, Service: serviceName, Mount: mountName}
	parsedTemplate, err := template.New("directory").Parse(GetValueForKey(dirFile, templateType))
	if err != nil {
		log.Fatal(err)
	}
	var result bytes.Buffer
	err = parsedTemplate.Execute(&result, data)
	if err != nil {
		log.Fatal(err)
	}
	return result.String()
}
