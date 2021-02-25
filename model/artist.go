package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Artist struct {
	ID    primitive.ObjectID    `json:"_id" bson:"_id"`
	Name  string                `json:"name" bson:"name"`
	Cover string                `json:"cover" bson:"cover"`
	Songs *[]primitive.ObjectID `json:"songs" bson:"songs"`
}
