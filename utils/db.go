package utils

import (
	"cms-api/config"
	"cms-api/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

var DB *gorm.DB

func init() {
	c := config.Get()
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			c.DB.Host,
			c.DB.Port,
			c.DB.User,
			c.DB.Name,
			c.DB.Pass,
			c.DB.SSL,
		),
	)
	if err != nil {
		log.Panic("failed to connect database")
	}
	DB = db
	DB.LogMode(c.Debug)

	DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.PostMeta{},
		&models.Lang{},
	)
}
