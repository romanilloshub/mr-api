package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"orov.io/siempreAbierto/plugin/auth"
)

type AdminHandler struct {
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (handler *AdminHandler) New(c *gin.Context) {
	email := c.Query("email")
	authUtil := auth.NewUtil(context.Background())
	err := authUtil.NewAdmin(email)
	log.Printf("The error: %s\n", err)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin created with mail: " + email})

}

func (handler *AdminHandler) GetUser(c *gin.Context) {
	email := c.Query("email")
	authUtil := auth.NewUtil(context.Background())
	user, err := authUtil.GetUser(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
