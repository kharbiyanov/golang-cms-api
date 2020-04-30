package utils

import (
	"cms-api/config"
	"errors"
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

type UploadedFile struct {
	Header   *multipart.FileHeader
	Name     string
	Path     string
	Original string
	MimeType *mimetype.MIME
}

func (f *UploadedFile) Save() error {
	file, err := f.Header.Open()

	if err != nil {
		return err
	}

	defer file.Close()

	if err := f.checkMimeType(); err != nil {
		return err
	}

	f.setFileInfo()

	if _, err := os.Stat(f.Path); err != nil {
		if os.IsNotExist(err) {
			if errDir := os.MkdirAll(f.Path, 0755); errDir != nil {
				return err
			}
		}
	}

	ff, err := os.OpenFile(f.GetPath(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer ff.Close()

	if _, err := io.Copy(ff, file); err != nil {
		return err
	}

	return nil
}

func (f *UploadedFile) GetPath() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

func (f *UploadedFile) setFileInfo() {
	rand.Seed(time.Now().UnixNano())
	currentTime := time.Now()

	ext := filepath.Ext(f.Header.Filename)
	wExt := strings.TrimSuffix(f.Header.Filename, ext)
	name := fmt.Sprintf("%s-%d%s", wExt, rand.Intn(100), ext)
	path := currentTime.Format("files/2006/01")

	f.Name = name
	f.Path = path
	f.Original = f.Header.Filename
}

func (f *UploadedFile) checkMimeType() error {
	reader, _ := f.Header.Open()
	mime, err := mimetype.DetectReader(reader)
	defer reader.Close()
	if err != nil {
		return err
	}
	f.MimeType = mime
	if !mimetype.EqualsAny(mime.String(), config.Get().MimeTypes...) {
		return errors.New(fmt.Sprintf("%s is now allowed", mime))
	}
	return nil
}
