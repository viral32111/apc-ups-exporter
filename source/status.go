package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Structure to hold the status response - apcupsd.org/manual/manual.html#status-report-fields
// NOTE: Integer values are stored as floats because Prometheus requires floats
type Status struct {

	// When information was last obtained from the UPS
	Date time.Time // DATE

	// Data reported by the UPS
	UPS struct {

		// Name (this can be from EEPROM or configuration file)
		Name string // UPSNAME

		// Status of the UPS
		StatusText string // STATUS
		StatusFlag int64 // STATFLAG

		// Information about the UPS
		ModelName string // MODEL
		FirmwareRevision string // FIRMWARE
		SerialNumber string // SERIALNO
		ManufacturedAt time.Time // MANDATE - SmartUPS X 3000

		// Load
		LoadPercent float64 // LOADPCT

		// Line voltage
		LineVoltage float64 // LINEV
		LineVoltageFluctuationSensitivity string // SENSE
		MaximumLineVoltage float64 // MAXLINEV - SmartUPS X 3000
		MinimumLineVoltage float64 // MINLINEV - SmartUPS X 3000
		OutputVoltage float64 // OUTPUTV - SmartUPS X 3000
		LineFrequency float64 // LINEFREQ - SmartUPS X 3000

		// Delay between alarm beeps
		AlarmIntervalSeconds float64 // ALARMDEL

		// Results of the last self-test
		SelfTestResult string // SELFTEST
		SelfTestInterval float64 // STESTI - SmartUPS X 3000

		// Internal temperature (in Celsius)
		Temperature float64 // ITEMP - SmartUPS X 3000

		// Data about the battery
		Battery struct {
			ChargePercent float64 // BCHARGE
			RemainingRuntimeMinutes float64 // TIMELEFT
			OutputVoltage float64 // BATTV
			LastReplacementDate time.Time // BATTDATE
			LowBatterySignalThreshold float64 // DLOWBATT - SmartUPS X 3000
			ExternalCount float64 // EXTBATTS - SmartUPS X 3000
		}

		// Expected power values
		Expect struct {
			MainsInputVoltage float64 // NOMINV
			BatteryOutputVoltage float64 // NOMBATTV
			PowerOutputWattage float64 // NOMPOWER
		}

	}

	// Data reported by the daemon
	Daemon struct {

		// Hostname of the system running the daemon
		SystemName string // HOSTNAME

		// Version of the daemon
		Version string // VERSION

		// When the daemon was started
		StartupTime time.Time // STARTTIME

		// Communication driver in use
		Driver string // DRIVER

		// Values from the configuration file
		Configuration struct {

			// Type of cable
			ManagementCable string // CABLE

			// Thresholds
			MinimumBatteryChargePercent float64 // MBATTCHG
			MinimumBatteryRemainingRuntimeMinutes float64 // MINTIMEL

			// Timeout
			MaximumTimeoutMinutes float64 // MAXTIME

			// UPS operating mode
			OperatingMode string // UPSMODE

		}

		// About the battery
		Battery struct {

			// Transfer to battery
			Transfer struct {

				// Total number of transfers
				Total float64 // NUMXFERS

				// Reason for the last transfer
				LastReason string // LASTXFER
				LastAt time.Time // XOFFBATT - SmartUPS X 3000

				// Line voltage below & above to trigger a transfer to battery
				LowLineVoltage float64 // LOTRANS
				HighLineVoltage float64 // HITRANS

			}

			// Time spent on battery
			TimeSpent struct {
				Current float64 // TONBATT
				Total float64 // CUMONBATT
			}

		}

	}

}

// Parses the status response into a structure
func ParseStatusText( text string ) ( status Status, err error ) {

	// Split the response into lines
	lines := strings.Split( text, "\n" )

	// Loop through all the lines...
	for _, line := range lines {

		// Skip lines that are empty
		line = strings.TrimSpace( line )
		if line == "" { continue }

		// Parse the line into key & value
		key, value, parseError := ParseLine( line )
		if parseError != nil { return Status{}, parseError }

		// Remove any trailing labels from the value
		value = strings.TrimSuffix( value, " Volts" ) // e.g., LINEV
		value = strings.TrimSuffix( value, " Seconds" ) // e.g., MAXTIME
		value = strings.TrimSuffix( value, " Minutes" ) // e.g., TIMELEFT
		value = strings.TrimSuffix( value, " Percent" ) // e.g., LOADPCT
		value = strings.TrimSuffix( value, " Watts" )
		value = strings.TrimSuffix( value, " Hz" ) // e.g., LINEFREQ
		value = strings.TrimSuffix( value, " C" ) // e.g., ITEMP

		// Assign the value to the correct property in the structure
		switch key {

			// "The date and time that the information was last obtained from the UPS"
			case "DATE": {
				parsedDate, dateParseError := time.Parse( "2006-01-02 15:04:05 -0700", value )
				if dateParseError != nil { return Status{}, dateParseError }

				status.Date = parsedDate
			}

			// "The name of the machine that collected the UPS data"
			case "HOSTNAME": status.Daemon.SystemName = value

			// "The apcupsd release number, build date, and platform"
			case "VERSION": status.Daemon.Version = value

			// "The name of the UPS as stored in the EEPROM or in the UPSNAME directive in the configuration file"
			case "UPSNAME": status.UPS.Name = value

			// "The cable as specified in the configuration file (UPSCABLE)"
			case "CABLE": status.Daemon.Configuration.ManagementCable = value
			case "DRIVER": status.Daemon.Driver = value

			// "The mode in which apcupsd is operating as specified in the configuration file (UPSMODE)"
			case "UPSMODE": status.Daemon.Configuration.OperatingMode = value

			// "The time/date that apcupsd was started"
			case "STARTTIME": {
				parsedDate, dateParseError := time.Parse( "2006-01-02 15:04:05 -0700", value )
				if dateParseError != nil { return Status{}, dateParseError }

				status.Daemon.StartupTime = parsedDate
			}

			// "The UPS model as derived from information from the UPS"
			case "MODEL": status.UPS.ModelName = value

			// "The current status of the UPS (ONLINE, ONBATT, etc.)"
			case "STATUS": status.UPS.StatusText = value

			// "The current line voltage as returned by the UPS"
			case "LINEV": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.LineVoltage = parsedFloat
			}

			// "The percentage of load capacity as estimated by the UPS"
			case "LOADPCT": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.LoadPercent = parsedFloat
			}

			// "The percentage charge on the batteries"
			case "BCHARGE": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Battery.ChargePercent = parsedFloat
			}

			// "The remaining runtime left on batteries as estimated by the UPS"
			case "TIMELEFT": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Battery.RemainingRuntimeMinutes = parsedFloat
			}

			// "If the battery charge percentage (BCHARGE) drops below this value, apcupsd will shutdown your system. Value is set in the configuration file (BATTERYLEVEL)"
			case "MBATTCHG": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Configuration.MinimumBatteryChargePercent = parsedFloat
			}

			// "apcupsd will shutdown your system if the remaining runtime equals or is below this point. Value is set in the configuration file (MINUTES)"
			case "MINTIMEL": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Configuration.MinimumBatteryRemainingRuntimeMinutes = parsedFloat
			}

			// "apcupsd will shutdown your system if the time on batteries exceeds this value. A value of zero disables the feature. Value is set in the configuration file (TIMEOUT)"
			case "MAXTIME": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Configuration.MaximumTimeoutMinutes = parsedFloat
			}

			// SmartUPS X 3000 - "The maximum line voltage since the last STATUS as returned by the UPS."
			case "MAXLINEV": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.MaximumLineVoltage = parsedFloat
			}

			// SmartUPS X 3000 - "The minimum line voltage since the last STATUS as returned by the UPS."
			case "MINLINEV": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.MinimumLineVoltage = parsedFloat
			}

			// SmartUPS X 3000 - "The voltage the UPS is supplying to your equipment."
			case "OUTPUTV": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.OutputVoltage = parsedFloat
			}

			// "The sensitivity level of the UPS to line voltage fluctuations"
			case "SENSE": status.UPS.LineVoltageFluctuationSensitivity = value

			// SmartUPS X 3000 - "The remaining runtime below which the UPS sends the low battery signal. At this point apcupsd will force an immediate emergency shutdown. "
			case "DLOWBATT": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Battery.LowBatterySignalThreshold = parsedFloat
			}

			// "The line voltage below which the UPS will switch to batteries"
			case "LOTRANS": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Battery.Transfer.LowLineVoltage = parsedFloat
			}

			// "The line voltage above which the UPS will switch to batteries"
			case "HITRANS": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Battery.Transfer.HighLineVoltage = parsedFloat
			}

			// SmartUPS X 3000 - "The internal UPS temperature as supplied by the UPS."
			case "ITEMP": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Temperature = parsedFloat
			}

			// "The delay period for the UPS alarm"
			case "ALARMDEL": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.AlarmIntervalSeconds = parsedFloat
			}

			// "Battery voltage as supplied by the UPS"
			case "BATTV": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Battery.OutputVoltage = parsedFloat
			}

			// SmartUPS X 3000 - "The line frequency in Hertz as given by the UPS."
			case "LINEFREQ": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.LineFrequency = parsedFloat
			}

			// "The reason for the last transfer to batteries"
			case "LASTXFER": status.Daemon.Battery.Transfer.LastReason = value

			// "The number of transfers to batteries since apcupsd startup"
			case "NUMXFERS": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Battery.Transfer.Total = parsedFloat
			}

			// "Time in seconds currently on batteries, or 0"
			case "TONBATT": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Battery.TimeSpent.Current = parsedFloat
			}

			// "Total (cumulative) time on batteries in seconds since apcupsd startup"
			case "CUMONBATT": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.Daemon.Battery.TimeSpent.Total = parsedFloat
			}

			// SmartUPS X 3000 - "Time and date of last transfer from batteries, or N/A."
			case "XOFFBATT": {
				if (value == "N/A") {
					status.Daemon.Battery.Transfer.LastAt = time.Unix(0, 0)
				} else {
					parsedDate, dateParseError := time.Parse( "2006-01-02 15:04:05 -0700", value )
					if dateParseError != nil { return Status{}, dateParseError }

					status.Daemon.Battery.Transfer.LastAt = parsedDate
				}
			}

			// "The results of the last self test"
			case "SELFTEST": status.UPS.SelfTestResult = value

			// SmartUPS X 3000 - "The interval in hours between automatic self tests."
			case "STESTI": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.SelfTestInterval = parsedFloat
			}

			// "Status flag. English version is given by STATUS"
			case "STATFLAG": {
				parsedInt, intParseError := strconv.ParseInt( strings.Replace( value, "0x", "", 1 ), 16, 64 )
				if intParseError != nil { return Status{}, intParseError }

				status.UPS.StatusFlag = parsedInt
			}

			// SmartUPS X 3000 - "The date the UPS was manufactured."
			case "MANDATE": {
				parsedDate, dateParseError := time.Parse( "2006-01-02", value )
				if dateParseError != nil {
					parsedDate, dateParseError = time.Parse( "01/02/2006", value ) // SmartUPS X 3000 reports date in MM/DD/YYYY format - https://github.com/viral32111/apc-ups-exporter/issues/30
          if dateParseError != nil {
            parsedDate, dateParseError = time.Parse( "01/02/06", value ) //SmartUPS 3000 reports date in MM/DD/YY format (see above link)
            if dateParseError != nil { return Status{}, dateParseError }
          }
				}

				status.UPS.ManufacturedAt = parsedDate
			}

			// "The UPS serial number"
			case "SERIALNO": status.UPS.SerialNumber = value

			// "The date that batteries were last replaced"
			case "BATTDATE": {
				parsedDate, dateParseError := time.Parse( "2006-01-02", value )
				if dateParseError != nil {
					parsedDate, dateParseError = time.Parse( "01/02/2006", value ) // SmartUPS X 3000 reports date in MM/DD/YYYY format - https://github.com/viral32111/apc-ups-exporter/issues/30
					if dateParseError != nil { return Status{}, dateParseError }
				}

				status.UPS.Battery.LastReplacementDate = parsedDate
			}

			// "The input voltage that the UPS is configured to expect"
			case "NOMINV": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Expect.MainsInputVoltage = parsedFloat
			}

			// "The nominal battery voltage"
			case "NOMBATTV": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Expect.BatteryOutputVoltage = parsedFloat
			}

			// SmartUPS X 3000 - "The number of external batteries as defined by the user. A correct number here helps the UPS compute the remaining runtime more accurately.""
			case "EXTBATTS": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Battery.ExternalCount = parsedFloat
			}

			// "The maximum power in Watts that the UPS is designed to supply"
			case "NOMPOWER": {
				parsedFloat, floatParseError := parseAsFloat( value, -1 )
				if floatParseError != nil { return Status{}, floatParseError }

				status.UPS.Expect.PowerOutputWattage = parsedFloat
			}

			// "The firmware revision number as reported by the UPS"
			case "FIRMWARE": status.UPS.FirmwareRevision = value

		}

	}

	// Return the populated structure
	return status, nil

}

// Parses a line from the status response
func ParseLine( line string ) ( key string, value string, err error ) {

	// Split the line into key & value
	pair := strings.SplitN( line, ":", 2 )
	if len( pair ) != 2 { return "", "", errors.New( "line does not contain separator character" ) }

	// Trim spaces from the key & value
	key = strings.TrimSpace( pair[ 0 ] )
	value = strings.TrimSpace( pair[ 1 ] )

	// Return the key & value
	return key, value, nil

}
