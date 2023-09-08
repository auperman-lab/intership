package db

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	API struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"api"`
	Environment string `yaml:"environment"`
	MongoDB     struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"mongodb"`
	PostgreSQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"postgresql"`
}

func InitMongo() {
	configData, err := os.ReadFile("/Users/nmacrii/Desktop/intership/configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config

	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Access configuration variables based on the environment
	if config.Environment == "dev" {
		mongodbConfig := config.MongoDB
		// Use MongoDB configurations for the development environment
		fmt.Printf("Connecting to MongoDB: %s:%d/%s\n", mongodbConfig.Host, mongodbConfig.Port, mongodbConfig.Database)
	} else {
		postgresqlConfig := config.PostgreSQL
		// Use PostgreSQL configurations for production
		fmt.Printf("Connecting to PostgreSQL: %s:%d/%s\n", postgresqlConfig.Host, postgresqlConfig.Port, postgresqlConfig.Database)
	}
}
