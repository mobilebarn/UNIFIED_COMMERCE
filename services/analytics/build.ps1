# Build script for analytics service
Set-Location -Path $PSScriptRoot
go build -o analytics-service.exe ./cmd/server/main.go