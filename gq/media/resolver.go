package main

import (
	"cms-api/utils"
	"github.com/graphql-go/graphql"
	"log"
	"mime/multipart"
)

func UploadMedia(params graphql.ResolveParams) (interface{}, error) {
	fh := params.Args["file"].(*multipart.FileHeader)
	uploadedFile, err := utils.SaveFile(fh)

	if err != nil {
		return nil, err
	}

	log.Printf("%+v", uploadedFile)

	return nil, nil
}
