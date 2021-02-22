package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Artist struct {
	Name  string               `json:"name" bson:"name"`
	Cover string               `json:"cover" bson:"cover"`
	Songs []primitive.ObjectID `json:"artist" bson:"artist"`
}
