#!/bin/bash
set -eu

# Setup script for ebiten-agent-example
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

install_xvfb() {
    local os=$1

    case "$os" in
        linux)
            echo "Installing Xvfb..."
            if command -v apt-get &> /dev/null; then
                sudo apt-get update
                sudo apt-get install -y xvfb
                echo "Xvfb installed."
            else
                echo "Warning: apt-get not found. Please install Xvfb manually."
            fi
            ;;
        macos)
            echo "Xvfb is not required on macOS (native display available)."
            echo "For headless testing, consider using XQuartz if needed."
            ;;
        *)
            echo "Warning: Unknown OS. Skipping Xvfb installation."
            ;;
    esac
}

main() {
    echo "=== ebiten-agent-example setup ==="

    local os
    os=$(detect_os)
    echo "Detected OS: $os"

    install_go_deps
    install_xvfb "$os"

    echo "=== Setup complete ==="
}

main
