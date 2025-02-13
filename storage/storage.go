package storage

import (
	"platform-go-challenge/models"
	"sync"
)

// Mutex for safe concurrent access
var (
	mu        sync.Mutex
	favorites = make(map[string][]models.Asset) // UserID → List of Assets
)

// GetFavoritesByType retrieves all favorites of a specific type for a given user
func GetFavoritesByType(userID string, assetType string) []models.Asset {
	mu.Lock()
	defer mu.Unlock()
	userFavorites, exists := favorites[userID]
	if !exists {
		return []models.Asset{}
	}
	// Filter by asset type
	var filtered []models.Asset
	for _, fav := range userFavorites {
		if string(fav.Type) == assetType { // Convert AssetType to string
			filtered = append(filtered, fav)
		}
	}
	return filtered
}

// Add asset to favorites
func AddFavorite(userID string, asset models.Asset) {
	mu.Lock()
	defer mu.Unlock()
	favorites[userID] = append(favorites[userID], asset)
}

// Get all favorites for a user
func GetFavorites(userID string) []models.Asset {
	mu.Lock()
	defer mu.Unlock()
	return favorites[userID]
}

// Remove asset from favorites
func RemoveFavorite(userID, assetID string) {
	mu.Lock()
	defer mu.Unlock()
	userAssets := favorites[userID]
	for i, asset := range userAssets {
		if asset.ID == assetID {
			favorites[userID] = append(userAssets[:i], userAssets[i+1:]...)
			return
		}
	}
}

// Edit asset description
func EditFavorite(userID, assetID, newDescription string) {
	mu.Lock()
	defer mu.Unlock()
	for i, asset := range favorites[userID] {
		if asset.ID == assetID {
			favorites[userID][i].Description = newDescription
			return
		}
	}
}
