package main

import (
	"os"
	"archive/zip"
	"path/filepath"
	"strings"
	"io"
	"io/ioutil"
	"log"
	"fmt"
	"time"
	"strconv"
)

func RecursiveZip(pathToZip, destinationPath string) error {
	os.MkdirAll(destinationPath, os.ModePerm)
	zipName := genDestFileName(pathToZip, destinationPath)
	destinationFile, err := os.Create(zipName)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func genDestFileName(pathToZip, destinationPath string) string{
	zippedName := strings.Split(pathToZip, "/")[:]
	fmt.Println(zippedName)
	t := time.Now().Format("20060102150405")
	var pathToZipSlice []string
	pathToZipSlice = strings.Split(pathToZip, "/")
	dirName := pathToZipSlice[len(pathToZipSlice)-1]

	return destinationPath + "/" + dirName + "_backup_" + t + ".zip"
}

func main() {
	args := os.Args
	fmt.Println(args[1:])
	src := args[1]
	dest := args[2]
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}
	var actionsMap = make(map[int]string)
	fmt.Println("Creating backup for the following files/dirs")
	fmt.Println(0, "- All subdirectoried")
	actionsMap[0] = ""
	i := 1
	for _, f := range files {
		fmt.Println(i, "-", f.Name())
		actionsMap[i] = f.Name()
		i++
	}
	fmt.Println("Choose directory for backup...")
	fmt.Print("Enter text: ")
	input := ""
	fmt.Scanln(&input)
	chosenNumber, _ := strconv.ParseInt(input, 10 , 0)
	chosenPath := actionsMap[int(chosenNumber)]
	if chosenPath != "" {
		fmt.Println(chosenPath)
		src += "/" + chosenPath
		dest += "/" + chosenPath
	} else {
		fmt.Println("Key does not exists. exiting...")
		return
	}
	fmt.Println("Starting backup", chosenPath)
	err = RecursiveZip(src, dest)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Backup completed successfully!")
	}
}
