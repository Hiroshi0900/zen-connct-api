package domain

import (
	"time"
)

// Session セッション値オブジェクト
type Session struct {
	id        string
	userID    string
	expiresAt time.Time
	createdAt time.Time
}

// NewSession セッションのコンストラクタ
func NewSession(id, userID string, duration time.Duration) *Session {
	now := time.Now()
	return &Session{
		id:        id,
		userID:    userID,
		expiresAt: now.Add(duration),
		createdAt: now,
	}
}

// ID ゲッター
func (s *Session) ID() string {
	return s.id
}

// UserID ゲッター
func (s *Session) UserID() string {
	return s.userID
}

// ExpiresAt ゲッター
func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

// CreatedAt ゲッター
func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

// IsExpired セッションの有効性チェック
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// IsValid セッションの有効性チェック
func (s *Session) IsValid() bool {
	return !s.IsExpired()
}