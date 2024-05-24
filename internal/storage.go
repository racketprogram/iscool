// internal/storage.go
package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var dataFile = "data.json"
var users = make(map[string]*User)

// SaveData saves the current state of users to a JSON file
func SaveData() error {
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
