package utils

import (
	"io"
	"mime/multipart"
	"os"
	"strings"
)

const (
	BaseUrl = "127.0.0.1:8080"
)

func SaveToFiles(f multipart.FileHeader, path string, name string) (string, error) {
	src, err := f.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	format := strings.Split(f.Filename, ".")
	mediaPath := path + name + "." + format[len(format)-1]
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
