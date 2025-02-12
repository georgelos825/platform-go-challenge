package routes

import (
	"github.com/gin-gonic/gin"
	"platform-go-challenge/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/favorites/:user_id", handlers.GetFavoritesHandler)
	r.POST("/favorites", handlers.AddFavoriteHandler)
	r.DELETE("/favorites/:user_id/:asset_id", handlers.RemoveFavoriteHandler)
	r.PUT("/favorites/:user_id/:asset_id", handlers.EditFavoriteHandler)

	return r
}
