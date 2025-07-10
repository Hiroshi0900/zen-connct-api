package service

import (
	"context"
	"zen-connect/internal/user/application/dto"
	"zen-connect/internal/user/domain"
)

//go:generate mockgen -source=user_service.go -destination=../../../mocks/user_service_mock.go -package=mocks

// UserService ユーザーサービスのインターフェース
type UserService interface {
	// RegisterOrUpdateUserFromAuth 認証後のユーザー登録・更新
	RegisterOrUpdateUserFromAuth(ctx context.Context, cmd dto.RegisterFromAuthCommand) (*domain.User, error)
	
	// GetUserByAuth0ID Auth0 IDでユーザーを取得
	GetUserByAuth0ID(ctx context.Context, auth0UserID string) (*domain.User, error)
	
	// GetUserByID ユーザーIDでユーザーを取得
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	
	// UpdateUser ユーザー情報を更新
	UpdateUser(ctx context.Context, user *domain.User) error
}