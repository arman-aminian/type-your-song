package model

type Song struct {
	Name        string `json:"name" bson:"name"`
	Cover       string `json:"cover" bson:"cover"`
	Genre       Genre  `json:"genre" bson:"genre"`
	Artist      Artist `json:"artist" bson:"artist"`
	Duration    int    `json:"duration" bson:"duration"`
	PassedUsers []User `json:"passed_users" bson:"passed_users"`
	Url         string `json:"url" bson:"url"`
	//todo add lyric

	MaxWPM     int `json:"max_wpm" bson:"max_wpm"`
	AvgWPM     int `json:"avg_wpm" bson:"avg_wpm"`
	WordsCount int `json:"words_count" bson:"words_count"`
	Score      int `json:"score" bson:"score"`
}