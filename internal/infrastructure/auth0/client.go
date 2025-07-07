package auth0

import (
	"fmt"
	"os"
)

// Config holds Auth0 configuration
type Config struct {
	Domain       string
	ClientID     string
	ClientSecret string
	Audience     string
}

// NewConfig creates a new Auth0 configuration from environment variables
func NewConfig() (*Config, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN environment variable is required")
	}

	clientID := os.Getenv("AUTH0_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("AUTH0_CLIENT_ID environment variable is required")
	}

	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("AUTH0_CLIENT_SECRET environment variable is required")
	}

	audience := os.Getenv("AUTH0_AUDIENCE")
	if audience == "" {
		return nil, fmt.Errorf("AUTH0_AUDIENCE environment variable is required")
	}

	return &Config{
		Domain:       domain,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
	}, nil
}

// JWKSUrl returns the JWKS URL for the Auth0 tenant
func (c *Config) JWKSUrl() string {
	return fmt.Sprintf("https://%s/.well-known/jwks.json", c.Domain)
}

// IssuerURL returns the issuer URL for JWT validation
func (c *Config) IssuerURL() string {
	return fmt.Sprintf("https://%s/", c.Domain)
}