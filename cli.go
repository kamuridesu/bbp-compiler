package main

import (
	"flag"
	// "fmt"
	"log"
	"os"
)

func getCliArgs() [2]string {
	options := [2]string{
		"pack",
	}
	if len(os.Args) < 2 || len(os.Args) > 2 {
		log.Fatal("Invalid argument numbers")
		os.Exit(1)
	}
	option := os.Args[1]
	var app_path string
	flag.StringVar(&app_path, "path", ".", "Application path")

	for _, op := range options {
		if option == op {
			return [2]string{option, app_path}
		}
	}
	log.Fatal("Argument not expected!")
	os.Exit(127)
	return [2]string{""}
}

func cli() {
	arguments := getCliArgs()
	command := arguments[0]
	path := arguments[1]
	if command == "pack" {
		createZipFromBundle(path)
		compile(path)
	}
}
