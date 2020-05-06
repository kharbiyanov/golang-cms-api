package main

import (
	"cms-api/config"
	"cms-api/models"
	"cms-api/utils"
	"github.com/graphql-go/graphql"
	"mime/multipart"
)

func UploadFile(params graphql.ResolveParams) (interface{}, error) {
	lang, _ := params.Args["lang"].(string)
	fh, _ := params.Args["file"].(*multipart.FileHeader)
	authUser := utils.GetAuthUser(params)

	uploadedFile := utils.UploadedFile{
		Header: fh,
	}
	if err := uploadedFile.Save(); err != nil {
		return nil, err
	}

	file := &models.File{
		Title:    uploadedFile.Original,
		MimeType: uploadedFile.MimeType.String(),
		File:     uploadedFile.GetPath(),
		AuthorID: authUser.ID,
	}

	if err := utils.DB.Create(file).Scan(file).Error; err != nil {
		return nil, err
	}

	translation := models.Translation{
		ElementType: "file",
		ElementID:   file.ID,
		Lang:        lang,
	}

	if err := utils.DB.Create(&translation).Scan(&translation).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func GetFiles(params graphql.ResolveParams) (interface{}, error) {
	first := config.Get().DefaultPostsLimit
	result := struct {
		Data  []models.File `json:"data"`
		Count int           `json:"count"`
	}{}

	lang, _ := params.Args["lang"].(string)

	tx := utils.DB.Table("files").
		Select("files.*").
		Joins("left join translations on translations.element_id = files.id").
		Where("translations.lang = ? and translations.element_type = ?", lang, "file")

	if err := tx.Count(&result.Count).Error; err != nil {
		return nil, err
	}
	if pFirst, ok := params.Args["first"].(int); ok {
		first = pFirst
	}
	tx = tx.Limit(first)
	if offset, ok := params.Args["offset"].(int); ok {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&result.Data).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func GetTranslationsInFile(params graphql.ResolveParams) (interface{}, error) {
	file, fileExist := params.Source.(models.File)

	if !fileExist {
		return nil, nil
	}

	innerSql := utils.DB.Table("translations t").Select("t.group_id").Where("t.element_id = ? AND t.element_type = ?", file.ID, "file").QueryExpr()

	tx := utils.DB.Table("translations").
		Select("translations.*").
		Joins("LEFT JOIN lang l ON translations.lang = l.code").
		Where("l.deleted_at IS NULL AND translations.element_type = ? AND translations.group_id = (?)", "file", innerSql)

	if err := tx.Find(&file.Translations).Error; err != nil {
		return nil, err
	}

	return file.Translations, nil
}
