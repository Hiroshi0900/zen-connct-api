package session

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
)

// SessionData represents the data stored in the session cookie
type SessionData struct {
	UserID      string    `json:"uid"`      // Local DB user ID
	Auth0UserID string    `json:"sub"`      // Auth0 user ID
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	ExpiresAt   time.Time `json:"exp"`
}

// CookieStore manages encrypted session cookies
type CookieStore struct {
	secureCookie   *securecookie.SecureCookie
	cookieName     string
	cookieDomain   string
	cookiePath     string
	cookieSecure   bool
	cookieHTTPOnly bool
	cookieSameSite http.SameSite
	maxAge         int
}

// NewCookieStore creates a new cookie-based session store
func NewCookieStore() (*CookieStore, error) {
	// Get secret key from environment
	secretKey := os.Getenv("SESSION_SECRET")
	if secretKey == "" {
		return nil, fmt.Errorf("SESSION_SECRET environment variable is required")
	}

	// Ensure secret key is 32 bytes for AES-256
	if len(secretKey) != 32 {
		return nil, fmt.Errorf("SESSION_SECRET must be exactly 32 bytes")
	}

	// Create secure cookie with encryption and authentication
	sc := securecookie.New([]byte(secretKey), nil)
	sc.SetSerializer(securecookie.JSONEncoder{})

	// Get cookie settings from environment
	cookieName := getEnvOrDefault("SESSION_COOKIE_NAME", "zen_session")
	cookieDomain := getEnvOrDefault("SESSION_COOKIE_DOMAIN", "localhost")
	cookiePath := getEnvOrDefault("SESSION_COOKIE_PATH", "/")
	cookieSecure := getEnvBoolOrDefault("SESSION_COOKIE_SECURE", false)
	cookieHTTPOnly := getEnvBoolOrDefault("SESSION_COOKIE_HTTP_ONLY", true)
	
	// Parse SameSite setting
	sameSiteStr := getEnvOrDefault("SESSION_COOKIE_SAME_SITE", "lax")
	var sameSite http.SameSite
	switch sameSiteStr {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	maxAge := getEnvIntOrDefault("SESSION_MAX_AGE", 86400) // 24 hours default

	return &CookieStore{
		secureCookie:   sc,
		cookieName:     cookieName,
		cookieDomain:   cookieDomain,
		cookiePath:     cookiePath,
		cookieSecure:   cookieSecure,
		cookieHTTPOnly: cookieHTTPOnly,
		cookieSameSite: sameSite,
		maxAge:         maxAge,
	}, nil
}

// SetSession encrypts and stores session data in a cookie
func (cs *CookieStore) SetSession(c echo.Context, data SessionData) error {
	// Set expiration time
	data.ExpiresAt = time.Now().Add(time.Duration(cs.maxAge) * time.Second)

	// Encode session data
	encoded, err := cs.secureCookie.Encode(cs.cookieName, data)
	if err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	// Create cookie
	cookie := &http.Cookie{
		Name:     cs.cookieName,
		Value:    encoded,
		Domain:   cs.cookieDomain,
		Path:     cs.cookiePath,
		MaxAge:   cs.maxAge,
		Secure:   cs.cookieSecure,
		HttpOnly: cs.cookieHTTPOnly,
		SameSite: cs.cookieSameSite,
	}

	c.SetCookie(cookie)
	return nil
}

// GetSession retrieves and decrypts session data from cookie
func (cs *CookieStore) GetSession(c echo.Context) (*SessionData, error) {
	// Get cookie
	cookie, err := c.Cookie(cs.cookieName)
	if err != nil {
		return nil, fmt.Errorf("session cookie not found: %w", err)
	}

	// Decode session data
	var data SessionData
	err = cs.secureCookie.Decode(cs.cookieName, cookie.Value, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode session: %w", err)
	}

	// Check expiration
	if time.Now().After(data.ExpiresAt) {
		return nil, fmt.Errorf("session has expired")
	}

	return &data, nil
}

// ClearSession removes the session cookie
func (cs *CookieStore) ClearSession(c echo.Context) {
	cookie := &http.Cookie{
		Name:     cs.cookieName,
		Value:    "",
		Domain:   cs.cookieDomain,
		Path:     cs.cookiePath,
		MaxAge:   -1,
		Secure:   cs.cookieSecure,
		HttpOnly: cs.cookieHTTPOnly,
		SameSite: cs.cookieSameSite,
	}

	c.SetCookie(cookie)
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}