package graphql

// CustomerAction represents the type of customer action
type CustomerAction string

const (
	CustomerActionView           CustomerAction = "VIEW"
	CustomerActionClick          CustomerAction = "CLICK"
	CustomerActionAddToCart      CustomerAction = "ADD_TO_CART"
	CustomerActionRemoveFromCart CustomerAction = "REMOVE_FROM_CART"
	CustomerActionPurchase       CustomerAction = "PURCHASE"
	CustomerActionSearch         CustomerAction = "SEARCH"
	CustomerActionFilter         CustomerAction = "FILTER"
	CustomerActionSort           CustomerAction = "SORT"
	CustomerActionShare          CustomerAction = "SHARE"
	CustomerActionWishlist       CustomerAction = "WISHLIST"
	CustomerActionReview         CustomerAction = "REVIEW"
	CustomerActionReturn         CustomerAction = "RETURN"
)

// RecommendationType represents the type of recommendation
type RecommendationType string

const (
	RecommendationTypePopularity    RecommendationType = "POPULARITY"
	RecommendationTypeCollaborative RecommendationType = "COLLABORATIVE"
	RecommendationTypeContent       RecommendationType = "CONTENT"
	RecommendationTypeTrending      RecommendationType = "TRENDING"
	RecommendationTypeSimilar       RecommendationType = "SIMILAR"
	RecommendationTypePersonalized  RecommendationType = "PERSONALIZED"
	RecommendationTypeCategory      RecommendationType = "CATEGORY"
	RecommendationTypeBrand         RecommendationType = "BRAND"
)

// AnalyticsReportType represents the type of analytics report
type AnalyticsReportType string

const (
	AnalyticsReportTypeSales      AnalyticsReportType = "SALES"
	AnalyticsReportTypeTraffic    AnalyticsReportType = "TRAFFIC"
	AnalyticsReportTypeConversion AnalyticsReportType = "CONVERSION"
	AnalyticsReportTypeCustomer   AnalyticsReportType = "CUSTOMER"
	AnalyticsReportTypeInventory  AnalyticsReportType = "INVENTORY"
	AnalyticsReportTypeFinancial  AnalyticsReportType = "FINANCIAL"
	AnalyticsReportTypeMarketing  AnalyticsReportType = "MARKETING"
)
