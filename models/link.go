package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters link newLink
type Link struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Order int                `json:"order" bson:"order"`
	Icon  string             `json:"icon" bson:"icon"`
	Href  string             `json:"href" bson:"href"`
	Label string             `json:"label" bson:"label"`
	IsNew string             `json:"isNew" bson:"isNew"`
}
