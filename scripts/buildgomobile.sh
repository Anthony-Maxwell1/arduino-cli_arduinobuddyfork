#!/usr/bin/env bash
set -euo pipefail

# =========================
# Arduinobuddy GoMobile Build Script (Linux)
# =========================

info()    { echo "[INFO] $*"; }
warn()    { echo "[WARNING] $*"; }
error()   { echo "[ERROR] $*" >&2; }

# 1. Check if gomobile is installed
info "Checking for gomobile..."
if ! command -v "$(go env GOPATH)/bin/gomobile" >/dev/null 2>&1; then
    warn "gomobile not found. Installing..."
    go install golang.org/x/mobile/cmd/gomobile@latest || {
        error "Failed to install gomobile."
        exit 1
    }
    info "gomobile installed successfully."
    info "Adding go bin to PATH."
    export PATH="$(go env GOPATH)/bin:$PATH"
    info "Post-install gomobile initialization"
    "$(go env GOPATH)/bin/gomobile" init
fi

# 2. Check if input file was provided
if [[ $# -lt 1 ]]; then
    echo "[USAGE] ./buildgomobile.sh input.go [ApiLevel] [OutputName]"
    exit 1
fi

INPUT="$1"

# Uncomment if you want strict input checking
# if [[ ! -f "$INPUT" ]]; then
#     error "Input file \"$INPUT\" not found."
#     exit 1
# fi

# 3. Get API Version argument or default to 25
API="${2:-25}"

# 4. Get output name argument or default to input filename (without extension)
if [[ $# -ge 3 ]]; then
    OUTPUT="$3"
else
    OUTPUT="$(basename "${INPUT%.go}.aar")"
fi

# 5. Build using gomobile
info "Building $INPUT with API $API..."
if ! "$(go env GOPATH)/bin/gomobile" bind -target=android -androidapi "$API" -o "$OUTPUT" "$INPUT"; then
    warn "Initial gomobile build failed. Initializing gomobile..."
    "$(go env GOPATH)/bin/gomobile" init

    info "Building $INPUT with API $API..."
    if ! "$(go env GOPATH)/bin/gomobile" bind -target=android -androidapi "$API" -o "$OUTPUT" "$INPUT"; then
        warn "Second gomobile build failed. Installing bind..."
        go get golang.org/x/mobile/bind

        info "Building $INPUT with API $API..."
        if ! "$(go env GOPATH)/bin/gomobile" bind -target=android -androidapi "$API" -o "$OUTPUT" "$INPUT"; then
            error "Gomobile build failed."
            exit 1
        fi
    fi
fi

echo "[SUCCESS] Build complete: $OUTPUT"
