package posts

import (
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
)

func GetSchemaConfig(postConfig PostConfig) (graphql.Fields, graphql.Fields) {
	postType := GetPostType(postConfig)
	query := GetQuery(postType, postConfig)
	mutation := GetMutation(postType, postConfig)
	return query, mutation
}

func GetQuery(postType *graphql.Object, postConfig PostConfig) graphql.Fields {
	fields := graphql.Fields{}

	postGet := postConfig.Slug + "Get"
	postList := postConfig.Slug + "List"

	fields[postGet] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Get %s by id", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return GetPost(p, postConfig)
		},
	}

	fields[postList] = &graphql.Field{
		Type:        graphql.NewList(postType),
		Description: fmt.Sprintf("Get %s list", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
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

	return fields
}

func GetMutation(postType *graphql.Object, postConfig PostConfig) graphql.Fields {
	fields := graphql.Fields{}

	postCreate := postConfig.Slug + "Create"
	postUpdate := postConfig.Slug + "Update"
	postDelete := postConfig.Slug + "Delete"

	fields[postCreate] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Create new %s", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
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

	fields[postUpdate] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Update %s by id", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"info": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"price": &graphql.ArgumentConfig{
				Type: graphql.Float,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			post := Post{}

			return post, nil
		},
	}

	fields[postDelete] = &graphql.Field{
		Type:        postType,
		Description: fmt.Sprintf("Delete %s by id", postConfig.Slug),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			post := Post{}

			return post, nil
		},
	}

	return fields
}
