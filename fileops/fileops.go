package fileops

import (
	"archive/tar"
	"fmt"
	"io"
	log "kvwmap-backup/logging"
	"os"
	"path/filepath"
	"strings"
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

	tarfile, err := os.Create(tarball)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tw := tar.NewWriter(tarfile)
	defer tw.Close()

	//prüfen ob stat möglich ist
	info, err := os.Stat(file)
	if err != nil {
		return nil
	}

	//header erstellen
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

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
		fmt.Println(dir)
		fmt.Println(base)
		fmt.Println(header.Name)
	}

	//header schreiben
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	//Verzeichnisse überspringen
	if info.IsDir() {
		return nil
	}

	//Datei öffnen
	copyfile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer copyfile.Close()
	//Datei ins Tar kopieren
	_, err = io.Copy(tw, copyfile)
	return err
}

// copypaste from https://golangdocs.com/tar-gzip-in-golang
func Tar(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
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
