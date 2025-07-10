package interfaces

import (
	"context"
	"net/http"
	"strings"
	
	"github.com/labstack/echo/v4"
	"zen-connect/internal/auth/application/service"
	"zen-connect/internal/auth/domain"
)

// AuthMiddleware 認証ミドルウェア
type AuthMiddleware struct {
	authService service.AuthService
}

// NewAuthMiddleware コンストラクタ
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth 認証を必要とするミドルウェア
func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーからトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "authorization header required",
				})
			}
			
			// "Bearer "プレフィックスを除去
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization format",
				})
			}
			
			// トークンを検証
			user, err := m.authService.VerifyToken(c.Request().Context(), token)
			if err != nil {
				if err == domain.ErrInvalidToken || err == domain.ErrTokenExpired {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "invalid or expired token",
					})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "authentication failed",
				})
			}
			
			// ユーザー情報をコンテキストに設定
			ctx := context.WithValue(c.Request().Context(), "user", user)
			c.SetRequest(c.Request().WithContext(ctx))
			
			return next(c)
		}
	}
}

// GetAuthenticatedUser コンテキストから認証済みユーザーを取得
func GetAuthenticatedUser(c echo.Context) (*domain.AuthenticatedUser, error) {
	user, ok := c.Request().Context().Value("user").(*domain.AuthenticatedUser)
	if !ok {
		return nil, domain.ErrUnauthorized
	}
	return user, nil
}