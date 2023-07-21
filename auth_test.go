package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	auth "github.com/go-app-service/api/auth"
	"github.com/stretchr/testify/assert"
)

func TestGoogleLoginHandler(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a router with the GoogleLoginHandler route
	router := gin.Default()
	router.GET("/google/login", auth.GoogleLogin)

	// Create a test server with the router as the handler
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Make a request to the test server
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(testServer.URL + "/google/login")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Ensure the correct Location header points to the Google login URL
	expectedURL := "https://accounts.google.com/o/oauth2/auth"
	assert.Contains(t, resp.Header.Get("Location"), expectedURL, "Redirect URL should match")
}
