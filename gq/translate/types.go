package main

import (
	"github.com/graphql-go/graphql"
)

var LangType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Lang",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"full_name": &graphql.Field{
				Type: graphql.String,
			},
			"code": &graphql.Field{
				Type: graphql.String,
			},
			"flag": &graphql.Field{
				Type: graphql.String,
			},
			"default": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)
