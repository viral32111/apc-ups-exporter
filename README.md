# APC UPS Exporter

[![CI](https://github.com/viral32111/apc-ups-exporter/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/viral32111/apc-ups-exporter/actions/workflows/ci.yml)
[![CodeQL](https://github.com/viral32111/apc-ups-exporter/actions/workflows/codeql.yml/badge.svg)](https://github.com/viral32111/apc-ups-exporter/actions/workflows/codeql.yml)
![GitHub tag (with filter)](https://img.shields.io/github/v/tag/viral32111/apc-ups-exporter?label=Latest)
![GitHub repository size](https://img.shields.io/github/repo-size/viral32111/apc-ups-exporter?label=Size)
![GitHub release downloads](https://img.shields.io/github/downloads/viral32111/apc-ups-exporter/total?label=Downloads)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/viral32111/apc-ups-exporter?label=Commits)

This is a [Prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for the status data reported by [APC's Uninterruptible Power Supplies](https://www.apc.com/uk/en/).

The status data is fetched from [apcupsd](http://www.apcupsd.org/)'s *Network Information Server* (NIS), so [ensure it is enabled in your `apcupsd.conf` configuration file](http://www.apcupsd.org/manual/manual.html#configuration-directives-used-by-the-network-information-server).

I test this against my [APC Back-UPS 850VA (BE850G2-UK)](https://www.apc.com/shop/uk/en/products/APC-Back-UPS-850VA-230V-USB-Type-C-and-A-charging-ports-8-BS-1363-outlets-2-surge-/P-BE850G2-UK).

## 📥 Usage

Download the [latest release](https://github.com/viral32111/apc-ups-exporter/releases/latest) for your platform. There are builds available for Linux and Windows, on 32-bit and 64-bit architectures of x86 and ARM. There are extra Linux builds to accommodate glibc and musl libraries.

The utility does not expect any command-line arguments. There are sensible defaults in place, so it *should* run without any configuration. However, functionality can be changed using the optional command-line flags below.

* `--nis-address <string>`: The Network Information Server's IPv4 address. Defaults to `127.0.0.1`.
* `--nis-port <number>`: The Network Information Server's TCP port number. Defaults to `3551`.
* `--metrics-address <string>`: The listening IPv4 address for the Prometheus HTTP metrics server. Defaults to `127.0.0.1`.
* `--metrics-port <number>`: The listening TCP port number for the Promtheus HTTP metrics server. Defaults to `5000`.
* `--metrics-path <string>`: The HTTP path to the metrics page. Defaults to `/metrics`.
* `--metrics-interval <string>`: The number of seconds to wait between collecting metrics. Defaults to `15`.

These flags can be prefixed with either a single (`-`) or double (`--`) hyphen.

Use the `--help` (`-h`) flag for more information.

### 🐳 Docker

Alternatively, there is a [Docker image](https://github.com/users/viral32111/packages/container/package/apc-ups-exporter) available for Linux.

The image is based on [Ubuntu](https://ubuntu.com) (`ghcr.io/viral32111/apc-ups-exporter:latest-ubuntu`) weighing in at roughly 130 MiB. However, there is a variant based on [Alpine Linux](https://alpinelinux.org) (`ghcr.io/viral32111/apc-ups-exporter:latest-alpine`) which is much lighter at roughly 20 MiB.

Run the command below to download the image and create a container. Replace the `:latest` tag with your desired variant (e.g., `:1.1.3-ubuntu`, `:main-alpine`, etc.).

```bash
docker container run \
	--name apc-ups-exporter \
	--network host \
	--detach \
	ghcr.io/viral32111/apc-ups-exporter:1
```

The host's networking stack is often required to connect to the daemon's Network Information Server.

## 🖼️ Examples

Serve metrics at `/metrics` the default loopback port `5000` using data fetched every 15 seconds from the Network Information Server at `192.168.0.5` on the default port `3551`:

```
$ apc-ups-exporter --nis-address 192.168.0.5
The configured Network Information Server is: 192.168.0.5:3551.

Resetting all metrics...
Starting background metrics collection...
Serving metrics page at http://127.0.0.1:5000/metrics...

Connected to the Network Information Server.
 Fetched status from the Network Information Server.
	Updated the status metric.
	Updated the power metrics.
	Updated the battery metrics.
	Updated the daemon metrics.
 Disconnected from the Network Information Server.
 Waiting 15 seconds for next collection..
```

## 📰 Metrics

The following Prometheus metrics are exported:

### Status

* `ups_status`

### Power

* `ups_power_input_expect_voltage`
* `ups_power_output_maximum_wattage`
* `ups_power_line_voltage`
* `ups_power_load_percent`

### Battery

* `ups_battery_output_actual_voltage`
* `ups_battery_time_spent_latest_seconds`
* `ups_battery_time_spent_total_seconds`
* `ups_battery_remaining_charge_percent`
* `ups_battery_remaining_time_minutes`

### Daemon

* `ups_daemon_remaining_charge_percent`
* `ups_daemon_remaining_time_minutes`
* `ups_daemon_timeout_minutes`
* `ups_daemon_transfer_count`
* `ups_daemon_start_timestamp`

## ⚖️ License

Copyright (C) 2022-2023 [viral32111](https://viral32111.com).

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see https://www.gnu.org/licenses.
