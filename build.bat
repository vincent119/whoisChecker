@echo off
REM WhoisChecker Cross-Platform Build Script for Windows

setlocal EnableDelayedExpansion

REM Project information
set APP_NAME=whoisChecker
set VERSION=v1.0.0
set BUILD_DIR=build
set SOURCE_FILE=./cmd/main.go

echo [92m🔨 Building WhoisChecker for multiple platforms...[0m

REM Create build directory
if exist "%BUILD_DIR%" (
    echo [93m📁 Cleaning existing build directory...[0m
    rmdir /s /q "%BUILD_DIR%"
)
mkdir "%BUILD_DIR%"

REM Get current timestamp
for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set "BUILD_TIME=%dt:~0,4%-%dt:~4,2%-%dt:~6,2%_%dt:~8,2%:%dt:~10,2%:%dt:~12,2%"

echo [93m📦 Version: %VERSION%[0m
echo [93m🕒 Build Time: %BUILD_TIME%[0m

echo.
echo [94m🚀 Starting cross-platform compilation...[0m
echo.

REM macOS Intel
echo [94m🏗️  Building for darwin/amd64...[0m
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o "%BUILD_DIR%/%APP_NAME%_darwin_amd64" "%SOURCE_FILE%"
if !errorlevel! equ 0 (
    echo [92m✅ Successfully built: %APP_NAME%_darwin_amd64[0m
) else (
    echo [91m❌ Failed to build for darwin/amd64[0m
)

REM macOS Apple Silicon
echo [94m🏗️  Building for darwin/arm64...[0m
set GOOS=darwin
set GOARCH=arm64
set CGO_ENABLED=0
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o "%BUILD_DIR%/%APP_NAME%_darwin_arm64" "%SOURCE_FILE%"
if !errorlevel! equ 0 (
    echo [92m✅ Successfully built: %APP_NAME%_darwin_arm64[0m
) else (
    echo [91m❌ Failed to build for darwin/arm64[0m
)

REM Windows 64-bit
echo [94m🏗️  Building for windows/amd64...[0m
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o "%BUILD_DIR%/%APP_NAME%_windows_amd64.exe" "%SOURCE_FILE%"
if !errorlevel! equ 0 (
    echo [92m✅ Successfully built: %APP_NAME%_windows_amd64.exe[0m
) else (
    echo [91m❌ Failed to build for windows/amd64[0m
)

REM Windows 32-bit
echo [94m🏗️  Building for windows/386...[0m
set GOOS=windows
set GOARCH=386
set CGO_ENABLED=0
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o "%BUILD_DIR%/%APP_NAME%_windows_386.exe" "%SOURCE_FILE%"
if !errorlevel! equ 0 (
    echo [92m✅ Successfully built: %APP_NAME%_windows_386.exe[0m
) else (
    echo [91m❌ Failed to build for windows/386[0m
)

REM Linux 64-bit
echo [94m🏗️  Building for linux/amd64...[0m
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o "%BUILD_DIR%/%APP_NAME%_linux_amd64" "%SOURCE_FILE%"
if !errorlevel! equ 0 (
    echo [92m✅ Successfully built: %APP_NAME%_linux_amd64[0m
) else (
    echo [91m❌ Failed to build for linux/amd64[0m
)

REM Linux ARM64
echo [94m🏗️  Building for linux/arm64...[0m
set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=0
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME%" -o "%BUILD_DIR%/%APP_NAME%_linux_arm64" "%SOURCE_FILE%"
if !errorlevel! equ 0 (
    echo [92m✅ Successfully built: %APP_NAME%_linux_arm64[0m
) else (
    echo [91m❌ Failed to build for linux/arm64[0m
)

echo.
echo [92m🎉 All builds completed![0m

REM List built files
echo.
echo [94m📋 Built files:[0m
dir "%BUILD_DIR%"

echo.
echo [92m✨ Build process completed![0m
echo [94m📁 All binaries are available in the '%BUILD_DIR%' directory[0m

pause
