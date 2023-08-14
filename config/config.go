package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

// Configurations Configurations from config file
type Configurations struct {
	Server  ServerConfigurations   `yaml:"server"`
	Service ServiceConfigurations  `yaml:"service"`
	DB      DatabaseConfigurations `yaml:"database"`
}

// ServerConfigurations Server configurations
type ServerConfigurations struct {
	Port int `yaml:"port"`
}

// ServiceConfigurations Service configurations
type ServiceConfigurations struct {
	Name string `yaml:"name"`
}

// DatabaseConfigurations Database configurations
type DatabaseConfigurations struct {
	DNS     string `yaml:"dns"`
	Pool    int    `yaml:"pool"`
	Timeout int    `yaml:"timeout"`
}

// LoadConfigFile Load the yml config file and environment variables
func LoadConfigFile(log *zap.SugaredLogger) *Configurations {
	configFile, err := os.ReadFile("./resources/config.yml")
	if err != nil {
		log.Fatalf("Failed to load the config file: %v", err)
	}

	// expand environment variables
	confContent := []byte(os.ExpandEnv(string(configFile)))

	configurations := &Configurations{}
	err = yaml.Unmarshal(confContent, configurations)
	if err != nil {
		log.Fatalf("Failed to unmarshall the config file: %v", err)
	}

	return configurations
}
