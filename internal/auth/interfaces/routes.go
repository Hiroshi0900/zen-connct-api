package interfaces

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"zen-connect/internal/infrastructure/auth0"
	"zen-connect/internal/infrastructure/logger"
	"zen-connect/internal/infrastructure/session"
	userservice "zen-connect/internal/user/application/service"
)

// AuthHandler 認証関連のHTTPハンドラー
type AuthHandler struct {
	authService     *auth0.AuthService
	callbackHandler *auth0.CallbackHandler
	sessionStore    *session.CookieStore
}

// NewAuthHandler コンストラクタ
func NewAuthHandler(authService *auth0.AuthService, userService userservice.UserService, sessionStore *session.CookieStore, provider *oidc.Provider, config *auth0.Config) (*AuthHandler, error) {
	// Create OAuth2 config
	oauth2Config := authService.GetOAuth2Config()

	// Create callback handler with new UserService
	callbackHandler := auth0.NewCallbackHandler(oauth2Config, provider, userService)

	return &AuthHandler{
		authService:     authService,
		callbackHandler: callbackHandler,
		sessionStore:    sessionStore,
	}, nil
}

// SetupRoutes 認証関連のルーティング設定
func (h *AuthHandler) SetupRoutes(e *echo.Echo) {
	auth := e.Group("/auth")

	// Authentication flow endpoints
	auth.GET("/login", h.SignIn)      // Redirects to Auth0
	auth.GET("/callback", h.Callback) // Handles Auth0 callback
	auth.GET("/logout", h.Logout)     // Clears session and redirects to Auth0 logout
	auth.GET("/me", h.Me)             // Returns current user information
	
	// Keep API endpoints too
	apiAuth := e.Group("/api/auth")
	apiAuth.GET("/login-url", h.GetLoginURL) // Returns login URL for AJAX
}

// SignIn handles GET /auth/login - redirects to Auth0 Universal Login
func (h *AuthHandler) SignIn(c echo.Context) error {
	logCtx := logger.WithContext(c.Request().Context()).WithComponent(logger.ComponentAuth)
	
	// Generate Auth0 login URL with state
	state := "state" // TODO: Generate secure random state
	loginURL := h.authService.GetLoginURL(state)
	
	logCtx.Info("User initiated login",
		zap.String("action", "login_start"),
		zap.String("provider", "auth0"),
		zap.String("login_url", loginURL),
		zap.String("user_agent", c.Request().UserAgent()),
		zap.String("remote_addr", c.RealIP()),
	)

	// Redirect to Auth0
	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// Callback handles GET /auth/callback - processes Auth0 callback
func (h *AuthHandler) Callback(c echo.Context) error {
	start := time.Now()
	logCtx := logger.WithContext(c.Request().Context()).WithComponent(logger.ComponentAuth)
	
	// Get authorization code and state from query parameters
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errorParam := c.QueryParam("error")
	errorDesc := c.QueryParam("error_description")
	
	logCtx.Info("Processing Auth0 callback",
		zap.String("action", "callback_start"),
		zap.String("state", state),
		zap.Bool("has_code", code != ""),
		zap.String("error_param", errorParam),
		zap.String("remote_addr", c.RealIP()),
	)

	// Check for authentication errors
	if errorParam != "" {
		logCtx.Warn("Auth0 authentication error",
			zap.String("error", errorParam),
			zap.String("error_description", errorDesc),
			zap.Duration("duration", time.Since(start)),
		)
		// Redirect to frontend with error
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		errorURL := fmt.Sprintf("%s/?error=%s&error_description=%s", frontendURL, errorParam, errorDesc)
		return c.Redirect(http.StatusTemporaryRedirect, errorURL)
	}

	// Check if we have an authorization code
	if code == "" {
		logCtx.Warn("Auth0 callback missing authorization code",
			zap.Duration("duration", time.Since(start)),
		)
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		errorURL := fmt.Sprintf("%s/?error=no_code", frontendURL)
		return c.Redirect(http.StatusTemporaryRedirect, errorURL)
	}

	// Process callback
	logCtx.Info("Processing Auth0 token exchange")
	sessionData, err := h.callbackHandler.HandleCallback(c.Request().Context(), code, state)
	if err != nil {
		logCtx.Error("Auth0 callback processing failed",
			zap.Error(err),
			zap.Duration("duration", time.Since(start)),
		)
		// Redirect to frontend with error
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		errorURL := fmt.Sprintf("%s/?error=auth_failed", frontendURL)
		return c.Redirect(http.StatusTemporaryRedirect, errorURL)
	}

	// Create session cookie
	logCtx.Info("Creating user session",
		zap.String("user_id", sessionData.UserID),
		zap.String("auth0_user_id", sessionData.Auth0UserID),
	)
	if err := h.sessionStore.SetSession(c, *sessionData); err != nil {
		logCtx.Error("Failed to create session cookie",
			zap.Error(err),
			zap.String("user_id", sessionData.UserID),
			zap.Duration("duration", time.Since(start)),
		)
		// Redirect to frontend with error
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}
		errorURL := fmt.Sprintf("%s/?error=session_failed", frontendURL)
		return c.Redirect(http.StatusTemporaryRedirect, errorURL)
	}

	logCtx.Info("Auth0 authentication completed successfully",
		zap.String("action", "login_success"),
		zap.String("user_id", sessionData.UserID),
		zap.String("auth0_user_id", sessionData.Auth0UserID),
		zap.String("user_email", sessionData.Email),
		zap.Duration("duration", time.Since(start)),
	)

	// Redirect to frontend with success
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	successURL := fmt.Sprintf("%s/", frontendURL)
	return c.Redirect(http.StatusTemporaryRedirect, successURL)
}

// GetLoginURL handles GET /api/auth/login-url - returns login URL for AJAX calls
func (h *AuthHandler) GetLoginURL(c echo.Context) error {
	// Generate random state for CSRF protection
	state := "random-state-string" // TODO: Generate proper random state
	
	loginURL := h.authService.GetLoginURL(state)

	return c.JSON(http.StatusOK, map[string]string{
		"login_url": loginURL,
		"state":     state,
	})
}

// Me handles GET /auth/me - returns current user information
func (h *AuthHandler) Me(c echo.Context) error {
	logCtx := logger.WithContext(c.Request().Context()).WithComponent(logger.ComponentAuth)
	
	// Get session from cookie
	sessionData, err := h.sessionStore.GetSession(c)
	if err != nil {
		logCtx.Info("User info request without valid session",
			zap.String("action", "user_info_unauthenticated"),
			zap.String("remote_addr", c.RealIP()),
			zap.Error(err),
		)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	logCtx.Debug("User info request successful",
		zap.String("action", "user_info_success"),
		zap.String("user_id", sessionData.UserID),
		zap.String("auth0_user_id", sessionData.Auth0UserID),
	)

	// Return user information
	response := map[string]interface{}{
		"user_id":       sessionData.UserID,
		"auth0_user_id": sessionData.Auth0UserID,
		"email":         sessionData.Email,
		"name":          sessionData.Name,
		"expires_at":    sessionData.ExpiresAt,
		"authenticated": true,
	}

	return c.JSON(http.StatusOK, response)
}

// Logout handles GET /auth/logout - clears session and redirects to Auth0 logout
func (h *AuthHandler) Logout(c echo.Context) error {
	logCtx := logger.WithContext(c.Request().Context()).WithComponent(logger.ComponentAuth)
	
	// Get current session data for logging
	sessionData, sessionErr := h.sessionStore.GetSession(c)
	if sessionErr == nil {
		logCtx = logCtx.WithUserInfo(sessionData.UserID, sessionData.Auth0UserID, sessionData.Email, sessionData.Name)
	}
	
	logCtx.Info("User initiated logout",
		zap.String("action", "logout_start"),
		zap.Bool("had_valid_session", sessionErr == nil),
		zap.String("remote_addr", c.RealIP()),
	)

	// Clear session cookie
	h.sessionStore.ClearSession(c)

	// Get Auth0 logout URL
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Construct Auth0 logout URL
	auth0Config, _ := auth0.NewConfig()
	logoutURL := fmt.Sprintf("https://%s/v2/logout?client_id=%s&returnTo=%s",
		auth0Config.Domain,
		auth0Config.ClientID,
		frontendURL,
	)
	
	logCtx.Info("User logout completed",
		zap.String("action", "logout_success"),
		zap.String("redirect_url", logoutURL),
	)

	return c.Redirect(http.StatusTemporaryRedirect, logoutURL)
}