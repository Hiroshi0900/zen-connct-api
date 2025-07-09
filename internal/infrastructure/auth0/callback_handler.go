package auth0

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"zen-connect/internal/infrastructure/session"
	"zen-connect/internal/user/domain"
)

// CallbackHandler handles OAuth2 callback from Auth0
type CallbackHandler struct {
	oauth2Config *oauth2.Config
	provider     *oidc.Provider
	userRepo     domain.UserRepository
}

// NewCallbackHandler creates a new callback handler
func NewCallbackHandler(oauth2Config *oauth2.Config, provider *oidc.Provider, userRepo domain.UserRepository) *CallbackHandler {
	return &CallbackHandler{
		oauth2Config: oauth2Config,
		provider:     provider,
		userRepo:     userRepo,
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
	// Try to find existing user by Auth0 user ID
	existingUser, err := h.userRepo.FindByAuth0UserID(claims.Sub)
	if err == nil {
		// User exists, update information if needed
		// Update profile if changed
		if existingUser.Profile().DisplayName() != claims.Name {
			existingUser.UpdateProfile(claims.Name, existingUser.Profile().Bio(), claims.Picture)
		}

		// Update email verification status
		if claims.EmailVerified && !existingUser.EmailVerified() {
			existingUser.VerifyEmail()
		}

		// Save updated user
		if err := h.userRepo.Save(existingUser); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		return existingUser, nil
	}

	// User doesn't exist, create new user
	email, err := domain.NewEmail(claims.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Create new user
	newUser := domain.NewUser(claims.Sub, email, claims.Name, claims.EmailVerified)
	
	// Set profile image if provided
	if claims.Picture != "" {
		newUser.UpdateProfile(claims.Name, "", claims.Picture)
	}

	// Save new user
	if err := h.userRepo.Save(newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return newUser, nil
}