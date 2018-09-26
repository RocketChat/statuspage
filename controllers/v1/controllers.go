package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func internalErrorHandler(c *gin.Context, err error) {
	log.Println(err)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}

func internalErrorHandlerDetailed(c *gin.Context, err error) {
	log.Println(err)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error", "details": err.Error()})
}
