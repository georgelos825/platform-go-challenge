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

var (
	router    *gin.Engine
	testToken string
)

// Setup router and get token once before tests
func init() {
	router = routes.SetupRouter()
	testToken = getTestToken()
}

// getTestToken retrieves a JWT token only once
func getTestToken() string {
	reqBody := "user_id=123"
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	return response["token"]
}

// Test adding a favorite asset
func TestAddFavorite(t *testing.T) {
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
	req.Header.Set("Authorization", "Bearer "+testToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Asset added to favorites")
}

// Test retrieving favorites
func TestGetFavorites(t *testing.T) {
	req, _ := http.NewRequest("GET", "/favorites/123", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test removing a favorite
func TestRemoveFavorite(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/favorites/123/1", nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test editing a favorite's description
func TestEditFavorite(t *testing.T) {
	reqBody := map[string]string{"new_description": "Updated Stock trends"}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/favorites/123/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Asset updated")
}
