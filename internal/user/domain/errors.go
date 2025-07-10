package domain

import "errors"

// User domain errors
var (
	// ErrInvalidInput 無効な入力
	ErrInvalidInput = errors.New("invalid input")
	
	// ErrEmailAlreadyExists メールアドレスが既に存在
	ErrEmailAlreadyExists = errors.New("email already exists")
	
	// ErrInvalidEmail 無効なメールアドレス
	ErrInvalidEmail = errors.New("invalid email")
)