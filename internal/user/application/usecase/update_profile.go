package usecase

import (
	"context"
	"zen-connect/internal/user/application/dto"
	"zen-connect/internal/user/application/service"
)

// UpdateProfileUseCase プロフィール更新ユースケース
type UpdateProfileUseCase struct {
	userService service.UserService
}

// NewUpdateProfileUseCase コンストラクタ
func NewUpdateProfileUseCase(userService service.UserService) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		userService: userService,
	}
}

// Execute プロフィール更新を実行
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	// ユーザーを取得
	user, err := uc.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	// プロフィールを更新
	user.UpdateProfile(req.DisplayName, req.Bio, req.ProfileImageURL)
	
	// 更新を保存
	if err := uc.userService.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	
	// レスポンスを作成
	response := &dto.UpdateProfileResponse{
		UserID: user.ID(),
		Profile: dto.ProfileDTO{
			DisplayName:     user.Profile().DisplayName(),
			Bio:             user.Profile().Bio(),
			ProfileImageURL: user.Profile().ProfileImageURL(),
		},
		UpdatedAt: user.UpdatedAt(),
	}
	
	return response, nil
}