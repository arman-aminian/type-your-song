package model

type Artist struct {
	Name  string `json:"name" bson:"name"`
	Songs []Song `json:"artist" bson:"artist"`
}
