package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//NotImplemented sends a status and message about it not being implemented yet
func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented yet"})
}
