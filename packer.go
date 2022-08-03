package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sanderhahn/gozip"
)

func cvtFileInfoToStrArray(files []fs.FileInfo) []string {
	var str_files []string
	for _, file := range files {
		str_files = append(str_files, file.Name())
	}
	return str_files
}

func getFiles(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	str_files := cvtFileInfoToStrArray(files)
	for index, file := range str_files {
		str_files[index] = filepath.Join(path, file)
	}
	return str_files
}

func createBundle(path string) {
	cmd := exec.Command("cmd.exe", "/c", "BatchBuildPython\\build.bat", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // add error checking
}

func createZipFromBundle(path string) {
	createBundle(path)
	files := getFiles(path)
	// fmt.Println(str_files)
	err := gozip.Zip("file.zip", files)
	if err != nil {
		log.Fatal(err)
	}
}
