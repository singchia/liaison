#!/bin/bash

# 简化的E2E测试运行脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${GREEN}=== Liaison Simple E2E Test Suite ===${NC}"
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

echo -e "${GREEN}Build successful${NC}"

# 检查服务器是否已经在运行
echo -e "${YELLOW}Checking if server is already running...${NC}"
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}Server is already running${NC}"
    SERVER_RUNNING=true
else
    echo -e "${YELLOW}Starting server...${NC}"
    # 启动服务器
    ./bin/liaison -c etc/liaison.yaml &
    SERVER_PID=$!
    SERVER_RUNNING=false
    
    # 等待服务器启动
    echo "Waiting for server to start..."
    for i in {1..30}; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            echo -e "${GREEN}Server started successfully${NC}"
            break
        fi
        if [ $i -eq 30 ]; then
            echo -e "${RED}Server failed to start within 30 seconds${NC}"
            kill $SERVER_PID 2>/dev/null || true
            exit 1
        fi
        sleep 1
    done
fi

# 运行简化测试
echo -e "${YELLOW}Running simple E2E tests...${NC}"
cd test/e2e

if go test -v -run "TestHealthCheck|TestIAMEndpoints|TestProtectedEndpoints" .; then
    echo -e "${GREEN}=== Simple E2E Tests Passed ===${NC}"
    TEST_RESULT=0
else
    echo -e "${RED}=== Simple E2E Tests Failed ===${NC}"
    TEST_RESULT=1
fi

# 清理
if [ "$SERVER_RUNNING" = false ]; then
    echo -e "${YELLOW}Stopping server...${NC}"
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
fi

exit $TEST_RESULT
