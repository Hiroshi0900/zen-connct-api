package service

import (
	"context"
	"zen-connect/internal/user/application/dto"
	"zen-connect/internal/user/domain"
)

// userServiceImpl UserServiceの実装
type userServiceImpl struct {
	userRepo domain.UserRepository
}

// NewUserService UserServiceのコンストラクタ
func NewUserService(userRepo domain.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// RegisterOrUpdateUserFromAuth 認証後のユーザー登録・更新
func (s *userServiceImpl) RegisterOrUpdateUserFromAuth(ctx context.Context, cmd dto.RegisterFromAuthCommand) (*domain.User, error) {
	// 既存ユーザーの検索
	user, err := s.userRepo.FindByAuth0UserID(cmd.Auth0UserID)
	if err != nil && err != domain.ErrUserNotFound {
		return nil, err
	}
	
	if user != nil {
		// 既存ユーザーの更新
		user.UpdateProfile(cmd.Name, "", cmd.Picture)
		
		if err := s.userRepo.Save(user); err != nil {
			return nil, err
		}
		
		return user, nil
	}
	
	// 新規ユーザーの作成
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}
	
	user = domain.NewUser(cmd.Auth0UserID, email, cmd.Name, cmd.EmailVerified)
	
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetUserByAuth0ID Auth0 IDでユーザーを取得
func (s *userServiceImpl) GetUserByAuth0ID(ctx context.Context, auth0UserID string) (*domain.User, error) {
	return s.userRepo.FindByAuth0UserID(auth0UserID)
}

// GetUserByID ユーザーIDでユーザーを取得
func (s *userServiceImpl) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateUser ユーザー情報を更新
func (s *userServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Save(user)
}