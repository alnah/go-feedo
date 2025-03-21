package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// configFilePath represents the JSON configuration file path
const configFilePath = ".config/go-feedo/config.json"

// DatabaseConfig holds configuration values for the database
type DatabaseConfig struct {
	URL             string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// SetUser configures the current user name for the database
func (db *DatabaseConfig) SetUser(currentUserName string) error {
	db.CurrentUserName = currentUserName
	return write(*db)
}

// Read decodes the JSON database configuration from the configuration file path
func Read() (DatabaseConfig, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return DatabaseConfig{}, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return DatabaseConfig{}, err
	}
	defer func() { _ = file.Close() }()
	dec := json.NewDecoder(file)
	var db DatabaseConfig
	if err = dec.Decode(&db); err != nil {
		return DatabaseConfig{}, err
	}
	return db, nil
}

// write encodes the JSON database configuration to the configuration file path
func write(db DatabaseConfig) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close }()
	enc := json.NewEncoder(file)
	enc.SetIndent("", " ") // pretty format
	if err = enc.Encode(db); err != nil {
		return err
	}
	return nil
}

// getConfigFilePath constructs the complete configuration file path
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configFilePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil { // Ensure the dir exists
		return "", err
	}
	return fullPath, nil
}
