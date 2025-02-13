package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"platform-go-challenge/models"
	"platform-go-challenge/storage"
	"sync"
)

// convertAssetsToInterface converts a slice of models.Asset to a slice of interface{}
func convertAssetsToInterface(assets []models.Asset) []interface{} {
	result := make([]interface{}, len(assets))
	for i, asset := range assets {
		result[i] = asset
	}
	return result
}

// Request structure for adding a favorite
type AddFavoriteRequest struct {
	UserID string       `json:"user_id"`
	Asset  models.Asset `json:"asset"`
}

// GetFavoritesHandler retrieves a user's favorites concurrently
func GetFavoritesHandler(c *gin.Context) {
	userID := c.Param("user_id")
	var wg sync.WaitGroup
	results := make([][]interface{}, 3)
	// Launch concurrent tasks with indexed results
	wg.Add(3)
	go func() {
		defer wg.Done()
		charts := storage.GetFavoritesByType(userID, "chart")
		results[0] = convertAssetsToInterface(charts)
	}()
	go func() {
		defer wg.Done()
		insights := storage.GetFavoritesByType(userID, "insight")
		results[1] = convertAssetsToInterface(insights)
	}()
	go func() {
		defer wg.Done()
		audiences := storage.GetFavoritesByType(userID, "audience")
		results[2] = convertAssetsToInterface(audiences)
	}()
	wg.Wait()
	// Combine all results into a single slice
	var favorites []interface{}
	for _, res := range results {
		favorites = append(favorites, res...)
	}
	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

// Add a favorite
func AddFavoriteHandler(c *gin.Context) {
	var req AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	// Validate required fields
	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}
	if req.Asset.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Asset ID is required"})
		return
	}
	if req.Asset.Type != models.ChartType && req.Asset.Type != models.InsightType && req.Asset.Type != models.AudienceType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset type"})
		return
	}
	storage.AddFavorite(req.UserID, req.Asset)
	c.JSON(http.StatusOK, gin.H{"message": "Asset added to favorites"})
}

// Remove a favorite
func RemoveFavoriteHandler(c *gin.Context) {
	userID := c.Param("user_id")
	assetID := c.Param("asset_id")
	if userID == "" || assetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Asset ID are required"})
		return
	}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	if req.NewDescription == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New description cannot be empty"})
		return
	}
	storage.EditFavorite(userID, assetID, req.NewDescription)
	c.JSON(http.StatusOK, gin.H{"message": "Asset updated"})
}
