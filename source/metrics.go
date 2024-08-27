package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Create the metrics
var (

	// Status (as number) - STATUS
	metricStatus = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Name: "status",
		Help: "The current status.",
	} )

	// Current internal temperature (as celsius) - ITEMP - SmartUPS X 3000
	metricTemperature = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Name: "temperature_celsius",
		Help: "The current internal temperature of the UPS.",
	} )

	/*************************************/

	// Expected power input (as voltage) - NOMPOWER
	metricPowerInputExpectVoltage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "input_expect_voltage",
		Help: "The expected input voltage.",
	} )

	// Maximum power output (as wattage) - NOMPOWER
	metricPowerOutputWattage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "output_maximum_wattage",
		Help: "The maximum power the UPS can output.",
	} )

	// Current line voltage (as voltage) - LINEV
	metricPowerLineVoltage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "line_voltage",
		Help: "The current line voltage as returned by the UPS.",
	} )

	// Maximum line voltage (as voltage) - MAXLINEV - SmartUPS X 3000
	metricPowerMaximumLineVoltage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "line_maximum_voltage",
		Help: "The maximum line voltage as returned by the UPS.",
	} )

	// Minimum line voltage (as voltage) - MINLINEV - SmartUPS X 3000
	metricPowerMinimumLineVoltage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "line_minimum_voltage",
		Help: "The minimum line voltage as returned by the UPS.",
	} )

	// Current line frequency (as hertz) - LINEFREQ - SmartUPS X 3000
	metricPowerLineFrequency = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "line_frequency_hertz",
		Help: "The current line frequency as returned by the UPS.",
	} )

	// Current output voltage (as voltage) - OUTPUTV - SmartUPS X 3000
	metricPowerOutputVoltage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "output_voltage",
		Help: "The current output voltage as returned by the UPS.",
	} )

	// Current load capacity (as percentage) - LOADPCT
	metricPowerLoadPercent = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "power",
		Name: "load_percent",
		Help: "The current load capacity as estimated by the UPS, as a percentage.",
	} )

	/*************************************/

	// Expected power output of the battery (as voltage) - NOMBATTV
	metricBatteryExpectVoltage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "output_expect_voltage",
		Help: "The expected output voltage of the battery.",
	} )

	// Actual power output of the battery (as voltage) - BATTV
	metricBatteryActualVoltage = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "output_actual_voltage",
		Help: "The actual output voltage of the battery.",
	} )

	// Latest time spent on battery (in seconds) - TONBATT
	metricBatteryTimeSpentLatestSeconds = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "time_spent_latest_seconds",
		Help: "The latest time spent on battery.",
	} )

	// Total time spent on battery (in seconds) - CUMONBATT
	metricBatteryTimeSpentTotalSeconds = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "time_spent_total_seconds",
		Help: "The total time spent on battery.",
	} )

	// Remaining charge of the battery (as percentage) - BCHARGE
	metricBatteryRemainingChargePercent = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "remaining_charge_percent",
		Help: "The remaining charge on the battery, as a percentage.",
	} )

	// Remaining time of the battery (in minutes) - TIMELEFT
	metricBatteryRemainingTimeMinutes = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "remaining_time_minutes",
		Help: "The remaining runtime left on the battery as estimated by the UPS, in minutes.",
	} )

	// Low battery threshold (in minutes) - DLOWBATT - SmartUPS X 3000
	metricBatteryLowThreshold = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "low_threshold_minutes",
		Help: "The low battery threshold, in minutes.",
	} )

	// Number of external batteries - EXTBATTS - SmartUPS X 3000
	metricBatteryCount = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "battery",
		Name: "count",
		Help: "The number of external batteries in the UPS.",
	} )

	/*************************************/

	// Configured minimum battery charge (as percentage) - MBATTCHG
	metricDaemonRemainingChargePercent = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "daemon",
		Name: "remaining_charge_percent",
		Help: "The configured minimum remaining charge on the battery to trigger a system shutdown, as a percentage.",
	} )

	// Configured minimum battery remaining time (in minutes) - MINTIMEL
	metricDaemonRemainingTimeMinutes = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "daemon",
		Name: "remaining_time_minutes",
		Help: "The configured minimum remaining runtime left on the battery to trigger a system shutdown, in minutes.",
	} )

	// Configured maximum timeout (in minutes) - MAXTIME
	metricDaemonTimeoutMinutes = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "daemon",
		Name: "timeout_minutes",
		Help: "The configured maximum time running on the battery to trigger a system shutdown, in minutes.",
	} )

	// Number of transfers to battery - NUMXFERS
	metricDaemonTransferCount = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "daemon",
		Name: "transfer_count",
		Help: "The number of transfers to the battery.",
	} )

	// Daemon startup time (as unix timestamp) - STARTTIME
	metricDaemonStartTimestamp = promauto.NewGauge( prometheus.GaugeOpts {
		Namespace: "ups",
		Subsystem: "daemon",
		Name: "start_timestamp",
		Help: "The date & time the daemon was started.",
	} )

)

// Sets all of the metrics to zero
func ResetMetrics() {

	// Status
	metricStatus.Set( 0 )
	metricTemperature.Set( 0 )

	// Power
	metricPowerInputExpectVoltage.Set( 0 )
	metricPowerOutputWattage.Set( 0 )
	metricPowerLineVoltage.Set( 0 )
	metricPowerMaximumLineVoltage.Set( 0 )
	metricPowerMinimumLineVoltage.Set( 0 )
	metricPowerLineFrequency.Set( 0 )
	metricPowerOutputVoltage.Set( 0 )
	metricPowerLoadPercent.Set( 0 )

	// Battery
	metricBatteryExpectVoltage.Set( 0 )
	metricBatteryActualVoltage.Set( 0 )
	metricBatteryTimeSpentLatestSeconds.Set( 0 )
	metricBatteryTimeSpentTotalSeconds.Set( 0 )
	metricBatteryRemainingChargePercent.Set( 0 )
	metricBatteryRemainingTimeMinutes.Set( 0 )
	metricBatteryLowThreshold.Set( 0 )
	metricBatteryCount.Set( 0 )

	// Daemon
	metricDaemonRemainingChargePercent.Set( 0 )
	metricDaemonRemainingTimeMinutes.Set( 0 )
	metricDaemonTimeoutMinutes.Set( 0 )
	metricDaemonTransferCount.Set( 0 )
	metricDaemonStartTimestamp.Set( 0 )

}

// Serves the metrics page over HTTP
func ServeMetrics( address net.IP, port int, path string ) ( err error ) {

	// Handle requests to the metrics path using the Prometheus HTTP handler
	http.Handle( path, promhttp.Handler() )

	// Listen for HTTP requests
	listenError := http.ListenAndServe( fmt.Sprintf( "%s:%d" , address, port ), nil )
	if listenError != nil { return listenError }

	// No error, all was good
	return nil

}
