#!/bin/bash

# Liaison Service Installation Script
# This script installs and configures the Liaison systemd service

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="liaison"
SERVICE_USER="liaison"
SERVICE_GROUP="liaison"
DEFAULT_EMAIL="default@liaison.com"
INSTALL_DIR="/opt/liaison"
CONFIG_DIR="/opt/liaison/conf"
DATA_DIR="/opt/liaison/data"
LOG_DIR="/opt/liaison/logs"
BIN_DIR="/opt/liaison/bin"
WEB_DIR="/opt/liaison/web"
EDGE_DIR="/opt/liaison/edge"
CERTS_DIR="/opt/liaison/certs"

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
mkdir -p "$INSTALL_DIR" "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR" "$BIN_DIR" "$WEB_DIR" "$EDGE_DIR" "$CERTS_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR"
chmod 755 "$INSTALL_DIR"

# Set specific permissions
chmod 750 "$DATA_DIR" "$LOG_DIR"
chmod 755 "$BIN_DIR" "$CONFIG_DIR" "$WEB_DIR" "$EDGE_DIR"
chmod 750 "$CERTS_DIR"

# Copy service files
echo -e "${YELLOW}Installing systemd service files...${NC}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
if [[ -d "$SCRIPT_DIR/systemd" ]]; then
    # Copy liaison.service
    if [[ -f "$SCRIPT_DIR/systemd/liaison.service" ]]; then
        cp "$SCRIPT_DIR/systemd/liaison.service" "/etc/systemd/system/"
        chmod 644 "/etc/systemd/system/liaison.service"
        echo -e "${GREEN}liaison.service installed${NC}"
    fi
    # Copy frontier.service
    if [[ -f "$SCRIPT_DIR/systemd/frontier.service" ]]; then
        cp "$SCRIPT_DIR/systemd/frontier.service" "/etc/systemd/system/"
        chmod 644 "/etc/systemd/system/frontier.service"
        echo -e "${GREEN}frontier.service installed${NC}"
    fi
elif [[ -f "$SCRIPT_DIR/liaison.service" ]]; then
    cp "$SCRIPT_DIR/liaison.service" "/etc/systemd/system/"
    chmod 644 "/etc/systemd/system/liaison.service"
else
    echo -e "${RED}Error: liaison.service not found${NC}"
    exit 1
fi

# Copy binaries
echo -e "${YELLOW}Copying binaries...${NC}"
if [[ -d "$SCRIPT_DIR/bin" ]]; then
    cp -f "$SCRIPT_DIR/bin/"* "$BIN_DIR/"
    chown "$SERVICE_USER:$SERVICE_GROUP" "$BIN_DIR"/*
    chmod 755 "$BIN_DIR"/*
    echo -e "${GREEN}Binaries copied${NC}"
    # Check if frontier binary exists
    if [[ -f "$BIN_DIR/frontier" ]]; then
        echo -e "${GREEN}  - frontier binary found${NC}"
    fi
else
    echo -e "${YELLOW}Warning: bin directory not found, skipping binary copy${NC}"
fi

# Get public IP address
echo -e "${YELLOW}Detecting public IP address...${NC}"
PUBLIC_ADDR=$(curl -s --max-time 5 ifconfig.me 2>/dev/null || curl -s --max-time 5 ifconfig.co 2>/dev/null || echo "localhost")
if [ -z "$PUBLIC_ADDR" ] || [ "$PUBLIC_ADDR" = "localhost" ]; then
    echo -e "${YELLOW}Warning: Could not detect public IP automatically.${NC}"
    echo -e "${YELLOW}Please manually set the public IP in $CONFIG_DIR/liaison.yaml after installation.${NC}"
    PUBLIC_ADDR="localhost"
else
    echo -e "${GREEN}Auto-detected public IP: ${BOLD}${CYAN}$PUBLIC_ADDR${NC}${GREEN}${NC}"
    echo -e "${YELLOW}⚠️  If this IP is incorrect, you can edit $CONFIG_DIR/liaison.yaml later${NC}"
fi

# Generate random JWT secret (32 characters minimum for security)
echo -e "${YELLOW}Generating JWT secret key...${NC}"
JWT_SECRET=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-32 2>/dev/null)
if [[ -z "$JWT_SECRET" ]]; then
    JWT_SECRET=$(head -c 32 /dev/urandom | base64 | tr -d "=+/" | cut -c1-32 2>/dev/null)
fi
if [[ -z "$JWT_SECRET" ]]; then
    JWT_SECRET=$(date +%s | sha256sum | base64 | tr -d "=+/" | cut -c1-32)
fi
if [[ ${#JWT_SECRET} -lt 32 ]]; then
    # Pad to 32 characters if needed
    JWT_SECRET="${JWT_SECRET}$(openssl rand -hex 16 | head -c $((32 - ${#JWT_SECRET})))"
fi
echo -e "${GREEN}JWT secret key generated${NC}"

# Render configuration templates
echo -e "${YELLOW}Rendering configuration files from templates...${NC}"
if [[ -d "$SCRIPT_DIR/conf" ]]; then
    # Render liaison.yaml from template
    if [[ -f "$SCRIPT_DIR/conf/liaison.yaml.template" ]]; then
        # Replace ${PUBLIC_ADDR} and ${JWT_SECRET} in the template
        sed -e "s|\${PUBLIC_ADDR}|${PUBLIC_ADDR}|g" \
            -e "s|\${JWT_SECRET}|${JWT_SECRET}|g" \
            "$SCRIPT_DIR/conf/liaison.yaml.template" > "$CONFIG_DIR/liaison.yaml"
        chown "$SERVICE_USER:$SERVICE_GROUP" "$CONFIG_DIR/liaison.yaml"
        chmod 644 "$CONFIG_DIR/liaison.yaml"
        echo -e "${GREEN}liaison.yaml rendered with public IP: ${BOLD}${CYAN}$PUBLIC_ADDR${NC}${GREEN}${NC}"
        echo -e "${GREEN}JWT secret key saved to configuration${NC}"
    fi
    
    # Render frontier.yaml from template (if it has PUBLIC_ADDR variable)
    if [[ -f "$SCRIPT_DIR/conf/frontier.yaml.template" ]]; then
        # Replace ${PUBLIC_ADDR} if it exists in the template
        sed "s|\${PUBLIC_ADDR}|${PUBLIC_ADDR}|g" "$SCRIPT_DIR/conf/frontier.yaml.template" > "$CONFIG_DIR/frontier.yaml"
        chown "$SERVICE_USER:$SERVICE_GROUP" "$CONFIG_DIR/frontier.yaml"
        chmod 644 "$CONFIG_DIR/frontier.yaml"
        echo -e "${GREEN}frontier.yaml rendered${NC}"
    fi
fi

# Copy configuration files (fallback if templates don't exist)
if [[ -d "$SCRIPT_DIR/etc" ]]; then
    for yaml_file in "$SCRIPT_DIR/etc/"*.yaml; do
        if [[ -f "$yaml_file" ]]; then
            filename=$(basename "$yaml_file")
            # Only copy if not already rendered from template
            if [[ ! -f "$CONFIG_DIR/$filename" ]]; then
                cp -f "$yaml_file" "$CONFIG_DIR/"
                chown "$SERVICE_USER:$SERVICE_GROUP" "$CONFIG_DIR/$filename"
                chmod 644 "$CONFIG_DIR/$filename"
            fi
        fi
    done
    if [[ -f "$CONFIG_DIR"/*.yaml ]]; then
        echo -e "${GREEN}Configuration files copied${NC}"
    fi
fi

# Copy web frontend files
echo -e "${YELLOW}Copying web frontend files...${NC}"
if [[ -d "$SCRIPT_DIR/web" ]]; then
    cp -r "$SCRIPT_DIR/web/"* "$WEB_DIR/" 2>/dev/null || true
    chown -R "$SERVICE_USER:$SERVICE_GROUP" "$WEB_DIR"
    chmod -R 755 "$WEB_DIR"
    echo -e "${GREEN}Web frontend files copied${NC}"
else
    echo -e "${YELLOW}Warning: web directory not found, skipping web files copy${NC}"
fi

# Copy edge binaries and scripts for all platforms
echo -e "${YELLOW}Copying edge binaries and scripts for all platforms...${NC}"
if [[ -d "$SCRIPT_DIR/edge" ]]; then
    cp -f "$SCRIPT_DIR/edge/"* "$EDGE_DIR/" 2>/dev/null || true
    chown "$SERVICE_USER:$SERVICE_GROUP" "$EDGE_DIR"/* 2>/dev/null || true
    # Set executable permission for scripts, read-only for tar.gz files
    chmod 755 "$EDGE_DIR"/install.sh 2>/dev/null || true
    chmod 755 "$EDGE_DIR"/uninstall.sh 2>/dev/null || true
    chmod 644 "$EDGE_DIR"/*.tar.gz 2>/dev/null || true
    echo -e "${GREEN}Edge files copied${NC}"
else
    echo -e "${YELLOW}Warning: edge directory not found, skipping edge files copy${NC}"
fi

# Generate TLS certificates
echo -e "${YELLOW}Generating TLS certificates...${NC}"
if command -v openssl >/dev/null 2>&1; then
    if [[ ! -f "$CERTS_DIR/server.crt" ]] || [[ ! -f "$CERTS_DIR/server.key" ]]; then
        openssl req -x509 -newkey rsa:4096 \
            -keyout "$CERTS_DIR/server.key" \
            -out "$CERTS_DIR/server.crt" \
            -days 365 \
            -nodes \
            -subj "/C=CN/ST=Beijing/L=Beijing/O=Liaison/OU=IT/CN=localhost" \
            2>/dev/null
        chmod 600 "$CERTS_DIR/server.key"
        chmod 644 "$CERTS_DIR/server.crt"
        chown "$SERVICE_USER:$SERVICE_GROUP" "$CERTS_DIR/server.key" "$CERTS_DIR/server.crt"
        echo -e "${GREEN}TLS certificates generated in $CERTS_DIR${NC}"
    else
        echo -e "${YELLOW}Certificates already exist, skipping generation${NC}"
    fi
else
    echo -e "${YELLOW}Warning: openssl not found, skipping certificate generation${NC}"
    echo -e "${YELLOW}You may need to generate certificates manually or install openssl${NC}"
fi

# Generate initial password
# Reload systemd
echo -e "${YELLOW}Reloading systemd daemon...${NC}"
systemctl daemon-reload

# Enable services
echo -e "${YELLOW}Enabling services...${NC}"
systemctl enable "$SERVICE_NAME"
if [[ -f "/etc/systemd/system/frontier.service" ]]; then
    systemctl enable frontier
    echo -e "${GREEN}frontier service enabled${NC}"
fi

# Start services in order: frontier first, then liaison
echo -e "${YELLOW}Starting services...${NC}"
LIAISON_STARTED=false
FRONTIER_STARTED=false
if [[ -f "/etc/systemd/system/frontier.service" ]]; then
    echo -e "${YELLOW}Starting frontier service...${NC}"
    # Reload systemd to pick up any service file changes
    systemctl daemon-reload
    systemctl start frontier
    sleep 2  # Give service a moment to start
    if systemctl is-active --quiet frontier; then
        echo -e "${GREEN}frontier service started${NC}"
        FRONTIER_STARTED=true
        echo -e "${YELLOW}Waiting 3 seconds for frontier to initialize...${NC}"
        sleep 3
    else
        echo -e "${RED}Warning: frontier service failed to start${NC}"
        systemctl status frontier --no-pager -l || true
        echo -e "${YELLOW}Trying to check logs: journalctl -u frontier -n 20 --no-pager${NC}"
        journalctl -u frontier -n 20 --no-pager || true
    fi
fi

echo -e "${YELLOW}Starting liaison service...${NC}"
# Reload systemd to pick up any service file changes
systemctl daemon-reload
systemctl start "$SERVICE_NAME"
sleep 2  # Give service a moment to start

# Wait for liaison service to be fully active and database initialized
MAX_WAIT=30
WAIT_COUNT=0
while [[ $WAIT_COUNT -lt $MAX_WAIT ]]; do
    if systemctl is-active --quiet "$SERVICE_NAME"; then
        # Check if database file exists (liaison creates it on startup)
        DB_FILE="/opt/liaison/data/liaison.db"
        if [[ -f "$DB_FILE" ]]; then
            # Give it a bit more time to ensure database schema is initialized
            sleep 2
            break
        fi
    fi
    sleep 1
    WAIT_COUNT=$((WAIT_COUNT + 1))
done

if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo -e "${GREEN}liaison service started${NC}"
    LIAISON_STARTED=true
    
    # Now create default user and password after liaison service is running
    echo -e "${YELLOW}Generating initial password...${NC}"
    INITIAL_PASSWORD=""
    # Temporarily disable exit on error for password generation
    set +e
    if [[ -f "$BIN_DIR/password-generator" ]]; then
        # Generate random password directly in the script
        # Using /dev/urandom for better randomness
        INITIAL_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16 2>/dev/null)
        if [[ -z "$INITIAL_PASSWORD" ]]; then
            # Fallback: use /dev/urandom if openssl is not available
            INITIAL_PASSWORD=$(head -c 16 /dev/urandom | base64 | tr -d "=+/" | cut -c1-16 2>/dev/null)
        fi
        if [[ -z "$INITIAL_PASSWORD" ]]; then
            # Last resort: use date + random number
            INITIAL_PASSWORD="$(date +%s | sha256sum | base64 | tr -d "=+/" | cut -c1-16)"
        fi
        
        if [[ -n "$INITIAL_PASSWORD" ]]; then
            # Use password-generator to hash and store the password
            # Wait a bit more to ensure database is fully ready
            sleep 1
            "$BIN_DIR/password-generator" -password "$INITIAL_PASSWORD" -email "$DEFAULT_EMAIL" -create >/dev/null 2>&1
            PASSWORD_EXIT_CODE=$?
            set -e  # Re-enable exit on error
            
            if [[ $PASSWORD_EXIT_CODE -eq 0 ]]; then
                echo -e "${GREEN}Default user created successfully${NC}"
            else
                echo -e "${YELLOW}Warning: Failed to store password in database${NC}"
                echo -e "${YELLOW}You can manually set the password later using:${NC}"
                echo -e "${YELLOW}  $BIN_DIR/password-generator -password <password> -email $DEFAULT_EMAIL -create${NC}"
                INITIAL_PASSWORD=""  # Clear password if creation failed
            fi
        else
            set -e  # Re-enable exit on error
            echo -e "${YELLOW}Warning: Failed to generate random password${NC}"
        fi
    else
        set -e  # Re-enable exit on error
        echo -e "${YELLOW}Warning: password-generator not found, skipping password generation${NC}"
    fi
else
    echo -e "${RED}Warning: liaison service failed to start${NC}"
    systemctl status "$SERVICE_NAME" --no-pager -l || true
    echo -e "${YELLOW}Trying to check logs: journalctl -u $SERVICE_NAME -n 20 --no-pager${NC}"
    journalctl -u "$SERVICE_NAME" -n 20 --no-pager || true
fi

# Check frontier status
if [[ -f "/etc/systemd/system/frontier.service" ]]; then
    if systemctl is-active --quiet frontier; then
        FRONTIER_STARTED=true
    fi
fi

echo -e "${GREEN}Installation completed successfully!${NC}"
echo ""
echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}${GREEN}  Installation Summary${NC}"
echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo ""
if [[ -n "$INITIAL_PASSWORD" ]]; then
    echo -e "${BOLD}${YELLOW}  Default Credentials:${NC}"
    echo -e "${BOLD}${CYAN}    Username: ${DEFAULT_EMAIL}${NC}"
    echo -e "${BOLD}${CYAN}    Password: ${INITIAL_PASSWORD}${NC}"
    echo ""
fi
echo -e "${BOLD}${YELLOW}  Public IP Address (Auto-detected):${NC}"
echo -e "${BOLD}${CYAN}    ${PUBLIC_ADDR}${NC}"
if [[ "$PUBLIC_ADDR" != "localhost" ]]; then
    echo -e "${YELLOW}    ⚠️  If this IP is incorrect, please edit:${NC}"
    echo -e "${YELLOW}       $CONFIG_DIR/liaison.yaml${NC}"
    echo -e "${YELLOW}       and update the 'server_url' field${NC}"
else
    echo -e "${YELLOW}    ⚠️  Public IP could not be auto-detected.${NC}"
    echo -e "${YELLOW}       Please edit $CONFIG_DIR/liaison.yaml and set 'server_url' manually${NC}"
fi
echo ""
echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo ""
# Service status
echo -e "${BOLD}${YELLOW}  Service Status:${NC}"
if [[ "$FRONTIER_STARTED" == "true" ]]; then
    echo -e "${GREEN}    ✓ frontier: Running${NC}"
else
    echo -e "${YELLOW}    ✗ frontier: Not running${NC}"
fi
if [[ "$LIAISON_STARTED" == "true" ]]; then
    echo -e "${GREEN}    ✓ liaison: Running${NC}"
else
    echo -e "${YELLOW}    ✗ liaison: Not running${NC}"
fi
echo ""
# Access information
if [[ "$LIAISON_STARTED" == "true" ]]; then
    if [[ "$PUBLIC_ADDR" != "localhost" ]]; then
        echo -e "${BOLD}${YELLOW}  Access URL:${NC}"
        echo -e "${BOLD}${CYAN}    https://${PUBLIC_ADDR}${NC}"
        echo ""
    fi
fi
echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review configuration files in $CONFIG_DIR/"
if [[ "$PUBLIC_ADDR" != "localhost" ]]; then
    echo "2. If the detected public IP (${PUBLIC_ADDR}) is incorrect,"
    echo "   edit $CONFIG_DIR/liaison.yaml and update the 'server_url' field"
else
    echo "2. The public IP could not be auto-detected."
    echo "   Please edit $CONFIG_DIR/liaison.yaml and set the 'server_url' field manually"
fi
if [[ -n "$INITIAL_PASSWORD" ]]; then
    echo "3. ⚠️  Please change the default password after first login!"
fi
if [[ -f "/etc/systemd/system/frontier.service" ]]; then
    echo "4. Check service status: systemctl status $SERVICE_NAME frontier"
    echo "5. View logs: journalctl -u $SERVICE_NAME -f (or journalctl -u frontier -f)"
else
    echo "3. Check service status: systemctl status $SERVICE_NAME"
    echo "4. View logs: journalctl -u $SERVICE_NAME -f"
fi
echo ""
echo -e "${GREEN}Service management commands:${NC}"
if [[ -f "/etc/systemd/system/frontier.service" ]]; then
    echo "  Start:   systemctl start $SERVICE_NAME frontier"
    echo "  Stop:    systemctl stop $SERVICE_NAME frontier"
    echo "  Restart: systemctl restart $SERVICE_NAME frontier"
    echo "  Status:  systemctl status $SERVICE_NAME frontier"
    echo "  Logs:    journalctl -u $SERVICE_NAME -f (or journalctl -u frontier -f)"
else
    echo "  Start:   systemctl start $SERVICE_NAME"
    echo "  Stop:    systemctl stop $SERVICE_NAME"
    echo "  Restart: systemctl restart $SERVICE_NAME"
    echo "  Status:  systemctl status $SERVICE_NAME"
    echo "  Logs:    journalctl -u $SERVICE_NAME -f"
fi
