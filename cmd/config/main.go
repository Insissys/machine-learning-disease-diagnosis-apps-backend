package main

import (
	"log"
	"os"

	"github.com/sefazi/machine-learning-disease-diagnosis-apps-backend/internal/config"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := yaml.Marshal(&config.GlobalConfig)
	if err != nil {
		panic(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("failed to get working directory:", err)
	}

	filePath := pwd + "/internal/config/config.yml"
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		panic(err)
	}

	log.Println("YAML data has been written to 'config.yml'")
}
