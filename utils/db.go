package utils

import (
	"cms-api/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			*Config.DB.Host,
			*Config.DB.Port,
			*Config.DB.User,
			*Config.DB.Name,
			*Config.DB.Pass,
			*Config.DB.SSL,
		),
	)
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	DB.LogMode(*Config.Debug)

	DB.AutoMigrate(&models.User{})
}
