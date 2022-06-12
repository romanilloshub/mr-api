package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

const bearerPrefix = "bearer "

type AuthMiddleware struct {
	firebaseClient *auth.Client
	ctx            context.Context
}

func NewAuthMiddleware(ctx context.Context) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseClient: client,
		ctx:            ctx,
	}
}

func (middleware *AuthMiddleware) AllowAnonymous(c *gin.Context) {
	c.Next()
}

func (middleware *AuthMiddleware) LoggedUser(c *gin.Context) {
	jwt, err := getJWT(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "You must supply an Authorization header with a valid jwt, prefixed by bearer.",
		})
		log.Printf("Can't obtain JWT: %v\n", err)
		return
	}
	_, err = client.VerifyIDToken(middleware.ctx, jwt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid jwt in Authorization header.",
		})
		log.Printf("Can't parse JWT: %v\n", err)
		return
	}
	c.Next()
}

func (middleware *AuthMiddleware) WithRol(rol string) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := getJWT(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You must supply an Authorization header with a valid jwt, prefixed by bearer.",
			})
			log.Printf("Can't obtain JWT: %v\n", err)
			return
		}

		idToken, err := client.VerifyIDToken(middleware.ctx, jwt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid jwt in Authorization header.",
			})
			log.Printf("Can't parse JWT: %v\n", err)
			return
		}

		if !idToken.Claims[rol].(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "You have no permission to do this operation",
			})
			return
		}

		c.Next()
	}
}

func getJWT(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader("Authorization")

	if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
		return "", errors.New("authorization info bad formatted")
	}

	jwt := strings.TrimPrefix(authorizationHeader, bearerPrefix)
	return jwt, nil
}
