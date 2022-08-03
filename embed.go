package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/akavel/rsrc/rsrc"
)

func getFilesFromUser(message string) string {
	data := input(message)
	if _, err := os.Stat(data); err == nil {
		return data
	} else if errors.Is(err, os.ErrNotExist) {
		log.Fatal("File does not exists")
		os.Exit(127)
	}
	log.Fatal("Unknowk error")
	os.Exit(1)
	return ""
}

func getArch() string {
	arches := [2]string{"amd64", "386"}
	user_inp := input("Choose the arch [386/amd64], default amd64:  ")
	for _, arch := range arches {
		if user_inp == arch {
			return arch
		}
	}
	return "amd64"
}

func getVersion() string {
	version := input("Enter the version (default: 0.0.1): ")
	if version == "" {
		return "0.0.1"
	}
	return version
}

func getPrivileges() string {
	privs := [3]string{"asInvoker", "highestAvailable", "requireAdministrator"}
	user_inp := input("Choose the arch [asInvoker/highestAvailable/requireAdministrator], default asInvoker: ")
	for _, priv := range privs {
		if user_inp == priv {
			return priv
		}
	}
	return "asInvoker"
}
func gatherMetadata() map[string]string {
	items := make(map[string]string)
	items["appname"] = input("Enter the app name: ")
	items["entrypoint"] = input("Enter the main file name (entrypoint): ")
	items["version"] = getVersion()
	items["icon"] = getFilesFromUser("Enter the icon path: ")
	items["arch"] = getArch()
	items["privilege"] = getPrivileges()
	return items
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}
}

func formatManifestFile() map[string]string {
	info := gatherMetadata()
	manifest, err := os.Open("./manifest.xml.template")
	check(err)
	scanner := bufio.NewScanner(manifest)
	manifest_text := ""
	for scanner.Scan() {
		manifest_text += scanner.Text() + "\n"
	}
	manifest_text = strings.Replace(manifest_text, "%{APPNAME}%", info["appname"], 2)
	manifest_text = strings.Replace(manifest_text, "%{VERSION}%", info["version"], 1)
	manifest_text = strings.Replace(manifest_text, "%{ARCH}%", info["arch"], 1)
	manifest_text = strings.Replace(manifest_text, "%{PRIVILEGE_LEVEL}%", info["privilege"], 1)
	out_filename := "./manifest.xml"
	e := os.WriteFile(out_filename, []byte(manifest_text), 0644)
	check(e)
	info["manifest"] = out_filename
	return info
}

func createManifestSyso(info map[string]string) {
	err := rsrc.Embed(info["appname"]+".syso", info["arch"], info["manifest"], info["icon"])
	check(err)
}

func formatMainFile(entrypoint string, path string) {
	main, err := os.Open("./main.go.template")
	check(err)
	scanner := bufio.NewScanner(main)
	main_text := ""
	for scanner.Scan() {
		main_text += scanner.Text() + "\n"
	}
	main_text = strings.Replace(main_text, "%{ENTRYPOINT}%", entrypoint, 1)
	e := os.WriteFile(filepath.Join(path, "main.go"), []byte(main_text), 0644)
	check(e)
}

func generatePackageInfo(path string) {
	info := gatherMetadata()
	formatMainFile(info["entrypoint"], path)
	createManifestSyso(info)
}

func compile(path string) {
	generatePackageInfo(path)
	cmd := exec.Command("cmd.exe", "/c", "cd", path, "&&", "go", "build", "&&", "del", "main.go")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // add error checking
}
