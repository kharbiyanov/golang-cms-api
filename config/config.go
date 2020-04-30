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
	MimeTypes         []string                `json:"mimeTypes"`
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

	viper.SetDefault("mimeTypes", mimeTypes)

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

var mimeTypes = []string{
	//.jpg .jpeg
	"image/jpeg",
	"image/pjpeg",
	//.png
	"image/png",
	//.gif
	"image/gif",
	//.ico
	"image/x-icon",
	//.pdf
	"application/pdf",
	//.doc
	"application/msword",
	//.docx
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	//.ppt
	"application/mspowerpoint",
	"application/powerpoint",
	"application/vnd.ms-powerpoint",
	"application/x-mspowerpoint",
	//.pptx
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	//.pps
	"application/mspowerpoint",
	"application/vnd.ms-powerpoint",
	//.ppsx
	"application/vnd.openxmlformats-officedocument.presentationml.slideshow",
	//.odt
	"application/vnd.oasis.opendocument.text",
	//.xls
	"application/excel",
	"application/vnd.ms-excel",
	"application/x-excel",
	"application/x-msexcel",
	//.xlsx
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	//.psd
	"application/octet-stream",
	//.mp3
	"audio/mpeg3",
	"audio/x-mpeg-3",
	"video/mpeg",
	"video/x-mpeg",
	//.m4a
	"audio/m4a",
	//.ogg
	"audio/ogg",
	//.wav
	"audio/wav",
	"audio/x-wav",
	//.mp4
	"video/mp4",
	//.m4v
	"video/x-m4v",
	//.mov
	"video/quicktime",
	//.wmv
	"video/x-ms-asf",
	"video/x-ms-wmv",
	//.avi
	"application/x-troff-msvideo",
	"video/avi",
	"video/msvideo",
	"video/x-msvideo",
	//.mpg
	"audio/mpeg",
	"video/mpeg",
	//.ogv
	"video/ogg",
	//.3gp
	"video/3gpp",
	"audio/3gpp",
	//.3g2
	"video/3gpp2",
	"audio/3gpp2",
	//.tar
	"application/x-tar",
	//.zip
	"application/zip",
	//.gz .gzip
	"application/x-zip",
	//.rar
	"application/rar",
	//.7z
	"application/x-7z-compressed",
}
