package model

type Genre struct {
	Name  string `json:"name" bson:"name"`
	Cover string `json:"cover" bson:"cover"`
	Songs []Song `json:"songs" bson:"songs"`
}
