package env

import (
	"os"
	"strconv"
)

func GetString(key string, defaultValue string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return
}

func GetInt(key string, defaultValue int) (value int) {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		value = defaultValue
	}

	return
}
