package auth0

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"zen-connect/internal/infrastructure/session"
	"zen-connect/internal/user/application/dto"
	"zen-connect/internal/user/application/service"
	"zen-connect/internal/user/domain"
)

// CallbackHandler handles OAuth2 callback from Auth0
type CallbackHandler struct {
	oauth2Config *oauth2.Config
	provider     *oidc.Provider
	userService  service.UserService
}

// NewCallbackHandler creates a new callback handler
func NewCallbackHandler(oauth2Config *oauth2.Config, provider *oidc.Provider, userService service.UserService) *CallbackHandler {
	return &CallbackHandler{
		oauth2Config: oauth2Config,
		provider:     provider,
		userService:  userService,
	}
}

// Claims represents ID token claims
type Claims struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
}

// HandleCallback processes the OAuth2 callback and returns session data
func (h *CallbackHandler) HandleCallback(ctx context.Context, code, state string) (*session.SessionData, error) {
	// Exchange authorization code for tokens
	oauth2Token, err := h.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Extract the ID Token from OAuth2 token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	// Parse and verify ID Token
	verifier := h.provider.Verifier(&oidc.Config{ClientID: h.oauth2Config.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	// Extract claims
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Create or update user
	user, err := h.createOrUpdateUser(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update user: %w", err)
	}

	// Create session data
	sessionData := &session.SessionData{
		UserID:       user.ID(),
		Auth0UserID:  user.Auth0UserID(),
		Email:        user.Email().String(),
		Name:         user.Profile().DisplayName(),
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		ExpiresAt:    oauth2Token.Expiry,
	}

	return sessionData, nil
}

// createOrUpdateUser creates a new user or updates existing user
func (h *CallbackHandler) createOrUpdateUser(claims Claims) (*domain.User, error) {
	// Use new UserService to register or update user from auth
	cmd := dto.RegisterFromAuthCommand{
		Auth0UserID:   claims.Sub,
		Email:         claims.Email,
		Name:          claims.Name,
		Picture:       claims.Picture,
		EmailVerified: claims.EmailVerified,
	}
	
	user, err := h.userService.RegisterOrUpdateUserFromAuth(context.Background(), cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to register or update user: %w", err)
	}
	
	return user, nil
}