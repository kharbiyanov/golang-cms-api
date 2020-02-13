package auth

import (
	"cms-api/utils"
	"errors"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type User struct {
	ID        int64  `json:"id"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token"`
}

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"userName": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var user User

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type:        userType,
				Description: "Login in dashboard",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"remember": &graphql.ArgumentConfig{
						Type:         graphql.Boolean,
						DefaultValue: false,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					username, usernameOK := params.Args["username"].(string)
					password, passwordOK := params.Args["password"].(string)
					_, rememberOK := params.Args["remember"].(bool)
					sessionToken := uuid.NewV4().String()

					if usernameOK && passwordOK && rememberOK {
						if username != "kharbiyanov" {
							return user, utils.GQError{
								Err:  errors.New("invalid login"),
								Code: 13,
							}
							//return user, gqlerrors.
						}
						if password != "123321" {
							return user, errors.New("wrong password")
						}

						user.ID = 7
						user.UserName = username
						user.FirstName = "Marat"
						user.LastName = "Kharbiyanov"
						user.Token = sessionToken

						_, err := utils.Redis.Do("SETEX", sessionToken, "120", user)
						if err != nil {
							return "", errors.New(http.StatusText(http.StatusInternalServerError))
						}
					}
					return user, nil
				},
			},

			"logout": &graphql.Field{
				Type:        userType,
				Description: "Logout from admin",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					token := params.Context.Value("authToken")
					_, err := utils.Redis.Do("DEL", token)
					if err != nil {
						return "", errors.New(http.StatusText(http.StatusInternalServerError))
					}

					return nil, nil
				},
			},
		},
	},
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)
