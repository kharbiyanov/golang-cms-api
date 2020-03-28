package server

import (
	"cms-api/config"
	"cms-api/services/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func Run() {
	SetupPlugins()

	c := config.Get()
	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	authGroup := router.Group("/auth/")
	authGroup.POST("/login/", auth.LoginHandler)
	authGroup.POST("/logout/", auth.LogoutHandler)

	graphqlGroup := router.Group("/graphql/")
	if c.Debug {
		graphqlGroup.GET("/", GetHandler())
	}
	graphqlGroup.POST("/", GetHandler())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000"}
	router.Use(cors.New(corsConfig))

	log.Panic(router.Run(c.ServerAddr))
}
