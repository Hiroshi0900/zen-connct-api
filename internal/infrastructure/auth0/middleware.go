package auth0

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware provides JWT authentication middleware for Auth0
type AuthMiddleware struct {
	validator *validator.Validator
}

// NewAuthMiddleware creates a new Auth0 authentication middleware
func NewAuthMiddleware(config *Config) (*AuthMiddleware, error) {
	// Parse JWKS URL
	jwksURL, err := url.Parse(config.JWKSUrl())
	if err != nil {
		return nil, err
	}

	// Create a JWKS provider with cache duration
	keyFunc := jwks.NewCachingProvider(jwksURL, 5*time.Minute)

	// Set up the validator
	jwtValidator, err := validator.New(
		keyFunc.KeyFunc,
		validator.RS256,
		config.IssuerURL(),
		[]string{config.Audience},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
	)
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		validator: jwtValidator,
	}, nil
}

// CustomClaims represents the custom claims in the JWT token
type CustomClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Validate validates the custom claims
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// RequireAuth returns a middleware function that validates JWT tokens
func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// Check Bearer prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header must start with 'Bearer '",
				})
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Token is required",
				})
			}

			// Validate token
			validatedClaims, err := m.validator.ValidateToken(c.Request().Context(), token)
			if err != nil {
				log.Printf("Token validation failed: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			// Extract custom claims
			customClaims, ok := validatedClaims.(*validator.ValidatedClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token claims",
				})
			}

			// Add user info to context
			ctx := context.WithValue(c.Request().Context(), "user_id", customClaims.RegisteredClaims.Subject)
			if customClaims.CustomClaims != nil {
				if claims, ok := customClaims.CustomClaims.(*CustomClaims); ok {
					ctx = context.WithValue(ctx, "user_email", claims.Email)
					ctx = context.WithValue(ctx, "user_name", claims.Name)
				}
			}

			// Update request context
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
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