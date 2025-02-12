package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"platform-go-challenge/models"
	"platform-go-challenge/storage"
)

// Request structure for adding a favorite
type AddFavoriteRequest struct {
	UserID string       `json:"user_id"`
	Asset  models.Asset `json:"asset"`
}

// Get all favorites for a user
func GetFavoritesHandler(c *gin.Context) {
	userID := c.Param("user_id")
	c.JSON(http.StatusOK, gin.H{"favorites": storage.GetFavorites(userID)})
}

// Add a favorite
func AddFavoriteHandler(c *gin.Context) {
	var req AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	storage.AddFavorite(req.UserID, req.Asset)
	c.JSON(http.StatusOK, gin.H{"message": "Asset added to favorites"})
}

// Remove a favorite
func RemoveFavoriteHandler(c *gin.Context) {
	userID := c.Param("user_id")
	assetID := c.Param("asset_id")
	storage.RemoveFavorite(userID, assetID)
	c.JSON(http.StatusOK, gin.H{"message": "Asset removed"})
}

// Edit favorite description
func EditFavoriteHandler(c *gin.Context) {
	userID := c.Param("user_id")
	assetID := c.Param("asset_id")

	var req struct {
		NewDescription string `json:"new_description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storage.EditFavorite(userID, assetID, req.NewDescription)
	c.JSON(http.StatusOK, gin.H{"message": "Asset updated"})
}
