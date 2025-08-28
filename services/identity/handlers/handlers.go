package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"unified-commerce/services/identity/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// IdentityHandler handles HTTP requests for identity operations
type IdentityHandler struct {
	service   *service.IdentityService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewIdentityHandler creates a new identity handler
func NewIdentityHandler(service *service.IdentityService, logger *logger.Logger) *IdentityHandler {
	return &IdentityHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all identity routes
func (h *IdentityHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/logout", h.Logout)
			auth.POST("/refresh", h.RefreshToken)
			auth.POST("/forgot-password", h.ForgotPassword)
			auth.POST("/reset-password", h.ResetPassword)
			auth.POST("/verify-email", h.VerifyEmail)
		}

		// Protected routes (authentication required)
		authConfig := middleware.DefaultAuthConfig()
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(authConfig))
		{
			// User profile management
			users := protected.Group("/users")
			{
				users.GET("/me", h.GetProfile)
				users.PUT("/me", h.UpdateProfile)
				users.POST("/change-password", h.ChangePassword)
				users.DELETE("/me", h.DeleteAccount)
			}

			// Admin routes (require admin role)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin", "super_admin"))
			{
				admin.GET("/users", h.ListUsers)
				admin.GET("/users/:id", h.GetUser)
				admin.PUT("/users/:id", h.UpdateUser)
				admin.DELETE("/users/:id", h.DeleteUser)
				admin.POST("/users/:id/activate", h.ActivateUser)
				admin.POST("/users/:id/deactivate", h.DeactivateUser)
			}
		}
	}
}

// Register handles user registration
func (h *IdentityHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	response, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrEmailExists:
			httputil.Conflict(c, "Email already exists")
		case service.ErrUsernameExists:
			httputil.Conflict(c, "Username already exists")
		default:
			h.logger.WithError(err).Error("Registration failed")
			httputil.InternalServerError(c, "Registration failed")
		}
		return
	}

	httputil.Created(c, response, "User registered successfully")
}

// Login handles user authentication
func (h *IdentityHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	response, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			httputil.Unauthorized(c, "Invalid credentials")
		case service.ErrUserInactive:
			httputil.Forbidden(c, "Account is inactive")
		case service.ErrEmailNotVerified:
			httputil.Forbidden(c, "Email not verified")
		default:
			h.logger.WithError(err).Error("Login failed")
			httputil.InternalServerError(c, "Login failed")
		}
		return
	}

	httputil.Success(c, response, "Login successful")
}

// Logout handles user logout
func (h *IdentityHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		httputil.BadRequest(c, "Invalid authorization header")
		return
	}
	token := authHeader[7:]

	if err := h.service.Logout(c.Request.Context(), userID.(string), token); err != nil {
		h.logger.WithError(err).Error("Logout failed")
		httputil.InternalServerError(c, "Logout failed")
		return
	}

	httputil.Success(c, nil, "Logout successful")
}

// RefreshToken handles token refresh
func (h *IdentityHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	// TODO: Implement refresh token logic
	httputil.Success(c, gin.H{"message": "Refresh token endpoint - to be implemented"})
}

// GetProfile returns the current user's profile
func (h *IdentityHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	// TODO: Get user profile from service
	httputil.Success(c, gin.H{"user_id": userID, "message": "Profile endpoint - to be implemented"})
}

// UpdateProfile updates the current user's profile
func (h *IdentityHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	// TODO: Implement profile update logic
	httputil.Success(c, gin.H{"user_id": userID, "message": "Profile update endpoint - to be implemented"})
}

// ChangePassword handles password change
func (h *IdentityHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID.(string), req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			httputil.BadRequest(c, "Current password is incorrect")
		default:
			h.logger.WithError(err).Error("Password change failed")
			httputil.InternalServerError(c, "Password change failed")
		}
		return
	}

	httputil.Success(c, nil, "Password changed successfully")
}

// ForgotPassword initiates password reset process
func (h *IdentityHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	// TODO: Implement forgot password logic
	httputil.Success(c, nil, "Password reset email sent if account exists")
}

// ResetPassword resets password using reset token
func (h *IdentityHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	// TODO: Implement password reset logic
	httputil.Success(c, nil, "Password reset successful")
}

// VerifyEmail verifies email using verification token
func (h *IdentityHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Token string `json:"token" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	// TODO: Implement email verification logic
	httputil.Success(c, nil, "Email verified successfully")
}

// DeleteAccount handles account deletion
func (h *IdentityHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	// TODO: Implement account deletion logic
	httputil.Success(c, gin.H{"user_id": userID, "message": "Account deletion endpoint - to be implemented"})
}

// Admin endpoints

// ListUsers lists all users (admin only)
func (h *IdentityHandler) ListUsers(c *gin.Context) {
	pagination := httputil.GetPaginationParams(c)

	// TODO: Implement user listing logic
	httputil.SuccessWithMeta(c, []interface{}{}, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   0,
	}, "Users retrieved successfully")
}

// GetUser gets a specific user (admin only)
func (h *IdentityHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	// TODO: Implement get user logic
	httputil.Success(c, gin.H{"user_id": userID, "message": "Get user endpoint - to be implemented"})
}

// UpdateUser updates a specific user (admin only)
func (h *IdentityHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	// TODO: Implement user update logic
	httputil.Success(c, gin.H{"user_id": userID, "message": "Update user endpoint - to be implemented"})
}

// DeleteUser deletes a specific user (admin only)
func (h *IdentityHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	// TODO: Implement user deletion logic
	httputil.Success(c, gin.H{"user_id": userID, "message": "Delete user endpoint - to be implemented"})
}

// ActivateUser activates a user account (admin only)
func (h *IdentityHandler) ActivateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	// TODO: Implement user activation logic
	httputil.Success(c, gin.H{"user_id": userID, "message": "Activate user endpoint - to be implemented"})
}

// DeactivateUser deactivates a user account (admin only)
func (h *IdentityHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	// TODO: Implement user deactivation logic
	httputil.Success(c, gin.H{"user_id": userID, "message": "Deactivate user endpoint - to be implemented"})
}
