package domain

import "time"

// AuthenticatedUser 認証済みユーザーを表す
type AuthenticatedUser struct {
	id       string
	email    string
	sub      string    // Auth0のsub claim
	name     string
	picture  string
	verified bool
	authTime time.Time
}

// NewAuthenticatedUser コンストラクタ
func NewAuthenticatedUser(id, email, sub, name, picture string, verified bool) *AuthenticatedUser {
	return &AuthenticatedUser{
		id:       id,
		email:    email,
		sub:      sub,
		name:     name,
		picture:  picture,
		verified: verified,
		authTime: time.Now(),
	}
}

// IsSessionValid セッションの有効性チェック
func (u *AuthenticatedUser) IsSessionValid() bool {
	return time.Since(u.authTime) < 24*time.Hour
}

// ID ゲッター
func (u *AuthenticatedUser) ID() string {
	return u.id
}

// Email ゲッター
func (u *AuthenticatedUser) Email() string {
	return u.email
}

// Sub Auth0のsub claimを取得
func (u *AuthenticatedUser) Sub() string {
	return u.sub
}

// Name ゲッター
func (u *AuthenticatedUser) Name() string {
	return u.name
}

// Picture ゲッター
func (u *AuthenticatedUser) Picture() string {
	return u.picture
}

// Verified ゲッター
func (u *AuthenticatedUser) Verified() bool {
	return u.verified
}

// AuthTime ゲッター
func (u *AuthenticatedUser) AuthTime() time.Time {
	return u.authTime
}