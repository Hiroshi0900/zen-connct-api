package service

import (
	"context"
	"zen-connect/internal/auth/domain"
)

//go:generate mockgen -source=auth_service.go -destination=../../../mocks/auth_service_mock.go -package=mocks

// AuthService 認証サービスのインターフェース
type AuthService interface {
	// VerifyToken トークンを検証し、認証済みユーザーを返す
	VerifyToken(ctx context.Context, token string) (*domain.AuthenticatedUser, error)
	
	// GetUserInfo ユーザー情報を取得する
	GetUserInfo(ctx context.Context, userID string) (*domain.AuthenticatedUser, error)
	
	// RefreshToken トークンを更新する
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}