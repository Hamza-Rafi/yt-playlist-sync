package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/goccy/go-yaml"
)

var configFilePath string = "./config.yml"

func main() {
	// load config from config file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalln(err)
	}

	spew.Dump(config)
}
