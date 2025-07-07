package main

import (
	"log"
	"zen-connect/internal/infrastructure/auth0"
	"zen-connect/internal/shared/interfaces"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	// Initialize Auth0 configuration
	auth0Config, err := auth0.NewConfig()
	if err != nil {
		log.Fatal("Failed to create Auth0 config:", err)
	}

	// Initialize Auth0 middleware
	authMiddleware, err := auth0.NewAuthMiddleware(auth0Config)
	if err != nil {
		log.Fatal("Failed to create Auth0 middleware:", err)
	}

	// Initialize Auth0 service and handler
	authService, err := auth0.NewAuthService(auth0Config)
	if err != nil {
		log.Fatal("Failed to create Auth0 service:", err)
	}
	auth0Handler := interfaces.NewAuth0Handler(authService)

	// Setup Auth0 routes
	auth0Handler.SetupRoutes(e)

	// Setup API documentation
	routesHandler := interfaces.NewRoutesHandler()
	routesHandler.SetupRoutes(e)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "OK",
			"service": "zen-connect-api",
		})
	})

	// Protected health check endpoint to test Auth0
	e.GET("/health/protected", func(c echo.Context) error {
		userID, ok := auth0.GetUserIDFromContext(c.Request().Context())
		if !ok {
			return c.JSON(401, map[string]string{
				"error": "User not authenticated",
			})
		}

		email, _ := auth0.GetUserEmailFromContext(c.Request().Context())
		name, _ := auth0.GetUserNameFromContext(c.Request().Context())

		return c.JSON(200, map[string]interface{}{
			"status":  "OK",
			"service": "zen-connect-api",
			"user_id": userID,
			"email":   email,
			"name":    name,
		})
	}, authMiddleware.RequireAuth())

	// Start server
	log.Println("Starting zen-connect API server on :8080")
	log.Printf("Auth0 Domain: %s", auth0Config.Domain)
	log.Printf("Auth0 Audience: %s", auth0Config.Audience)
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}