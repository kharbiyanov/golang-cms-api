package main

import (
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

func GetUser(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)

	var user = models.User{}
	user.ID = id

	if err := utils.DB.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
