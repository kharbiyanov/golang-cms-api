package main

import (
	"cms-api/config"
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"fmt"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/gomail.v2"
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

	if user.State == 0 {
		return nil, &errors.ErrorWithCode{
			Message: errors.UserNotActivatedMessage,
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
	email, _ := params.Args["email"].(string)
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

	if !utils.DB.Where(&models.User{UserName: userName, State: 1}).First(&user).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.UserNameExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if !utils.DB.Where(&models.User{Email: email, State: 1}).First(&user).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.EmailExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	// TODO: Check Inactive Users

	user.State = 0
	user.ActivationCode = uuid.NewV4().String()
	user.Email = email

	if err := utils.DB.Create(user).Scan(user).Error; err != nil {
		return nil, err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Get().SMTP.UserName)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Your Activation Code!")
	m.SetBody("text/html", fmt.Sprintf("Activations Code: <code>%s</code>", user.ActivationCode))

	if err := utils.SendMail(m); err != nil {
		return nil, err
	}

	return user, nil
}

func Activate(params graphql.ResolveParams) (interface{}, error) {
	activationCode, _ := params.Args["activation_code"].(string)

	user := &models.User{
		ActivationCode: activationCode,
		State:          0,
	}

	if utils.DB.Where(&user).First(&user).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.ActivationCodeNotFoundMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	user.State = 1
	user.ActivationCode = ""

	if err := utils.DB.Save(&user).Find(&user).Error; err != nil {
		return nil, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.InternalServerErrorCode,
		}
	}

	return user, nil
}
