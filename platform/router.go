package platform

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}
