#!/bin/bash

# Liaison Service Installation Script
# This script installs and configures the Liaison systemd service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="liaison"
SERVICE_USER="liaison"
SERVICE_GROUP="liaison"
INSTALL_DIR="/opt/liaison"
CONFIG_DIR="/opt/liaison/conf"
DATA_DIR="/opt/liaison/data"
LOG_DIR="/opt/liaison/logs"
BIN_DIR="/opt/liaison/bin"

echo -e "${GREEN}Installing Liaison Service...${NC}"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}This script must be run as root${NC}"
   exit 1
fi

# Create service user and group
echo -e "${YELLOW}Creating service user and group...${NC}"
if ! id "$SERVICE_USER" &>/dev/null; then
    useradd --system --no-create-home --shell /bin/false "$SERVICE_USER"
    echo -e "${GREEN}Created user: $SERVICE_USER${NC}"
else
    echo -e "${YELLOW}User $SERVICE_USER already exists${NC}"
fi

# Create directories
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p "$INSTALL_DIR" "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR" "$BIN_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR"
chmod 755 "$INSTALL_DIR"

# Set specific permissions
chmod 750 "$DATA_DIR" "$LOG_DIR"
chmod 755 "$BIN_DIR" "$CONFIG_DIR"

# Copy service file
echo -e "${YELLOW}Installing systemd service file...${NC}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
if [[ -f "$SCRIPT_DIR/systemd/liaison.service" ]]; then
    cp "$SCRIPT_DIR/systemd/liaison.service" "/etc/systemd/system/"
elif [[ -f "$SCRIPT_DIR/liaison.service" ]]; then
    cp "$SCRIPT_DIR/liaison.service" "/etc/systemd/system/"
else
    echo -e "${RED}Error: liaison.service not found${NC}"
    exit 1
fi
chmod 644 "/etc/systemd/system/liaison.service"

# Copy binaries
echo -e "${YELLOW}Copying binaries...${NC}"
if [[ -d "$SCRIPT_DIR/bin" ]]; then
    cp -f "$SCRIPT_DIR/bin/"* "$BIN_DIR/"
    chown "$SERVICE_USER:$SERVICE_GROUP" "$BIN_DIR"/*
    chmod 755 "$BIN_DIR"/*
    echo -e "${GREEN}Binaries copied${NC}"
else
    echo -e "${YELLOW}Warning: bin directory not found, skipping binary copy${NC}"
fi

# Copy configuration files
echo -e "${YELLOW}Copying configuration files...${NC}"
if [[ -d "$SCRIPT_DIR/etc" ]]; then
    cp -f "$SCRIPT_DIR/etc/"*.yaml "$CONFIG_DIR/" 2>/dev/null || true
    chown "$SERVICE_USER:$SERVICE_GROUP" "$CONFIG_DIR"/*.yaml 2>/dev/null || true
    chmod 644 "$CONFIG_DIR"/*.yaml 2>/dev/null || true
    echo -e "${GREEN}Configuration files copied${NC}"
else
    echo -e "${YELLOW}Warning: etc directory not found, skipping config copy${NC}"
fi

# Reload systemd
echo -e "${YELLOW}Reloading systemd daemon...${NC}"
systemctl daemon-reload

# Enable service
echo -e "${YELLOW}Enabling liaison service...${NC}"
systemctl enable "$SERVICE_NAME"

echo -e "${GREEN}Installation completed successfully!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review and edit configuration files in $CONFIG_DIR/ if needed"
echo "2. Start the service: systemctl start $SERVICE_NAME"
echo "3. Check status: systemctl status $SERVICE_NAME"
echo "4. View logs: journalctl -u $SERVICE_NAME -f"
echo ""
echo -e "${GREEN}Service management commands:${NC}"
echo "  Start:   systemctl start $SERVICE_NAME"
echo "  Stop:    systemctl stop $SERVICE_NAME"
echo "  Restart: systemctl restart $SERVICE_NAME"
echo "  Status:  systemctl status $SERVICE_NAME"
echo "  Logs:    journalctl -u $SERVICE_NAME -f"
