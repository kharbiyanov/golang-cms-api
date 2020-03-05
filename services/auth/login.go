package auth

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
	UserName *string `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Remember bool    `json:"remember"`
}

func LoginHandler(c *gin.Context) {
	var args Login

	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		UserName: args.UserName,
	}

	if err := utils.DB.Where(&user).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.InvalidLoginErrorMessage})
		return
	}

	if !utils.CheckPasswordHash(args.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.WrongPasswordErrorMessage})
		return
	}

	token, expiresAt, err := utils.GenerateToken(user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "expiresAt": expiresAt})
}
