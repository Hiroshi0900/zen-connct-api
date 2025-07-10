package interfaces

import (
	"github.com/labstack/echo/v4"
	"zen-connect/internal/infrastructure/session"
	"zen-connect/internal/user/application/dto"
	"zen-connect/internal/user/application/usecase"
)

// UserHandler ユーザー関連のHTTPハンドラー
type UserHandler struct {
	registerUserUseCase     *usecase.RegisterUserUseCase
	updateProfileUseCase    *usecase.UpdateProfileUseCase
	getUserProfileUseCase   *usecase.GetUserProfileUseCase
}

// NewUserHandler コンストラクタ
func NewUserHandler(
	registerUserUseCase *usecase.RegisterUserUseCase,
	updateProfileUseCase *usecase.UpdateProfileUseCase,
	getUserProfileUseCase *usecase.GetUserProfileUseCase,
) *UserHandler {
	return &UserHandler{
		registerUserUseCase:     registerUserUseCase,
		updateProfileUseCase:    updateProfileUseCase,
		getUserProfileUseCase:   getUserProfileUseCase,
	}
}

// SetupRoutes ユーザー関連のルーティング設定
func (h *UserHandler) SetupRoutes(e *echo.Echo, sessionMiddleware *session.Middleware) {
	userGroup := e.Group("/users")
	
	// 現在のユーザー情報取得（認証が必要）
	userGroup.GET("/me", h.GetCurrentUser, sessionMiddleware.RequireAuth())
}

// GetCurrentUser 現在のユーザー情報取得
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	// セッションからユーザーIDを取得
	userID, ok := session.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return c.JSON(401, map[string]string{
			"error": "User not authenticated",
		})
	}
	
	// ユーザー情報を取得
	req := &dto.GetUserProfileRequest{
		UserID: userID,
	}
	
	response, err := h.getUserProfileUseCase.Execute(c.Request().Context(), req)
	if err != nil {
		return c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	}
	
	return c.JSON(200, response)
}