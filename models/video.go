package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters video newVideo
type Video struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description" bson:"description"`
	YoutubeLink  string             `json:"youtubeLink" bson:"youtubeLink"`
	DownloadLink string             `json:"downloadLink" bson:"downloadLink"`
	Meta         Meta               `json:"meta" bson:"meta"`
}

type Meta struct {
	Geo         Geo    `json:"geo" bson:"geo"`
	GeoSnapshot string `json:"geoSnapshot" bson:"geoSnapshot"`
}

type Geo struct {
	Lat string `json:"lat" bson:"lat"`
	Lng string `json:"lng" bson:"lng"`
}
