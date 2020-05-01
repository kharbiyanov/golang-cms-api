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
