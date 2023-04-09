# APC UPS Exporter

This is a [Prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for the status data reported by [APC's Uninterruptible Power Supplies](https://www.apc.com/uk/en/).

The status data is fetched from [apcupsd](http://www.apcupsd.org/)'s *network information server* (NIS), so [ensure it is enabled in your `apcupsd.conf` configuration file](http://www.apcupsd.org/manual/manual.html#configuration-directives-used-by-the-network-information-server).

I test this against my [APC Back-UPS 850VA (BE850G2-UK)](https://www.apc.com/shop/uk/en/products/APC-Back-UPS-850VA-230V-USB-Type-C-and-A-charging-ports-8-BS-1363-outlets-2-surge-/P-BE850G2-UK).

## Usage

### Release

Download the pre-built executable from the latest stable release for your platform:

* [Windows (64-bit)](https://github.com/viral32111/apc-ups-exporter/releases/latest/download/apc-ups-exporter-windows-amd64.exe)
* [Linux (glibc) (64-bit)](https://github.com/viral32111/apc-ups-exporter/releases/latest/download/apc-ups-exporter-linux-amd64-glibc)
* [Linux (musl) (64-bit)](https://github.com/viral32111/apc-ups-exporter/releases/latest/download/apc-ups-exporter-linux-amd64-musl)

Checksums are available on [the release page](https://github.com/viral32111/apc-ups-exporter/releases/latest).

### Docker

Alternatively, a [Docker image](https://github.com/users/viral32111/packages/container/package/apc-ups-exporter) is available for Linux.

* Ubuntu 22.10: `ghcr.io/viral32111/apc-ups-exporter:main-ubuntu`.
* Alpine Linux v3.17: `ghcr.io/viral32111/apc-ups-exporter:main-alpine`.

Run the following command to download the image and create a Docker container:

```
docker run \
    --name apc-ups-exporter \
    --network host \
    --detach \
    ghcr.io/viral32111/apc-ups-exporter:latest
```

Replace the `:latest` tag on the image with your desired version/flavour (e.g., `:1.2.0-ubuntu`, `:1.2-alpine`, etc.).

The host's network stack is usually required to connect to the network information server.

### Flags

The flags can be prefixed with either a single hyphen (`-`) or a double hyphen (`--`).

Use the `--help` flag to show usage along with a list of flags with descriptions and default values.

All of the flags are optional, with sensible default values.

* `--nis-address <string>`
  * The network information server's IPv4 address.
  * Defaults to `127.0.0.1`.
  * Example: `--nis-address 192.168.0.5`.
* `--nis-port <number>`
  * The network information server's TCP port number.
  * Defaults to `3551`.
  * Example: `--nis-port 1234`.
* `--metrics-address <string>`
  * The listening IPv4 address for the Prometheus HTTP metrics server.
  * Defaults to `127.0.0.1`.
  * Example: `--metrics-address 192.168.0.2`.
* `--metrics-port <number>`
  * The listening TCP port number for the Promtheus HTTP metrics server.
  * Defaults to `5000`.
  * Example: `--metrics-port 8080`.
* `--metrics-path <string>`
  * The HTTP path to the metrics page.
  * Defaults to `/metrics`.
  * Example: `--metrics-path /stats`.
* `--metrics-interval <string>`
  * The number of seconds to wait between collecting metrics.
  * Defaults to `15`.
  * Example: `--metrics-interval 5`.

### Examples

Fetch data from the network information server at `192.168.0.5` on port `3551`, and serve metrics at `/metrics` on loopback port `5000` every 15 seconds.

```
$ apc-ups-exporter --nis-address 192.168.0.5
The configured Network Information Server is: 192.168.0.5:3551.

Resetting all metrics...
Starting background metrics collection...
Serving metrics page at http://127.0.0.1:5000/metrics...

Connected to the network information server.
 Fetched status from the network information server.
  Updated the status metric.
  Updated the power metrics.
  Updated the battery metrics.
  Updated the daemon metrics.
 Disconnected from the network information server.
 Waiting 15 seconds for next collection..
```

## Metrics

The following Prometheus metrics are exported:

* Status
  * `ups_status`

* Power
  * `ups_power_input_expect_voltage`
  * `ups_power_output_maximum_wattage`
  * `ups_power_line_voltage`
  * `ups_power_load_percent`

* Battery
  * `ups_battery_output_actual_voltage`
  * `ups_battery_time_spent_latest_seconds`
  * `ups_battery_time_spent_total_seconds`
  * `ups_battery_remaining_charge_percent`
  * `ups_battery_remaining_time_minutes`

* Daemon
  * `ups_daemon_remaining_charge_percent`
  * `ups_daemon_remaining_time_minutes`
  * `ups_daemon_timeout_minutes`
  * `ups_daemon_transfer_count`
  * `ups_daemon_start_timestamp`

## License

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
