package main

import (
	"fmt"
	"net"
	"os"

	//"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promauto"
)

/*var batteryChargeGague = promauto.NewGauge( prometheus.GaugeOpts {
	Namespace: "ups",
	Subsystem: "battery",
	Name: "charge",
	Help: "The percentage charge on the batteries.",
} )*/

/*func doMetrics() {
	fmt.Println( "Starting metrics collection..." )

	go func() {
		for {
			exampleCounter.Inc()
			fmt.Println( "Incremented counter" )

			time.Sleep( 1 * time.Second )
		}
	}()
}*/

func main() {

	//prometheus.MustRegister( batteryChargeGague )

	// Configuration
	nisAddress := net.IPv4( 192, 168, 0, 10 )
	nisPort := 3551

	// Create
	var networkInformationServer NetworkInformationServer

	// Connect
	connectError := networkInformationServer.Connect( nisAddress, nisPort, 5000 )
	if connectError != nil { exitWithErrorMessage( connectError.Error() ) }
	defer networkInformationServer.Disconnect()
	//fmt.Println( "Connected" )

	// Status
	status, statusError := networkInformationServer.FetchStatus()
	if statusError != nil { exitWithErrorMessage( statusError.Error() ) }

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
	fmt.Printf( "Startup Time: %s\n", status.StartupTime.Format( "2006-01-02 15:04:05" ) )
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
	fmt.Printf( "Total Transfers: %d\n", status.TransferCount )
	fmt.Printf( "Last Transfer Reason: '%s'\n", status.LastTransferReason )
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
