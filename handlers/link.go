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

var linkCollectionName = "link"

type LinkHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewLinkHandler(ctx context.Context, database *mongo.Database, redisClient *redis.Client) *LinkHandler {
	return &LinkHandler{
		collection:  database.Collection(linkCollectionName),
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// swagger:operation GET /link link listLink
// Returns list of link
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func (handler *LinkHandler) ListLinkHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(linkCollectionName).Result()
	if err == redis.Nil {
		log.Printf("Link list requested to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)

		linkList := make([]models.Link, 0)
		for cur.Next(handler.ctx) {
			var link models.Link
			cur.Decode(&link)
			linkList = append(linkList, link)
		}

		data, _ := json.Marshal(linkList)
		handler.redisClient.Set(linkCollectionName, string(data), 0)
		c.JSON(http.StatusOK, linkList)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Link list requested to Redis")
		linkList := make([]models.Link, 0)
		json.Unmarshal([]byte(val), &linkList)
		c.JSON(http.StatusOK, linkList)
	}
}

// swagger:operation POST /link link newLink
// Create a new link
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func (handler *LinkHandler) NewLinkHandler(c *gin.Context) {
	var link models.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link.ID = primitive.NewObjectID()
	_, err := handler.collection.InsertOne(handler.ctx, link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new link"})
		return
	}

	log.Println("Removing link data from Redis")
	handler.redisClient.Del(linkCollectionName)

	c.JSON(http.StatusOK, link)
}

// swagger:operation PUT /link/{id} link updateLink
// Update an existing link
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the link
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
//         description: Invalid link ID
func (handler *LinkHandler) UpdateLinkHandler(c *gin.Context) {
	id := c.Param("id")
	var link models.Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "order", Value: link.Order},
		{Key: "icon", Value: link.Icon},
		{Key: "href", Value: link.Href},
		{Key: "label", Value: link.Label},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.redisClient.Del(linkCollectionName)

	c.JSON(http.StatusOK, gin.H{"message": "Link has been updated"})
}

// swagger:operation DELETE /link/{id} link deleteLink
// Delete an existing link
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the link
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid link ID
func (handler *LinkHandler) DeleteLinkHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	handler.redisClient.Del(linkCollectionName)

	c.JSON(http.StatusOK, gin.H{"message": "Link has been deleted"})
}

// swagger:operation GET /link/{id} link getOneLink
// Get one link
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: link ID
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
func (handler *LinkHandler) GetOneLinkHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var link models.Link
	err := cur.Decode(&link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, link)
}
