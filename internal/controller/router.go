package controller

import (
	_ "Testovoe_5/docs"
	"Testovoe_5/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(app *gin.Engine, songsService *service.Services) {

	songs := app.Group("/songs")
	newSongsRoutes(songs, songsService)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
