// internal/storage.go
package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var dataFile = "data.json"
var users = make(map[string]*User)
var useMockData = false

// UseMockData sets the mock data for testing
func UseMockData(mockUsers map[string]*User) {
	users = mockUsers
	useMockData = true
}

// SaveData saves the current state of users to a JSON file
func SaveData() error {
	if useMockData {
		return nil
	}
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, data, 0644)
}

// LoadData loads the state of users from a JSON file
func LoadData() error {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return nil // No file, skip loading
	}

	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &users)
}

// SetDataFile sets the data file path and checks if it's a valid path
func SetDataFile(path string) error {
	if strings.HasSuffix(path, "/") {
		return errors.New("provided path ends with a '/', please provide a valid file path")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return errors.New("invalid file path")
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			dir := filepath.Dir(absPath)
			if _, dirErr := os.Stat(dir); os.IsNotExist(dirErr) {
				return errors.New("directory does not exist")
			}
		} else {
			return errors.New("error checking path")
		}
	} else if info.IsDir() {
		return errors.New("provided path is a directory, please provide a valid file path")
	}

	dataFile = absPath
	return nil
}
