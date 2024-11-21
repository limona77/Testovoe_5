package controller

import (
	"Testovoe_5/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *gin.Engine, songsService *service.Services) {
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	newSongsRoutes(app, songsService)
}
