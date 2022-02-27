package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var config Configuration // global configuration

type Configuration struct {
	Server struct {
		Host string `yaml:"host", envconfig:"SERVER_HOST"`
		Port string `yaml:"port", envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		ConnectionString string `yaml:"connectionstring", envconfig:"DB_CONNECTION_STRING`
	} `yaml:"database"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func SetupConfig() {
	readFile(&config)
	readEnviorement(&config) // overwrite dev enviorement at production
}

func GetConfig() Configuration {
	return config
}

func readFile(config *Configuration) {
	filepath, _ := filepath.Abs("src/config.yml")
	f, err := os.Open(filepath)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		processError(err)
	}
}

func readEnviorement(config *Configuration) {
	err := envconfig.Process("", config)
	if err != nil {
		processError(err)
	}
}
