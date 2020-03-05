package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetTokenExpDuration(remember bool) time.Duration {
	expDuration := time.Hour * 48 // 2 days

	if remember {
		expDuration = time.Hour * 336 // 14 days
	}

	return expDuration
}

func GetContextFromParams(p graphql.ResolveParams) *gin.Context {
	rootValue := p.Info.RootValue.(map[string]interface{})
	return rootValue["context"].(*gin.Context)
}
