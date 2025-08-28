package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"unified-commerce/services/merchant-account/models"
	"unified-commerce/services/merchant-account/repository"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/utils"
)

var (
	ErrMerchantNotFound     = errors.New("merchant not found")
	ErrMerchantExists       = errors.New("merchant already exists")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrInvalidSubscription  = errors.New("invalid subscription")
	ErrPlanNotFound         = errors.New("plan not found")
	ErrMemberNotFound       = errors.New("member not found")
	ErrInvalidRole          = errors.New("invalid role")
)

// MerchantService handles merchant account business logic
type MerchantService struct {
	repo   *repository.Repository
	logger *logger.Logger
}

// NewMerchantService creates a new merchant service
func NewMerchantService(repo *repository.Repository, logger *logger.Logger) *MerchantService {
	return &MerchantService{
		repo:   repo,
		logger: logger,
	}
}

// CreateMerchantRequest represents a merchant creation request
type CreateMerchantRequest struct {
	BusinessName string `json:"business_name" validate:"required"`
	LegalName    string `json:"legal_name"`
	BusinessType string `json:"business_type" validate:"required"`
	Industry     string `json:"industry"`
	PrimaryEmail string `json:"primary_email" validate:"required,email"`
	PrimaryPhone string `json:"primary_phone"`
	WebsiteURL   string `json:"website_url"`
	Description  string `json:"description"`
	
	// Owner information
	OwnerUserID string `json:"owner_user_id" validate:"required"`
	
	// Initial address
	Address *CreateAddressRequest `json:"address"`
	
	// Initial plan
	PlanID string `json:"plan_id"`
}

// CreateAddressRequest represents an address creation request
type CreateAddressRequest struct {
	Type       string `json:"type" validate:"required"`
	Label      string `json:"label"`
	Company    string `json:"company"`
	Address1   string `json:"address1" validate:"required"`
	Address2   string `json:"address2"`
	City       string `json:"city" validate:"required"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code" validate:"required"`
	Country    string `json:"country" validate:"required"`
	Phone      string `json:"phone"`
}

// UpdateMerchantRequest represents a merchant update request
type UpdateMerchantRequest struct {
	BusinessName string                 `json:"business_name"`
	LegalName    string                 `json:"legal_name"`
	Industry     string                 `json:"industry"`
	WebsiteURL   string                 `json:"website_url"`
	Description  string                 `json:"description"`
	LogoURL      string                 `json:"logo_url"`
	Settings     models.MerchantSettings `json:"settings"`
}

// InviteMemberRequest represents a member invitation request
type InviteMemberRequest struct {
	UserID      string   `json:"user_id" validate:"required"`
	Role        string   `json:"role" validate:"required"`
	Permissions []string `json:"permissions"`
}

// CreateMerchant creates a new merchant account
func (s *MerchantService) CreateMerchant(ctx context.Context, req *CreateMerchantRequest) (*models.Merchant, error) {
	// Check if merchant with email already exists
	existing, err := s.repo.Merchant.GetByEmail(ctx, req.PrimaryEmail)
	if err == nil && existing != nil {
		return nil, ErrMerchantExists
	}

	// Generate merchant ID and business-specific settings
	merchant := &models.Merchant{
		BusinessName: req.BusinessName,
		LegalName:    req.LegalName,
		BusinessType: req.BusinessType,
		Industry:     req.Industry,
		PrimaryEmail: req.PrimaryEmail,
		PrimaryPhone: req.PrimaryPhone,
		WebsiteURL:   req.WebsiteURL,
		Description:  req.Description,
		Status:       "pending",
		OnboardingStep: 1,
		Settings: models.MerchantSettings{
			Currency:      "USD",
			Timezone:      "UTC",
			WeightUnit:    "kg",
			DimensionUnit: "cm",
			OrderIDPrefix: utils.GenerateShortID(),
			NotificationPrefs: map[string]bool{
				"order_notifications":    true,
				"payment_notifications":  true,
				"inventory_notifications": true,
			},
			BusinessHours: map[string]string{
				"monday":    "09:00-17:00",
				"tuesday":   "09:00-17:00",
				"wednesday": "09:00-17:00",
				"thursday":  "09:00-17:00",
				"friday":    "09:00-17:00",
				"saturday":  "closed",
				"sunday":    "closed",
			},
		},
	}

	// Create merchant
	if err := s.repo.Merchant.Create(ctx, merchant); err != nil {
		s.logger.WithError(err).Error("Failed to create merchant")
		return nil, fmt.Errorf("failed to create merchant account")
	}

	// Add owner as first member
	ownerMember := &models.MerchantMember{
		MerchantID: merchant.ID,
		UserID:     req.OwnerUserID,
		Role:       "owner",
		Status:     "active",
		JoinedAt:   utils.TimePointer(time.Now()),
		Permissions: []string{"*"}, // Owner has all permissions
	}

	if err := s.repo.MerchantMember.Create(ctx, ownerMember); err != nil {
		s.logger.WithError(err).Error("Failed to create owner member")
		// Note: In production, you might want to rollback the merchant creation
	}

	// Add initial address if provided
	if req.Address != nil {
		address := &models.MerchantAddress{
			MerchantID: merchant.ID,
			Type:       req.Address.Type,
			Label:      req.Address.Label,
			Company:    req.Address.Company,
			Address1:   req.Address.Address1,
			Address2:   req.Address.Address2,
			City:       req.Address.City,
			State:      req.Address.State,
			PostalCode: req.Address.PostalCode,
			Country:    req.Address.Country,
			Phone:      req.Address.Phone,
			IsDefault:  true,
		}

		if err := s.repo.MerchantAddress.Create(ctx, address); err != nil {
			s.logger.WithError(err).Error("Failed to create initial address")
		}
	}

	// Create initial subscription if plan provided
	if req.PlanID != "" {
		plan, err := s.repo.Plan.GetByID(ctx, req.PlanID)
		if err == nil {
			subscription := &models.Subscription{
				MerchantID:         merchant.ID,
				PlanID:             plan.ID,
				Status:             "active",
				BillingCycle:       "monthly",
				Amount:             plan.MonthlyPrice,
				Currency:           plan.Currency,
				CurrentPeriodStart: time.Now(),
				CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
			}

			// Add trial period if plan has trial days
			if plan.TrialDays > 0 {
				trialEnd := time.Now().AddDate(0, 0, plan.TrialDays)
				subscription.TrialEndsAt = &trialEnd
			}

			if err := s.repo.Subscription.Create(ctx, subscription); err != nil {
				s.logger.WithError(err).Error("Failed to create initial subscription")
			}
		}
	}

	// Reload merchant with all relationships
	return s.repo.Merchant.GetByID(ctx, merchant.ID)
}

// GetMerchant retrieves a merchant by ID
func (s *MerchantService) GetMerchant(ctx context.Context, merchantID, userID string) (*models.Merchant, error) {
	merchant, err := s.repo.Merchant.GetByID(ctx, merchantID)
	if err != nil {
		return nil, ErrMerchantNotFound
	}

	// Check if user has access to this merchant
	if !merchant.HasMember(userID) {
		return nil, ErrUnauthorized
	}

	return merchant, nil
}

// GetMerchantsByUser retrieves all merchants for a user
func (s *MerchantService) GetMerchantsByUser(ctx context.Context, userID string) ([]models.Merchant, error) {
	return s.repo.Merchant.GetByUserID(ctx, userID)
}

// UpdateMerchant updates merchant information
func (s *MerchantService) UpdateMerchant(ctx context.Context, merchantID, userID string, req *UpdateMerchantRequest) (*models.Merchant, error) {
	merchant, err := s.repo.Merchant.GetByID(ctx, merchantID)
	if err != nil {
		return nil, ErrMerchantNotFound
	}

	// Check if user has permission to update
	memberRole := merchant.GetMemberRole(userID)
	if memberRole != "owner" && memberRole != "admin" {
		return nil, ErrUnauthorized
	}

	// Update fields
	if req.BusinessName != "" {
		merchant.BusinessName = req.BusinessName
	}
	if req.LegalName != "" {
		merchant.LegalName = req.LegalName
	}
	if req.Industry != "" {
		merchant.Industry = req.Industry
	}
	if req.WebsiteURL != "" {
		merchant.WebsiteURL = req.WebsiteURL
	}
	if req.Description != "" {
		merchant.Description = req.Description
	}
	if req.LogoURL != "" {
		merchant.LogoURL = req.LogoURL
	}

	// Update settings
	if req.Settings.Currency != "" {
		merchant.Settings.Currency = req.Settings.Currency
	}
	if req.Settings.Timezone != "" {
		merchant.Settings.Timezone = req.Settings.Timezone
	}
	if req.Settings.WeightUnit != "" {
		merchant.Settings.WeightUnit = req.Settings.WeightUnit
	}
	if req.Settings.DimensionUnit != "" {
		merchant.Settings.DimensionUnit = req.Settings.DimensionUnit
	}
	if len(req.Settings.NotificationPrefs) > 0 {
		merchant.Settings.NotificationPrefs = req.Settings.NotificationPrefs
	}
	if len(req.Settings.BusinessHours) > 0 {
		merchant.Settings.BusinessHours = req.Settings.BusinessHours
	}

	if err := s.repo.Merchant.Update(ctx, merchant); err != nil {
		s.logger.WithError(err).Error("Failed to update merchant")
		return nil, fmt.Errorf("failed to update merchant")
	}

	return merchant, nil
}

// ActivateMerchant activates a merchant account
func (s *MerchantService) ActivateMerchant(ctx context.Context, merchantID string) error {
	return s.repo.Merchant.UpdateStatus(ctx, merchantID, "active")
}

// SuspendMerchant suspends a merchant account
func (s *MerchantService) SuspendMerchant(ctx context.Context, merchantID string) error {
	return s.repo.Merchant.UpdateStatus(ctx, merchantID, "suspended")
}

// InviteMember invites a new member to the merchant account
func (s *MerchantService) InviteMember(ctx context.Context, merchantID, inviterUserID string, req *InviteMemberRequest) (*models.MerchantMember, error) {
	merchant, err := s.repo.Merchant.GetByID(ctx, merchantID)
	if err != nil {
		return nil, ErrMerchantNotFound
	}

	// Check if inviter has permission
	inviterRole := merchant.GetMemberRole(inviterUserID)
	if inviterRole != "owner" && inviterRole != "admin" {
		return nil, ErrUnauthorized
	}

	// Validate role
	validRoles := []string{"admin", "manager", "staff", "viewer"}
	if !utils.StringInSlice(req.Role, validRoles) {
		return nil, ErrInvalidRole
	}

	// Check if user is already a member
	existing, _ := s.repo.MerchantMember.GetByUserAndMerchant(ctx, req.UserID, merchantID)
	if existing != nil {
		return nil, fmt.Errorf("user is already a member")
	}

	member := &models.MerchantMember{
		MerchantID:  merchantID,
		UserID:      req.UserID,
		Role:        req.Role,
		Status:      "invited",
		Permissions: req.Permissions,
		InvitedBy:   inviterUserID,
		InvitedAt:   utils.TimePointer(time.Now()),
	}

	if err := s.repo.MerchantMember.Create(ctx, member); err != nil {
		s.logger.WithError(err).Error("Failed to create member invitation")
		return nil, fmt.Errorf("failed to invite member")
	}

	return member, nil
}

// AcceptInvitation accepts a member invitation
func (s *MerchantService) AcceptInvitation(ctx context.Context, memberID, userID string) error {
	member, err := s.repo.MerchantMember.GetByUserAndMerchant(ctx, userID, "")
	if err != nil {
		return ErrMemberNotFound
	}

	if member.Status != "invited" {
		return fmt.Errorf("invitation not pending")
	}

	return s.repo.MerchantMember.UpdateStatus(ctx, memberID, "active")
}

// RemoveMember removes a member from the merchant account
func (s *MerchantService) RemoveMember(ctx context.Context, merchantID, memberID, removerUserID string) error {
	merchant, err := s.repo.Merchant.GetByID(ctx, merchantID)
	if err != nil {
		return ErrMerchantNotFound
	}

	// Check if remover has permission
	removerRole := merchant.GetMemberRole(removerUserID)
	if removerRole != "owner" && removerRole != "admin" {
		return ErrUnauthorized
	}

	// Cannot remove owner
	for _, member := range merchant.Members {
		if member.ID == memberID && member.Role == "owner" {
			return fmt.Errorf("cannot remove owner")
		}
	}

	return s.repo.MerchantMember.UpdateStatus(ctx, memberID, "removed")
}

// ListPlans returns all available subscription plans
func (s *MerchantService) ListPlans(ctx context.Context) ([]models.Plan, error) {
	return s.repo.Plan.ListActive(ctx)
}

// CreateSubscription creates a new subscription for a merchant
func (s *MerchantService) CreateSubscription(ctx context.Context, merchantID, planID, userID string) (*models.Subscription, error) {
	merchant, err := s.repo.Merchant.GetByID(ctx, merchantID)
	if err != nil {
		return nil, ErrMerchantNotFound
	}

	// Check if user has permission
	if !merchant.IsOwner(userID) {
		return nil, ErrUnauthorized
	}

	plan, err := s.repo.Plan.GetByID(ctx, planID)
	if err != nil {
		return nil, ErrPlanNotFound
	}

	// Cancel existing active subscription
	if activeSubscription := merchant.GetActiveSubscription(); activeSubscription != nil {
		s.repo.Subscription.UpdateStatus(ctx, activeSubscription.ID, "cancelled")
	}

	subscription := &models.Subscription{
		MerchantID:         merchantID,
		PlanID:             planID,
		Status:             "active",
		BillingCycle:       "monthly",
		Amount:             plan.MonthlyPrice,
		Currency:           plan.Currency,
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
	}

	if err := s.repo.Subscription.Create(ctx, subscription); err != nil {
		s.logger.WithError(err).Error("Failed to create subscription")
		return nil, fmt.Errorf("failed to create subscription")
	}

	return subscription, nil
}

// GetSubscription retrieves merchant's subscription information
func (s *MerchantService) GetSubscription(ctx context.Context, merchantID, userID string) (*models.Subscription, error) {
	merchant, err := s.repo.Merchant.GetByID(ctx, merchantID)
	if err != nil {
		return nil, ErrMerchantNotFound
	}

	if !merchant.HasMember(userID) {
		return nil, ErrUnauthorized
	}

	return s.repo.Subscription.GetActiveByMerchantID(ctx, merchantID)
}