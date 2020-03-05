package server

import (
	"cms-api/graphql"
	"cms-api/services/auth"
	"cms-api/utils"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	if ! *utils.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	a := r.Group("/auth")
	a.POST("/login", auth.LoginHandler)
	a.POST("/logout", auth.LogoutHandler)

	g := r.Group("/graphql")
	if *utils.Config.Debug {
		g.GET("", graphql.GetHandler())
	}
	g.POST("", graphql.GetHandler())

	panic(r.Run(*utils.Config.ServerAddr))
}
