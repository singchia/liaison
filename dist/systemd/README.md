# Liaison Systemd Service

This directory contains systemd service files and installation scripts for the Liaison edge computing management service.

## Files

- `liaison.service` - Systemd service unit file
- `install.sh` - Installation script
- `uninstall.sh` - Uninstallation script
- `README.md` - This documentation

## Installation

### Prerequisites

- Linux system with systemd
- Root privileges
- Liaison binary compiled and ready

### Quick Installation

1. Make the installation script executable:
   ```bash
   chmod +x install.sh
   ```

2. Run the installation script as root:
   ```bash
   sudo ./install.sh
   ```

3. Copy the liaison binary to the installation directory:
   ```bash
   sudo cp /path/to/liaison /opt/liaison/bin/
   sudo chown liaison:liaison /opt/liaison/bin/liaison
   sudo chmod +x /opt/liaison/bin/liaison
   ```

4. Copy the configuration file:
   ```bash
   sudo cp /path/to/liaison.yaml /opt/liaison/etc/
   sudo chown liaison:liaison /opt/liaison/etc/liaison.yaml
   ```

5. Start the service:
   ```bash
   sudo systemctl start liaison
   ```

### Manual Installation

If you prefer to install manually:

1. Create the service user:
   ```bash
   sudo useradd --system --no-create-home --shell /bin/false liaison
   ```

2. Create directories:
   ```bash
   sudo mkdir -p /opt/liaison/{bin,etc,data,logs}
   sudo chown -R liaison:liaison /opt/liaison
   sudo chmod 750 /opt/liaison/{data,logs}
   ```

3. Copy the service file:
   ```bash
   sudo cp liaison.service /etc/systemd/system/
   sudo systemctl daemon-reload
   sudo systemctl enable liaison
   ```

## Service Management

### Basic Commands

```bash
# Start the service
sudo systemctl start liaison

# Stop the service
sudo systemctl stop liaison

# Restart the service
sudo systemctl restart liaison

# Check service status
sudo systemctl status liaison

# View service logs
sudo journalctl -u liaison -f

# Enable service to start on boot
sudo systemctl enable liaison

# Disable service from starting on boot
sudo systemctl disable liaison
```

### Log Management

The service logs to systemd journal by default. You can view logs using:

```bash
# View recent logs
sudo journalctl -u liaison

# Follow logs in real-time
sudo journalctl -u liaison -f

# View logs from today
sudo journalctl -u liaison --since today

# View logs with timestamps
sudo journalctl -u liaison -o short-iso
```

## Configuration

The service expects the configuration file at `/opt/liaison/etc/liaison.yaml`. Make sure to:

1. Update the configuration file according to your environment
2. Ensure the database path is writable by the liaison user
3. Configure the correct frontier server addresses
4. Set appropriate log levels and file paths

### Example Configuration

```yaml
manager:
  listen:
    addr: 0.0.0.0:8080
    network: tcp
    tls:
      enable: false
  db: /opt/liaison/data/liaison.db
frontier:
  dial:
    addrs:
      - 127.0.0.1:30011
    network: tcp
    tls:
      enable: false
log:
  level: info
  file: /opt/liaison/logs/liaison.log
  maxsize: 100
  maxrolls: 10
```

## Security Features

The service is configured with several security features:

- Runs as a dedicated system user (`liaison`)
- Uses `NoNewPrivileges` to prevent privilege escalation
- Has restricted file system access with `ProtectSystem=strict`
- Uses `PrivateTmp` for isolated temporary directories
- Has resource limits to prevent resource exhaustion
- Restricts various system capabilities

## Troubleshooting

### Service Won't Start

1. Check the service status:
   ```bash
   sudo systemctl status liaison
   ```

2. Check the logs:
   ```bash
   sudo journalctl -u liaison --no-pager
   ```

3. Verify the binary exists and is executable:
   ```bash
   ls -la /opt/liaison/bin/liaison
   ```

4. Check file permissions:
   ```bash
   ls -la /opt/liaison/
   ```

### Permission Issues

1. Ensure the liaison user owns the installation directory:
   ```bash
   sudo chown -R liaison:liaison /opt/liaison
   ```

2. Check that the configuration file is readable:
   ```bash
   sudo -u liaison cat /opt/liaison/etc/liaison.yaml
   ```

### Database Issues

1. Ensure the database directory is writable:
   ```bash
   sudo chmod 750 /opt/liaison/data
   sudo chown liaison:liaison /opt/liaison/data
   ```

2. Check disk space:
   ```bash
   df -h /opt/liaison/data
   ```

## Uninstallation

To remove the service:

1. Make the uninstallation script executable:
   ```bash
   chmod +x uninstall.sh
   ```

2. Run the uninstallation script as root:
   ```bash
   sudo ./uninstall.sh
   ```

The script will:
- Stop and disable the service
- Remove the systemd service file
- Optionally remove the installation directory
- Optionally remove the service user

## Directory Structure

After installation, the directory structure will be:

```
/opt/liaison/
├── bin/
│   └── liaison          # Liaison binary
├── etc/
│   └── liaison.yaml     # Configuration file
├── data/
│   └── liaison.db       # Database file
└── logs/
    └── liaison.log      # Log file
```

## Support

For issues and questions:
- Check the service logs: `sudo journalctl -u liaison -f`
- Verify configuration: `sudo -u liaison /opt/liaison/bin/liaison -h`
- Review system resources: `systemctl status liaison`
