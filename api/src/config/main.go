package config

import (
	"os"
)

var configMap = map[string]string{}

func init() {
	configMap["githubToken"] = os.Getenv("githubToken")
	configMap["listeningPort"] = os.Getenv("listeningPort")
}

func GetValue(key string) (string) {
	return configMap[key]
}