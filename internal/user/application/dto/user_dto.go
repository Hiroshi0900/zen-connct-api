package dto

import "time"

// RegisterFromAuthCommand 認証からのユーザー登録コマンド
type RegisterFromAuthCommand struct {
	Auth0UserID   string `json:"auth0_user_id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// RegisterUserRequest ユーザー登録リクエスト
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterUserResponse ユーザー登録レスポンス
type RegisterUserResponse struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateProfileRequest プロフィール更新リクエスト
type UpdateProfileRequest struct {
	DisplayName       string `json:"display_name"`
	Bio               string `json:"bio"`
	ProfileImageURL   string `json:"profile_image_url"`
}

// UpdateProfileResponse プロフィール更新レスポンス
type UpdateProfileResponse struct {
	UserID    string    `json:"user_id"`
	Profile   ProfileDTO `json:"profile"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUserProfileRequest ユーザープロフィール取得リクエスト
type GetUserProfileRequest struct {
	UserID string `json:"user_id"`
}

// GetUserProfileResponse ユーザープロフィール取得レスポンス
type GetUserProfileResponse struct {
	UserID    string     `json:"user_id"`
	Email     string     `json:"email"`
	Profile   ProfileDTO `json:"profile"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ProfileDTO プロフィール情報のDTO
type ProfileDTO struct {
	DisplayName     string `json:"display_name"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
}

// UserDTO ユーザー情報のDTO
type UserDTO struct {
	UserID      string     `json:"user_id"`
	Auth0UserID string     `json:"auth0_user_id"`
	Email       string     `json:"email"`
	Profile     ProfileDTO `json:"profile"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}