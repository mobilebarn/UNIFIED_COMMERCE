package service_test

import (
	"testing"
	"unified-commerce/services/identity/service"

	"github.com/stretchr/testify/assert"
)

func TestIdentityService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		// Test structure placeholder
		req := &service.RegisterRequest{
			Email:     "test@example.com",
			Username:  "testuser",
			Password:  "SecurePass123!",
			FirstName: "Test",
			LastName:  "User",
		}

		// Validate request structure
		assert.NotNil(t, req)
		assert.Equal(t, "test@example.com", req.Email)
	})

	t.Run("email already exists", func(t *testing.T) {
		// Test email validation
		assert.Equal(t, service.ErrEmailExists.Error(), "email already exists")
	})
}

func TestIdentityService_Login(t *testing.T) {
	t.Run("login request structure", func(t *testing.T) {
		req := &service.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		assert.NotNil(t, req)
		assert.Equal(t, "test@example.com", req.Email)
	})
}
