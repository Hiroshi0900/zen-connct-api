package auth0

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"zen-connect/internal/infrastructure/session"
	"zen-connect/internal/user/domain"
)

// CallbackHandler handles Auth0 OAuth2 callback
type CallbackHandler struct {
	config       *Config
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	userRepo     domain.UserRepository
	sessionStore *session.CookieStore
}

// NewCallbackHandler creates a new callback handler
func NewCallbackHandler(config *Config, oauth2Config *oauth2.Config, verifier *oidc.IDTokenVerifier, userRepo domain.UserRepository, sessionStore *session.CookieStore) *CallbackHandler {
	return &CallbackHandler{
		config:       config,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		userRepo:     userRepo,
		sessionStore: sessionStore,
	}
}

// Claims represents the claims in the ID token
type Claims struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// HandleCallback processes the OAuth2 callback from Auth0
func (h *CallbackHandler) HandleCallback(ctx context.Context, code, state string) (*session.SessionData, error) {
	// Exchange code for tokens
	token, err := h.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Extract and verify ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in token response")
	}

	idToken, err := h.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	// Extract claims from ID token
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	// Create or update user
	user, err := h.createOrUpdateUser(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update user: %w", err)
	}

	// Create session data
	sessionData := &session.SessionData{
		UserID:      user.ID(),
		Auth0UserID: claims.Sub,
		Email:       claims.Email,
		Name:        claims.Name,
	}

	return sessionData, nil
}

// createOrUpdateUser creates a new user or updates existing user
func (h *CallbackHandler) createOrUpdateUser(claims Claims) (domain.User, error) {
	// Try to find existing user by Auth0 user ID
	existingUser, err := h.userRepo.FindByAuth0UserID(claims.Sub)
	if err == nil {
		// User exists, update information if needed
		if activeUser, ok := domain.IsActive(existingUser); ok {
			// Update email if changed
			if activeUser.Email().Value() != claims.Email {
				email, err := domain.NewEmail(claims.Email)
				if err != nil {
					return nil, fmt.Errorf("invalid email: %w", err)
				}
				activeUser.UpdateEmail(email)
			}

			// Update profile if changed
			if activeUser.Profile().DisplayName() != claims.Name {
				activeUser.UpdateProfile(claims.Name, activeUser.Profile().Bio(), claims.Picture)
			}

			// Update email verification status
			if claims.EmailVerified && !activeUser.EmailVerified() {
				activeUser.VerifyEmail()
			}

			// Save updated user
			if err := h.userRepo.Save(activeUser); err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}

			return activeUser, nil
		}
	}

	// Try to find existing user by email (for migration from provisional to active)
	email, err := domain.NewEmail(claims.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	existingUser, err = h.userRepo.FindByEmail(email)
	if err == nil {
		// Found provisional user, activate it
		if provisionalUser, ok := domain.IsProvisional(existingUser); ok {
			activeUser := provisionalUser.ActivateUser(claims.Sub, claims.Name, claims.EmailVerified)
			
			// Update profile image
			if claims.Picture != "" {
				activeUser.UpdateProfile(claims.Name, "", claims.Picture)
			}

			if err := h.userRepo.Save(activeUser); err != nil {
				return nil, fmt.Errorf("failed to activate user: %w", err)
			}

			return activeUser, nil
		}
	}

	// Create new active user
	newUser := domain.NewActiveUser(claims.Sub, email, claims.Name, claims.EmailVerified)
	
	// Set profile image if provided
	if claims.Picture != "" {
		newUser.UpdateProfile(claims.Name, "", claims.Picture)
	}

	if err := h.userRepo.Save(newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}