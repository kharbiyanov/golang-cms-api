package main

import (
	"cms-api/models"
	"github.com/graphql-go/graphql"
)

var (
	mutationFields = graphql.Fields{}
)

func InitSchema(plugin *models.Plugin) {
	setupMutation()

	plugin.MutationFields = mutationFields
}

func setupMutation() {
	mutationFields["authLogin"] = &graphql.Field{
		Type:        TokenType,
		Description: "Login.",
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return Login(params)
		},
	}
	mutationFields["authLogout"] = &graphql.Field{
		Type:        TokenType,
		Description: "Logout.",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return Logout(params)
		},
	}
	mutationFields["authRegister"] = &graphql.Field{
		Type:        UserType,
		Description: "Registration.",
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return Register(params)
		},
	}
	mutationFields["authActivate"] = &graphql.Field{
		Type:        UserType,
		Description: "Activation.",
		Args: graphql.FieldConfigArgument{
			"activation_code": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return Activate(params)
		},
	}
}
