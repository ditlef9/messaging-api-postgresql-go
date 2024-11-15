// File: handlers/product_handlers.go

package handlers

import (
	"ekeberg.com/messaging-api-postgresql-go/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetMessages(c *gin.Context) {
	products, err := models.GetMessages(-1)
	if err != nil {
		log.Fatal(err)
	}

	if products == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": products})
	}
}

func GetMessageById(c *gin.Context) {
	id := c.Param("id")
	sqlMessage, err := models.GetMessageById(id)
	if err != nil {
		log.Fatal(err)
	}

	if sqlMessage.Id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": sqlMessage})
	}
}

func Options(c *gin.Context) {
	ourOptions := "HTTP/1.1 200 OK\n" +
		"Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Origin: http://locahost:8080\n" +
		"Access-Control-Allow-Methods: GET\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	c.String(200, ourOptions)
}
