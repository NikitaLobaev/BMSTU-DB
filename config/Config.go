package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port uint   `json:"port"`
	} `json:"server"`
	Database struct {
		Host     string `json:"host"`
		Port     uint   `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
}

func NewConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		return nil, err
	}

	config := new(Config)
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
