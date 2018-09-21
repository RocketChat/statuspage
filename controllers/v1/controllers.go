package v1

import (
	"log"

	"github.com/gin-gonic/gin"
)

func internalErrorHandler(c *gin.Context, err error) {
	log.Println(err)

	c.JSON(500, gin.H{"error": "internal error"})
}
