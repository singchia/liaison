#!/bin/bash

# Liaison Edge 安装脚本
# 此脚本会根据操作系统自动下载并安装对应的 Edge 安装包

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认配置
SERVER_ADDR="${SERVER_ADDR:-localhost:8080}"
PACKAGES_DIR="/opt/liaison/packages"
INSTALL_DIR="/opt/liaison/edge"

# 解析参数
ACCESS_KEY=""
SECRET_KEY=""

show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --access-key=KEY     Access key (required)"
    echo "  --secret-key=KEY     Secret key (required)"
    echo "  --server-addr=ADDR   Server address (default: localhost:8080)"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 --access-key=xxx --secret-key=yyy --server-addr=example.com:8080"
    exit 0
}

while [[ $# -gt 0 ]]; do
    case $1 in
        --access-key=*)
            ACCESS_KEY="${1#*=}"
            shift
            ;;
        --secret-key=*)
            SECRET_KEY="${1#*=}"
            shift
            ;;
        --server-addr=*)
            SERVER_ADDR="${1#*=}"
            shift
            ;;
        -h|--help)
            show_help
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

if [ -z "$ACCESS_KEY" ] || [ -z "$SECRET_KEY" ]; then
    echo -e "${RED}Error: --access-key and --secret-key are required${NC}"
    echo "Use --help for usage information"
    exit 1
fi

echo -e "${GREEN}Starting Liaison Edge installation...${NC}"

# 检测操作系统
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # 检测 Linux 架构
        ARCH=$(uname -m)
        if [ "$ARCH" == "x86_64" ]; then
            echo "linux-amd64"
        elif [ "$ARCH" == "aarch64" ] || [ "$ARCH" == "arm64" ]; then
            echo "linux-arm64"
        else
            echo "linux-$ARCH"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # 检测 macOS 架构
        ARCH=$(uname -m)
        if [ "$ARCH" == "x86_64" ]; then
            echo "darwin-amd64"
        elif [ "$ARCH" == "arm64" ]; then
            echo "darwin-arm64"
        else
            echo "darwin-$ARCH"
        fi
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
        echo "windows-amd64"
    else
        echo -e "${RED}Unsupported operating system: $OSTYPE${NC}"
        exit 1
    fi
}

OS_ARCH=$(detect_os)
echo -e "${GREEN}Detected OS/Arch: $OS_ARCH${NC}"

# 确定安装包文件名（统一使用 tar.gz 格式）
PACKAGE_NAME="liaison-edge-${OS_ARCH}.tar.gz"

# 创建临时目录
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# 下载安装包
echo -e "${YELLOW}Downloading installation package...${NC}"
# 构建下载 URL（使用 http:// 前缀，因为下载需要协议）
PACKAGE_URL="http://${SERVER_ADDR}/packages/edge/${PACKAGE_NAME}"
echo -e "${YELLOW}Server address: ${SERVER_ADDR}${NC}"
echo -e "${YELLOW}Package URL: ${PACKAGE_URL}${NC}"

if command -v curl >/dev/null 2>&1; then
    HTTP_CODE=$(curl -sSL -o "${TMP_DIR}/${PACKAGE_NAME}" -w "%{http_code}" "${PACKAGE_URL}")
elif command -v wget >/dev/null 2>&1; then
    wget -q -O "${TMP_DIR}/${PACKAGE_NAME}" "${PACKAGE_URL}" || HTTP_CODE="404"
    if [ $? -eq 0 ]; then
        HTTP_CODE="200"
    else
        HTTP_CODE="404"
    fi
else
    echo -e "${RED}Error: curl or wget is required${NC}"
    exit 1
fi

if [ "$HTTP_CODE" != "200" ]; then
    echo -e "${RED}Error: Failed to download package (HTTP $HTTP_CODE)${NC}"
    echo -e "${YELLOW}Package URL: $PACKAGE_URL${NC}"
    exit 1
fi

# 解压安装包
echo -e "${YELLOW}Extracting package...${NC}"
cd "$TMP_DIR"
if ! tar -xzf "${PACKAGE_NAME}"; then
    echo -e "${RED}Error: Failed to extract package${NC}"
    exit 1
fi

# 安装
echo -e "${YELLOW}Installing...${NC}"

# 查找解压后的二进制文件
BINARY_NAME="liaison-edge"
if [[ "$OS_ARCH" == "windows"* ]]; then
    BINARY_NAME="liaison-edge.exe"
fi

if [ ! -f "${TMP_DIR}/${BINARY_NAME}" ]; then
    echo -e "${RED}Error: Binary file ${BINARY_NAME} not found in package${NC}"
    exit 1
fi

# Linux/macOS/Windows 安装
mkdir -p "$INSTALL_DIR"
cp "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/liaison-edge"
chmod +x "${INSTALL_DIR}/liaison-edge"

# 创建配置文件
mkdir -p "${INSTALL_DIR}/etc"

cat > "${INSTALL_DIR}/etc/liaison-edge.yaml" <<EOF
manager:
  dial:
    addrs:
      - ${SERVER_ADDR}
    network: tcp
    tls:
      enable: false
  auth:
    access_key: "${ACCESS_KEY}"
    secret_key: "${SECRET_KEY}"
log:
  level: info
  file: ./logs/liaison-edge.log
  maxsize: 100
  maxrolls: 10
EOF

echo -e "${GREEN}Installation completed!${NC}"
echo -e "${GREEN}Edge binary: ${INSTALL_DIR}/liaison-edge${NC}"
echo -e "${GREEN}Config file: ${INSTALL_DIR}/etc/liaison-edge.yaml${NC}"
echo ""
echo -e "${YELLOW}To start the edge, run:${NC}"
echo -e "${YELLOW}  ${INSTALL_DIR}/liaison-edge -c ${INSTALL_DIR}/etc/liaison-edge.yaml${NC}"

