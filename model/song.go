package model

import (
	"github.com/arman-aminian/gosub/parsers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Song struct {
	ID          primitive.ObjectID    `json:"_id" bson:"_id"`
	Name        string                `json:"name" bson:"name"`
	Cover       string                `json:"cover" bson:"cover"`
	Genre       string                `json:"genre" bson:"genre"`
	Artist      primitive.ObjectID    `json:"artist" bson:"artist"`
	Duration    int                   `json:"duration" bson:"duration"`
	PassedUsers *[]primitive.ObjectID `json:"passed_users" bson:"passed_users"`
	Url         string                `json:"url" bson:"url"`
	Lyrics      parsers.Srt           `json:"lyrics" bson:"lyrics"`

	MaxWPM     int `json:"max_wpm" bson:"max_wpm"`
	AvgWPM     int `json:"avg_wpm" bson:"avg_wpm"`
	WordsCount int `json:"words_count" bson:"words_count"`
	Score      int `json:"score" bson:"score"`
}
