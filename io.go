package main

import (
	"bufio"
	"os"
	"strings"
)

func CreateFileIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, nil, 0666); err != nil {
			return err
		}
	}
	return nil
}

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0666); err != nil {
			return err
		}
	}
	return nil
}

func AppendToFile(path, content string) error {
	if err := CreateFileIfNotExists(path); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func ReadlinesFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []string
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
