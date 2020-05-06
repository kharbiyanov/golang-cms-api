package main

import (
	"cms-api/models"
	"cms-api/utils"
	graphqlmultipart "git.osg.uz/kharbiyanov/graphql-multipart-middleware"
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
	queryFields["fileList"] = &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name: "FileList",
			Fields: graphql.Fields{
				"data": &graphql.Field{
					Type: graphql.NewList(FileType),
				},
				"count": &graphql.Field{
					Type: graphql.Int,
				},
			},
		}),
		Description: "Get file list",
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"first": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"offset": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return GetFiles(params)
		},
	}
}

func setupMutation() {
	mutationFields["fileCreate"] = &graphql.Field{
		Type:        FileType,
		Description: "Create new file.",
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"file": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphqlmultipart.Upload),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "file", "create"); err != nil {
				return nil, err
			}

			return UploadFile(params)
		},
	}
	mutationFields["fileUpdate"] = &graphql.Field{
		Type:        FileType,
		Description: "Update file.",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "file", "update"); err != nil {
				return nil, err
			}

			return UpdateFile(params)
		},
	}
}
