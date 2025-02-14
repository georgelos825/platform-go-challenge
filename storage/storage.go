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

// GetFavoritesByType retrieves all favorites of a specific type for a given user using parallel filtering
func GetFavoritesByType(userID string, assetType string) []models.Asset {
	mu.Lock()
	userFavorites, exists := favorites[userID]
	mu.Unlock()
	if !exists {
		return []models.Asset{}
	}
	// Filter in parallel
	var filtered []models.Asset
	var wg sync.WaitGroup
	var muFilter sync.Mutex
	for _, fav := range userFavorites {
		wg.Add(1)
		go func(fav models.Asset) {
			defer wg.Done()
			if string(fav.Type) == assetType {
				muFilter.Lock()
				filtered = append(filtered, fav)
				muFilter.Unlock()
			}
		}(fav)
	}
	wg.Wait()
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
