package handlers

import (
	"encoding/json"
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
	UserID string          `json:"user_id"`
	Asset  json.RawMessage `json:"asset"`
}

// GetFavoritesHandler retrieves a user's favorites concurrently
func GetFavoritesHandler(c *gin.Context) {
	userID := c.Param("user_id")
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allFavorites []models.AssetInterface
	getAndAppend := func(assetType string) {
		defer wg.Done()
		assets := storage.GetFavoritesByType(userID, assetType)
		mu.Lock()
		allFavorites = append(allFavorites, assets...)
		mu.Unlock()
	}
	wg.Add(3)
	go getAndAppend("chart")
	go getAndAppend("insight")
	go getAndAppend("audience")
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"favorites": allFavorites})
}

// Add a favorite
func AddFavoriteHandler(c *gin.Context) {
	var req struct {
		UserID string          `json:"user_id"`
		Asset  json.RawMessage `json:"asset"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var base models.Asset
	if err := json.Unmarshal(req.Asset, &base); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset format"})
		return
	}
	switch base.Type {
	case models.ChartType:
		var chart models.Chart
		if err := json.Unmarshal(req.Asset, &chart); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chart format"})
			return
		}
		storage.AddFavorite(req.UserID, chart)
	case models.AudienceType:
		var audience models.Audience
		if err := json.Unmarshal(req.Asset, &audience); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audience format"})
			return
		}
		storage.AddFavorite(req.UserID, audience)
	case models.InsightType:
		var insight models.Insight
		if err := json.Unmarshal(req.Asset, &insight); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid insight format"})
			return
		}
		storage.AddFavorite(req.UserID, insight)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset type"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Asset added successfully"})
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
