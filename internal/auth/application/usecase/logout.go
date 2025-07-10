package usecase

import (
	"context"
	"zen-connect/internal/auth/application/dto"
	"zen-connect/internal/auth/application/service"
)

// LogoutUseCase ログアウトユースケース
type LogoutUseCase struct {
	sessionService service.SessionService
}

// NewLogoutUseCase コンストラクタ
func NewLogoutUseCase(sessionService service.SessionService) *LogoutUseCase {
	return &LogoutUseCase{
		sessionService: sessionService,
	}
}

// Execute ログアウトを実行
func (uc *LogoutUseCase) Execute(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error) {
	// セッションを無効化
	err := uc.sessionService.InvalidateSession(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}
	
	// レスポンスを作成
	response := &dto.LogoutResponse{
		RedirectURL: "/login",
	}
	
	return response, nil
}