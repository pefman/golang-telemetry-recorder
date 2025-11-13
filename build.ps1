# Build Script for F1 Telemetry Recorder
# Run this script to build the application

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "  Building F1 Telemetry Recorder" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

Write-Host "Go version:" -ForegroundColor Green
go version
Write-Host ""

# Tidy up dependencies
Write-Host "Resolving dependencies..." -ForegroundColor Yellow
go mod tidy

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Failed to resolve dependencies" -ForegroundColor Red
    exit 1
}

# Build the application
Write-Host "Building executable..." -ForegroundColor Yellow
go build -o f1-telemetry-recorder.exe

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Build failed" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "================================================" -ForegroundColor Green
Write-Host "  Build completed successfully!" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green
Write-Host ""
Write-Host "Executable: f1-telemetry-recorder.exe" -ForegroundColor Cyan
Write-Host ""
Write-Host "To run the application:" -ForegroundColor Yellow
Write-Host "  .\f1-telemetry-recorder.exe" -ForegroundColor White
Write-Host ""
Write-Host "For help, see README.md or QUICKSTART.md" -ForegroundColor Gray
Write-Host ""
