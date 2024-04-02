package main

import (
	"os"
	"os/exec"
)

const CONVERT_OUTPUT_FORMAT = "mp4"

func ConvertDirectory(directory string) {
	scriptPath := "./scripts/convert_to_mp4.sh"
	cmd := exec.Command(scriptPath, LIBRARY_PATH)
	cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
}
