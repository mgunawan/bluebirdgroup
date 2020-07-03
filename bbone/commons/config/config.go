package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

//Get get string configuration by defined environtment variable key and default value
func Get(envKey, defaultVal string) string {
	val := os.Getenv(envKey)
	if val == "" {
		val = defaultVal
	}
	return val
}

//GetInt get string configuration by defined environtment variable key and default value
func GetInt(envKey string, defaultVal int) int {
	str := os.Getenv(envKey)
	if str == "" {
		return defaultVal
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return defaultVal
	}
	return int(val)
}

//GetA get value in string by "," separator
func GetA(envKey string, defaultVal []string) []string {
	str := os.Getenv(envKey)
	if str == "" {
		return defaultVal
	}
	return strings.Split(str, ",")
}

//GetBool get boolean value configuration
func GetBool(envKey string, defaultValue bool) bool {
	str := os.Getenv(envKey)
	if str == "" {
		return defaultValue
	}
	val, err := strconv.ParseBool(str)
	if err != nil {
		return defaultValue
	}
	return val
}
