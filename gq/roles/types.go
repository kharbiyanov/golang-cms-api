package main

import (
	"github.com/graphql-go/graphql"
)

var AccessType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Access",
		Fields: graphql.Fields{
			"object": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"action": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)

var RoleType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Role",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"access": &graphql.Field{
				Type: graphql.NewList(AccessType),
			},
		},
	},
)
