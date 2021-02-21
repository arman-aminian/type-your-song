package utils

import (
	"io"
	"mime/multipart"
	"os"
)

const (
	BaseUrl = "127.0.0.1:8080"
)

func SaveToFiles(f multipart.FileHeader, path string) (string, error) {
	src, err := f.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mediaPath := path + f.Filename
	dst, err := os.Create(mediaPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return mediaPath, nil
}
