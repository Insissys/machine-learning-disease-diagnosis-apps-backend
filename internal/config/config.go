package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Config struct {
		Server   Server   `yaml:"server"`
		Client   Client   `yaml:"client"`
		Database Database `yaml:"database"`
	} `yaml:"config"`
}

var GlobalConfig Configuration

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Client struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Timeout  int32  `yaml:"timeout"`
}

type Database struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}

func LoadConfig() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to get working directory:", err)
	}

	filePath := pwd + "/internal/config/config.yml"
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read config file at %s: %v", filePath, err)
	}

	err = yaml.Unmarshal(file, &GlobalConfig)
	if err != nil {
		log.Fatalln("failed to unmarshal config:", err)
	}

	log.Printf("Config loaded: %+v\n", GlobalConfig)
}

func GetConfig() *Configuration {
	return &GlobalConfig
}
