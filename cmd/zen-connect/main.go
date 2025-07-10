package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zen-connect/internal/infrastructure/auth0"
	"zen-connect/internal/infrastructure/logger"
	"zen-connect/internal/infrastructure/postgres"
	"zen-connect/internal/infrastructure/session"
	"zen-connect/internal/shared/interfaces"
	"zen-connect/internal/user/infrastructure"
	userservice "zen-connect/internal/user/application/service"
	userusecase "zen-connect/internal/user/application/usecase"
	userinterfaces "zen-connect/internal/user/interfaces"
	authinterfaces "zen-connect/internal/auth/interfaces"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize logger
	loggerConfig := logger.NewConfig()
	if err := logger.Initialize(loggerConfig); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Close()

	// Use structured logging from here on
	logger.Info("Starting zen-connect application",
		zap.String("service", loggerConfig.ServiceName),
		zap.String("version", loggerConfig.ServiceVersion),
		zap.String("environment", loggerConfig.Environment),
	)

	ctx := context.Background()

	// Initialize PostgreSQL client
	logger.Info("Initializing PostgreSQL client")
	pgClient, err := postgres.NewClient(ctx)
	if err != nil {
		logger.Fatal("Failed to create PostgreSQL client", zap.Error(err))
	}
	defer func() {
		logger.Info("Closing PostgreSQL client")
		pgClient.Close()
	}()
	logger.Info("PostgreSQL client initialized successfully")

	// Initialize repositories
	logger.Info("Initializing repositories")
	userRepo := infrastructure.NewPostgresUserRepository(pgClient.Pool)

	// Initialize session store
	logger.Info("Initializing session store")
	sessionStore, err := session.NewCookieStore()
	if err != nil {
		logger.Fatal("Failed to create session store", zap.Error(err))
	}
	logger.Info("Session store initialized successfully")

	// Initialize Auth0 configuration
	logger.Info("Initializing Auth0 configuration")
	auth0Config, err := auth0.NewConfig()
	if err != nil {
		logger.Fatal("Failed to create Auth0 config", zap.Error(err))
	}
	logger.Info("Auth0 configuration loaded",
		zap.String("domain", auth0Config.Domain),
		zap.String("audience", auth0Config.Audience),
	)

	// Initialize OIDC provider
	logger.Info("Initializing OIDC provider")
	provider, err := oidc.NewProvider(ctx, auth0Config.IssuerURL())
	if err != nil {
		logger.Fatal("Failed to create OIDC provider", zap.Error(err))
	}
	logger.Info("OIDC provider initialized successfully")

	// Initialize Auth0 service
	logger.Info("Initializing Auth0 service")
	authService, err := auth0.NewAuthService(auth0Config)
	if err != nil {
		logger.Fatal("Failed to create Auth0 service", zap.Error(err))
	}
	logger.Info("Auth0 service initialized successfully")

	// Initialize session middleware
	logger.Info("Initializing session middleware")
	sessionMiddleware := session.NewMiddleware(sessionStore)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Global middleware with structured logging
	e.Use(logger.RequestLoggerMiddleware())
	e.Use(logger.SessionLoggerMiddleware())
	e.Use(logger.ErrorLoggerMiddleware())
	e.Use(logger.RecoveryLoggerMiddleware())
	
	// CORS middleware with proper configuration for authentication
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // デフォルト値
	}
	logger.Info("CORS configuration", zap.String("frontend_url", frontendURL))
	
	// CORS設定（開発環境用）
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET", 
			"POST", 
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Initialize new architecture components
	logger.Info("Initializing new architecture components")
	
	// User service
	userService := userservice.NewUserService(userRepo)

	// Initialize new auth handler with UserService
	logger.Info("Initializing new auth handler")
	newAuthHandler, err := authinterfaces.NewAuthHandler(authService, userService, sessionStore, provider, auth0Config)
	if err != nil {
		logger.Fatal("Failed to create new auth handler", zap.Error(err))
	}
	logger.Info("New auth handler initialized successfully")

	// Setup routes
	logger.Info("Setting up application routes")
	// Use new auth handler
	newAuthHandler.SetupRoutes(e)

	// User use cases
	getUserProfileUseCase := userusecase.NewGetUserProfileUseCase(userService)
	
	// User handler (only keeping GetCurrentUser endpoint)
	userHandler := userinterfaces.NewUserHandler(nil, nil, getUserProfileUseCase)
	userHandler.SetupRoutes(e, sessionMiddleware)

	// Setup API documentation
	routesHandler := interfaces.NewRoutesHandler()
	routesHandler.SetupRoutes(e)
	
	logger.Info("Routes configured successfully")

	// Health check endpoints
	e.GET("/health", func(c echo.Context) error {
		logCtx := logger.WithContext(c.Request().Context())
		logCtx.Debug("Health check requested")
		return c.JSON(200, map[string]string{
			"status":  "OK",
			"service": "zen-connect-api",
		})
	})

	// Protected health check endpoint to test session-based auth
	e.GET("/health/protected", func(c echo.Context) error {
		logCtx := logger.WithContext(c.Request().Context())
		userID, ok := session.GetUserIDFromContext(c.Request().Context())
		if !ok {
			logCtx.Warn("Protected health check accessed without authentication")
			return c.JSON(401, map[string]string{
				"error": "User not authenticated",
			})
		}

		email, _ := session.GetUserEmailFromContext(c.Request().Context())
		name, _ := session.GetUserNameFromContext(c.Request().Context())
		auth0UserID, _ := session.GetAuth0UserIDFromContext(c.Request().Context())

		logCtx.Info("Protected health check accessed by authenticated user",
			zap.String("user_id", userID),
			zap.String("auth0_user_id", auth0UserID),
		)

		return c.JSON(200, map[string]interface{}{
			"status":        "OK",
			"service":       "zen-connect-api",
			"user_id":       userID,
			"auth0_user_id": auth0UserID,
			"email":         email,
			"name":          name,
		})
	}, sessionMiddleware.RequireAuth())

	// Database health check
	e.GET("/health/db", func(c echo.Context) error {
		logCtx := logger.WithContext(c.Request().Context())
		logCtx.Debug("Database health check requested")
		
		start := time.Now()
		if err := pgClient.Health(c.Request().Context()); err != nil {
			logCtx.Error("Database health check failed",
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
			)
			return c.JSON(500, map[string]string{
				"status": "ERROR",
				"error":  err.Error(),
			})
		}
		
		logCtx.Info("Database health check successful",
			zap.Duration("duration", time.Since(start)),
		)
		return c.JSON(200, map[string]string{
			"status": "OK",
			"db":     "connected",
		})
	})

	// Start server with graceful shutdown
	go func() {
		logger.Info("Starting zen-connect API server",
			zap.String("port", "8080"),
			zap.String("auth0_domain", auth0Config.Domain),
			zap.String("auth0_audience", auth0Config.Audience),
			zap.String("database_url", "***MASKED***"),
			zap.String("environment", loggerConfig.Environment),
		)
		if err := e.Start(":8080"); err != nil {
			logger.Info("Server stopped", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown signal received, starting graceful shutdown")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	
	logger.Info("Server shutdown completed successfully")
}