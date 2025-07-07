package interfaces

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"zen-connect/internal/infrastructure/auth0"
)

// Auth0Handler handles Auth0 authentication flow
type Auth0Handler struct {
	authService *auth0.AuthService
}

// NewAuth0Handler creates a new Auth0Handler
func NewAuth0Handler(authService *auth0.AuthService) *Auth0Handler {
	return &Auth0Handler{
		authService: authService,
	}
}

// SignInRequest represents the signin request
type SignInRequest struct {
	RedirectURI string `json:"redirect_uri" query:"redirect_uri"`
}

// SignIn handles GET /auth/login - redirects to Auth0 Universal Login (like reference app)
func (h *Auth0Handler) SignIn(c echo.Context) error {
	// Generate Auth0 login URL with state
	loginURL := h.authService.GetLoginURL("state")

	// Redirect to Auth0 (307 Temporary Redirect like reference app)
	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// CallbackRequest represents callback data
type CallbackRequest struct {
	AccessToken string `json:"access_token" query:"access_token"`
	TokenType   string `json:"token_type" query:"token_type"`
	ExpiresIn   string `json:"expires_in" query:"expires_in"`
	Error       string `json:"error" query:"error"`
	ErrorDesc   string `json:"error_description" query:"error_description"`
}

// Callback handles GET /api/auth/callback - processes Auth0 callback
func (h *Auth0Handler) Callback(c echo.Context) error {
	var req CallbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid callback parameters",
		})
	}

	// Check for authentication errors
	if req.Error != "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":             req.Error,
			"error_description": req.ErrorDesc,
		})
	}

	// Check if we have an access token
	if req.AccessToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No access token received",
		})
	}

	// Return token to frontend
	response := map[string]interface{}{
		"access_token": req.AccessToken,
		"token_type":   req.TokenType,
		"expires_in":   req.ExpiresIn,
		"success":      true,
	}

	return c.JSON(http.StatusOK, response)
}

// GetLoginURL handles GET /api/auth/login-url - returns login URL for AJAX calls
func (h *Auth0Handler) GetLoginURL(c echo.Context) error {
	// Generate random state for CSRF protection
	state := "random-state-string" // TODO: Generate proper random state
	
	loginURL := h.authService.GetLoginURL(state)

	return c.JSON(http.StatusOK, map[string]string{
		"login_url": loginURL,
		"state":     state,
	})
}

// SetupRoutes sets up the Auth0 authentication routes (matching reference app pattern)
func (h *Auth0Handler) SetupRoutes(e *echo.Echo) {
	auth := e.Group("/auth")

	// Authentication flow endpoints (matching reference app)
	auth.GET("/login", h.SignIn)            // Redirects to Auth0
	auth.GET("/callback", h.Callback)       // Handles Auth0 callback
	
	// Keep API endpoints too
	apiAuth := e.Group("/api/auth")
	apiAuth.GET("/login-url", h.GetLoginURL)   // Returns login URL for AJAX
}