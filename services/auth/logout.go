package auth

import (
	"cms-api/config"
	"cms-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogoutHandler(c *gin.Context) {
	if err := utils.RemoveToken(c.GetHeader(config.Get().AuthTokenHeader)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
