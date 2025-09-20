package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system (can be merchant staff, customers, or admin)
type User struct {
	ID                string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email             string         `json:"email" gorm:"uniqueIndex;not null"`
	Username          string         `json:"username" gorm:"uniqueIndex"`
	PasswordHash      string         `json:"-" gorm:"not null"` // Never include in JSON responses
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	Phone             string         `json:"phone"`
	Avatar            string         `json:"avatar"`
	IsActive          bool           `json:"is_active" gorm:"default:true"`
	IsEmailVerified   bool           `json:"is_email_verified" gorm:"default:false"`
	EmailVerifiedAt   *time.Time     `json:"email_verified_at"`
	LastLoginAt       *time.Time     `json:"last_login_at"`
	PasswordChangedAt *time.Time     `json:"password_changed_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Roles    []UserRole    `json:"roles" gorm:"foreignKey:UserID"`
	Sessions []UserSession `json:"-" gorm:"foreignKey:UserID"`
	// Note: MerchantMembers relationship is managed by the Merchant Account service
}

// IsEntity marks User as a federation entity
func (u User) IsEntity() {}

// Role represents a role in the system
type Role struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Permissions []RolePermission `json:"permissions" gorm:"foreignKey:RoleID"`
	UserRoles   []UserRole       `json:"-" gorm:"foreignKey:RoleID"`
}

// Permission represents a permission in the system
type Permission struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"` // e.g., "products", "orders", "users"
	Action      string    `json:"action"`   // e.g., "create", "read", "update", "delete"
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	RolePermissions []RolePermission `json:"-" gorm:"foreignKey:PermissionID"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID        string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string     `json:"user_id" gorm:"not null;index"`
	RoleID    string     `json:"role_id" gorm:"not null;index"`
	GrantedBy string     `json:"granted_by"` // ID of user who granted this role
	GrantedAt time.Time  `json:"granted_at"`
	ExpiresAt *time.Time `json:"expires_at"` // Optional expiration

	// Relationships
	User Role `json:"user" gorm:"foreignKey:UserID"`
	Role Role `json:"role" gorm:"foreignKey:RoleID"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoleID       string    `json:"role_id" gorm:"not null;index"`
	PermissionID string    `json:"permission_id" gorm:"not null;index"`
	CreatedAt    time.Time `json:"created_at"`

	// Relationships
	Role       Role       `json:"role" gorm:"foreignKey:RoleID"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID"`
}

// UserSession represents an active user session
type UserSession struct {
	ID         string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     string     `json:"user_id" gorm:"not null;index"`
	Token      string     `json:"-" gorm:"uniqueIndex;not null"` // JWT token hash
	Type       string     `json:"type" gorm:"not null"`          // "access", "refresh"
	IsActive   bool       `json:"is_active" gorm:"default:true"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`

	// Additional session metadata
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
	Device    string `json:"device"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string     `json:"user_id" gorm:"not null;index"`
	Token     string     `json:"token" gorm:"uniqueIndex;not null"`
	IsUsed    bool       `json:"is_used" gorm:"default:false"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// EmailVerification represents an email verification request
type EmailVerification struct {
	ID         string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     string     `json:"user_id" gorm:"not null;index"`
	Email      string     `json:"email" gorm:"not null"`
	Token      string     `json:"token" gorm:"uniqueIndex;not null"`
	IsUsed     bool       `json:"is_used" gorm:"default:false"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	VerifiedAt *time.Time `json:"verified_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// AuditLog represents security and action audit logs
type AuditLog struct {
	ID         string                 `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     string                 `json:"user_id" gorm:"index"`
	Action     string                 `json:"action" gorm:"not null"` // "login", "logout", "password_change", etc.
	Resource   string                 `json:"resource"`               // What was affected
	ResourceID string                 `json:"resource_id"`            // ID of affected resource
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	Success    bool                   `json:"success"`
	ErrorMsg   string                 `json:"error_message"`
	Metadata   map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt  time.Time              `json:"created_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook for User model
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		// GORM will handle UUID generation via default:gen_random_uuid()
	}
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Username
	}
	return u.FirstName + " " + u.LastName
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, userRole := range u.Roles {
		if userRole.Role.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permissionName string) bool {
	for _, userRole := range u.Roles {
		for _, rolePermission := range userRole.Role.Permissions {
			if rolePermission.Permission.Name == permissionName {
				return true
			}
		}
	}
	return false
}

// GetRoleNames returns a slice of role names for the user
func (u *User) GetRoleNames() []string {
	roleNames := make([]string, len(u.Roles))
	for i, userRole := range u.Roles {
		roleNames[i] = userRole.Role.Name
	}
	return roleNames
}
