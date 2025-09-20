#!/bin/bash

# Liaison Service Uninstallation Script
# This script removes the Liaison systemd service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="liaison"
SERVICE_USER="liaison"
INSTALL_DIR="/opt/liaison"

echo -e "${GREEN}Uninstalling Liaison Service...${NC}"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}This script must be run as root${NC}"
   exit 1
fi

# Stop and disable service
echo -e "${YELLOW}Stopping and disabling service...${NC}"
if systemctl is-active --quiet "$SERVICE_NAME"; then
    systemctl stop "$SERVICE_NAME"
    echo -e "${GREEN}Service stopped${NC}"
else
    echo -e "${YELLOW}Service is not running${NC}"
fi

if systemctl is-enabled --quiet "$SERVICE_NAME"; then
    systemctl disable "$SERVICE_NAME"
    echo -e "${GREEN}Service disabled${NC}"
else
    echo -e "${YELLOW}Service is not enabled${NC}"
fi

# Remove service file
echo -e "${YELLOW}Removing systemd service file...${NC}"
if [[ -f "/etc/systemd/system/$SERVICE_NAME.service" ]]; then
    rm "/etc/systemd/system/$SERVICE_NAME.service"
    echo -e "${GREEN}Service file removed${NC}"
else
    echo -e "${YELLOW}Service file not found${NC}"
fi

# Reload systemd
echo -e "${YELLOW}Reloading systemd daemon...${NC}"
systemctl daemon-reload

# Ask about removing data
echo ""
read -p "Do you want to remove the installation directory ($INSTALL_DIR)? [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [[ -d "$INSTALL_DIR" ]]; then
        rm -rf "$INSTALL_DIR"
        echo -e "${GREEN}Installation directory removed${NC}"
    else
        echo -e "${YELLOW}Installation directory not found${NC}"
    fi
else
    echo -e "${YELLOW}Installation directory preserved${NC}"
fi

# Ask about removing user
echo ""
read -p "Do you want to remove the service user ($SERVICE_USER)? [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if id "$SERVICE_USER" &>/dev/null; then
        userdel "$SERVICE_USER"
        echo -e "${GREEN}Service user removed${NC}"
    else
        echo -e "${YELLOW}Service user not found${NC}"
    fi
else
    echo -e "${YELLOW}Service user preserved${NC}"
fi

echo -e "${GREEN}Uninstallation completed successfully!${NC}"
