package domain

import "errors"

// 認証関連のエラー定義
var (
	// ErrInvalidToken 無効なトークン
	ErrInvalidToken = errors.New("invalid token")
	
	// ErrTokenExpired トークンの有効期限切れ
	ErrTokenExpired = errors.New("token expired")
	
	// ErrSessionNotFound セッションが見つからない
	ErrSessionNotFound = errors.New("session not found")
	
	// ErrSessionExpired セッションの有効期限切れ
	ErrSessionExpired = errors.New("session expired")
	
	// ErrUnauthorized 認証が必要
	ErrUnauthorized = errors.New("unauthorized")
	
	// ErrInvalidClaims 無効なクレーム
	ErrInvalidClaims = errors.New("invalid claims")
)