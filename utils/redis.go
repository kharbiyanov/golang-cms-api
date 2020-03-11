package utils

import (
	"cms-api/config"
	"cms-api/errors"
	"cms-api/models"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

var (
	Redis redis.Conn
)

func init() {
	log.Println(config.Get().RedisAddr)
	conn, err := redis.Dial("tcp", config.Get().RedisAddr)
	if err != nil {
		log.Panic(err)
	}
	Redis = conn
}

func GenerateToken(user models.User) (string, time.Time, error) {
	token := uuid.NewV4().String()
	expDuration := GetTokenExpDuration(false)

	expTime := time.Now().Add(expDuration)

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return token, expTime, err
	}
	if _, err := Redis.Do("SETEX", token, fmt.Sprintf("%.0f", expDuration.Seconds()), jsonUser); err != nil {
		return token, expTime, err
	}
	return token, expTime, nil
}

func CheckToken(p graphql.ResolveParams) (models.User, error) {
	ctx := GetContextFromParams(p)
	user := models.User{}

	token, err := GetBearerToken(ctx.GetHeader("Authorization"))

	if err != nil {
		return user, err
	}

	response, err := Redis.Do("GET", token)
	if err != nil {
		return user, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.InternalServerErrorCode,
		}
	}
	if response == nil {
		return user, &errors.ErrorWithCode{
			Message: errors.InvalidTokenErrorCodeMessage,
			Code:    errors.ForbiddenCode,
		}
	}

	if err := json.Unmarshal([]byte(fmt.Sprintf("%s", response)), &user); err != nil {
		return user, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.InternalServerErrorCode,
		}
	}

	expDuration := GetTokenExpDuration(false)

	if _, err := Redis.Do("EXPIRE", token, fmt.Sprintf("%.0f", expDuration.Seconds())); err != nil {
		return user, &errors.ErrorWithCode{
			Message: err.Error(),
			Code:    errors.InternalServerErrorCode,
		}
	}
	return user, nil
}

func CheckPermission(userName string, object string, action string) error {
	if err := Roles.LoadPolicy(); err != nil {
		return err
	}

	ok := Roles.Enforce(userName, object, action)

	if ! ok {
		return &errors.ErrorWithCode{
			Message: errors.ForbiddenCodeMessage,
			Code:    errors.ForbiddenCode,
		}
	}

	return nil
}

func ValidateUser(p graphql.ResolveParams, object string, action string) error {
	if user, err := CheckToken(p); err != nil {
		return err
	} else {
		if err := CheckPermission(*user.UserName, object, action); err != nil {
			return err
		}
	}

	return nil
}

func RemoveToken(token string) error {
	if _, err := Redis.Do("DEL", token); err != nil {
		return err
	}
	return nil
}
