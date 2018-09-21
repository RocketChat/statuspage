package middleware

import (
	"github.com/gin-gonic/gin"
)

//NotImplemented sends a status and message about it not being implemented yet
func NotImplemented(c *gin.Context) {
	c.JSON(501, gin.H{"message": "Not implemented yet."})
}
