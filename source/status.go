package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// http://www.apcupsd.org/manual/manual.html#status-report-fields
type Status struct {

	Date time.Time // DATE

	Daemon struct {
		Hostname string // HOSTNAME
		Version string // VERSION
		Mode string // UPSMODE

		Configuration struct {
			MinimumBatteryCharge int64 // MBATTCHG
			MinimumBatteryTimeLeft int64 // MINTIMEL
			MaximumTimeout int64 // MAXTIME
		}
	}

	Name string // UPSNAME
	Model string // MODEL
	Firmware string // FIRMWARE
	Serial string // SERIALNO

	Management struct {
		Cable string // CABLE
		Driver string // DRIVER
	}

	Status string // STATUS
	StartupTime time.Time // STARTTIME

	LoadPercent float64 // LOADPCT
	LineVoltage float64 // LINEV

	Battery struct {
		Charge float64 // BCHARGE
		TimeLeft float64 // TIMELEFT
		Voltage float64 // BATTV
		LastReplacementDate time.Time // BATTDATE
	}

	Sensitivity string // SENSE

	LowTransferVoltage float64 // LOTRANS
	HighTransferVoltage float64 // HITRANS

	AlarmDelayInterval int64 // ALARMDEL

	LastTransferReason string // LASTXFER
	TransferCount int64 // NUMXFERS

	TimeOnBattery time.Duration // TONBATT
	TotalTimeOnBattery time.Duration // CUMONBATT

	SelfTestResult string // SELFTEST

	StatusFlag int64 // STATFLAG

	NormalInputVoltage float64 // NOMINV
	NormalBatteryVoltage float64 // NOMBATTV
	NormalPowerOutputWattage int64 // NOMPOWER

}

func ParseStatusText( text string ) ( status Status, err error ) {
	lines := strings.Split( text, "\n" )

	for _, line := range lines {
		if line == "" { continue }

		key, value, parseError := ParseLine( line )
		if parseError != nil { return status, parseError }

		// Remove labels
		value = strings.TrimSuffix( value, " Volts" )
		value = strings.TrimSuffix( value, " Seconds" )
		value = strings.TrimSuffix( value, " Minutes" )
		value = strings.TrimSuffix( value, " Percent" )
		value = strings.TrimSuffix( value, " Watts" )

		//fmt.Printf( "'%s' = '%s'\n", key, value )

		switch key {
			// APC

			case "DATE": status.Date, _ = time.Parse( "2006-01-02 15:04:05 -0700", value )

			case "HOSTNAME": status.Daemon.Hostname = value
			case "VERSION": status.Daemon.Version = value

			case "UPSNAME": status.Name = value

			case "CABLE": status.Management.Cable = value
			case "DRIVER": status.Management.Driver = value

			case "UPSMODE": status.Daemon.Mode = value

			case "STARTTIME": status.StartupTime, _ = time.Parse( "2006-01-02 15:04:05 -0700", value )

			case "MODEL": status.Model = value

			case "STATUS": status.Status = value

			case "LINEV": status.LineVoltage, _ = strconv.ParseFloat( value, 64 )
			case "LOADPCT": status.LoadPercent, _ = strconv.ParseFloat( value, 64 )

			case "BCHARGE": status.Battery.Charge, _ = strconv.ParseFloat( value, 64 )
			case "TIMELEFT": status.Battery.TimeLeft, _ = strconv.ParseFloat( value, 64 )

			case "MBATTCHG": status.Daemon.Configuration.MinimumBatteryCharge, _ = strconv.ParseInt( value, 10, 64 )
			case "MINTIMEL": status.Daemon.Configuration.MinimumBatteryTimeLeft, _ = strconv.ParseInt( value, 10, 64 )
			case "MAXTIME": status.Daemon.Configuration.MaximumTimeout, _ = strconv.ParseInt( value, 10, 64 )

			case "SENSE": status.Sensitivity = value

			case "LOTRANS": status.LowTransferVoltage, _ = strconv.ParseFloat( value, 64 )
			case "HITRANS": status.HighTransferVoltage, _ = strconv.ParseFloat( value, 64 )

			case "ALARMDEL": status.AlarmDelayInterval, _ = strconv.ParseInt( value, 10, 64 )

			case "LASTXFER": status.LastTransferReason = value
			case "NUMXFERS": status.TransferCount, _ = strconv.ParseInt( value, 10, 64 )

			case "TONBATT": status.TimeOnBattery, _ = time.ParseDuration( value + "s" )
			case "CUMONBATT": status.TotalTimeOnBattery, _ = time.ParseDuration( value + "s" )

			case "SELFTEST": status.SelfTestResult = value

			case "STATFLAG": status.StatusFlag, _ = strconv.ParseInt( value, 16, 64 )

			case "SERIALNO": status.Serial = value

			case "BATTDATE": status.Battery.LastReplacementDate, _ = time.Parse( "2006-01-02", value )

			case "NOMINV": status.NormalInputVoltage, _ = strconv.ParseFloat( value, 64 )
			case "NOMBATTV": status.NormalBatteryVoltage, _ = strconv.ParseFloat( value, 64 )
			case "NOMPOWER": status.NormalPowerOutputWattage, _ = strconv.ParseInt( value, 10, 64 )

			case "FIRMWARE": status.Firmware = value

			// END APC
		}
	}

	return status, nil
}

func ParseLine( line string ) ( key string, value string, err error ) {

	// Split
	pair := strings.SplitN( line, ":", 2 )
	if len( pair ) != 2 { return "", "", errors.New( "line does not contain separator character" ) }

	// Trim
	key = strings.TrimSpace( pair[ 0 ] )
	value = strings.TrimSpace( pair[ 1 ] )

	return key, value, nil
}
