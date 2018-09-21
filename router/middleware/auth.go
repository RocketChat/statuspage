package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//IsAuthorized checks to ensure the request can be made
func IsAuthorized(c *gin.Context) {
	token := c.GetHeader("Authorization")
	validToken := os.Getenv("API_TOKEN")

	if !(len(validToken) > 0 && token == validToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}
