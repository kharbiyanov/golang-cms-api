package server

import (
	"cms-api/config"
	"cms-api/services/auth"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	SetupPlugins()

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
		g.GET("", GetHandler())
	}
	g.POST("", GetHandler())

	log.Panic(r.Run(c.ServerAddr))
}
