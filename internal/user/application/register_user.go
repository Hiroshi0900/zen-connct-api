package application

import (
	"errors"
	"zen-connect/internal/user/domain"
)

// RegisterUserRequest represents the input for user registration
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterUserResponse represents the output of user registration
type RegisterUserResponse struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
}

// RegisterUserUseCase handles user registration business logic
type RegisterUserUseCase struct {
	userRepo domain.UserRepository
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase
func NewRegisterUserUseCase(userRepo domain.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: userRepo,
	}
}

// Execute performs user registration
func (uc *RegisterUserUseCase) Execute(req RegisterUserRequest) (*RegisterUserResponse, error) {
	// Validate and create Email value object
	email, err := domain.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	_, err = uc.userRepo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user (unverified)
	user := domain.NewUser(email.Value(), email, "", false)

	// Save user
	if err := uc.userRepo.Save(user); err != nil {
		return nil, err
	}

	// Create response
	response := &RegisterUserResponse{
		UserID:    user.ID(),
		Email:     user.Email().Value(),
		State:     "UnverifiedUser",
		CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}