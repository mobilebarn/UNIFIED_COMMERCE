#!/usr/bin/env powershell

Write-Host "🧹 Cleaning up duplicate federation scripts..." -ForegroundColor Yellow

$projectRoot = "c:\Users\dane\OneDrive\Desktop\UNIFIED_COMMERCE"

# List of files to remove (keeping the unified one)
$filesToRemove = @(
    "scripts\setup-federation.ps1",
    "scripts\setup-federation-simple.ps1", 
    "scripts\setup-federation-fixed.ps1",
    "scripts\setup-all-federation.ps1",
    "test-graphql-federation.ps1",
    "test-final-federation.ps1",
    "generate-federation.sh",
    "federation.go",
    "generated.go"
)

foreach ($file in $filesToRemove) {
    $fullPath = Join-Path $projectRoot $file
    if (Test-Path $fullPath) {
        Write-Host "  🗑️  Removing $file..." -NoNewline
        try {
            Remove-Item $fullPath -Force
            Write-Host " ✅" -ForegroundColor Green
        }
        catch {
            Write-Host " ❌ Failed to remove" -ForegroundColor Red
        }
    }
}

Write-Host "`n✨ Cleanup complete! Use the federation script:" -ForegroundColor Green
Write-Host "   .\scripts\setup-federation.ps1 -All" -ForegroundColor Cyan
