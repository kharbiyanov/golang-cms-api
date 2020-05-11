package main

import (
	"cms-api/config"
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

var (
	queryFields    = graphql.Fields{}
	mutationFields = graphql.Fields{}
)

func InitSchema(plugin *models.Plugin) {
	setupQuery()
	setupMutation()

	plugin.QueryFields = queryFields
	plugin.MutationFields = mutationFields
}

func setupQuery() {
	queryFields["userGet"] = &graphql.Field{
		Type:        UserType,
		Description: "Get user by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "user", "get"); err != nil {
				return nil, err
			}
			return GetUser(params)
		},
	}
	queryFields["userList"] = &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "UserList",
			Fields: graphql.Fields{
				"data": &graphql.Field{
					Type: graphql.NewList(UserType),
				},
				"count": &graphql.Field{
					Type: graphql.Int,
				},
			},
		}),
		Description: "Get user list.",
		Args: graphql.FieldConfigArgument{
			"first": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: config.Get().DefaultPostsLimit,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "user", "list"); err != nil {
				return nil, err
			}
			return GetUserList(params)
		},
	}
}

func setupMutation() {
	mutationFields["userUpdate"] = &graphql.Field{
		Type:        UserType,
		Description: "Update user by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"last_name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"first_name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"middle_name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"phone": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "user", "update"); err != nil {
				return nil, err
			}
			return UpdateUser(params)
		},
	}
	mutationFields["userDelete"] = &graphql.Field{
		Type:        UserType,
		Description: "Delete user by id.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "user", "delete"); err != nil {
				return nil, err
			}
			return DeleteUser(params)
		},
	}
}
