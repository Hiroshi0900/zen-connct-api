package usecase

import (
	"context"
	"zen-connect/internal/auth/application/dto"
	"zen-connect/internal/auth/application/service"
)

// VerifyTokenUseCase トークン検証ユースケース
type VerifyTokenUseCase struct {
	authService service.AuthService
}

// NewVerifyTokenUseCase コンストラクタ
func NewVerifyTokenUseCase(authService service.AuthService) *VerifyTokenUseCase {
	return &VerifyTokenUseCase{
		authService: authService,
	}
}

// Execute トークン検証を実行
func (uc *VerifyTokenUseCase) Execute(ctx context.Context, req *dto.VerifyTokenRequest) (*dto.VerifyTokenResponse, error) {
	// トークンを検証
	user, err := uc.authService.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	
	// DTOに変換して返す
	response := &dto.VerifyTokenResponse{
		UserID:   user.ID(),
		Email:    user.Email(),
		Sub:      user.Sub(),
		Name:     user.Name(),
		Picture:  user.Picture(),
		Verified: user.Verified(),
	}
	
	return response, nil
}