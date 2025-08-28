package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	httputil "unified-commerce/services/shared/http"
)

// Logger creates a structured logging middleware
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			"timestamp":  param.TimeStamp.Format(time.RFC3339),
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"user_agent": param.Request.UserAgent(),
			"request_id": param.Keys["request_id"],
		}).Info("HTTP Request")

		return ""
	})
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// CORS middleware for Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID     string   `json:"user_id"`
	MerchantID string   `json:"merchant_id,omitempty"`
	Email      string   `json:"email"`
	Roles      []string `json:"roles"`
	jwt.RegisteredClaims
}

// AuthConfig holds JWT authentication configuration
type AuthConfig struct {
	SecretKey     string
	TokenLookup   string // "header:Authorization" or "query:token" or "cookie:jwt"
	TokenHeadName string // "Bearer"
	SkipPaths     []string
}

// DefaultAuthConfig returns default authentication configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		SecretKey:     "unified-commerce-jwt-secret-key",
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		SkipPaths:     []string{"/health", "/metrics"},
	}
}

// JWTAuth creates a JWT authentication middleware
func JWTAuth(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for certain paths
		for _, path := range config.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		token := extractToken(c, config)
		if token == "" {
			httputil.Unauthorized(c, "Missing or invalid token")
			c.Abort()
			return
		}

		claims, err := validateToken(token, config.SecretKey)
		if err != nil {
			httputil.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("merchant_id", claims.MerchantID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			httputil.Forbidden(c, "No roles found")
			c.Abort()
			return
		}

		roles, ok := userRoles.([]string)
		if !ok {
			httputil.Forbidden(c, "Invalid roles format")
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, required := range requiredRoles {
			for _, userRole := range roles {
				if userRole == required {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			httputil.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireMerchant ensures the user is associated with a merchant
func RequireMerchant() gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID, exists := c.Get("merchant_id")
		if !exists || merchantID == "" {
			httputil.Forbidden(c, "Merchant context required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractToken extracts JWT token from request based on configuration
func extractToken(c *gin.Context, config *AuthConfig) string {
	parts := strings.Split(config.TokenLookup, ":")
	if len(parts) != 2 {
		return ""
	}

	switch parts[0] {
	case "header":
		auth := c.GetHeader(parts[1])
		if len(auth) > len(config.TokenHeadName)+1 && strings.HasPrefix(auth, config.TokenHeadName+" ") {
			return auth[len(config.TokenHeadName)+1:]
		}
	case "query":
		return c.Query(parts[1])
	case "cookie":
		if cookie, err := c.Cookie(parts[1]); err == nil {
			return cookie
		}
	}

	return ""
}

// validateToken validates a JWT token and returns claims
func validateToken(tokenString, secretKey string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

// Recovery middleware recovers from panics and returns a 500 error
func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get("request_id")

		logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"panic":      recovered,
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
		}).Error("Panic recovered")

		httputil.InternalServerError(c, "Internal server error")
	})
}

// RateLimit creates a simple rate limiting middleware
func RateLimit() gin.HandlerFunc {
	// This is a simple implementation. In production, you would use Redis or similar
	return func(c *gin.Context) {
		// TODO: Implement proper rate limiting with Redis
		c.Next()
	}
}
