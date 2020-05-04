package main

import (
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
	lang, _ := params.Args["lang"].(string)
	// TODO: добавить кол-во найденных результатов
	var files []models.File

	tx := utils.DB.Table("files").
		Select("files.*").
		Joins("left join translations on translations.element_id = files.id").
		Where("translations.lang = ? and translations.element_type = ?", lang, "file")

	if first, ok := params.Args["first"].(int); ok {
		tx = tx.Limit(first)
	}
	if offset, ok := params.Args["offset"].(int); ok {
		tx = tx.Offset(offset)
	}
	if err := tx.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
