# Identity Service API Testing Guide

## Testing the Identity Service

Once the Identity Service is running, you can test it using these API calls.

### Prerequisites
- Identity Service running on port 8001
- PostgreSQL database initialized
- Infrastructure services running via Docker Compose

### API Endpoints

#### 1. Health Check
```bash
curl -X GET http://localhost:8001/health
```

Expected Response:
```json
{
  "service": "identity",
  "status": "healthy",
  "time": "2024-01-01T12:00:00Z",
  "checks": {
    "postgres": "healthy",
    "redis": "healthy"
  }
}
```

#### 2. User Registration
```bash
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "username": "johndoe",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }'
```

Expected Response:
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid-here",
      "email": "john.doe@example.com",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "is_active": true,
      "created_at": "2024-01-01T12:00:00Z"
    },
    "access_token": "jwt-token-here",
    "refresh_token": "refresh-token-here",
    "expires_in": 3600
  }
}
```

#### 3. User Login
```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
  }'
```

#### 4. Get User Profile (Authenticated)
```bash
curl -X GET http://localhost:8001/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 5. Change Password (Authenticated)
```bash
curl -X POST http://localhost:8001/api/v1/users/change-password \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "SecurePass123!",
    "new_password": "NewSecurePass456!"
  }'
```

#### 6. Logout (Authenticated)
```bash
curl -X POST http://localhost:8001/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### Admin Endpoints (Require Admin Role)

#### List All Users
```bash
curl -X GET http://localhost:8001/api/v1/admin/users \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE"
```

#### Get Specific User
```bash
curl -X GET http://localhost:8001/api/v1/admin/users/USER_ID_HERE \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE"
```

### Error Responses

#### Validation Error (422)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "validation": "Password must be at least 8 characters long"
    }
  }
}
```

#### Unauthorized (401)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid credentials"
  }
}
```

#### Conflict (409)
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT", 
    "message": "Email already exists"
  }
}
```

### Database Verification

You can verify the data is being stored correctly by connecting to PostgreSQL:

```bash
# Connect to PostgreSQL container
docker exec -it unified-commerce-postgres psql -U identity_user -d identity_service

# Check users table
SELECT id, email, username, first_name, last_name, is_active, created_at FROM users;

# Check roles table
SELECT * FROM roles;

# Check audit logs
SELECT user_id, action, success, created_at FROM audit_logs ORDER BY created_at DESC LIMIT 10;
```

### Testing Workflow

1. **Start Infrastructure**: `docker-compose up -d postgres redis`
2. **Start Identity Service**: `go run services/identity/cmd/server/main.go`
3. **Test Health Check**: Verify service is running
4. **Register User**: Create a new user account
5. **Login**: Authenticate and get JWT token
6. **Test Protected Endpoints**: Use JWT token for authenticated requests
7. **Test Admin Features**: Create admin user and test admin endpoints

This testing guide demonstrates the complete authentication and authorization flow of our unified commerce platform's Identity Service.