#!/bin/bash

# E2E测试运行脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${GREEN}=== Liaison E2E Test Suite ===${NC}"
echo "Project root: $PROJECT_ROOT"

# 检查依赖
echo -e "${YELLOW}Checking dependencies...${NC}"

# 检查Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# 检查测试依赖
echo "Installing test dependencies..."
go mod tidy
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require

# 构建项目
echo -e "${YELLOW}Building project...${NC}"
if ! go build -o ./bin/liaison cmd/manager/main.go; then
    echo -e "${RED}Error: Failed to build project${NC}"
    exit 1
fi

# 检查二进制文件
if [ ! -f "./bin/liaison" ]; then
    echo -e "${RED}Error: Binary file not found${NC}"
    exit 1
fi

echo -e "${GREEN}Build successful${NC}"

# 运行E2E测试
echo -e "${YELLOW}Running E2E tests...${NC}"
cd test/e2e

# 设置测试环境变量
export LIAISON_DB_PATH="test_liaison.db"
export LIAISON_LOG_LEVEL="debug"

# 运行测试
if go test -v -timeout 5m .; then
    echo -e "${GREEN}=== E2E Tests Passed ===${NC}"
    exit 0
else
    echo -e "${RED}=== E2E Tests Failed ===${NC}"
    exit 1
fi
