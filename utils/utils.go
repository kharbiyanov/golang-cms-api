package utils

import (
	"bytes"
	"cms-api/config"
	"cms-api/errors"
	"cms-api/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"text/template"
	"time"
)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(b), err
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

func GetBearerToken(bearer string) (string, error) {
	token := strings.Split(bearer, "Bearer ")

	if len(token) == 1 {
		return "", &errors.ErrorWithCode{
			Message: errors.ForbiddenCodeMessage,
			Code:    errors.ForbiddenCode,
		}
	}

	return token[1], nil
}

func GetPermalink(object interface{}) (string, error) {
	c := config.Get()
	tmpl := template.New("Permalink")
	buf := &bytes.Buffer{}
	var permalink models.Permalink
	tmplText := ""
	switch obj := object.(type) {
	case models.Post:
		postConfig := GetPostConfig(obj.Type)
		permalink = models.Permalink{
			Id:     obj.ID,
			Object: obj.Type,
			Slug:   obj.Slug,
		}
		tmplText = postConfig.SingleUrl
	case models.Term:
		permalink = models.Permalink{
			Id:     obj.ID,
			Object: obj.Taxonomy,
			Slug:   obj.Slug,
		}
		tmplText = "/{{.Object}}/{{.Slug}}/"
	}
	if _, err := tmpl.Parse(tmplText); err != nil {
		return "", err
	}
	if err := tmpl.Execute(buf, permalink); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", c.SiteUrl, buf.String()), nil
}

func GetPostConfig(postType string) models.PostConfig {
	c := config.Get()
	var postConfig models.PostConfig
	for _, conf := range c.PostTypes {
		if postType == conf.Type {
			postConfig = conf
			break
		}
	}
	return postConfig
}
