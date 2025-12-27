package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine) {

	//ver si está vivo
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "server say: pong")
	})

}
