package controller

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(app *gin.Engine) {
	// TODO: прокинуть сервисы
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
