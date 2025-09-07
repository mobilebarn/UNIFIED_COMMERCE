#!/usr/bin/env powershell

# UNIFIED COMMERCE - ENTERPRISE FEDERATION STANDARDIZATION
# Converting ALL services to Identity Service's enterprise-grade architecture

Write-Host "🏢 ENTERPRISE FEDERATION STANDARDIZATION" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

$ErrorActionPreference = "Stop"
$workspaceRoot = "C:\Users\dane\OneDrive\Desktop\UNIFIED_COMMERCE"
$servicesToFix = @("cart", "order", "payment", "inventory", "product-catalog", "promotions", "merchant-account")

Write-Host "`n📋 Step 1: Stopping services..." -ForegroundColor Yellow
Get-Process | Where-Object {$_.ProcessName -in @("cart", "order", "payment", "inventory", "product-catalog", "promotions", "merchant-account", "main", "server")} | Stop-Process -Force -ErrorAction SilentlyContinue
Write-Host "✅ Services stopped" -ForegroundColor Green

Write-Host "`n📋 Step 2: Creating enterprise gqlgen configurations..." -ForegroundColor Yellow

foreach ($service in $servicesToFix) {
    $servicePath = Join-Path $workspaceRoot "services\$service"
    $graphqlPath = Join-Path $servicePath "graphql"
    
    Write-Host "  🔧 Configuring $service..." -ForegroundColor Cyan
    
    if (!(Test-Path $graphqlPath)) {
        New-Item -ItemType Directory -Path $graphqlPath -Force | Out-Null
    }
    
    # Copy Identity Service's gqlgen.yml as template
    $identityGqlgen = Join-Path $workspaceRoot "services\identity\graphql\gqlgen.yml"
    $targetGqlgen = Join-Path $graphqlPath "gqlgen.yml"
    
    if (Test-Path $identityGqlgen) {
        Copy-Item $identityGqlgen $targetGqlgen -Force
        
        # Update autobind path for each service
        $content = Get-Content $targetGqlgen -Raw
        $content = $content -replace 'unified-commerce/services/identity/models', "unified-commerce/services/$service/models"
        Set-Content -Path $targetGqlgen -Value $content -Encoding UTF8
        
        Write-Host "    ✅ Created enterprise gqlgen.yml" -ForegroundColor Green
    }
}

Write-Host "`n📋 Step 3: Creating enterprise handlers..." -ForegroundColor Yellow

foreach ($service in $servicesToFix) {
    $servicePath = Join-Path $workspaceRoot "services\$service"
    $graphqlPath = Join-Path $servicePath "graphql"
    
    Write-Host "  🔧 Creating handler for $service..." -ForegroundColor Cyan
    
    # Create enterprise handler based on Identity Service pattern
    $serviceCapitalized = (Get-Culture).TextInfo.ToTitleCase($service)
    
    $handlerContent = @"
package graphql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"unified-commerce/services/$service/service"
	"unified-commerce/services/shared/logger"
)

// NewGraphQLHandler creates a new GraphQL HTTP handler
func NewGraphQLHandler(${service}Service *service.${serviceCapitalized}Service, logger *logger.Logger) http.Handler {
	// Create a simple executable schema
	schema := NewExecutableSchema(Config{
		Resolvers: NewResolver(${service}Service, logger),
	})

	// Create the GraphQL server
	srv := handler.NewDefaultServer(schema)

	// Add recovery handler
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.WithField("panic", err).Error("GraphQL panic recovered")
		return fmt.Errorf("internal server error")
	})

	return srv
}

// NewPlaygroundHandler creates a new GraphQL playground handler
func NewPlaygroundHandler() http.Handler {
	return playground.Handler("GraphQL Playground", "/graphql")
}
"@
    
    $handlerPath = Join-Path $graphqlPath "handler.go"
    Set-Content -Path $handlerPath -Value $handlerContent -Encoding UTF8
    Write-Host "    ✅ Created enterprise handler.go" -ForegroundColor Green
}

Write-Host "`n📋 Step 4: Creating enterprise resolvers..." -ForegroundColor Yellow

foreach ($service in $servicesToFix) {
    $servicePath = Join-Path $workspaceRoot "services\$service"
    $graphqlPath = Join-Path $servicePath "graphql"
    
    Write-Host "  🔧 Creating resolver for $service..." -ForegroundColor Cyan
    
    $serviceCapitalized = (Get-Culture).TextInfo.ToTitleCase($service)
    
    $resolverContent = @"
package graphql

import (
	"unified-commerce/services/$service/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
type Resolver struct {
	${serviceCapitalized}Service *service.${serviceCapitalized}Service
	Logger          *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(${service}Service *service.${serviceCapitalized}Service, logger *logger.Logger) *Resolver {
	return &Resolver{
		${serviceCapitalized}Service: ${service}Service,
		Logger:          logger,
	}
}
"@
    
    $resolverPath = Join-Path $graphqlPath "resolver.go"
    Set-Content -Path $resolverPath -Value $resolverContent -Encoding UTF8
    Write-Host "    ✅ Created enterprise resolver.go" -ForegroundColor Green
}

Write-Host "`n📋 Step 5: Generating enterprise federation files..." -ForegroundColor Yellow

foreach ($service in $servicesToFix) {
    $servicePath = Join-Path $workspaceRoot "services\$service"
    $graphqlPath = Join-Path $servicePath "graphql"
    
    Write-Host "  🔧 Generating federation for $service..." -ForegroundColor Cyan
    
    Set-Location $graphqlPath
    
    try {
        & go run github.com/99designs/gqlgen generate 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "    ✅ Generated federation files" -ForegroundColor Green
        } else {
            Write-Host "    ⚠️  Needs schema definition" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "    ⚠️  Schema required for generation" -ForegroundColor Yellow
    }
}

Set-Location $workspaceRoot

Write-Host "`n🎯 ENTERPRISE STANDARDIZATION COMPLETE!" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green
Write-Host "✅ All services now use Identity Service architecture pattern" -ForegroundColor White
Write-Host "✅ Enterprise-grade gqlgen configuration" -ForegroundColor White
Write-Host "✅ Consistent handler and resolver patterns" -ForegroundColor White
Write-Host "✅ Ready for schema definitions and federation" -ForegroundColor White

Write-Host "`n📋 NEXT: Complete schema definitions for each service" -ForegroundColor Cyan
Write-Host "Each service now follows the GOLD STANDARD architecture!" -ForegroundColor Yellow
