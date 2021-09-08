package config

import (
	"os"
)

var (
	ServerAddr  string
	BaseURL     string
	FileStorage string
)

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func InitEnv() {
	ServerAddr = getEnv("SERVER_ADDRESS", "localhost:8080")
	BaseURL = getEnv("BASE_URL", "http://localhost:8080")
	FileStorage = getEnv("FILE_STORAGE_PATH", "urlStorage.gob")
}
