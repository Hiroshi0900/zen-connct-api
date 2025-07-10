package dto

import "time"

// SessionData セッション情報のDTO
type SessionData struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CallbackRequest コールバックリクエストDTO
type CallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// CallbackResponse コールバックレスポンスDTO
type CallbackResponse struct {
	SessionData *SessionData `json:"session_data"`
	RedirectURL string       `json:"redirect_url"`
}

// VerifyTokenRequest トークン検証リクエストDTO
type VerifyTokenRequest struct {
	Token string `json:"token"`
}

// VerifyTokenResponse トークン検証レスポンスDTO
type VerifyTokenResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Sub      string `json:"sub"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Verified bool   `json:"verified"`
}

// LogoutRequest ログアウトリクエストDTO
type LogoutRequest struct {
	SessionID string `json:"session_id"`
}

// LogoutResponse ログアウトレスポンスDTO
type LogoutResponse struct {
	RedirectURL string `json:"redirect_url"`
}