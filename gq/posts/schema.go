package main

import (
	"cms-api/config"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
)

var (
	queryFields    = graphql.Fields{}
	mutationFields = graphql.Fields{}
)

func InitSchema(plugin *models.Plugin) {
	for _, postConfig := range config.Get().PostTypes {
		postType := GetPostType(postConfig)
		setupPostsQuery(postType, postConfig)
		setupPostsMutation(postType, postConfig)
	}
	setupMetaQuery()
	setupMetaMutations()

	plugin.QueryFields = queryFields
	plugin.MutationFields = mutationFields
}

func setupPostsQuery(postType *graphql.Object, postConfig models.PostConfig) {
	queryFields[fmt.Sprintf("%sGet", postConfig.Slug)] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Get %s by id.", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return GetPost(params, postConfig)
		},
	}

	queryFields[fmt.Sprintf("%sList", postConfig.Slug)] = &graphql.Field{
		Type:        graphql.NewList(postType),
		Description: fmt.Sprintf("Get %s list.", postConfig.Slug),
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
			return GetPosts(params, postConfig)
		},
	}
}

func setupPostsMutation(postType *graphql.Object, postConfig models.PostConfig) {
	mutationFields[fmt.Sprintf("%sCreate", postConfig.Slug)] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Create new %s.", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
			"lang": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"content": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"excerpt": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"status": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"slug": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, postConfig.PluralSlug, "create"); err != nil {
				return nil, err
			}

			return CreatePost(params, postConfig)
		},
	}

	mutationFields[fmt.Sprintf("%sUpdate", postConfig.Slug)] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Update %s by id.", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"title": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"content": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"excerpt": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"status": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"slug": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, postConfig.PluralSlug, "update"); err != nil {
				return nil, err
			}

			return UpdatePost(params, postConfig)
		},
	}

	mutationFields[fmt.Sprintf("%sDelete", postConfig.Slug)] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Delete %s by id.", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, postConfig.PluralSlug, "delete"); err != nil {
				return nil, err
			}

			return DeletePost(params, postConfig)
		},
	}
}

func setupMetaQuery() {
	queryFields["metaGet"] = &graphql.Field{
		Type:        graphql.NewList(MetaType),
		Description: "Get meta by post_id and meta_keys.",
		Args: graphql.FieldConfigArgument{
			"post_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"keys": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return GetMeta(params)
		},
	}
}

func setupMetaMutations() {
	mutationFields["metaUpdate"] = &graphql.Field{
		Type:        MetaType,
		Description: "Update meta by post_id and meta_key. If the meta field for the post does not exist, it will be added.",
		Args: graphql.FieldConfigArgument{
			"post_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"value": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "meta", "update"); err != nil {
				return nil, err
			}

			return UpdateMeta(params)
		},
	}

	mutationFields["metaDelete"] = &graphql.Field{
		Type:        MetaType,
		Description: "Delete %s by post_id and meta_key.",
		Args: graphql.FieldConfigArgument{
			"post_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"key": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			if err := utils.ValidateUser(params, "meta", "delete"); err != nil {
				return nil, err
			}

			return DeleteMeta(params)
		},
	}
}
