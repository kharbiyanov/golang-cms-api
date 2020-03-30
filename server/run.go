package server

import (
	"cms-api/config"
	"cms-api/services/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"time"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Authorization", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//router.Use(cors.Default())

	log.Panic(router.Run(c.ServerAddr))
}
