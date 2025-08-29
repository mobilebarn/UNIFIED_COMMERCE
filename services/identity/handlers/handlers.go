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

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	response, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrInvalidToken, service.ErrTokenExpired:
			httputil.Unauthorized(c, "Invalid or expired refresh token")
		default:
			h.logger.WithError(err).Error("Token refresh failed")
			httputil.InternalServerError(c, "Token refresh failed")
		}
		return
	}

	httputil.Success(c, response, "Token refreshed successfully")
}

// GetProfile returns the current user's profile
func (h *IdentityHandler) GetProfile(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	user, err := h.service.ValidateToken(c.Request.Context(), c.GetHeader("Authorization")[7:]) // Remove "Bearer "
	if err != nil {
		switch err {
		case service.ErrInvalidToken, service.ErrTokenExpired:
			httputil.Unauthorized(c, "Invalid or expired token")
		default:
			h.logger.WithError(err).Error("Failed to get user profile")
			httputil.InternalServerError(c, "Failed to get user profile")
		}
		return
	}

	// Remove sensitive information
	user.PasswordHash = ""

	httputil.Success(c, user, "Profile retrieved successfully")
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

	// Get current user
	user, err := h.service.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user")
		httputil.InternalServerError(c, "Failed to update profile")
		return
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	// Update user in database
	if err := h.service.UpdateUser(c.Request.Context(), user); err != nil {
		h.logger.WithError(err).Error("Failed to update user")
		httputil.InternalServerError(c, "Failed to update profile")
		return
	}

	// Remove sensitive information
	user.PasswordHash = ""

	httputil.Success(c, user, "Profile updated successfully")
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

	if err := h.service.ForgotPassword(c.Request.Context(), req.Email); err != nil {
		// Don't reveal if user exists for security reasons
		h.logger.WithError(err).Debug("Forgot password request failed")
	}

	// Always return success to prevent email enumeration
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

	if err := h.service.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		switch err {
		case service.ErrInvalidToken:
			httputil.BadRequest(c, "Invalid or expired reset token")
		case service.ErrUserNotFound:
			httputil.BadRequest(c, "Invalid reset token")
		default:
			h.logger.WithError(err).Error("Password reset failed")
			httputil.InternalServerError(c, "Password reset failed")
		}
		return
	}

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

	if err := h.service.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		switch err {
		case service.ErrInvalidToken:
			httputil.BadRequest(c, "Invalid or expired verification token")
		case service.ErrUserNotFound:
			httputil.BadRequest(c, "Invalid verification token")
		default:
			h.logger.WithError(err).Error("Email verification failed")
			httputil.InternalServerError(c, "Email verification failed")
		}
		return
	}

	httputil.Success(c, nil, "Email verified successfully")
}

// DeleteAccount handles account deletion
func (h *IdentityHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), userID.(string)); err != nil {
		switch err {
		case service.ErrUserNotFound:
			httputil.NotFound(c, "User not found")
		default:
			h.logger.WithError(err).Error("Failed to delete user account")
			httputil.InternalServerError(c, "Failed to delete account")
		}
		return
	}

	httputil.Success(c, nil, "Account deleted successfully")
}

// Admin endpoints

// ListUsers lists all users (admin only)
func (h *IdentityHandler) ListUsers(c *gin.Context) {
	pagination := httputil.GetPaginationParams(c)

	users, total, err := h.service.ListUsers(c.Request.Context(), pagination.CalculateOffset(), pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list users")
		httputil.InternalServerError(c, "Failed to retrieve users")
		return
	}

	// Remove sensitive information from all users
	for i := range users {
		users[i].PasswordHash = ""
	}

	meta := &httputil.MetaInfo{
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total),
	}

	httputil.SuccessWithMeta(c, users, meta, "Users retrieved successfully")
}

// GetUser gets a specific user (admin only)
func (h *IdentityHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			httputil.NotFound(c, "User not found")
		default:
			h.logger.WithError(err).Error("Failed to get user")
			httputil.InternalServerError(c, "Failed to retrieve user")
		}
		return
	}

	// Remove sensitive information
	user.PasswordHash = ""

	httputil.Success(c, user, "User retrieved successfully")
}

// UpdateUser updates a specific user (admin only)
func (h *IdentityHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		IsActive  *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	// Get current user
	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			httputil.NotFound(c, "User not found")
		default:
			h.logger.WithError(err).Error("Failed to get user")
			httputil.InternalServerError(c, "Failed to update user")
		}
		return
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update user in database
	if err := h.service.UpdateUser(c.Request.Context(), user); err != nil {
		h.logger.WithError(err).Error("Failed to update user")
		httputil.InternalServerError(c, "Failed to update user")
		return
	}

	// Remove sensitive information
	user.PasswordHash = ""

	httputil.Success(c, user, "User updated successfully")
}

// DeleteUser deletes a specific user (admin only)
func (h *IdentityHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	// Prevent admin from deleting themselves
	currentUserID, exists := c.Get("user_id")
	if exists && currentUserID.(string) == userID {
		httputil.BadRequest(c, "Cannot delete your own account")
		return
	}

	if err := h.service.DeleteUser(c.Request.Context(), userID); err != nil {
		switch err {
		case service.ErrUserNotFound:
			httputil.NotFound(c, "User not found")
		default:
			h.logger.WithError(err).Error("Failed to delete user")
			httputil.InternalServerError(c, "Failed to delete user")
		}
		return
	}

	httputil.Success(c, nil, "User deleted successfully")
}

// ActivateUser activates a user account (admin only)
func (h *IdentityHandler) ActivateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	if err := h.service.ActivateUser(c.Request.Context(), userID); err != nil {
		switch err {
		case service.ErrUserNotFound:
			httputil.NotFound(c, "User not found")
		default:
			h.logger.WithError(err).Error("Failed to activate user")
			httputil.InternalServerError(c, "Failed to activate user")
		}
		return
	}

	httputil.Success(c, nil, "User activated successfully")
}

// DeactivateUser deactivates a user account (admin only)
func (h *IdentityHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		httputil.BadRequest(c, "User ID is required")
		return
	}

	// Prevent admin from deactivating themselves
	currentUserID, exists := c.Get("user_id")
	if exists && currentUserID.(string) == userID {
		httputil.BadRequest(c, "Cannot deactivate your own account")
		return
	}

	if err := h.service.DeactivateUser(c.Request.Context(), userID); err != nil {
		switch err {
		case service.ErrUserNotFound:
			httputil.NotFound(c, "User not found")
		default:
			h.logger.WithError(err).Error("Failed to deactivate user")
			httputil.InternalServerError(c, "Failed to deactivate user")
		}
		return
	}

	httputil.Success(c, nil, "User deactivated successfully")
}
