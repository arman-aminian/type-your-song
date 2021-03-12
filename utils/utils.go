package utils

import (
	"errors"
	"github.com/arman-aminian/type-your-song/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

const (
	BaseUrl   = "127.0.0.1:8080"
	Guest     = "guest"
	NotPassed = "none"
	Simple    = "simple"
	Medium    = "medium"
	Hard      = "hard"
)

func ValidPassedLevel(s string) bool {
	if s == Simple || s == Medium || s == Hard {
		return true
	}
	return false
}

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

func FindPassedSong(slice []model.PassedSong, val primitive.ObjectID) (model.PassedSong, error) {
	if slice == nil {
		return model.PassedSong{}, errors.New("null")
	}
	for _, item := range slice {
		if item.SID == val {
			return item, nil
		}
	}
	return model.PassedSong{}, errors.New("not found")
}
