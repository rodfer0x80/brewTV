package main

import (
	"log"
	"os"
)

const LOGFILE_PATH = "./brewTV.log"

func StartLogger() {
	CreateFileIfNotExists(LOGFILE_PATH)
	log_file, err := os.OpenFile(LOGFILE_PATH, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("[StartLogger]::Error creating log file:", err)
		os.Exit(4)
	}
	defer log_file.Close()
	log.SetOutput(log_file)
}
