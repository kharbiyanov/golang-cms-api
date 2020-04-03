package main

import (
	"cms-api/errors"
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
)

func CreateLang(params graphql.ResolveParams) (interface{}, error) {
	lang := &models.Lang{
		FullName: params.Args["full_name"].(string),
		Code:     params.Args["code"].(string),
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

	if !utils.DB.Where(&models.Lang{Code: lang.Code}).First(&lang).RecordNotFound() {
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

func DeleteLang(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)

	lang := models.Lang{}

	if utils.DB.First(&lang, id).RecordNotFound() {
		return nil, &errors.ErrorWithCode{
			Message: errors.LangNotFoundMessage,
			Code:    errors.InvalidParamsCode,
		}
	}

	if err := utils.DB.Delete(lang).Error; err != nil {
		return nil, err
	}

	if lang.Default {
		newDefaultLang := models.Lang{}
		if utils.DB.First(&newDefaultLang).RecordNotFound() {
			return nil, &errors.ErrorWithCode{
				Message: errors.LangNotFoundMessage,
				Code:    errors.InvalidParamsCode,
			}
		}
		if err := utils.DB.Model(&newDefaultLang).Updates(models.Lang{Default: true}).Error; err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func GetLangList() (interface{}, error) {
	var langList []models.Lang

	if err := utils.DB.Find(&langList).Error; err != nil {
		return nil, err
	}
	return langList, nil
}

func UpdateLang(params graphql.ResolveParams) (interface{}, error) {
	id, _ := params.Args["id"].(int)
	def, defExist := params.Args["default"].(bool)
	fields := make(map[string]interface{})

	var lang models.Lang

	if fullName, ok := params.Args["full_name"].(string); ok {
		fields["full_name"] = fullName
	}
	if flag, ok := params.Args["flag"].(string); ok {
		fields["flag"] = flag
	}
	if code, ok := params.Args["code"].(string); ok {
		fields["code"] = code
		if !utils.DB.Where(&models.Lang{Code: code}).Not(&models.Lang{ID: id}).First(&lang).RecordNotFound() {
			return nil, &errors.ErrorWithCode{
				Message: errors.LangCodeExistMessage,
				Code:    errors.InvalidParamsCode,
			}
		}
	}
	if defExist {
		fields["default"] = def
		if def {
			utils.DB.Model(&models.Lang{}).Where(&models.Lang{Default: true}).Update("default", false)
		}
	}

	lang.ID = id

	if err := utils.DB.Model(&lang).Updates(fields).Scan(&lang).Error; err != nil {
		return nil, err
	}
	return lang, nil
}
