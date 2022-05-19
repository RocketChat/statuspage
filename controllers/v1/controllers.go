package v1

import (
	"log"
	"net/http"

	"github.com/RocketChat/statuscentral/core"
	"github.com/gin-gonic/gin"
)

func internalErrorHandler(c *gin.Context, err error) {
	log.Println(err)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}

func badRequestHandlerDetailed(c *gin.Context, err error) {
	log.Println(err)

	c.JSON(http.StatusBadRequest, gin.H{"error": "bad request", "details": err.Error()})
}

func internalErrorHandlerDetailed(c *gin.Context, err error) {
	log.Println(err)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error", "details": err.Error()})
}

// LivenessCheckHandler checks to see whether the database responds to a ping
func LivenessCheckHandler(c *gin.Context) {
	if err := core.LivenessCheck(); err != nil {
		internalErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
