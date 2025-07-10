package usecase

import (
	"context"
	"zen-connect/internal/auth/application/dto"
	"zen-connect/internal/auth/application/service"
	"zen-connect/internal/auth/domain"
)

// CreateSessionUseCase セッション作成ユースケース
type CreateSessionUseCase struct {
	authService    service.AuthService
	sessionService service.SessionService
}

// NewCreateSessionUseCase コンストラクタ
func NewCreateSessionUseCase(
	authService service.AuthService,
	sessionService service.SessionService,
) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		authService:    authService,
		sessionService: sessionService,
	}
}

// Execute セッション作成を実行
func (uc *CreateSessionUseCase) Execute(ctx context.Context, user *domain.AuthenticatedUser) (*dto.SessionData, error) {
	// セッションを作成
	session, err := uc.sessionService.CreateSession(ctx, user)
	if err != nil {
		return nil, err
	}
	
	// DTOに変換して返す
	sessionData := &dto.SessionData{
		SessionID: session.ID(),
		UserID:    session.UserID(),
		ExpiresAt: session.ExpiresAt(),
	}
	
	return sessionData, nil
}