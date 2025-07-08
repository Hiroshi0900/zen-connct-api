package auth0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// AuthService provides Auth0 authentication operations
type AuthService struct {
	config       *Config
	oauth2Config *oauth2.Config
	provider     *oidc.Provider
}

// NewAuthService creates a new Auth0 authentication service
func NewAuthService(config *Config) (*AuthService, error) {
	// Initialize OIDC provider
	provider, err := oidc.NewProvider(context.Background(), config.IssuerURL())
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// Configure OAuth2
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  os.Getenv("API_URL") + "/auth/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &AuthService{
		config:       config,
		oauth2Config: oauth2Config,
		provider:     provider,
	}, nil
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents Auth0 login response
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Login authenticates user with email/password using Resource Owner Password Grant
// Note: This requires enabling "Password" grant type in Auth0 Application settings
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	url := fmt.Sprintf("https://%s/oauth/token", s.config.Domain)
	
	payload := map[string]string{
		"grant_type": "password",
		"username":   req.Email,
		"password":   req.Password,
		"audience":   s.config.Audience,
		"client_id":  s.config.ClientID,
		"client_secret": s.config.ClientSecret,
		"scope":      "openid profile email",
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return nil, fmt.Errorf("login failed with status: %d, error: %v", resp.StatusCode, errorResp)
		}
		return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &loginResp, nil
}

// GetLoginURL generates Universal Login URL using oauth2 library with audience
func (s *AuthService) GetLoginURL(state string) string {
	// Test without audience first
	return s.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	
	// Add audience parameter for Auth0 API access (commented out for testing)
	// return s.oauth2Config.AuthCodeURL(state, 
	// 	oauth2.AccessTypeOffline,
	// 	oauth2.SetAuthURLParam("audience", s.config.Audience),
	// )
}

// GetOAuth2Config returns the OAuth2 configuration
func (s *AuthService) GetOAuth2Config() *oauth2.Config {
	return s.oauth2Config
}