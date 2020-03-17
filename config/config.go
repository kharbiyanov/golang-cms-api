package config

import (
	"cms-api/models"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type config struct {
	ServerAddr        string              `json:"serverAddr"`
	RedisAddr         string              `json:"redisAddr"`
	DB                db                  `json:"db"`
	Debug             bool                `json:"debug"`
	DefaultPostsLimit int                 `json:"defaultPostsLimit"`
	PostTypes         []models.PostConfig `json:"postTypes"`
}

type db struct {
	Host string `json:"host"`
	Port string `json:"port"`
	SSL  string `json:"ssl"`
	Name string `json:"name"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

var c = config{}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	viper.SetDefault("serverAddr", ":3000")
	viper.SetDefault("redisAddr", ":6379")
	viper.SetDefault("debug", false)

	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.ssl", "disable")

	viper.SetDefault("defaultPostsLimit", 10)

	if err := viper.Unmarshal(&c); err != nil {
		log.Panicf("Unable to decode into struct, %v", err)
	}

	log.Println(c)

	if c.DB.Name == "" {
		log.Panic("db.name is not specified in config file")
	}
	if c.DB.User == "" {
		log.Panic("db.user is not specified in config file")
	}
	if c.DB.Pass == "" {
		log.Panic("db.pass is not specified in config file")
	}
}

func Get() config {
	return c
}
