#!/bin/bash

# Liaison Service Uninstallation Script
# This script removes the Liaison systemd service

# Don't exit on error for commands that may fail (like pgrep when no process found)
set +e

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

# Stop and disable services
echo -e "${YELLOW}Stopping and disabling services...${NC}"

# Stop and disable liaison service
if systemctl is-active --quiet "$SERVICE_NAME" 2>/dev/null; then
    systemctl stop "$SERVICE_NAME"
    echo -e "${GREEN}liaison service stopped${NC}"
else
    echo -e "${YELLOW}liaison service is not running${NC}"
fi

if systemctl is-enabled --quiet "$SERVICE_NAME" 2>/dev/null; then
    systemctl disable "$SERVICE_NAME"
    echo -e "${GREEN}liaison service disabled${NC}"
else
    echo -e "${YELLOW}liaison service is not enabled${NC}"
fi

# Stop and disable frontier service
if systemctl is-active --quiet frontier 2>/dev/null; then
    systemctl stop frontier
    echo -e "${GREEN}frontier service stopped${NC}"
else
    echo -e "${YELLOW}frontier service is not running${NC}"
fi

if systemctl is-enabled --quiet frontier 2>/dev/null; then
    systemctl disable frontier
    echo -e "${GREEN}frontier service disabled${NC}"
else
    echo -e "${YELLOW}frontier service is not enabled${NC}"
fi

# Force kill any remaining processes
echo -e "${YELLOW}Checking for remaining processes...${NC}"

# Kill liaison processes
LIAISON_PIDS=$(pgrep -f "/opt/liaison/bin/liaison" 2>/dev/null)
if [[ -n "$LIAISON_PIDS" ]]; then
    echo -e "${YELLOW}Killing remaining liaison processes (PIDs: $LIAISON_PIDS)...${NC}"
    kill -TERM $LIAISON_PIDS 2>/dev/null
    sleep 2
    # Force kill if still running
    LIAISON_PIDS=$(pgrep -f "/opt/liaison/bin/liaison" 2>/dev/null)
    if [[ -n "$LIAISON_PIDS" ]]; then
        echo -e "${YELLOW}Force killing liaison processes (PIDs: $LIAISON_PIDS)...${NC}"
        kill -KILL $LIAISON_PIDS 2>/dev/null
    fi
else
    echo -e "${GREEN}No liaison processes found${NC}"
fi

# Kill frontier processes
FRONTIER_PIDS=$(pgrep -f "/opt/liaison/bin/frontier" 2>/dev/null)
if [[ -n "$FRONTIER_PIDS" ]]; then
    echo -e "${YELLOW}Killing remaining frontier processes (PIDs: $FRONTIER_PIDS)...${NC}"
    kill -TERM $FRONTIER_PIDS 2>/dev/null
    sleep 2
    # Force kill if still running
    FRONTIER_PIDS=$(pgrep -f "/opt/liaison/bin/frontier" 2>/dev/null)
    if [[ -n "$FRONTIER_PIDS" ]]; then
        echo -e "${YELLOW}Force killing frontier processes (PIDs: $FRONTIER_PIDS)...${NC}"
        kill -KILL $FRONTIER_PIDS 2>/dev/null
    fi
else
    echo -e "${GREEN}No frontier processes found${NC}"
fi

# Wait a bit for processes to fully terminate
sleep 1

# Remove service files
echo -e "${YELLOW}Removing systemd service files...${NC}"
if [[ -f "/etc/systemd/system/$SERVICE_NAME.service" ]]; then
    rm "/etc/systemd/system/$SERVICE_NAME.service"
    echo -e "${GREEN}liaison.service removed${NC}"
else
    echo -e "${YELLOW}liaison.service not found${NC}"
fi

if [[ -f "/etc/systemd/system/frontier.service" ]]; then
    rm "/etc/systemd/system/frontier.service"
    echo -e "${GREEN}frontier.service removed${NC}"
else
    echo -e "${YELLOW}frontier.service not found${NC}"
fi

# Reload systemd
echo -e "${YELLOW}Reloading systemd daemon...${NC}"
systemctl daemon-reload

# Check if any processes are still using the installation directory
echo -e "${YELLOW}Checking for processes using installation directory...${NC}"
REMAINING_PIDS=$(pgrep -f "/opt/liaison" 2>/dev/null)
if [[ -n "$REMAINING_PIDS" ]]; then
    echo -e "${RED}Warning: Found processes still using /opt/liaison:${NC}"
    ps -fp $REMAINING_PIDS 2>/dev/null
    echo -e "${YELLOW}Attempting to force kill these processes...${NC}"
    kill -KILL $REMAINING_PIDS 2>/dev/null
    sleep 2
    # Verify they're gone
    REMAINING_PIDS=$(pgrep -f "/opt/liaison" 2>/dev/null)
    if [[ -n "$REMAINING_PIDS" ]]; then
        echo -e "${RED}Warning: Some processes could not be killed. You may need to manually stop them.${NC}"
    else
        echo -e "${GREEN}All processes terminated${NC}"
    fi
else
    echo -e "${GREEN}No processes using installation directory${NC}"
fi

# Ask about removing data
echo ""
read -p "Do you want to remove the installation directory ($INSTALL_DIR)? [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [[ -d "$INSTALL_DIR" ]]; then
        # Try to unmount any filesystems in the directory first
        umount "$INSTALL_DIR"/* 2>/dev/null || true
        
        # Remove directory
        if rm -rf "$INSTALL_DIR" 2>/dev/null; then
            echo -e "${GREEN}Installation directory removed${NC}"
        else
            echo -e "${RED}Warning: Could not remove installation directory${NC}"
            echo -e "${YELLOW}Some files may still be in use. You may need to manually remove: $INSTALL_DIR${NC}"
            # Try to remove what we can
            find "$INSTALL_DIR" -type f -delete 2>/dev/null || true
            find "$INSTALL_DIR" -type d -empty -delete 2>/dev/null || true
        fi
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
