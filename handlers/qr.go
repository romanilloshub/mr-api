package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"orov.io/siempreAbierto/models"
)

const blogURL = "https://blog.romanillos.net"

type QRHandler struct {
	video       *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewQRHandler(ctx context.Context, database *mongo.Database, redisClient *redis.Client) *QRHandler {
	return &QRHandler{
		video:       database.Collection(videoCollectionName),
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// swagger:operation GET /siempre-abierto/video/qr/{qr} video getOneVideo
// Get one video
// ---
// produces:
// - application/json
// parameters:
//   - name: qr
//     in: path
//     description: QR ID
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
func (handler *QRHandler) RedirectToVideoByQR(c *gin.Context) {
	QR := c.Param("qr")
	cur := handler.video.FindOne(handler.ctx, bson.M{
		"qr": QR,
	})

	if cur.Err() == mongo.ErrNoDocuments {
		c.Redirect(http.StatusTemporaryRedirect, blogURL)
		return
	}
	var video models.Video
	err := cur.Decode(&video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if video.PostURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, video.PostURL)
	}
	c.Redirect(http.StatusTemporaryRedirect, blogURL)
}
