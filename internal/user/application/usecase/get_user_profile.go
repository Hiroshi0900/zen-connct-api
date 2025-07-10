package usecase

import (
	"context"
	"zen-connect/internal/user/application/dto"
	"zen-connect/internal/user/application/service"
)

// GetUserProfileUseCase ユーザープロフィール取得ユースケース
type GetUserProfileUseCase struct {
	userService service.UserService
}

// NewGetUserProfileUseCase コンストラクタ
func NewGetUserProfileUseCase(userService service.UserService) *GetUserProfileUseCase {
	return &GetUserProfileUseCase{
		userService: userService,
	}
}

// Execute ユーザープロフィール取得を実行
func (uc *GetUserProfileUseCase) Execute(ctx context.Context, req *dto.GetUserProfileRequest) (*dto.GetUserProfileResponse, error) {
	// ユーザーを取得
	user, err := uc.userService.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	
	// レスポンスを作成
	response := &dto.GetUserProfileResponse{
		UserID: user.ID(),
		Email:  user.Email().String(),
		Profile: dto.ProfileDTO{
			DisplayName:     user.Profile().DisplayName(),
			Bio:             user.Profile().Bio(),
			ProfileImageURL: user.Profile().ProfileImageURL(),
		},
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
	
	return response, nil
}