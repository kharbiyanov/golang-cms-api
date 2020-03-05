package utils

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	ServerAddr  *string `json:"SERVER_ADDR"`
	RedisAddr   *string `json:"REDIS_ADDR"`
	DB          db      `json:"DB"`
	Debug       *bool   `json:"DEBUG"`
	TokenHeader string  `json:"-"`
}

type db struct {
	Host *string `json:"HOST"`
	Port *string `json:"PORT"`
	SSL  *string `json:"SSL_MODE"`
	Name *string `json:"NAME"`
	User *string `json:"USER"`
	Pass *string `json:"PASS"`
}

var Config = &config{
	TokenHeader: "Auth-Token",
}

func init() {
	fileName := "config.json"
	configFile, err := os.Open(fileName)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&Config); err != nil {
		panic(err)
	}

	if Config.Debug == nil {
		debug := false
		Config.Debug = &debug
	}

	if Config.ServerAddr == nil {
		log.Fatal("Env SERVER_ADDR does not exist")
	}

	if Config.RedisAddr == nil {
		log.Fatal("Env REDIS_ADDR does not exist")
	}

	if Config.DB.Host == nil {
		log.Fatal("Env DB_HOST does not exist")
	}

	if Config.DB.Port == nil {
		log.Fatal("Env DB_PORT does not exist")
	}

	if Config.DB.SSL == nil {
		log.Fatal("Env DB_SSL_MODE does not exist")
	}

	if Config.DB.Name == nil {
		log.Fatal("Env DB_NAME does not exist")
	}

	if Config.DB.User == nil {
		log.Fatal("Env DB_USER does not exist")
	}

	if Config.DB.Pass == nil {
		log.Fatal("Env DB_PASS does not exist")
	}
}
