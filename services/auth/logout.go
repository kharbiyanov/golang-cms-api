package auth

import (
	"cms-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogoutHandler(c *gin.Context) {
	if err := utils.RemoveToken(c.GetHeader(utils.Config.TokenHeader)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
