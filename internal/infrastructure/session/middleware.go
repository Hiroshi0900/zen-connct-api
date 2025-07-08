package session

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Middleware provides session-based authentication middleware
type Middleware struct {
	cookieStore *CookieStore
}

// NewMiddleware creates a new session middleware
func NewMiddleware(cookieStore *CookieStore) *Middleware {
	return &Middleware{
		cookieStore: cookieStore,
	}
}

// RequireAuth returns a middleware function that validates session cookies
func (m *Middleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session from cookie
			sessionData, err := m.cookieStore.GetSession(c)
			if err != nil {
				log.Printf("Session validation failed: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authentication required",
				})
			}

			// Add session data to context
			ctx := context.WithValue(c.Request().Context(), "session", sessionData)
			ctx = context.WithValue(ctx, "user_id", sessionData.UserID)
			ctx = context.WithValue(ctx, "auth0_user_id", sessionData.Auth0UserID)
			ctx = context.WithValue(ctx, "user_email", sessionData.Email)
			ctx = context.WithValue(ctx, "user_name", sessionData.Name)

			// Update request context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// OptionalAuth returns a middleware function that optionally validates session cookies
// This middleware doesn't fail if no session is found, but adds session data to context if available
func (m *Middleware) OptionalAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get session from cookie (ignore errors)
			sessionData, err := m.cookieStore.GetSession(c)
			if err == nil {
				// Add session data to context
				ctx := context.WithValue(c.Request().Context(), "session", sessionData)
				ctx = context.WithValue(ctx, "user_id", sessionData.UserID)
				ctx = context.WithValue(ctx, "auth0_user_id", sessionData.Auth0UserID)
				ctx = context.WithValue(ctx, "user_email", sessionData.Email)
				ctx = context.WithValue(ctx, "user_name", sessionData.Name)

				// Update request context
				c.SetRequest(c.Request().WithContext(ctx))
			}

			return next(c)
		}
	}
}

// Context helper functions to extract session data from context

// GetSessionFromContext extracts session data from context
func GetSessionFromContext(ctx context.Context) (*SessionData, bool) {
	session, ok := ctx.Value("session").(*SessionData)
	return session, ok
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}

// GetAuth0UserIDFromContext extracts Auth0 user ID from context
func GetAuth0UserIDFromContext(ctx context.Context) (string, bool) {
	auth0UserID, ok := ctx.Value("auth0_user_id").(string)
	return auth0UserID, ok
}

// GetUserEmailFromContext extracts user email from context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("user_email").(string)
	return email, ok
}

// GetUserNameFromContext extracts user name from context
func GetUserNameFromContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value("user_name").(string)
	return name, ok
}