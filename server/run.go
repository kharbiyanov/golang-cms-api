package server

import (
	"cms-api/config"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	graphqlGroup := router.Group("/graphql/")
	if c.Debug {
		graphqlGroup.GET("/", GetHandler())
	}
	graphqlGroup.POST("/", GetHandler())

	log.Panic(router.Run(c.ServerAddr))
}
