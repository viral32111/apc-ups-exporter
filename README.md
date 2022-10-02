# APC UPS Exporter

This is a [Prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for the status data reported by [APC's Uninterruptible Power Supplies](https://www.apc.com/uk/en/).

The status data is fetched from [apcupsd](http://www.apcupsd.org/)'s *Network Information Server*, so [ensure this is enabled in your `apcupsd.conf` configuration file](http://www.apcupsd.org/manual/manual.html#configuration-directives-used-by-the-network-information-server).

I test this against my [APC Back-UPS 850VA (BE850G2-UK)](https://www.apc.com/shop/uk/en/products/APC-Back-UPS-850VA-230V-USB-Type-C-and-A-charging-ports-8-BS-1363-outlets-2-surge-/P-BE850G2-UK).

## Usage

Download the [latest release](https://github.com/viral32111/apc-ups-exporter/releases/latest) for your platform, both Linux (glibc & musl) and Windows builds are available.

### Flags

The flags can be prefixed with either a single hyphen (`-`) or a double hyphen (`--`).

Use the `-help` flag to show usage, and a list of these flags with descriptions and default values.

All of the flags are optional, with sensible default values.

* `-nis-address <string>` to specify the Network Information Server's IPv4 address (e.g. `-nis-address 192.168.0.5`, defaults to `127.0.0.1`).
* `-nis-port <number>` to specify the Network Information Server's port number (e.g. `-nis-port 1234`, defaults to `3551`).
* `-metrics-address <string>` to specify the listening IPv4 address for the Prometheus HTTP metrics server (e.g. `-metrics-address 192.168.0.2`, defaults to `127.0.0.1`).
* `-metrics-port <number>` to specify the listening port number for the Promtheus HTTP metrics server (e.g. `-metrics-port 8080`, defaults to `5000`).
* `-metrics-path <string>` to specify the HTTP path to the metrics page (e.g. `-metrics-path /stats`, defaults to `/metrics`).
* `-metrics-interval <string>` to specify the time to wait (in seconds) between collecting metrics (e.g. `-metrics-interval 5`, defaults to `15`).

### Examples

Fetch status data from the Network Information Server at `192.168.0.5` on port `3551`, and serve metrics at `/metrics` on loopback port `5000` every 15 seconds.

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

### Metrics

The following metrics are exported:

* `ups_status`

* `ups_power_input_expect_voltage`
* `ups_power_output_maximum_wattage`
* `ups_power_line_voltage`
* `ups_power_load_percent`

* `ups_battery_output_actual_voltage`
* `ups_battery_time_spent_latest_seconds`
* `ups_battery_time_spent_total_seconds`
* `ups_battery_remaining_charge_percent`
* `ups_battery_remaining_time_minutes`

* `ups_daemon_remaining_charge_percent`
* `ups_daemon_remaining_time_minutes`
* `ups_daemon_timeout_minutes`
* `ups_daemon_transfer_count`
* `ups_daemon_start_timestamp`

## License

Copyright (C) 2022 [viral32111](https://viral32111.com).

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
