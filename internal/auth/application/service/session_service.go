package service

import (
	"context"
	"zen-connect/internal/auth/domain"
)

//go:generate mockgen -source=session_service.go -destination=../../../mocks/session_service_mock.go -package=mocks

// SessionService セッション管理サービスのインターフェース
type SessionService interface {
	// CreateSession セッションを作成する
	CreateSession(ctx context.Context, user *domain.AuthenticatedUser) (*domain.Session, error)
	
	// GetSession セッションを取得する
	GetSession(ctx context.Context, sessionID string) (*domain.Session, error)
	
	// InvalidateSession セッションを無効化する
	InvalidateSession(ctx context.Context, sessionID string) error
	
	// RefreshSession セッションを更新する
	RefreshSession(ctx context.Context, sessionID string) (*domain.Session, error)
}