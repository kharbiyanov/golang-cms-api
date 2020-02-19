package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func init() {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			Config.DB.Host,
			Config.DB.Port,
			Config.DB.User,
			Config.DB.Name,
			Config.DB.Pass,
			Config.DB.SSL,
		),
	)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	DB = db
	DB.LogMode(true)

	//db.Create(&Product{Code: "L1212", Price: 1000})

	var product Product
	//db.First(&product, "code = ?", "L1212")
	db.Model(&product).Where("code = ?", "L1212").Update("Price", 5000)

}
