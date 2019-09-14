package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/rs/xid"
)

func main() {
	flag.Parse()
	if err := cmd(flag.Arg(0)); err != nil {
		log.Fatal(err)
	}
}

func cmd(inputDir string) error {
	if inputDir == "" {
		return fmt.Errorf("No argument error")
	}

	files, err := getFileList(inputDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println(file.Name)
		if _, err := xid.FromString(file.Name); err == nil {
			continue
		}

		id := xid.NewWithTime(file.Time)
		newName := id.String() + filepath.Ext(file.Name)
		err := os.Rename(path.Join(inputDir, file.Name), path.Join(inputDir, newName))
		if err != nil {
			return err
		}
	}

	return nil
}

type file struct {
	Name string
	Time time.Time
}

func getFileList(dir string) ([]file, error) {
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []file
	for _, path := range paths {
		if path.IsDir() {
			continue
		}
		if !isFileImage(path.Name()) {
			continue
		}

		file := file{Name: path.Name(), Time: path.ModTime()}
		files = append(files, file)
	}
	return files, nil
}

func isFileImage(name string) bool {
	exts := []string{".jpeg", ".jpg", ".png"}
	for _, ext := range exts {
		if filepath.Ext(name) == ext {
			return true
		}
	}
	return false
}
