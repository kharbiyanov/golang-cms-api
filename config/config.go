package config

import (
	"cms-api/models"
	"github.com/spf13/viper"
	"log"
)

type config struct {
	ServerAddr        string                  `json:"serverAddr"`
	RedisAddr         string                  `json:"redisAddr"`
	SiteUrl           string                  `json:"siteUrl"`
	DB                db                      `json:"db"`
	Debug             bool                    `json:"debug"`
	DefaultPostsLimit int                     `json:"defaultPostsLimit"`
	PostTypes         []models.PostConfig     `json:"postTypes"`
	Taxonomies        []models.TaxonomyConfig `json:"taxonomies"`
	SMTP              smtp                    `json:"smtp"`
}

type smtp struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
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
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	viper.SetDefault("siteUrl", "http://localhost:8080")
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
