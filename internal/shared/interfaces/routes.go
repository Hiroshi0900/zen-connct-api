package interfaces

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RoutesHandler provides API documentation and route listing
type RoutesHandler struct{}

// NewRoutesHandler creates a new routes handler
func NewRoutesHandler() *RoutesHandler {
	return &RoutesHandler{}
}

// ListRoutes returns all available API endpoints
func (h *RoutesHandler) ListRoutes(c echo.Context) error {
	routes := map[string]interface{}{
		"service": "zen-connect-api",
		"version": "1.0.0",
		"endpoints": map[string]interface{}{
			"health": map[string]string{
				"GET /health":           "Basic health check",
				"GET /health/protected": "Protected health check (requires Auth0 token)",
			},
			"authentication": map[string]string{
				"GET /auth/login":       "Redirect to Auth0 login",
				"GET /auth/callback":    "Auth0 callback handler",
				"GET /api/auth/login-url": "Get Auth0 login URL (AJAX)",
			},
			"documentation": map[string]string{
				"GET /api/routes": "This endpoint - list all routes",
			},
		},
		"auth": map[string]interface{}{
			"provider": "Auth0",
			"domain":   "dev-ie6tlg1ol8xjemia.us.auth0.com",
			"flow":     "Authorization Code with PKCE",
			"scopes":   []string{"openid", "profile", "email"},
		},
	}

	return c.JSON(http.StatusOK, routes)
}

// SetupRoutesDocumentation sets up the routes documentation endpoint
func (h *RoutesHandler) SetupRoutes(e *echo.Echo) {
	api := e.Group("/api")
	api.GET("/routes", h.ListRoutes)
}