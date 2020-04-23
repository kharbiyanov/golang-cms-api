package main

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

func Login(params graphql.ResolveParams) (interface{}, error) {
	userName, _ := params.Args["username"].(string)
	password, _ := params.Args["password"].(string)

	user := models.User{
		UserName: userName,
	}

	if err := utils.DB.Where(&user).First(&user).Error; err != nil {
		return nil, &errors.ErrorWithCode{
			Message: errors.InvalidLoginErrorMessage,
			Code:    errors.ForbiddenCode,
		}
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, &errors.ErrorWithCode{
			Message: errors.WrongPasswordErrorMessage,
			Code:    errors.ForbiddenCode,
		}
	}

	token, expiresAt, err := utils.GenerateToken(user)

	if err != nil {
		return nil, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.ForbiddenCode,
		}
	}

	return models.Token{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func Logout(params graphql.ResolveParams) (interface{}, error) {
	ctx := utils.GetContextFromParams(params)
	token, err := utils.GetBearerToken(ctx.GetHeader("Authorization"))

	if err != nil {
		return nil, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.InternalServerErrorCode,
		}
	}

	if err := utils.RemoveToken(token); err != nil {
		return nil, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.InternalServerErrorCode,
		}
	}

	return models.Token{
		Token: token,
	}, nil
}

func Register(params graphql.ResolveParams) (interface{}, error) {
	userName, _ := params.Args["username"].(string)
	password, _ := params.Args["password"].(string)

	hashPassword, err := utils.HashPassword(password)

	if err != nil {
		return nil, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.InternalServerErrorCode,
		}
	}

	user := &models.User{
		UserName: userName,
		Password: hashPassword,
	}

	if !utils.DB.Where(&models.User{UserName: userName}).First(&user).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.UserNameExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Create(user).Scan(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
