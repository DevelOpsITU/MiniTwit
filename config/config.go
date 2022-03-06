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
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	Database struct {
		Type             string `yaml:"type" envconfig:"DB_TYPE"`
		ConnectionString string `yaml:"connectionString" envconfig:"DB_CONNECTION_STRING"`
	} `yaml:"database"`
	Development struct {
		GenerateMockData bool `yaml:"generateMockData" envconfig:"DEVELOPMENT_GENERATE_MOCK_DATA"`
	} `yaml:"development"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func SetupConfig() {
	readFile(&config)
	readEnviorement(&config) // overwrite dev enviorement at production
}

func SetupTestConfig() {
	readTestConfigFile(&config)
	readEnviorement(&config) // overwrite dev enviorement at production
}

func GetConfig() Configuration {
	return config
}

func readFile(config *Configuration) {
	filePath, _ := filepath.Abs("./config.yml")
	f, err := os.Open(filePath)
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

func readTestConfigFile(config *Configuration) {
	filePath, _ := filepath.Abs("../../config.yml")
	f, err := os.Open(filePath)
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
