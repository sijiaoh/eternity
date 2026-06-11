#!/bin/bash
set -eu

# Setup script for Eternity
# Supports Ubuntu and macOS

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "macos" ;;
        *)       echo "unknown" ;;
    esac
}

install_go_deps() {
    echo "Installing Go dependencies..."
    go mod download
    go mod tidy
    echo "Go dependencies installed."
}

install_linux_deps() {
    echo "Installing Linux dependencies..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y gcc xvfb xorg-dev
        echo "Linux dependencies installed."
    else
        echo "Warning: apt-get not found. Please install gcc, xvfb and xorg-dev manually."
    fi
}

install_system_deps() {
    local os=$1

    case "$os" in
        linux)
            install_linux_deps
            ;;
        macos)
            echo "No additional system dependencies required on macOS."
            ;;
        *)
            echo "Warning: Unknown OS. Skipping system dependency installation."
            ;;
    esac
}

main() {
    echo "=== Eternity setup ==="

    local os
    os=$(detect_os)
    echo "Detected OS: $os"

    install_go_deps
    install_system_deps "$os"

    echo "=== Setup complete ==="
}

main
