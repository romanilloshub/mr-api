package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:parameters video newVideo
type Video struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Image        string             `json:"image" bson:"image"`
	Description  string             `json:"description" bson:"description"`
	URL          string             `json:"url" bson:"url"`
	LatLng       []float64          `json:"latlng"`
	DownloadLink string             `json:"downloadLink" bson:"downloadLink"`
	Meta         Meta               `json:"meta" bson:"meta"`
	Link         []VideoLink        `json:"link" bson:"link"`
	QR           string             `json:"qr" bson:"qr"`
	PostURL      string             `json:"postURL" bson:"postURL"`
}

type VideoLink struct {
	Type        string `json:"type" bson:"type"`
	Name        string `json:"name" bson:"name"`
	Icon        string `json:"icon" bson:"icon"`
	Href        string `json:"href" bson:"href"`
	Description string `json:"description" bson:"description"`
}

type Meta struct {
	Geo         Geo    `json:"geo" bson:"geo"`
	GeoSnapshot string `json:"geoSnapshot" bson:"geoSnapshot"`
}

type Geo struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}
