package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters link newLink
type Link struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Type        string             `json:"type" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	Icon        string             `json:"icon" bson:"icon"`
	Href        string             `json:"href" bson:"href"`
	Description string             `json:"description" bson:"description"`
}
