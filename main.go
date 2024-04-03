package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const DEBUG = 0

const LOGFILE_PATH = "./brewTV.log"

func main() {
	CreateFileIfNotExists(LOGFILE_PATH)
	log_file, err := os.OpenFile(LOGFILE_PATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("[main]::Error setting up logger")
		os.Exit(4)
	}
	defer log_file.Close()
	log.SetOutput(log_file)

	startFlag := flag.Bool("start", false, "Run the server")
	setupFlag := flag.Bool("setup", false, "Add allowed MAC address")
	convertFlag := flag.Bool("convert", false, "Convert library to mp4")

	flag.Parse()

	if *startFlag {
		StartServer()
	}

	if *convertFlag {
		ConvertDirectory(LIBRARY_PATH)
	}

	if *setupFlag {
		args := flag.Args()
		if len(args) > 0 {
			AddToMacAddressAllowlist(args[0]) // Pass the first argument to AddMacAddress
		} else {
			fmt.Println("Missing MAC address argument for setup.")
		}
	}
}
