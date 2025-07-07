package interfaces

import (
	"net/http"
	"zen-connect/internal/user/application"

	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	registerUserUseCase *application.RegisterUserUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(registerUserUseCase *application.RegisterUserUseCase) *UserHandler {
	return &UserHandler{
		registerUserUseCase: registerUserUseCase,
	}
}

// RegisterUser handles POST /api/users/register
func (h *UserHandler) RegisterUser(c echo.Context) error {
	var req application.RegisterUserRequest

	// Bind JSON request to struct
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email is required",
		})
	}

	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Password is required",
		})
	}

	// Execute use case
	response, err := h.registerUserUseCase.Execute(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Return successful response
	return c.JSON(http.StatusCreated, response)
}

// SetupRoutes sets up the user-related routes
func (h *UserHandler) SetupRoutes(e *echo.Echo) {
	api := e.Group("/api")
	users := api.Group("/users")

	users.POST("/register", h.RegisterUser)
}