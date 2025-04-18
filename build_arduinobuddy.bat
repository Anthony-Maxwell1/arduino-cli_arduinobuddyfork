@echo off
REM Compiles any golang script to an android-compatible binary.
REM This script is designed to be run in a Windows environment.
REM === Arduino CLI Cross-Compile for Android ===

REM === CONFIGURATION ===
set GOOS=android
set GOARCH=arm64
set CGO_ENABLED=0
set OUTPUT_NAME=arduino-cli-android
set REPO_DIR=arduino-cli

echo === Setting up environment ===
set GOOS=%GOOS%
set GOARCH=%GOARCH%
set CGO_ENABLED=%CGO_ENABLED%

echo === Building Arduino CLI for Android (%GOARCH%) ===
go build -o %OUTPUT_NAME% ./
if errorlevel 1 (
    echo Build failed.
    exit /b 1
)

echo === Verifying Binary ===
where file >nul 2>nul
if errorlevel 1 (
    echo Could not find `file` utility to verify binary.
    echo You can install it using Git Bash or WSL.
) else (
    file %OUTPUT_NAME%
)

echo.
echo Done! Output: %OUTPUT_NAME%
