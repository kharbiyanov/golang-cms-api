package main

import (
	"cms-api/config"
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

func GetUserList(params graphql.ResolveParams) (interface{}, error) {
	first := config.Get().DefaultPostsLimit
	result := struct {
		Data  []models.User `json:"data"`
		Count int           `json:"count"`
	}{}

	tx := utils.DB.Order("created_at desc")

	if err := tx.Model(&models.User{}).Count(&result.Count).Error; err != nil {
		return nil, err
	}
	if pFirst, ok := params.Args["first"].(int); ok {
		first = pFirst
	}
	tx = tx.Limit(first)
	if offset, ok := params.Args["offset"].(int); ok {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}
