package repository

import (
	"context"
	"time"

	"unified-commerce/services/identity/models"
	"unified-commerce/services/shared/database"
)

// UserRepository handles user data operations
type UserRepository struct {
	db *database.PostgresDB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.PostgresDB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.DB.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.DB.WithContext(ctx).
		Preload("Roles.Role.Permissions.Permission").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.DB.WithContext(ctx).
		Preload("Roles.Role.Permissions.Permission").
		First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.DB.WithContext(ctx).
		Preload("Roles.Role.Permissions.Permission").
		First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.DB.WithContext(ctx).Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.DB.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Get total count
	if err := r.db.DB.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.DB.WithContext(ctx).
		Preload("Roles.Role").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, total, err
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", now).Error
}

// UpdatePassword updates the user's password hash and timestamp
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_hash":       passwordHash,
			"password_changed_at": now,
		}).Error
}

// VerifyEmail marks a user's email as verified
func (r *UserRepository) VerifyEmail(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_email_verified": true,
			"email_verified_at": now,
		}).Error
}

// RoleRepository handles role data operations
type RoleRepository struct {
	db *database.PostgresDB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *database.PostgresDB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create creates a new role
func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	return r.db.DB.WithContext(ctx).Create(role).Error
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(ctx context.Context, id string) (*models.Role, error) {
	var role models.Role
	err := r.db.DB.WithContext(ctx).
		Preload("Permissions.Permission").
		First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName retrieves a role by name
func (r *RoleRepository) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.DB.WithContext(ctx).
		Preload("Permissions.Permission").
		First(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List retrieves all active roles
func (r *RoleRepository) List(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.DB.WithContext(ctx).
		Preload("Permissions.Permission").
		Where("is_active = ?", true).
		Find(&roles).Error
	return roles, err
}

// UserSessionRepository handles user session data operations
type UserSessionRepository struct {
	db *database.PostgresDB
}

// NewUserSessionRepository creates a new user session repository
func NewUserSessionRepository(db *database.PostgresDB) *UserSessionRepository {
	return &UserSessionRepository{db: db}
}

// Create creates a new user session
func (r *UserSessionRepository) Create(ctx context.Context, session *models.UserSession) error {
	return r.db.DB.WithContext(ctx).Create(session).Error
}

// GetByToken retrieves a session by token hash
func (r *UserSessionRepository) GetByToken(ctx context.Context, tokenHash string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.DB.WithContext(ctx).
		Preload("User").
		First(&session, "token = ? AND is_active = ? AND expires_at > ?", tokenHash, true, time.Now()).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateLastUsed updates the session's last used timestamp
func (r *UserSessionRepository) UpdateLastUsed(ctx context.Context, sessionID string) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("id = ?", sessionID).
		Update("last_used_at", now).Error
}

// DeactivateSession deactivates a session
func (r *UserSessionRepository) DeactivateSession(ctx context.Context, sessionID string) error {
	return r.db.DB.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("id = ?", sessionID).
		Update("is_active", false).Error
}

// DeactivateUserSessions deactivates all sessions for a user
func (r *UserSessionRepository) DeactivateUserSessions(ctx context.Context, userID string) error {
	return r.db.DB.WithContext(ctx).
		Model(&models.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
}

// CleanupExpiredSessions removes expired sessions
func (r *UserSessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	return r.db.DB.WithContext(ctx).
		Where("expires_at < ? OR is_active = ?", time.Now(), false).
		Delete(&models.UserSession{}).Error
}

// PasswordResetRepository handles password reset operations
type PasswordResetRepository struct {
	db *database.PostgresDB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *database.PostgresDB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset request
func (r *PasswordResetRepository) Create(ctx context.Context, reset *models.PasswordReset) error {
	return r.db.DB.WithContext(ctx).Create(reset).Error
}

// GetByToken retrieves a password reset by token
func (r *PasswordResetRepository) GetByToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.DB.WithContext(ctx).
		Preload("User").
		First(&reset, "token = ? AND is_used = ? AND expires_at > ?", token, false, time.Now()).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

// MarkAsUsed marks a password reset as used
func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).
		Model(&models.PasswordReset{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_used": true,
			"used_at": now,
		}).Error
}

// AuditLogRepository handles audit log operations
type AuditLogRepository struct {
	db *database.PostgresDB
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *database.PostgresDB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	return r.db.DB.WithContext(ctx).Create(log).Error
}

// GetByUserID retrieves audit logs for a specific user
func (r *AuditLogRepository) GetByUserID(ctx context.Context, userID string, offset, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := r.db.DB.WithContext(ctx).Model(&models.AuditLog{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, total, err
}

// Repository aggregates all repositories
type Repository struct {
	User          *UserRepository
	Role          *RoleRepository
	UserSession   *UserSessionRepository
	PasswordReset *PasswordResetRepository
	AuditLog      *AuditLogRepository
}

// NewRepository creates a new repository with all sub-repositories
func NewRepository(db *database.PostgresDB) *Repository {
	return &Repository{
		User:          NewUserRepository(db),
		Role:          NewRoleRepository(db),
		UserSession:   NewUserSessionRepository(db),
		PasswordReset: NewPasswordResetRepository(db),
		AuditLog:      NewAuditLogRepository(db),
	}
}

// DB returns underlying gorm DB for one-off internal operations (avoid in regular code paths)
func (r *Repository) DB() *database.PostgresDB {
	return r.User.db
}

// Migrate runs database migrations for all models
func (r *Repository) Migrate() error {
	db := r.User.db.DB
	
	// Drop merchant_members table if it exists (moved to Merchant Account service)
	db.Exec("DROP TABLE IF EXISTS merchant_members CASCADE")
	
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.UserSession{},
		// Note: MerchantMember is managed by the Merchant Account service
		&models.PasswordReset{},
		&models.EmailVerification{},
		&models.AuditLog{},
	)
}
