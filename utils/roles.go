package utils

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Roles *casbin.Enforcer

func init() {
	casbinGormAdapter := gormadapter.NewAdapter("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		*Config.DB.Host,
		*Config.DB.Port,
		*Config.DB.User,
		*Config.DB.Name,
		*Config.DB.Pass,
		*Config.DB.SSL,
	), true)
	Roles = casbin.NewEnforcer("rbac_model.conf", casbinGormAdapter)
}
