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

# 根据操作系统设置正确的路径
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    # Windows 系统：使用标准 Windows 路径
    if command -v cygpath >/dev/null 2>&1; then
        INSTALL_DIR=$(cygpath -w "/c/Program Files/Liaison" 2>/dev/null || echo "C:\\Program Files\\Liaison")
        BIN_DIR=$(cygpath -w "/c/Program Files/Liaison/bin" 2>/dev/null || echo "C:\\Program Files\\Liaison\\bin")
        CONFIG_DIR=$(cygpath -w "/c/Program Files/Liaison/conf" 2>/dev/null || echo "C:\\Program Files\\Liaison\\conf")
        LOG_DIR=$(cygpath -w "/c/Program Files/Liaison/logs" 2>/dev/null || echo "C:\\Program Files\\Liaison\\logs")
    else
        INSTALL_DIR="C:\\Program Files\\Liaison"
        BIN_DIR="C:\\Program Files\\Liaison\\bin"
        CONFIG_DIR="C:\\Program Files\\Liaison\\conf"
        LOG_DIR="C:\\Program Files\\Liaison\\logs"
    fi
    BINARY_NAME="liaison-edge.exe"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS 系统：使用标准 macOS 路径
    BIN_DIR="/usr/local/bin"
    CONFIG_DIR="${HOME}/Library/Application Support/liaison"
    LOG_DIR="${HOME}/Library/Logs/liaison"
    INSTALL_DIR="/usr/local"
    BINARY_NAME="liaison-edge"
    DATA_DIR=""
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux 系统：使用标准 Linux 路径
    BIN_DIR="/usr/local/bin"
    CONFIG_DIR="/etc/liaison"
    DATA_DIR="/var/lib/liaison"
    LOG_DIR="/var/log/liaison"
    INSTALL_DIR="/usr/local"
    BINARY_NAME="liaison-edge"
else
    # 其他系统：使用默认路径（兼容旧版本）
    INSTALL_DIR="/opt/liaison"
    BIN_DIR="/opt/liaison/bin"
    CONFIG_DIR="/opt/liaison/conf"
    LOG_DIR="/opt/liaison/logs"
    BINARY_NAME="liaison-edge"
    DATA_DIR="/opt/liaison/data"
fi

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
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    # Windows 系统：使用 taskkill 或 pgrep
    if command -v taskkill >/dev/null 2>&1; then
        taskkill //F //IM "${BINARY_NAME}" 2>/dev/null || true
    else
        EDGE_PIDS=$(pgrep -f "liaison-edge" 2>/dev/null || true)
        if [ -n "$EDGE_PIDS" ]; then
            echo "$EDGE_PIDS" | xargs kill -TERM 2>/dev/null || true
            sleep 2
            REMAINING_PIDS=$(pgrep -f "liaison-edge" 2>/dev/null || true)
            if [ -n "$REMAINING_PIDS" ]; then
                echo "$REMAINING_PIDS" | xargs kill -KILL 2>/dev/null || true
            fi
        fi
    fi
    echo -e "${GREEN}All liaison-edge processes stopped${NC}"
else
    # Linux/macOS 系统
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
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS: 二进制在 /usr/local/bin 下，需要 sudo
    if [ -f "${BIN_DIR}/${BINARY_NAME}" ]; then
        echo -e "${YELLOW}Removing binary file...${NC}"
        sudo rm -f "${BIN_DIR}/${BINARY_NAME}"
        echo -e "${GREEN}Binary file removed${NC}"
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux 需要 sudo 权限删除 /usr/local/bin 下的文件
    if [ -f "${BIN_DIR}/${BINARY_NAME}" ]; then
        echo -e "${YELLOW}Removing binary file...${NC}"
        sudo rm -f "${BIN_DIR}/${BINARY_NAME}"
        echo -e "${GREEN}Binary file removed${NC}"
    fi
else
    if [ -f "${BIN_DIR}/${BINARY_NAME}" ]; then
        echo -e "${YELLOW}Removing binary file...${NC}"
        rm -f "${BIN_DIR}/${BINARY_NAME}"
        echo -e "${GREEN}Binary file removed${NC}"
    fi
fi

# 删除配置文件
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS: 配置文件在用户目录下，不需要 sudo
    if [ -f "${CONFIG_DIR}/liaison-edge.yaml" ]; then
        echo -e "${YELLOW}Removing configuration file...${NC}"
        rm -f "${CONFIG_DIR}/liaison-edge.yaml"
        echo -e "${GREEN}Configuration file removed${NC}"
    fi
    # 如果配置目录为空，也删除目录
    if [ -d "${CONFIG_DIR}" ] && [ -z "$(ls -A "${CONFIG_DIR}" 2>/dev/null)" ]; then
        rmdir "${CONFIG_DIR}" 2>/dev/null || true
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux 需要 sudo 权限删除 /etc 下的文件
    if [ -f "${CONFIG_DIR}/liaison-edge.yaml" ]; then
        echo -e "${YELLOW}Removing configuration file...${NC}"
        sudo rm -f "${CONFIG_DIR}/liaison-edge.yaml"
        echo -e "${GREEN}Configuration file removed${NC}"
    fi
    # 如果配置目录为空，也删除目录
    if [ -d "${CONFIG_DIR}" ] && [ -z "$(ls -A "${CONFIG_DIR}" 2>/dev/null)" ]; then
        sudo rmdir "${CONFIG_DIR}" 2>/dev/null || true
    fi
else
    if [ -f "${CONFIG_DIR}/liaison-edge.yaml" ]; then
        echo -e "${YELLOW}Removing configuration file...${NC}"
        rm -f "${CONFIG_DIR}/liaison-edge.yaml"
        echo -e "${GREEN}Configuration file removed${NC}"
    fi
fi

# 删除数据文件（Linux）
if [[ "$OSTYPE" == "linux-gnu"* ]] && [ -n "$DATA_DIR" ] && [ -d "$DATA_DIR" ]; then
    DATA_FILES=$(find "$DATA_DIR" -name "*edge*" -o -name "*liaison-edge*" 2>/dev/null || true)
    if [ -n "$DATA_FILES" ]; then
        echo -e "${YELLOW}Removing data files...${NC}"
        sudo find "$DATA_DIR" -name "*edge*" -o -name "*liaison-edge*" -delete 2>/dev/null || true
        echo -e "${GREEN}Data files removed${NC}"
    fi
    # 如果数据目录为空，也删除目录
    if [ -z "$(ls -A "${DATA_DIR}" 2>/dev/null)" ]; then
        sudo rmdir "${DATA_DIR}" 2>/dev/null || true
    fi
fi

# 删除日志文件
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS: 日志文件在用户目录下，不需要 sudo
    if [ -d "$LOG_DIR" ]; then
        LOG_FILES=$(find "$LOG_DIR" -name "*liaison-edge*" 2>/dev/null || true)
        if [ -n "$LOG_FILES" ]; then
            echo -e "${YELLOW}Removing log files...${NC}"
            find "$LOG_DIR" -name "*liaison-edge*" -delete 2>/dev/null || true
            echo -e "${GREEN}Log files removed${NC}"
        fi
        # 如果日志目录为空，也删除目录
        if [ -z "$(ls -A "${LOG_DIR}" 2>/dev/null)" ]; then
            rmdir "${LOG_DIR}" 2>/dev/null || true
        fi
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux 需要 sudo 权限删除 /var/log 下的文件
    if [ -d "$LOG_DIR" ]; then
        LOG_FILES=$(find "$LOG_DIR" -name "*liaison-edge*" 2>/dev/null || true)
        if [ -n "$LOG_FILES" ]; then
            echo -e "${YELLOW}Removing log files...${NC}"
            sudo find "$LOG_DIR" -name "*liaison-edge*" -delete 2>/dev/null || true
            echo -e "${GREEN}Log files removed${NC}"
        fi
        # 如果日志目录为空，也删除目录
        if [ -z "$(ls -A "${LOG_DIR}" 2>/dev/null)" ]; then
            sudo rmdir "${LOG_DIR}" 2>/dev/null || true
        fi
    fi
else
    if [ -d "$LOG_DIR" ]; then
        LOG_FILES=$(find "$LOG_DIR" -name "*liaison-edge*" 2>/dev/null || true)
        if [ -n "$LOG_FILES" ]; then
            echo -e "${YELLOW}Removing log files...${NC}"
            find "$LOG_DIR" -name "*liaison-edge*" -delete 2>/dev/null || true
            echo -e "${GREEN}Log files removed${NC}"
        fi
    fi
fi

# 询问是否删除数据文件（如果有的话，仅限非 Linux 系统）
if [[ "$OSTYPE" != "linux-gnu"* ]] && [ -n "$DATA_DIR" ] && [ -d "$DATA_DIR" ]; then
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

# 静默清理历史 /opt/liaison 目录（如果存在且为空或只包含 edge 相关文件）
if [ -d "/opt/liaison" ]; then
    # 检查目录是否为空
    OPT_CONTENTS=$(find "/opt/liaison" -mindepth 1 -maxdepth 1 2>/dev/null | wc -l)
    if [ "$OPT_CONTENTS" -eq 0 ]; then
        echo -e "${YELLOW}Removing empty legacy directory /opt/liaison...${NC}"
        sudo rmdir "/opt/liaison" 2>/dev/null || true
    else
        # 检查是否只有 edge 相关文件
        OPT_EDGE_ONLY=$(find "/opt/liaison" -mindepth 1 -maxdepth 1 ! -name "*edge*" ! -name "*liaison-edge*" 2>/dev/null | wc -l)
        if [ "$OPT_EDGE_ONLY" -eq 0 ]; then
            echo -e "${YELLOW}Removing legacy /opt/liaison directory (only edge files found)...${NC}"
            sudo rm -rf "/opt/liaison" 2>/dev/null || true
        fi
    fi
fi

echo ""
echo -e "${GREEN}✅ Liaison Edge uninstallation completed!${NC}"
