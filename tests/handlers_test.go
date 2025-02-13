package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/models"
	"platform-go-challenge/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Global JWT token for tests
var testToken string

// SetupRouter initializes a test Gin server
func SetupRouter() *gin.Engine {
	return routes.SetupRouter()
}

// GetTestToken retrieves a JWT token for testing
func GetTestToken(t *testing.T) string {
	if testToken != "" {
		return testToken
	}

	// Prepare request to get token
	reqBody := "user_id=123"
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	w := httptest.NewRecorder()
	router := SetupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Extract token from response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	token, exists := response["token"]
	assert.True(t, exists, "Token not found in response")

	testToken = token
	return token
}

// Test adding a favorite asset
func TestAddFavorite(t *testing.T) {
	router := SetupRouter()
	token := GetTestToken(t)

	reqBody := models.Asset{
		ID:          "1",
		Type:        models.ChartType,
		Description: "Stock trends",
	}
	requestData := map[string]interface{}{
		"user_id": "123",
		"asset":   reqBody,
	}
	jsonValue, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/favorites", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Asset added to favorites")
}

// Test retrieving favorites
func TestGetFavorites(t *testing.T) {
	router := SetupRouter()
	token := GetTestToken(t)

	req, _ := http.NewRequest("GET", "/favorites/123", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test removing a favorite
func TestRemoveFavorite(t *testing.T) {
	router := SetupRouter()
	token := GetTestToken(t)

	req, _ := http.NewRequest("DELETE", "/favorites/123/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test editing a favorite's description
func TestEditFavorite(t *testing.T) {
	router := SetupRouter()
	token := GetTestToken(t)

	reqBody := map[string]string{"new_description": "Updated Stock trends"}
	jsonValue, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/favorites/123/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Asset updated")
}
