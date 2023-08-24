package main

import (
	"log"
	"os"
	"os/exec"
)

const CONVERT_OUTPUT_FORMAT = "mp4"

func ConvertDirectory(directory string) {
	videoFormats := []string{"avi", "wmv", "mov", "flv", "mkv", "mpg", "mpeg", "webm", "3gp", "m4v", "ogv"}

	scriptPath := "./scripts/convert_all_files_in_dir_by_format.sh"

	for _, input_format := range videoFormats {
		cmd := exec.Command(scriptPath, LIBRARY_PATH, input_format, CONVERT_OUTPUT_FORMAT)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("[ConvertDirectory]::Error running script for format %s: %v\n", input_format, err)
		} else {
			log.Printf("[ConvertDirectory]::Conversion successful for format %s\n", input_format)
		}
	}
}
