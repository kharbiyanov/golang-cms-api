package run

import (
	"cms-api/config"
	"cms-api/graphql"
	"cms-api/services/auth"
	"github.com/gin-gonic/gin"
	"log"
)

func Server() {
	c := config.Get()
	if ! c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	a := r.Group("/auth")
	a.POST("/login", auth.LoginHandler)
	a.POST("/logout", auth.LogoutHandler)

	g := r.Group("/graphql")
	if c.Debug {
		g.GET("", graphql.GetHandler())
	}
	g.POST("", graphql.GetHandler())

	log.Panic(r.Run(c.ServerAddr))
}
