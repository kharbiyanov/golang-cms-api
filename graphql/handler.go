package graphql

import (
	"cms-api/graphql/posts"
	"cms-api/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
)

func GetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := handler.New(&handler.Config{
			Schema:     getSchema(),
			Pretty:     true,
			Playground: *utils.Config.Debug,
			RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
				return map[string]interface{}{
					"context": c,
				}
			},
		})
		h.ServeHTTP(c.Writer, c.Request)
	}
}

var postTypes = []posts.PostConfig{
	{
		Slug:       "page",
		PluralSlug: "pages",
	},
	{
		Slug:       "news",
		PluralSlug: "news",
	},
}

func getSchema() *graphql.Schema {

	query := graphql.Fields{}
	mutation := graphql.Fields{}

	for _, postConfig := range postTypes {
		queryFields, mutationFields := posts.GetSchemaConfig(postConfig)
		for key, val := range queryFields {
			query[key] = val
		}
		for key, val := range mutationFields {
			mutation[key] = val
		}
	}

	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    getQuery(query),
			Mutation: getMutation(mutation),
		},
	)
	return &schema
}

func getQuery(postsQuery graphql.Fields) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: postsQuery,
		},
	)
}

func getMutation(postsMutation graphql.Fields) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: postsMutation,
	})
}
