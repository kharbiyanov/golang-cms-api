package main

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

func CreateLang(params graphql.ResolveParams) (interface{}, error) {
	lang := &models.Lang{
		FullName: params.Args["fullName"].(string),
		Code:     params.Args["code"].(string),
		Flag:     params.Args["flag"].(string),
	}

	if flag, ok := params.Args["flag"].(string); ok {
		lang.Flag = flag
	}

	if def, ok := params.Args["default"].(bool); ok {
		lang.Default = def
	}

	if utils.DB.Where(&models.Lang{Default: true}).First(&models.Lang{}).RecordNotFound() {
		lang.Default = true
	}

	if ! utils.DB.Where(&models.Lang{Code: lang.Code}).First(&lang).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.LangCodeExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if ! utils.DB.Where(&models.Lang{Code: lang.Code}).First(&lang).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.LangCodeExistMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Create(lang).Scan(lang).Error; err != nil {
		return nil, err
	}
	return lang, nil
}
