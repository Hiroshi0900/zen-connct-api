package usecase

import (
	"context"
	"zen-connect/internal/user/application/dto"
	"zen-connect/internal/user/application/service"
	"zen-connect/internal/user/domain"
)

// RegisterUserUseCase ユーザー登録ユースケース
type RegisterUserUseCase struct {
	userService service.UserService
}

// NewRegisterUserUseCase コンストラクタ
func NewRegisterUserUseCase(userService service.UserService) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userService: userService,
	}
}

// Execute ユーザー登録を実行
func (uc *RegisterUserUseCase) Execute(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	// バリデーション
	if req.Email == "" || req.Password == "" {
		return nil, domain.ErrInvalidInput
	}
	
	// メールアドレスの値オブジェクトを作成
	email, err := domain.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}
	
	// ユーザーエンティティを作成
	user := domain.NewUser("", email, "", false)
	
	// ユーザーを保存
	if err := uc.userService.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	
	// レスポンスを作成
	response := &dto.RegisterUserResponse{
		UserID:    user.ID(),
		Email:     user.Email().String(),
		CreatedAt: user.CreatedAt(),
	}
	
	return response, nil
}