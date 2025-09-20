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
cp "$(dirname "$0")/liaison.service" "/etc/systemd/system/"
chmod 644 "/etc/systemd/system/liaison.service"

# Reload systemd
echo -e "${YELLOW}Reloading systemd daemon...${NC}"
systemctl daemon-reload

# Enable service
echo -e "${YELLOW}Enabling liaison service...${NC}"
systemctl enable "$SERVICE_NAME"

echo -e "${GREEN}Installation completed successfully!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Copy the liaison binary to $BIN_DIR/"
echo "2. Copy the configuration file to $CONFIG_DIR/liaison.yaml"
echo "3. Start the service: systemctl start $SERVICE_NAME"
echo "4. Check status: systemctl status $SERVICE_NAME"
echo "5. View logs: journalctl -u $SERVICE_NAME -f"
echo ""
echo -e "${GREEN}Service management commands:${NC}"
echo "  Start:   systemctl start $SERVICE_NAME"
echo "  Stop:    systemctl stop $SERVICE_NAME"
echo "  Restart: systemctl restart $SERVICE_NAME"
echo "  Status:  systemctl status $SERVICE_NAME"
echo "  Logs:    journalctl -u $SERVICE_NAME -f"
