package utils

import (
	"cms-api/config"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Roles *casbin.Enforcer

func init() {
	c := config.Get()
	casbinGormAdapter := gormadapter.NewAdapter("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Name,
		c.DB.Pass,
		c.DB.SSL,
	), true)
	Roles = casbin.NewEnforcer("rbac_model.conf", casbinGormAdapter)
}
