package config

import "os"

var (
	ServerAddr string
	BaseURL    string
)

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func Init() {
	ServerAddr = getEnv("SERVER_ADDRESS", "localhost:8080")
	BaseURL = getEnv("BASE_URL", "http://localhost:8080")
}
