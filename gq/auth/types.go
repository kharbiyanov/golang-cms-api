package main

import (
	"github.com/graphql-go/graphql"
)

var TokenType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Token",
		Fields: graphql.Fields{
			"token": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"expires_at": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AuthUser",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"created_at": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"updated_at": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"middle_name": &graphql.Field{
				Type: graphql.String,
			},
			"avatar": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.Int,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
