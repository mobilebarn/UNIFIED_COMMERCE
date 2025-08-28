package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// MetaInfo contains pagination and other metadata
type MetaInfo struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// PaginationParams holds pagination parameters from query string
type PaginationParams struct {
	Page    int `form:"page" binding:"omitempty,min=1"`
	PerPage int `form:"per_page" binding:"omitempty,min=1,max=100"`
}

// DefaultPagination returns default pagination parameters
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:    1,
		PerPage: 20,
	}
}

// GetPaginationParams extracts and validates pagination parameters from request
func GetPaginationParams(c *gin.Context) PaginationParams {
	var params PaginationParams

	// Set defaults
	params.Page = 1
	params.PerPage = 20

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse per_page
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 && perPage <= 100 {
			params.PerPage = perPage
		}
	}

	return params
}

// CalculateOffset calculates database offset from pagination params
func (p PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.PerPage
}

// CalculateTotalPages calculates total pages from total count
func (p PaginationParams) CalculateTotalPages(total int64) int {
	if total == 0 {
		return 0
	}
	return int((total + int64(p.PerPage) - 1) / int64(p.PerPage))
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}, message ...string) {
	response := Response{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(http.StatusOK, response)
}

// SuccessWithMeta sends a successful response with metadata
func SuccessWithMeta(c *gin.Context, data interface{}, meta *MetaInfo, message ...string) {
	response := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(http.StatusOK, response)
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}, message ...string) {
	response := Response{
		Success: true,
		Data:    data,
	}

	if len(message) > 0 {
		response.Message = message[0]
	} else {
		response.Message = "Resource created successfully"
	}

	c.JSON(http.StatusCreated, response)
}

// NoContent sends a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, details ...map[string]interface{}) {
	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	}

	if len(details) > 0 {
		response.Error.Details = details[0]
	}

	c.JSON(http.StatusBadRequest, response)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message ...string) {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: msg,
		},
	}

	c.JSON(http.StatusUnauthorized, response)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message ...string) {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "FORBIDDEN",
			Message: msg,
		},
	}

	c.JSON(http.StatusForbidden, response)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message ...string) {
	msg := "Resource not found"
	if len(message) > 0 {
		msg = message[0]
	}

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "NOT_FOUND",
			Message: msg,
		},
	}

	c.JSON(http.StatusNotFound, response)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string, details ...map[string]interface{}) {
	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "CONFLICT",
			Message: message,
		},
	}

	if len(details) > 0 {
		response.Error.Details = details[0]
	}

	c.JSON(http.StatusConflict, response)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message ...string) {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}

	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_ERROR",
			Message: msg,
		},
	}

	c.JSON(http.StatusInternalServerError, response)
}

// ValidationError sends a 422 Unprocessable Entity response for validation errors
func ValidationError(c *gin.Context, details map[string]interface{}) {
	response := Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: details,
		},
	}

	c.JSON(http.StatusUnprocessableEntity, response)
}
