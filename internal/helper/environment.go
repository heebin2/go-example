package helper

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const DefaultEnvFileName = ".env"

func LoadDotEnv(envFile string) error {
	return godotenv.Load(envFile)
}

func EnvInt(key string) (int, error) {
	v := os.Getenv(key)
	if len(v) <= 0 {
		return 0, fmt.Errorf("not found env %s", key)
	}
	return strconv.Atoi(v)
}

func EnvString(key string) (string, error) {
	v := os.Getenv(key)
	if len(v) <= 0 {
		return "", fmt.Errorf("not found env %s", key)
	}
	return v, nil
}
