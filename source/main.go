package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	PROJECT_NAME = "APC UPS Exporter"
	PROJECT_VERSION = "1.0.0"

	AUTHOR_NAME = "viral32111"
	AUTHOR_WEBSITE = "https://viral32111.com"
)

func main() {

	// Values of the command-line flags, and the defaults
	flagNisAddress := "127.0.0.1"
	flagNisPort := 3551
	flagMetricsAddress := "127.0.0.1"
	flagMetricsPort := 5000
	flagMetricsPath := "/metrics"
	flagMetricsInterval := 15 // Default Prometheus scrape interval

	// Setup the command-line flags
	flag.StringVar( &flagNisAddress, "nis-address", flagNisAddress, "The IPv4 address of the apcupsd Network Information Server." )
	flag.IntVar( &flagNisPort, "nis-port", flagNisPort, "The port number of the apcupsd Network Information Server." )
	flag.StringVar( &flagMetricsAddress, "metrics-address", flagMetricsAddress, "The IPv4 address to listen on for the Prometheus HTTP metrics server." )
	flag.IntVar( &flagMetricsPort, "metrics-port", flagMetricsPort, "The port number to listen on for the Prometheus HTTP metrics server." )
	flag.StringVar( &flagMetricsPath, "metrics-path", flagMetricsPath, "The full HTTP path to the metrics page." )
	flag.IntVar( &flagMetricsInterval, "metrics-interval", flagMetricsInterval, "The time in seconds to wait between collecting metrics." )

	// Set a custom help message
	flag.Usage = func() {
		fmt.Printf( "%s, v%s, by %s (%s).\n", PROJECT_NAME, PROJECT_VERSION, AUTHOR_NAME, AUTHOR_WEBSITE )

		fmt.Printf( "\nUsage: %s [-h/-help] [-nis-address <IPv4 address>] [-nis-port <number>] [-metrics-address <IPv4 address>] [-metrics-port <number>] [-metrics-path <string>] [-metrics-interval <seconds>]\n", os.Args[ 0 ] )
		flag.PrintDefaults()

		os.Exit( 1 ) // By default it exits with code 2
	}

	// Parse the command-line flags
	flag.Parse()

	// Require a valid IP address for the Network Information Server
	nisAddress := net.ParseIP( flagNisAddress )
	if ( flagNisAddress == "" || nisAddress == nil || nisAddress.To4() == nil ) { exitWithErrorMessage( "Invalid IPv4 address for apcupsd's Nnetwork Information Server." ) }

	// Require a valid port number for the Network Information Server
	if ( flagNisPort <= 0 || flagNisPort >= 65536 ) { exitWithErrorMessage( "Invalid port number for apcupsd's Network Information Server." ) }

	// Require a valid IP address for the Prometheus HTTP metrics server
	metricsAddress := net.ParseIP( flagMetricsAddress )
	if ( flagMetricsAddress == "" || metricsAddress == nil || metricsAddress.To4() == nil ) { exitWithErrorMessage( "Invalid listening IPv4 address for the Prometheus HTTP metrics server." ) }

	// Require a valid port number for the Prometheus HTTP metrics server
	if ( flagMetricsPort <= 0 || flagMetricsPort >= 65536 ) { exitWithErrorMessage( "Invalid listening port number for the Prometheus HTTP metrics server." ) }

	// Require a valid HTTP path for the metrics page
	if ( flagMetricsPath == "" || flagMetricsPath[ 0 : 1 ] != "/" || flagMetricsPath[ 1 : ] == "/" ) { exitWithErrorMessage( "Invalid path for the metrics page, must have a leading slash and no trailing slash." ) }

	// Require a valid interval for collecting metrics
	if ( flagMetricsInterval <= 0 ) { exitWithErrorMessage( "Invalid interval to wait between collecting metrics, must be greater than 0." ) }

	// Reset all metrics - probably not needed
	ResetMetrics()
	fmt.Println( "Reset all metrics to zero." )

	// Start collecting metrics in the background
	go collectMetricsInBackground( flagMetricsInterval, nisAddress, flagNisPort )

	// Serve metrics page
	fmt.Printf( "Serving metrics at http://%s:%d%s...\n\n", flagMetricsAddress, flagMetricsPort, flagMetricsPath )
	ServeMetrics( metricsAddress, flagMetricsPort, flagMetricsPath )

}

// Runs in the background to periodically collect metrics
func collectMetricsInBackground( interval int, nisAddress net.IP, nisPort int ) {
	fmt.Printf( "Starting background metrics collection...\n\n" )

	for {
		updateMetrics( nisAddress, nisPort )

		fmt.Printf( " Waiting %d seconds until next collection...\n\n", interval )
		time.Sleep( time.Duration( interval ) * time.Second )
	}
}

// Updates the metrics with the latest status from the NIS
func updateMetrics( nisAddress net.IP, nisPort int ) {

	// Create
	var networkInformationServer NetworkInformationServer

	// Connect
	connectError := networkInformationServer.Connect( nisAddress, nisPort, 5000 )
	if connectError != nil { exitWithErrorMessage( connectError.Error() ) }
	defer networkInformationServer.Disconnect()
	fmt.Println( "Connected to Network Information Server." )

	// Fetch status
	status, statusError := networkInformationServer.FetchStatus()
	if statusError != nil { exitWithErrorMessage( statusError.Error() ) }
	fmt.Println( " Fetched status from Network Information Server." )

	// Update status metric
	switch status.Status {
		case "ONLINE": metricStatus.Set( 1 )
		case "ONBATT": metricStatus.Set( 2 )
		default: metricStatus.Set( -1 )
	}

	// Update power metrics
	metricPowerInputExpectVoltage.Set( status.NormalInputVoltage )
	metricPowerOutputWattage.Set( float64( status.NormalPowerOutputWattage ) )
	metricPowerLineVoltage.Set( status.LineVoltage )
	metricPowerLoadPercent.Set( status.LoadPercent )

	// Update battery metrics
	metricBatteryExpectVoltage.Set( status.NormalBatteryVoltage )
	metricBatteryActualVoltage.Set( status.Battery.Voltage )
	metricBatteryTimeSpentLatestSeconds.Set( float64( status.TimeOnBattery ) )
	metricBatteryTimeSpentTotalSeconds.Set( float64( status.TotalTimeOnBattery ) )
	metricBatteryRemainingChargePercent.Set( status.Battery.Charge )
	metricBatteryRemainingTimeMinutes.Set( status.Battery.TimeLeft )
	
	// Update daemon metrics
	metricDaemonRemainingChargePercent.Set( float64( status.Daemon.Configuration.MinimumBatteryCharge ) )
	metricDaemonRemainingTimeMinutes.Set( float64( status.Daemon.Configuration.MinimumBatteryTimeLeft ) )
	metricDaemonTimeoutMinutes.Set( float64( status.Daemon.Configuration.MaximumTimeout ) )
	metricDaemonTransferCount.Set( float64( status.Daemon.Transfer.Count ) )
	metricDaemonStartTimestamp.Set( float64( status.Daemon.StartupTime.Unix() ) )

	fmt.Println( " Disconnected from Network Information Server." )

	
	// Daemon
	fmt.Printf( "Daemon Hostname: '%s'\n", status.Daemon.Hostname )
	fmt.Printf( "Daemon Version: '%s'\n", status.Daemon.Version )
	fmt.Printf( "Daemon Mode: '%s'\n", status.Daemon.Mode )
	fmt.Println()

	// Configuration
	fmt.Printf( "Configured Minimum Battery Charge: %d %%\n", status.Daemon.Configuration.MinimumBatteryCharge )
	fmt.Printf( "Configured Minimum Battery Time Left: %d minutes\n", status.Daemon.Configuration.MinimumBatteryTimeLeft )
	fmt.Printf( "Configured Maximum Timeout: %d minutes\n", status.Daemon.Configuration.MaximumTimeout )
	fmt.Println()

	// Information
	fmt.Printf( "Name: '%s'\n", status.Name )
	fmt.Printf( "Model: '%s'\n", status.Model )
	fmt.Printf( "Firmware Version: '%s'\n", status.Firmware )
	fmt.Printf( "Serial Number: '%s'\n", status.Serial )
	fmt.Println()

	// Management
	fmt.Printf( "Management Cable: '%s'\n", status.Management.Cable )
	fmt.Printf( "Management Driver: '%s'\n", status.Management.Driver )
	fmt.Println()

	// Status
	fmt.Printf( "Status: '%s' (%d)\n", status.Status, status.StatusFlag )
	fmt.Printf( "Startup Time: %s\n", status.Daemon.StartupTime.Format( "2006-01-02 15:04:05" ) )
	fmt.Println()

	// Load
	fmt.Printf( "Current Load: %.1f %%\n", status.LoadPercent )
	fmt.Printf( "Maximum Load: %d watts\n", status.NormalPowerOutputWattage )
	fmt.Printf( "Line Voltage: %.1f volts (Expected: %.1f volts)\n", status.LineVoltage, status.NormalInputVoltage )

	// Battery
	fmt.Printf( "Battery Charge: %.1f %%\n", status.Battery.Charge )
	fmt.Printf( "Battery Time Left: %.1f minutes\n", status.Battery.TimeLeft )
	fmt.Printf( "Battery Output: %.1f volts (Expected: %.1f volts)\n", status.Battery.Voltage, status.NormalBatteryVoltage )
	fmt.Printf( "Battery Last Replaced: %s\n", status.Battery.LastReplacementDate.Format( "2006-01-02" ) )
	fmt.Println()

	// Sensitivity
	fmt.Printf( "Sensitivity: '%s'\n", status.Sensitivity )
	fmt.Println()

	// Transfer Voltage
	fmt.Printf( "Low Transfer: %.1f volts\n", status.LowTransferVoltage )
	fmt.Printf( "High Transfer: %.1f volts\n", status.HighTransferVoltage )
	fmt.Println()

	// Alarm
	fmt.Printf( "Alarm Interval: %d seconds\n", status.AlarmDelayInterval )
	fmt.Println()

	// Transfer to battery
	fmt.Printf( "Total Transfers: %d\n", status.Daemon.Transfer.Count )
	fmt.Printf( "Last Transfer Reason: '%s'\n", status.Daemon.Transfer.LastReason )
	fmt.Println()

	// Battery time
	fmt.Printf( "Time On Battery: %d seconds\n", status.TimeOnBattery )
	fmt.Printf( "Total Time On Battery: %d seconds\n", status.TotalTimeOnBattery )
	fmt.Println()

	// Self-test
	fmt.Printf( "Last Self-Test Result: '%s'\n", status.SelfTestResult )
	

}

func exitWithErrorMessage( message string ) {
	fmt.Println( message )
	os.Exit( 1 )
}
