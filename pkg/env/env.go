package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/praveennagaraj97/shopee/pkg/color"
	logger "github.com/praveennagaraj97/shopee/pkg/log"
)

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		logger.ErrorLogFatal("Failed to load env variables" + color.Reset)
	}

	logger.PrintLog("Environment Variables Loaded ⚙️ ", color.Purple)
}

func GetEnvVariable(key string) string {
	return os.Getenv(key)
}
