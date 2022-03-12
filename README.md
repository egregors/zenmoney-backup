# zenmoney-backup 

Backup your [zenmoney](zenmoney.ru) data by schedule.

- Backup all your data as `csv` files
- Different storage backends
- Configurable backup schedule
- Docker image

---
<div align="center">

[![Build Status](https://github.com/egregors/zenmoney-backup/actions/workflows/go.yml/badge.svg)](https://github.com/egregors/zenmoney-backup/actions) [![Coverage Status](https://coveralls.io/repos/github/egregors/zenmoney-backup/badge.svg?branch=main)](https://coveralls.io/github/egregors/zenmoney-backup?branch=main)

</div>

## Usage

### Binary

To make binary just download this repo and run `make build`.

```shell
git clone https://github.com/egregors/zenmoney-backup.git
make build
```

### Params

| short | long           | ENV          |                                                         |
|-------|----------------|--------------|---------------------------------------------------------|
| -l    | --zen_username | ZEN_USERNAME | Your zenmoney login                                     |
| -p    | --zen_password | ZEN_PASSWORD | Your zenmoney password                                  |
| -t    | --sleep_time   | SLEEP_TIME   | Backup performs every SLEEP_TIME minutes (default: 24h) |
|       | --dbg          | DEBUG        | Debug mode                                              |

### Docker

```shell
docker run --rm -e ZEN_USERNAME=*** -e ZEN_PASSWORD=*** -e SLEEP_TIME=24h -v $(pwd):/backups zenb
```

