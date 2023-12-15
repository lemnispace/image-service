package util

import (
	"log"

	"github.com/gin-gonic/gin"
)

// SendErrorResponse sends an error response to the client with the specified status code, message, and error.
// If the error is not nil, it logs the message and error. Otherwise, it only logs the message.
// It uses the gin.Context to send the JSON response.
func SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	} else {
		log.Printf("%s", message)
	}
	c.JSON(statusCode, gin.H{"error": message})
}
