# ZenMoney Backup

Automatically backup your [ZenMoney](https://zenmoney.ru) data using the official API with OAuth token authentication.

- 🔐 Secure API-based authentication with OAuth tokens
- 📦 Full data export as JSON files  
- ⏰ Configurable backup schedule
- 🐳 Lightweight Docker image (~8.5MB)
- 🔄 Automatic retry and error handling

---
<div align="center">

[![Build Status](https://github.com/egregors/zenmoney-backup/actions/workflows/go.yml/badge.svg)](https://github.com/egregors/zenmoney-backup/actions) 
[![Coverage Status](https://coveralls.io/repos/github/egregors/zenmoney-backup/badge.svg?branch=main)](https://coveralls.io/github/egregors/zenmoney-backup?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/egregors/zenmoney-backup)](https://goreportcard.com/report/github.com/egregors/zenmoney-backup)

</div>

## 🚀 Quick Start

### Get Your API Token

First, you need to obtain an OAuth token from ZenMoney:

1. Visit [https://zerro.app/token](https://zerro.app/token) (or [https://zenmoney.ru/api](https://zenmoney.ru/api))
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
```

## 📋 Command Line Options

| Short | Long | Environment | Description |
|-------|------|-------------|-------------|
| `-t` | `--zenmoney OAuth token` | `ZEN_TOKEN` | ZenMoney API Token (required) |
| `-p` | `--sleep_time` | `SLEEP_TIME` | Backup interval (default: 24h) |
| | `--dbg` | `DEBUG` | Enable debug mode |

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
- `DEBUG`: Set to "true" for debug logging

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
