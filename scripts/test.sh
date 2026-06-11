#!/bin/bash
# Run Go tests with Xvfb on Linux for headless execution.
# Usage: ./scripts/test.sh [test-target]

set -eu

OS="$(uname -s)"
TEST_TARGET="${1:-./...}"

run_tests() {
    CGO_ENABLED=1 go test -v "$TEST_TARGET"
}

case "$OS" in
    Linux)
        if ! command -v Xvfb &> /dev/null; then
            echo "Error: Xvfb is not installed. Install with: sudo apt-get install xvfb"
            exit 1
        fi

        DISPLAY_NUM=""
        for num in 99 98 97 96 95; do
            if ! [ -e "/tmp/.X${num}-lock" ]; then
                DISPLAY_NUM=$num
                break
            fi
        done

        if [ -z "$DISPLAY_NUM" ]; then
            echo "Error: No available display number (99-95 all in use)"
            exit 1
        fi

        export DISPLAY=":${DISPLAY_NUM}"

        Xvfb ":${DISPLAY_NUM}" -screen 0 1024x768x24 &
        XVFB_PID=$!

        trap "kill $XVFB_PID 2>/dev/null || true" EXIT

        sleep 1
        if ! kill -0 $XVFB_PID 2>/dev/null; then
            echo "Error: Failed to start Xvfb"
            exit 1
        fi

        run_tests
        ;;
    Darwin)
        run_tests
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac
