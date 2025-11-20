@echo off
setlocal enabledelayedexpansion

:: =========================
:: Arduinobuddy GoMobile Build Script
:: =========================

:: 1. Check if gomobile is installed
echo [INFO] Checking for gomobile...
where gomobile >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [WARNING] gomobile not found. Installing...
    go install golang.org/x/mobile/cmd/gomobile@latest
    if %ERRORLEVEL% neq 0 (
        echo [ERROR] Failed to install gomobile.
        exit /b 1
    )
    echo [INFO] Post-Install gomobile initialization
    gomobile init
)

:: 2. Check if input file was provided
if "%~1"=="" (
    echo [USAGE] buildgomobile.bat input.go [ApiLevel] [OutputName]
    exit /b 1
)

set INPUT=%~1

:: if not exist "%INPUT%" (
::     echo [ERROR] Input file "%INPUT%" not found.
::     exit /b 1
:: )

:: 3. Get API Version argument or default to 25
set API=%~2
if "%API%"=="" (
    set API=25
)

:: 4. Get output name argument or default to input filename (without extension)
set OUTPUT=%~3
if "%OUTPUT%"=="" (
    for %%F in ("%INPUT%.aar") do set OUTPUT=%%~nF
)

:: 5. Build using gomobile
echo [INFO] Building %INPUT% with API %API%...
gomobile bind -target=android -androidapi %API% -o "%OUTPUT%" "%INPUT%"

if %ERRORLEVEL% neq 0 (
    echo [WARNING] Initial gomobile build failed. Initializing gomobile...
    gomobile init
    echo [INFO] Building %INPUT% with API %API%...
    gomobile bind -target=android -androidapi %API% -o "%OUTPUT%" "%INPUT%"
    if %ERRORLEVEL% neq 0 (
        echo [WARNING] Second gomobile build failed. Installing bind...
        go get golang.org/x/mobile/bind
        echo [INFO] Building %INPUT% with API %API%...
        gomobile bind -target=android -androidapi %API% -o "%OUTPUT%" "%INPUT%"
        if %ERRORLEVEL% neq 0 (
            echo [ERROR] Gomobile build failed.
            exit /b 1
        )
    )
)

echo [SUCCESS] Build complete: %OUTPUT%
endlocal