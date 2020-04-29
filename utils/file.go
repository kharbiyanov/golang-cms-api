package utils

import (
	"cms-api/models"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
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

	uploadedFile := &models.UploadedFile{}
	setFileInfo(fh, uploadedFile)

	if _, err := os.Stat(uploadedFile.Path); err != nil {
		if os.IsNotExist(err) {
			if errDir := os.MkdirAll(uploadedFile.Path, 0755); errDir != nil {
				return nil, err
			}
		}
	}

	f, err := os.OpenFile(uploadedFile.GetPath(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return nil, err
	}

	mime, err := mimetype.DetectFile(uploadedFile.GetPath())

	if err != nil {
		return nil, err
	}

	uploadedFile.MimeType = mime.String()

	return uploadedFile, err
}

func setFileInfo(fh *multipart.FileHeader, file *models.UploadedFile) {
	rand.Seed(time.Now().UnixNano())
	currentTime := time.Now()

	ext := filepath.Ext(fh.Filename)
	wExt := strings.TrimSuffix(fh.Filename, ext)
	name := fmt.Sprintf("%s-%d%s", wExt, rand.Intn(100), ext)
	path := currentTime.Format("files/2006/01")

	file.Name = name
	file.Path = path
	file.Original = fh.Filename
}
