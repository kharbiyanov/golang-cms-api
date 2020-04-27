package utils

import (
	"cms-api/models"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFile(fh *multipart.FileHeader) (*models.UploadedFile, error) {
	file, err := fh.Open()

	if err != nil {
		return nil, err
	}

	defer file.Close()

	uploadedFile := getFileName(fh.Filename)

	if _, err := os.Stat(uploadedFile.Path); err != nil {
		if os.IsNotExist(err) {
			if errDir := os.MkdirAll(uploadedFile.Path, 0755); errDir != nil {
				return nil, err
			}
		}
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", uploadedFile.Path, uploadedFile.Name), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return nil, err
	}

	return uploadedFile, err
}

func getFileName(filename string) *models.UploadedFile {
	rand.Seed(time.Now().UnixNano())
	currentTime := time.Now()

	ext := filepath.Ext(filename)
	wExt := strings.TrimSuffix(filename, ext)
	name := fmt.Sprintf("%s-%d%s", wExt, rand.Intn(100), ext)
	path := currentTime.Format("uploads/2006/01/02")

	return &models.UploadedFile{
		Name:     name,
		Path:     path,
		Original: filename,
	}
}
