# zenmoney-backup 

Backup your [zenmoney](zenmoney.ru) data by schedule.

- Backup all your data as `csv` files
- Different storage backends
- Configurable backup schedule
- Docker image

---
<div align="center">

[![Build Status](https://github.com/egregors/zenmoney-backup/actions/workflows/go.yml/badge.svg)](https://github.com/egregors/zenmoney-backup/actions) 
[![Coverage Status](https://coveralls.io/repos/github/egregors/zenmoney-backup/badge.svg?branch=main)](https://coveralls.io/github/egregors/zenmoney-backup?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/egregors/zenmoney-backup)](https://goreportcard.com/report/github.com/egregors/zenmoney-backup)

</div>

## Usage

The simplest way to run backing is use [docker image](https://github.com/egregors/zenmoney-backup/pkgs/container/zenmoney-backup%2Fzenb)

### Docker
To start backups pulling container just run:

```shell
docker run --rm                         \
  -e ZEN_USERNAME=your_zenmoney_login   \
  -e ZEN_PASSWORD=your_zenmoney_pass    \
  -e SLEEP_TIME=24h                     \
  -v $(pwd):/backups                    \
  ghcr.io/egregors/zenmoney-backup/zenb
```

Don't forget change `your_zenmoney_login` and `your_zenmoney_pass` to your login and pass respectively. 
Backup files will be saved in your current directory. To change it define absolut path to the folder you need instead of `$(pwd)`.

![termtosvg_hcs4cfax](https://user-images.githubusercontent.com/2153895/158082850-47c1d4fd-0883-44ea-a246-729a60d7e51d.svg)

To build `image` locally pull this repo and run `make docker`.

### Binary

You cat use binary as well. To make binary just download this repo and run `make build`.

```shell
git clone https://github.com/egregors/zenmoney-backup.git
make build
```

Credentials and settings could be passed like a CLI arguments either ENV.

```shell
./zenb -l MyUsername -p MySuperSecretPass --sleep_time=24h
```

#### Params

| short | long           | ENV          |                                                         |
|-------|----------------|--------------|---------------------------------------------------------|
| -l    | --zen_username | ZEN_USERNAME | Your zenmoney login                                     |
| -p    | --zen_password | ZEN_PASSWORD | Your zenmoney password                                  |
| -t    | --sleep_time   | SLEEP_TIME   | Backup performs every SLEEP_TIME minutes (default: 24h) |
|       | --dbg          | DEBUG        | Debug mode                                              |

## Development

Use `Makefile` to development stuff. 

```shell
git:(main) ✗ make help
Usage: make [task]

task                 help
------               ----
build                Build binary
docker               Build Docker image
run                  Run in debug mode
lint                 Lint the files
test                 Run tests
                     
update-go-deps       Updating Go dependencies
                     
help                 Show help message
```

## Contributing
Bug reports, bug fixes and new features are always welcome.
Please open issues and submit pull requests for any new code.
