package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const envPath = "./build/.env"

func GetEnvVariable(key string) string {

	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func IsProd() bool {
	return GetEnvVariable("ENV") == "prod"
}

func CacheIsEnabled() bool {
	return GetEnvVariable("CACHE_MODE") == "on"
}
