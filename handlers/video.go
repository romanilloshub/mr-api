package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"orov.io/siempreAbierto/models"
)

var videoCollectionName = "video"

type VideoHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewVideoHandler(ctx context.Context, database *mongo.Database, redisClient *redis.Client) *VideoHandler {
	return &VideoHandler{
		collection:  database.Collection(videoCollectionName),
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// swagger:operation GET /video video listVideo
// Returns list of video
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func (handler *VideoHandler) ListVideoHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(videoCollectionName).Result()
	if err == redis.Nil {
		log.Printf("Video list requested to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)

		videoList := make([]models.Video, 0)
		for cur.Next(handler.ctx) {
			var video models.Video
			cur.Decode(&video)
			videoList = append(videoList, video)
		}

		data, _ := json.Marshal(videoList)
		handler.redisClient.Set(videoCollectionName, string(data), 0)
		c.JSON(http.StatusOK, videoList)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Video list requested to Redis")
		videoList := make([]models.Video, 0)
		json.Unmarshal([]byte(val), &videoList)
		c.JSON(http.StatusOK, videoList)
	}
}

// swagger:operation POST /video video newVideo
// Create a new video
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func (handler *VideoHandler) NewVideoHandler(c *gin.Context) {
	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	video.ID = primitive.NewObjectID()
	_, err := handler.collection.InsertOne(handler.ctx, video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new video"})
		return
	}

	log.Println("Removing video data from Redis")
	handler.redisClient.Del(videoCollectionName)

	c.JSON(http.StatusOK, video)
}

// swagger:operation PUT /video/{id} video updateVideo
// Update an existing video
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the video
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
//     '404':
//         description: Invalid video ID
func (handler *VideoHandler) UpdateVideoHandler(c *gin.Context) {
	id := c.Param("id")
	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: video.Title},
		{Key: "description", Value: video.Description},
		{Key: "youtubeLink", Value: video.YoutubeLink},
		{Key: "downloadLink", Value: video.DownloadLink},
		{Key: "meta", Value: video.Meta},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.redisClient.Del(videoCollectionName)

	c.JSON(http.StatusOK, gin.H{"message": "Video has been updated"})
}

// swagger:operation DELETE /video/{id} video deleteVideo
// Delete an existing video
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the video
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid video ID
func (handler *VideoHandler) DeleteVideoHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.redisClient.Del(videoCollectionName)

	c.JSON(http.StatusOK, gin.H{"message": "Video has been deleted"})
}

// swagger:operation GET /video/{id} video getOneVideo
// Get one video
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: video ID
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
func (handler *VideoHandler) GetOneVideoHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var video models.Video
	err := cur.Decode(&video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}
