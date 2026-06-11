#!/bin/bash
#
# Go代码格式化和静态分析脚本
# 支持Ubuntu和Mac
#

set -eu
cd "$(dirname "$0")/.." || exit 1

FAILED=""

echo "=== gofmt ==="
UNFORMATTED=$(gofmt -l . 2>/dev/null | grep -v '^vendor/' || true)
if [ -n "$UNFORMATTED" ]; then
    echo "以下文件需要格式化："
    echo "$UNFORMATTED"
    echo ""
    echo "运行 'gofmt -w .' 来格式化"
    FAILED=1
else
    echo "所有文件格式正确"
fi

echo ""
echo "=== go vet ==="
if ! go vet ./...; then
    FAILED=1
fi

echo ""
echo "=== golangci-lint ==="
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run ./...; then
        FAILED=1
    fi
else
    echo "golangci-lint未安装，跳过"
    echo "安装: https://golangci-lint.run/usage/install/"
fi

echo ""
if [ -n "$FAILED" ]; then
    echo "发现问题"
    exit 1
else
    echo "检查通过"
fi
