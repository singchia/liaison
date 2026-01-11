#!/bin/bash

# Liaison Edge 卸载脚本
# 此脚本会停止并卸载 Liaison Edge 服务

set +e  # 允许某些命令失败而不退出脚本

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
SERVICE_NAME="liaison-edge"
INSTALL_DIR="/opt/liaison"
BIN_DIR="/opt/liaison/bin"
CONFIG_DIR="/opt/liaison/conf"
LOG_DIR="/opt/liaison/logs"
DATA_DIR="/opt/liaison/data"

echo -e "${YELLOW}Uninstalling Liaison Edge...${NC}"

# 停止 systemd 服务（Linux）
if systemctl list-unit-files | grep -q "${SERVICE_NAME}.service"; then
    echo -e "${YELLOW}Stopping systemd service...${NC}"
    sudo systemctl stop "${SERVICE_NAME}" 2>/dev/null || true
    sudo systemctl disable "${SERVICE_NAME}" 2>/dev/null || true
    echo -e "${GREEN}systemd service stopped and disabled${NC}"
fi

# 停止 launchd 服务（macOS）
if [[ "$(uname)" == "Darwin" ]]; then
    PLIST_FILE="${HOME}/Library/LaunchAgents/com.liaison.edge.plist"
    if [ -f "$PLIST_FILE" ]; then
        echo -e "${YELLOW}Stopping launchd service...${NC}"
        # 使用新的 launchctl bootout API (macOS 10.11+)
        launchctl bootout "gui/$(id -u)/com.liaison.edge" 2>/dev/null || \
        launchctl unload "$PLIST_FILE" 2>/dev/null || true
        echo -e "${GREEN}launchd service stopped${NC}"
    fi
fi

# 停止所有运行中的 liaison-edge 进程
echo -e "${YELLOW}Stopping running processes...${NC}"
EDGE_PIDS=$(pgrep -f "liaison-edge" 2>/dev/null || true)
if [ -n "$EDGE_PIDS" ]; then
    echo "$EDGE_PIDS" | xargs kill -TERM 2>/dev/null || true
    sleep 2
    # 如果还有进程，强制杀死
    REMAINING_PIDS=$(pgrep -f "liaison-edge" 2>/dev/null || true)
    if [ -n "$REMAINING_PIDS" ]; then
        echo "$REMAINING_PIDS" | xargs kill -KILL 2>/dev/null || true
    fi
    echo -e "${GREEN}All liaison-edge processes stopped${NC}"
else
    echo -e "${YELLOW}No running liaison-edge processes found${NC}"
fi

# 删除 systemd 服务文件
if [ -f "/etc/systemd/system/${SERVICE_NAME}.service" ]; then
    echo -e "${YELLOW}Removing systemd service file...${NC}"
    sudo rm -f "/etc/systemd/system/${SERVICE_NAME}.service"
    sudo systemctl daemon-reload 2>/dev/null || true
    echo -e "${GREEN}systemd service file removed${NC}"
fi

# 删除 launchd plist 文件（macOS）
if [[ "$(uname)" == "Darwin" ]]; then
    PLIST_FILE="${HOME}/Library/LaunchAgents/com.liaison.edge.plist"
    if [ -f "$PLIST_FILE" ]; then
        echo -e "${YELLOW}Removing launchd plist file...${NC}"
        rm -f "$PLIST_FILE"
        echo -e "${GREEN}launchd plist file removed${NC}"
    fi
fi

# 删除二进制文件
if [ -f "${BIN_DIR}/liaison-edge" ]; then
    echo -e "${YELLOW}Removing binary file...${NC}"
    rm -f "${BIN_DIR}/liaison-edge"
    echo -e "${GREEN}Binary file removed${NC}"
fi

# 删除配置文件
if [ -f "${CONFIG_DIR}/liaison-edge.yaml" ]; then
    echo -e "${YELLOW}Removing configuration file...${NC}"
    rm -f "${CONFIG_DIR}/liaison-edge.yaml"
    echo -e "${GREEN}Configuration file removed${NC}"
fi

# 删除日志文件
if [ -d "$LOG_DIR" ]; then
    LOG_FILES=$(find "$LOG_DIR" -name "*liaison-edge*" 2>/dev/null || true)
    if [ -n "$LOG_FILES" ]; then
        echo -e "${YELLOW}Removing log files...${NC}"
        find "$LOG_DIR" -name "*liaison-edge*" -delete 2>/dev/null || true
        echo -e "${GREEN}Log files removed${NC}"
    fi
fi

# 询问是否删除数据文件（如果有的话）
if [ -d "$DATA_DIR" ]; then
    EDGE_DATA=$(find "$DATA_DIR" -name "*edge*" -o -name "*liaison-edge*" 2>/dev/null || true)
    if [ -n "$EDGE_DATA" ]; then
        echo ""
        read -p "Do you want to remove edge data files? [y/N]: " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}Removing edge data files...${NC}"
            find "$DATA_DIR" -name "*edge*" -o -name "*liaison-edge*" -delete 2>/dev/null || true
            echo -e "${GREEN}Edge data files removed${NC}"
        fi
    fi
fi

echo ""
echo -e "${GREEN}✅ Liaison Edge uninstallation completed!${NC}"
