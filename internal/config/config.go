package config

import (
	"flag"
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
	// Checking flags first
	flag.StringVar(&ServerAddr, "a", "", "Server address")
	flag.StringVar(&BaseURL, "b", "", "Base URL")
	flag.StringVar(&FileStorage, "f", "", "Storage filename")
	flag.Parse()

	// If no flags checking ENV. If no ENV - setting defaults
	if ServerAddr == "" {
		ServerAddr = getEnv("SERVER_ADDRESS", "localhost:8080")
	}
	if BaseURL == "" {
		BaseURL = getEnv("BASE_URL", "http://localhost:8080")
	}
	if FileStorage == "" {
		FileStorage = getEnv("FILE_STORAGE_PATH", "urlStorage.gob")
	}
}
