package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func CreateFileIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func AppendToFile(path string, content string) error {
	CreateFileIfNotExists(path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("[AppendToFile]::Error opening file %s", path)
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		log.Printf("[AppendToFile]::Error writting to file %s", path)
		return err
	}
	return nil
}

func ReadlinesFromFile(path string) ([]string, error) {
	var data []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		data = append(data, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}
