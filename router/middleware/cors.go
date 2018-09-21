package middleware

import (
	"github.com/gin-gonic/gin"
)

//CORSMiddleware provides some headers for CORS
func CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization, Cache-Control, Expires, Pragma, X-powered-by")
	c.Writer.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Expires", "-1")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("X-powered-by", "Rocket Fuel and Rocketeers")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}
