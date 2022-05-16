// Siempre Abierto API
//
// So we will need to add a description
//
//	Schemes: http
//  Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Javi <hi@orov.io> https://soriandev.orov.io
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"orov.io/siempreAbierto/handlers"
)

var linkHandler *handlers.LinkHandler
var videoHandler *handlers.VideoHandler
var pingHandler *handlers.PingHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Can't connect to database: " + err.Error())
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	database := client.Database(os.Getenv("MONGO_DATABASE"))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	status := redisClient.Ping()
	log.Printf("Redis: %s\n", status)

	linkHandler = handlers.NewLinkHandler(ctx, database, redisClient)
	videoHandler = handlers.NewVideoHandler(ctx, database, redisClient)
	pingHandler = handlers.NewPingHandler()
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:  []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "https://ngrok")
		},
	}))

	router.POST("/link", linkHandler.NewLinkHandler)
	router.GET("/link", linkHandler.ListLinkHandler)
	router.PUT("/link/:id", linkHandler.UpdateLinkHandler)
	router.DELETE("/link/:id", linkHandler.DeleteLinkHandler)
	router.GET("/link/:id", linkHandler.GetOneLinkHandler)

	router.POST("/video", videoHandler.NewVideoHandler)
	router.GET("/video", videoHandler.ListVideoHandler)
	router.PUT("/video/:id", videoHandler.UpdateVideoHandler)
	router.DELETE("/video/:id", videoHandler.DeleteVideoHandler)
	router.GET("/video/:id", videoHandler.GetOneVideoHandler)

	router.GET("/ping", pingHandler.Ping)

	router.Run()
}
