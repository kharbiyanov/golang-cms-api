package main

import (
	"cms-api/config"
	"cms-api/errors"
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

	if err := utils.Roles.LoadPolicy(); err != nil {
		return user, err
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

func UpdateUser(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	fields := make(map[string]interface{})

	var user models.User

	if val, ok := params.Args["last_name"].(string); ok {
		fields["last_name"] = val
	}
	if val, ok := params.Args["first_name"].(string); ok {
		fields["first_name"] = val
	}
	if val, ok := params.Args["middle_name"].(string); ok {
		fields["middle_name"] = val
	}
	if val, ok := params.Args["phone"].(int); ok {
		fields["phone"] = val
	}

	user.ID = id

	if err := utils.DB.Model(&user).Updates(fields).Scan(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func DeleteUser(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)

	user := models.User{}

	if utils.DB.First(&user, id).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.UserNotFoundMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	return user, utils.DB.Delete(user).Error
}

func GetRoles(params graphql.ResolveParams) (interface{}, error) {
	user, userExist := params.Source.(models.User)

	if !userExist {
		return nil, nil
	}

	return utils.Roles.GetRolesForUser(user.UserName)
}
