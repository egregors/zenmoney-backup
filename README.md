# ZenMoney Backup

Automatically backup your [ZenMoney](https://zenmoney.ru) data using the official API with OAuth token authentication.

- 🔐 Secure API-based authentication with OAuth tokens
- 📦 Full data export as JSON files  
- ⏰ Configurable backup schedule
- 🐳 Lightweight Docker image (~8.5MB)
- 🔄 Automatic retry and error handling
- 🔔 Error notifications via ntfy.sh

---
<div align="center">

[![Build Status](https://github.com/egregors/zenmoney-backup/actions/workflows/go.yml/badge.svg)](https://github.com/egregors/zenmoney-backup/actions) 
[![Coverage Status](https://coveralls.io/repos/github/egregors/zenmoney-backup/badge.svg?branch=main)](https://coveralls.io/github/egregors/zenmoney-backup?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/egregors/zenmoney-backup)](https://goreportcard.com/report/github.com/egregors/zenmoney-backup)

</div>

## 🚀 Quick Start

### Get Your API Token

First, you need to obtain an OAuth token from ZenMoney:

1. Visit [https://zerro.app/token](https://zerro.app/token)
2. Log in to your ZenMoney account
3. Generate an API token
4. Copy the token for use with the backup tool

### Docker (Recommended)

The easiest way to run backups is using the Docker image:

```bash
docker run --rm \
  -e ZEN_TOKEN="your_oauth_token_here" \
  -e SLEEP_TIME="24h" \
  -v $(pwd)/backups:/backups \
  zenb:latest
```

Replace `your_oauth_token_here` with your actual OAuth token from the step above.

**Parameters:**
- `ZEN_TOKEN`: Your ZenMoney OAuth token (required)
- `SLEEP_TIME`: Backup interval (default: 24h)
- `TIMEOUT`: Backup request timeout in seconds (default: 10)
- `NOTIFY_URL`: ntfy.sh notification URL for error alerts (optional)
- `DEBUG`: Set to `true` for debug logging

**Volume mounting:**
- The container saves backups to `/backups` directory
- Mount your local directory to persist backups: `-v $(pwd)/backups:/backups`

### Binary

You can also run the application as a standalone binary:

```bash
# Download and build
git clone https://github.com/egregors/zenmoney-backup.git
cd zenmoney-backup
make build-local

# Run with your token
./build/zenb -t "your_oauth_token_here" --sleep_time="24h"

# Run with custom timeout (useful for large backup files)
./build/zenb -t "your_oauth_token_here" -c 30
```

## 📋 Command Line Options

| Short | Long | Environment | Description |
|-------|------|-------------|-------------|
| `-t` | `--zenmoney OAuth token` | `ZEN_TOKEN` | ZenMoney API Token (required) |
| `-p` | `--sleep_time` | `SLEEP_TIME` | Backup interval (default: 24h) |
| `-c` | `--timeout` | `TIMEOUT` | Backup request timeout in seconds (default: 10) |
| `-n` | `--notify_url` | `NOTIFY_URL` | ntfy.sh notification URL (optional) |
| | `--dbg` | `DEBUG` | Enable debug mode |

## 🔔 Error Notifications

ZenMoney Backup supports error notifications via [ntfy.sh](https://ntfy.sh). When configured, you'll receive push notifications whenever a backup error occurs (such as API failures, network issues, or storage problems).

### Setting up ntfy.sh Notifications

1. **Choose a topic name** for your notifications (e.g., `zenmoney-backup-alerts`)
2. **Subscribe to the topic** on your device:
   - **Mobile**: Download the [ntfy app](https://ntfy.sh) and subscribe to your topic
   - **Desktop**: Use the [web app](https://ntfy.sh/app) or desktop client
3. **Configure the notification URL**:

```bash
# With Docker
docker run --rm \
  -e ZEN_TOKEN="your_token" \
  -e NOTIFY_URL="https://ntfy.sh/your_topic" \
  -v $(pwd)/backups:/backups \
  zenb:latest

# With binary
./build/zenb -t "your_token" -n "https://ntfy.sh/your_topic"
```

**Important:** Choose a unique topic name that others won't guess. Anyone who knows your topic name can read your notifications.

### Self-hosted ntfy

You can also use a self-hosted ntfy server:

```bash
-e NOTIFY_URL="https://your-ntfy-server.com/your_topic"
```

## 📁 Backup Format

The tool creates JSON backup files in the `backups/` directory with the following naming convention:

```
zen_2024-06-29_15-30-45.json
```

Each backup contains:
- All transactions
- Account information
- Categories
- Tags and labels
- Other ZenMoney data

The JSON format preserves all data structure and can be easily processed by other tools if needed.

## 🔧 Development

### Prerequisites

- Go 1.24+
- Docker (optional)
- Make

### Building

```bash
# Build for current OS
make build-local

# Build for Linux (Docker-compatible)
make build

# Build Docker image
make docker

# Run tests
make test

# View all available commands
make help
```

### Project Structure

```
├── cmd/           # Application entry point
├── srv/           # Backup server logic
├── store/         # Storage implementations
├── backups/       # Default backup directory (created automatically)
├── Dockerfile     # Docker build configuration
└── Makefile       # Build automation
```

## 🐳 Docker

### Build Image

```bash
make docker
```

### Run Container

```bash
# Interactive run
make docker-run

# Background with custom settings
docker run -d \
  --name zenmoney-backup \
  -e ZEN_TOKEN="your_token" \
  -e SLEEP_TIME="12h" \
  -v /host/path/to/backups:/backups \
  zenb:latest

# Check logs
docker logs zenmoney-backup
```

### Environment Variables

- `ZEN_TOKEN`: Your ZenMoney OAuth token (required)
- `SLEEP_TIME`: Backup interval (e.g., "1h", "30m", "24h")
- `TIMEOUT`: Backup request timeout in seconds (default: 10)
- `NOTIFY_URL`: ntfy.sh notification URL for error alerts (optional)
- `DEBUG`: Set to "true" for debug logging

## 🔄 Autostart on Linux (systemd)

You can configure ZenMoney Backup to run automatically as a systemd service on Linux. This allows the backup tool to start on boot and restart automatically if it fails.

### Binary Variant

Create a systemd unit file for running the compiled binary:

**File:** `/etc/systemd/system/zenmoney-backup.service`

```ini
[Unit]
# Description of the service
Description=ZenMoney Backup Service
# Start after the network is available
After=network-online.target
Wants=network-online.target

[Service]
# Run the service as a specific user (replace 'your_username' with actual system username)
User=your_username
# Working directory for the service
WorkingDirectory=/opt/zenmoney-backup
# Command to execute - adjust path to your binary location
ExecStart=/opt/zenmoney-backup/zenb -t "${ZEN_TOKEN}" --sleep_time="${SLEEP_TIME}"
# Restart policy - always restart if the service stops
Restart=always
# Wait 10 seconds before restarting after failure
RestartSec=10
# Environment variables for the service
Environment="ZEN_TOKEN=your_token_here"
Environment="SLEEP_TIME=24h"
Environment="TIMEOUT=10"
# Optional: Enable error notifications via ntfy.sh
# Environment="NOTIFY_URL=https://ntfy.sh/your_topic"
# Optional: Enable debug logging
# Environment="DEBUG=true"

[Install]
# Enable the service to start on boot (multi-user target)
WantedBy=multi-user.target
```

**Installation Steps:**

1. Create the working directory and copy your binary:
   ```bash
   sudo mkdir -p /opt/zenmoney-backup
   sudo cp /path/to/zenb /opt/zenmoney-backup/
   sudo chmod +x /opt/zenmoney-backup/zenb
   ```

2. Create the systemd unit file:
   ```bash
   sudo nano /etc/systemd/system/zenmoney-backup.service
   # Paste the configuration above, replacing placeholders with your values
   ```

3. Set proper permissions:
   ```bash
   sudo chmod 644 /etc/systemd/system/zenmoney-backup.service
   ```

4. Reload systemd to recognize the new service:
   ```bash
   sudo systemctl daemon-reload
   ```

5. Enable the service to start on boot:
   ```bash
   sudo systemctl enable zenmoney-backup.service
   ```

6. Start the service:
   ```bash
   sudo systemctl start zenmoney-backup.service
   ```

7. Check the service status:
   ```bash
   sudo systemctl status zenmoney-backup.service
   ```

8. View service logs:
   ```bash
   sudo journalctl -u zenmoney-backup.service -f
   ```

### Docker Variant

Create a systemd unit file for running the backup tool via Docker:

**File:** `/etc/systemd/system/zenmoney-backup.service`

```ini
[Unit]
# Description of the service
Description=ZenMoney Backup Docker Service
# Start after Docker is available
After=docker.service
Requires=docker.service

[Service]
# Type of service - simple for long-running containers
Type=simple
# Remove any existing container with the same name before starting
# The leading dash (-) means systemd will ignore failure if container doesn't exist
ExecStartPre=-/usr/bin/docker rm -f zenmoney-backup
# Start the Docker container with required parameters
# Use backslash (\) for line continuation in systemd unit files
ExecStart=/usr/bin/docker run --rm \
  --name zenmoney-backup \
  -e ZEN_TOKEN=your_token_here \
  -e SLEEP_TIME=24h \
  -e TIMEOUT=10 \
  -e NOTIFY_URL=https://ntfy.sh/your_topic \
  -v /opt/zenmoney-backup/backups:/backups \
  zenb:latest
# Stop the container gracefully
ExecStop=/usr/bin/docker stop zenmoney-backup
# Restart policy - always restart if the container stops
Restart=always
# Wait 10 seconds before restarting after failure
RestartSec=10

[Install]
# Enable the service to start on boot (multi-user target)
WantedBy=multi-user.target
```

**Note:** To enable error notifications, uncomment or add the `NOTIFY_URL` environment variable with your ntfy.sh topic URL. Remove the line or leave it empty to disable notifications.

**Installation Steps:**

1. Create the backup directory:
   ```bash
   sudo mkdir -p /opt/zenmoney-backup/backups
   sudo chmod 755 /opt/zenmoney-backup/backups
   ```

2. Pull or build the Docker image:
   ```bash
   # Option 1: Build locally
   cd /path/to/zenmoney-backup
   make docker
   
   # Option 2: Pull from registry (if available)
   # docker pull zenb:latest
   ```

3. Create the systemd unit file:
   ```bash
   sudo nano /etc/systemd/system/zenmoney-backup.service
   # Paste the configuration above, replacing placeholders with your values
   ```

4. Set proper permissions:
   ```bash
   sudo chmod 644 /etc/systemd/system/zenmoney-backup.service
   ```

5. Reload systemd to recognize the new service:
   ```bash
   sudo systemctl daemon-reload
   ```

6. Enable the service to start on boot:
   ```bash
   sudo systemctl enable zenmoney-backup.service
   ```

7. Start the service:
   ```bash
   sudo systemctl start zenmoney-backup.service
   ```

8. Check the service status:
   ```bash
   sudo systemctl status zenmoney-backup.service
   ```

9. View service logs:
   ```bash
   sudo journalctl -u zenmoney-backup.service -f
   # Or use Docker logs directly:
   docker logs -f zenmoney-backup
   ```

**Notes:**
- Replace `your_token_here` with your actual ZenMoney OAuth token
- Replace `your_topic` with your unique ntfy.sh topic name (for notifications)
- Adjust paths (`/opt/zenmoney-backup`) to match your preferred installation location
- For the binary variant, replace `your_username` with the user that should run the service
- You can customize `SLEEP_TIME`, `TIMEOUT`, `NOTIFY_URL`, and other environment variables as needed
- Both variants will automatically restart the service if it crashes
- Use `sudo systemctl stop zenmoney-backup.service` to stop the service
- Use `sudo systemctl disable zenmoney-backup.service` to prevent auto-start on boot

## 📊 Example Output

```
zenmoney-backup v1.2.3
~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=[,,_,,]:3
[INFO] login...
[INFO] downloading...
[DEBUG] downloading data ...
[DEBUG] downloaded
[INFO] zen_2024-06-29_15-30-45.json saved
[INFO] sleep for 24h0m0s
```

## 🔒 Security Notes

- Keep your OAuth token secure and never commit it to version control
- The application doesn't store your credentials permanently
- Backup files contain sensitive financial data - store them securely
- Use environment variables or secure secret management in production

## 🤝 Contributing

Bug reports, bug fixes, and new features are always welcome! Please open issues and submit pull requests for any new code.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make test` and `make lint`
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙋‍♂️ Support

If you encounter any issues:

1. Check that your OAuth token is valid and active
2. Ensure you have proper network connectivity
3. Review the logs for error messages
4. Open an issue on GitHub with details

---

Made with ❤️ for the ZenMoney community (thanks @nemirlev for ZenMoney go SDK)
