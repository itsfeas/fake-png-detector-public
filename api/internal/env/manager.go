package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var envMap map[string]string

func GetEnvMap() *map[string]string {
	return &envMap
}

func InitializeEnvMap(path string) error {
	env, err := godotenv.Read(path)
	if err != nil {
		return err
	}
	envMap = env
	return nil
}

func LoadEnvFile(path string) error {
	home, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := godotenv.Load(home + path); err != nil {
		return fmt.Errorf("error loading .env file")
	}
	return nil
}
