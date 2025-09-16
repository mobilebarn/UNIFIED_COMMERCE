package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"unified-commerce/services/identity/models"
	"unified-commerce/services/identity/repository"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
	"unified-commerce/services/shared/utils"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrEmailNotVerified   = errors.New("email not verified")
)

// IdentityService handles authentication and authorization business logic
type IdentityService struct {
	repo      *repository.Repository
	logger    *logger.Logger
	jwtSecret string
	tracer    interface{} // This would be *service.Tracer but we can"t import it here due to circular dependency
}

// NewIdentityService creates a new identity service
func NewIdentityService(repo *repository.Repository, logger *logger.Logger, jwtSecret string, tracer interface{}) *IdentityService {
	return &IdentityService{
		repo:      repo,
		logger:    logger,
		jwtSecret: jwtSecret,
		tracer:    tracer,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
}

// Register creates a new user account
func (s *IdentityService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// For now, we'll implement without tracing to avoid complexity
	// In a full implementation, you would add tracing here

	// Check if email already exists
	existingUser, err := s.repo.User.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailExists
	}

	// Check if username already exists (if provided)
	if req.Username != "" {
		existingUser, err = s.repo.User.GetByUsername(ctx, req.Username)
		if err == nil && existingUser != nil {
			return nil, ErrUsernameExists
		}
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("Failed to hash password")
		return nil, fmt.Errorf("failed to create user account")
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		IsActive:     true,
	}

	if err := s.repo.User.Create(ctx, user); err != nil {
		s.logger.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user account")
	}

	// Assign default role (customer)
	defaultRole, err := s.repo.Role.GetByName(ctx, "customer")
	if err == nil {
		userRole := &models.UserRole{
			UserID:    user.ID,
			RoleID:    defaultRole.ID,
			GrantedAt: time.Now(),
		}
		// Note: You'd need to add this method to the repository
		// s.repo.UserRole.Create(ctx, userRole)
		_ = userRole // Placeholder to avoid unused variable error
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate tokens")
		return nil, fmt.Errorf("failed to complete registration")
	}

	// Create session
	sessionToken := s.hashToken(accessToken)
	session := &models.UserSession{
		UserID:    user.ID,
		Token:     sessionToken,
		Type:      "access",
		IsActive:  true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.UserSession.Create(ctx, session); err != nil {
		s.logger.WithError(err).Error("Failed to create session")
	}

	// Log audit event
	s.logAuditEvent(ctx, user.ID, "user_register", "user", user.ID, true, "", nil)

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// Login authenticates a user and returns tokens
func (s *IdentityService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	var user *models.User
	var err error

	// Find user by email or username
	if req.Email != "" {
		user, err = s.repo.User.GetByEmail(ctx, req.Email)
	} else if req.Username != "" {
		user, err = s.repo.User.GetByUsername(ctx, req.Username)
	} else {
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		s.logAuditEvent(ctx, "", "user_login", "user", "", false, "User not found", nil)
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		s.logAuditEvent(ctx, user.ID, "user_login", "user", user.ID, false, "User inactive", nil)
		return nil, ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logAuditEvent(ctx, user.ID, "user_login", "user", user.ID, false, "Invalid password", nil)
		return nil, ErrInvalidCredentials
	}

	// Update last login
	if err := s.repo.User.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.WithError(err).Error("Failed to update last login")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate tokens")
		return nil, fmt.Errorf("failed to complete login")
	}

	// Create session
	sessionToken := s.hashToken(accessToken)
	session := &models.UserSession{
		UserID:    user.ID,
		Token:     sessionToken,
		Type:      "access",
		IsActive:  true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.UserSession.Create(ctx, session); err != nil {
		s.logger.WithError(err).Error("Failed to create session")
	}

	// Log successful login
	s.logAuditEvent(ctx, user.ID, "user_login", "user", user.ID, true, "", nil)

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// Logout invalidates user session
func (s *IdentityService) Logout(ctx context.Context, userID, accessToken string) error {
	// Hash token to find session
	sessionToken := s.hashToken(accessToken)

	// Find and deactivate session
	session, err := s.repo.UserSession.GetByToken(ctx, sessionToken)
	if err == nil {
		if err := s.repo.UserSession.DeactivateSession(ctx, session.ID); err != nil {
			s.logger.WithError(err).Error("Failed to deactivate session")
		}
	}

	// Log logout event
	s.logAuditEvent(ctx, userID, "user_logout", "user", userID, true, "", nil)

	return nil
}

// ValidateToken validates a JWT token and returns user information
func (s *IdentityService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims := &middleware.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token is expired
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrTokenExpired
	}

	// Get user from database
	user, err := s.repo.User.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if session is still active
	sessionToken := s.hashToken(tokenString)
	session, err := s.repo.UserSession.GetByToken(ctx, sessionToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Update session last used
	if err := s.repo.UserSession.UpdateLastUsed(ctx, session.ID); err != nil {
		s.logger.WithError(err).Error("Failed to update session last used")
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *IdentityService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Get user
	user, err := s.repo.User.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Validate new password
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	// Update password
	if err := s.repo.User.UpdatePassword(ctx, userID, string(passwordHash)); err != nil {
		return fmt.Errorf("failed to update password")
	}

	// Invalidate all user sessions (force re-login)
	if err := s.repo.UserSession.DeactivateUserSessions(ctx, userID); err != nil {
		s.logger.WithError(err).Error("Failed to deactivate user sessions")
	}

	// Log password change
	s.logAuditEvent(ctx, userID, "password_change", "user", userID, true, "", nil)

	return nil
}

// generateTokens generates access and refresh tokens for a user
func (s *IdentityService) generateTokens(user *models.User) (string, string, error) {
	// Access token (1 hour)
	accessClaims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  user.GetRoleNames(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token (7 days)
	refreshClaims := &middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// hashToken creates a hash of a token for storage
func (s *IdentityService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// logAuditEvent logs an audit event
func (s *IdentityService) logAuditEvent(ctx context.Context, userID, action, resource, resourceID string, success bool, errorMsg string, metadata map[string]interface{}) {
	auditLog := &models.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Success:    success,
		ErrorMsg:   errorMsg,
		Metadata:   metadata,
	}

	if err := s.repo.AuditLog.Create(ctx, auditLog); err != nil {
		s.logger.WithError(err).Error("Failed to create audit log")
	}
}

// GetUserByID retrieves a user by ID
func (s *IdentityService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.User.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUser updates user information
func (s *IdentityService) UpdateUser(ctx context.Context, user *models.User) error {
	if err := s.repo.User.Update(ctx, user); err != nil {
		s.logger.WithError(err).Error("Failed to update user")
		return fmt.Errorf("failed to update user")
	}

	// Log audit event
	s.logAuditEvent(ctx, user.ID, "user_update", "user", user.ID, true, "", nil)

	return nil
}

// ListUsers retrieves paginated list of users
func (s *IdentityService) ListUsers(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	users, total, err := s.repo.User.List(ctx, offset, limit)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list users")
		return nil, 0, fmt.Errorf("failed to list users")
	}

	return users, total, nil
}

// DeleteUser soft deletes a user
func (s *IdentityService) DeleteUser(ctx context.Context, userID string) error {
	// Check if user exists
	_, err := s.repo.User.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Deactivate all user sessions
	if err := s.repo.UserSession.DeactivateUserSessions(ctx, userID); err != nil {
		s.logger.WithError(err).Error("Failed to deactivate user sessions")
	}

	// Soft delete user
	if err := s.repo.User.Delete(ctx, userID); err != nil {
		s.logger.WithError(err).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user")
	}

	// Log audit event
	s.logAuditEvent(ctx, userID, "user_delete", "user", userID, true, "", nil)

	return nil
}

// RefreshToken validates a refresh token and returns new access/refresh tokens
func (s *IdentityService) RefreshToken(ctx context.Context, refreshTokenString string) (*AuthResponse, error) {
	claims := &middleware.JWTClaims{}

	token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if token is expired
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrTokenExpired
	}

	// Get user from database
	user, err := s.repo.User.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.generateTokens(user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate tokens")
		return nil, fmt.Errorf("failed to refresh tokens")
	}

	// Create new session
	sessionToken := s.hashToken(accessToken)
	session := &models.UserSession{
		UserID:    user.ID,
		Token:     sessionToken,
		Type:      "access",
		IsActive:  true,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.UserSession.Create(ctx, session); err != nil {
		s.logger.WithError(err).Error("Failed to create session")
	}

	// Log token refresh
	s.logAuditEvent(ctx, user.ID, "token_refresh", "user", user.ID, true, "", nil)

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

// ForgotPassword initiates password reset process
func (s *IdentityService) ForgotPassword(ctx context.Context, email string) error {
	// Get user by email
	user, err := s.repo.User.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists for security reasons
		s.logger.Debug("Forgot password request for non-existent email")
		return nil // Always return success to prevent email enumeration
	}

	// Generate reset token
	resetToken, err := utils.GenerateRandomString(32)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate reset token")
		return fmt.Errorf("failed to initiate password reset")
	}

	// Create password reset record
	passwordReset := &models.PasswordReset{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	if err := s.repo.PasswordReset.Create(ctx, passwordReset); err != nil {
		s.logger.WithError(err).Error("Failed to create password reset")
		return fmt.Errorf("failed to initiate password reset")
	}

	// TODO: Send password reset email
	// In a real implementation, you would send an email with the reset token
	s.logger.Info("Password reset token generated", "user_id", user.ID, "token", resetToken)

	// Log audit event
	s.logAuditEvent(ctx, user.ID, "password_reset_request", "user", user.ID, true, "", nil)

	return nil
}

// ResetPassword resets password using reset token
func (s *IdentityService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate new password
	if err := utils.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Get password reset record
	passwordReset, err := s.repo.PasswordReset.GetByToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	// Update user password
	if err := s.repo.User.UpdatePassword(ctx, passwordReset.UserID, string(passwordHash)); err != nil {
		return fmt.Errorf("failed to update password")
	}

	// Mark reset token as used
	if err := s.repo.PasswordReset.MarkAsUsed(ctx, passwordReset.ID); err != nil {
		s.logger.WithError(err).Error("Failed to mark password reset as used")
	}

	// Deactivate all user sessions (force re-login)
	if err := s.repo.UserSession.DeactivateUserSessions(ctx, passwordReset.UserID); err != nil {
		s.logger.WithError(err).Error("Failed to deactivate user sessions")
	}

	// Log audit event
	s.logAuditEvent(ctx, passwordReset.UserID, "password_reset", "user", passwordReset.UserID, true, "", nil)

	return nil
}

// VerifyEmail verifies email using verification token
func (s *IdentityService) VerifyEmail(ctx context.Context, token string) error {
	// TODO: Implement email verification logic
	// In a real implementation, you would:
	// 1. Look up email verification record by token
	// 2. Check if token is not expired
	// 3. Mark user's email as verified
	// 4. Mark verification token as used

	// For now, we'll implement a basic version
	// This is a placeholder - you'd need to add EmailVerification repository methods
	s.logger.Info("Email verification requested", "token", token)

	// TODO: Add proper implementation once EmailVerificationRepository methods are added
	return fmt.Errorf("email verification not fully implemented")
}

// ActivateUser activates a user account (admin only)
func (s *IdentityService) ActivateUser(ctx context.Context, userID string) error {
	// Get user
	user, err := s.repo.User.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Update user status
	user.IsActive = true
	if err := s.repo.User.Update(ctx, user); err != nil {
		s.logger.WithError(err).Error("Failed to activate user")
		return fmt.Errorf("failed to activate user")
	}

	// Log audit event
	s.logAuditEvent(ctx, userID, "user_activate", "user", userID, true, "", nil)

	return nil
}

// DeactivateUser deactivates a user account (admin only)
func (s *IdentityService) DeactivateUser(ctx context.Context, userID string) error {
	// Get user
	user, err := s.repo.User.GetByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Update user status
	user.IsActive = false
	if err := s.repo.User.Update(ctx, user); err != nil {
		s.logger.WithError(err).Error("Failed to deactivate user")
		return fmt.Errorf("failed to deactivate user")
	}

	// Deactivate all user sessions
	if err := s.repo.UserSession.DeactivateUserSessions(ctx, userID); err != nil {
		s.logger.WithError(err).Error("Failed to deactivate user sessions")
	}

	// Log audit event
	s.logAuditEvent(ctx, userID, "user_deactivate", "user", userID, true, "", nil)

	return nil
}

// startSpan is a helper method to start a tracing span
// Note: This is a simplified implementation. In a real implementation,
// you would pass the actual tracer from the base service
func (s *IdentityService) startSpan(ctx context.Context, name string) (context.Context, interface{}) {
	// This is a placeholder implementation
	// In a real implementation, you would use the actual tracer
	// For now, we just return the context as-is
	return ctx, nil
}
