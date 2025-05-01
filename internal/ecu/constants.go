package ecu

// Common fuel and ignition constants
const (
	// Fueling constants
	StoichiometricAFR = 14.7 // For gasoline
	MinimumAFR        = 11.0 // Very rich, risk of fouling
	MaximumAFR        = 16.0 // Very lean, risk of overheating

	// Ignition constants
	MinimumAdvance = 5.0  // Minimum ignition advance (degrees BTDC)
	MaximumAdvance = 45.0 // Maximum safe ignition advance (degrees BTDC)

	// Engine constants
	MinimumRPM = 800.0   // Minimum stable idle RPM
	MaximumRPM = 12000.0 // Absolute maximum RPM (hardware limit)

	// Temperature ranges
	MinEngineTemp = 20.0  // Minimum expected engine temperature (°C)
	OptEngineTemp = 90.0  // Optimal engine temperature (°C)
	MaxEngineTemp = 120.0 // Maximum safe engine temperature (°C)

	// Knock threshold
	KnockThreshold = 0.7 // Threshold for knock detection (0-1)
)

// TuningPresets defines common tuning configurations
var TuningPresets = map[string]TuningPreset{
	"Stock": {
		Name:         "Stock",
		Description:  "Factory stock tuning",
		FuelTrim:     0.0,
		IgnitionTrim: 0.0,
		Notes:        "Factory settings, safe for all conditions",
	},
	"Performance": {
		Name:         "Performance",
		Description:  "Increased power with premium fuel",
		FuelTrim:     2.0, // 2% more fuel
		IgnitionTrim: 2.0, // 2 degrees more advance
		Notes:        "Requires 91+ octane fuel, provides better throttle response",
	},
	"Economy": {
		Name:         "Economy",
		Description:  "Optimized for fuel economy",
		FuelTrim:     -3.0, // 3% less fuel
		IgnitionTrim: 1.0,  // 1 degree more advance
		Notes:        "May reduce power slightly but improves fuel economy",
	},
	"Yoshimura Exhaust": {
		Name:         "Yoshimura Alpha 2",
		Description:  "Tuned for Yoshimura Alpha 2 exhaust",
		FuelTrim:     3.0, // 3% more fuel
		IgnitionTrim: 1.5, // 1.5 degrees more advance
		Notes:        "Compensates for increased flow, optimal power with aftermarket exhaust",
	},
}

// TuningPreset represents a preset ECU configuration
type TuningPreset struct {
	Name         string
	Description  string
	FuelTrim     float64
	IgnitionTrim float64
	Notes        string
}

// Sensor limits define the normal operating ranges for sensors
var SensorLimits = map[string]SensorRange{
	"RPM": {
		Min:       0.0,
		Max:       12000.0,
		LowWarn:   800.0,
		HighWarn:  11000.0,
		Unit:      "RPM",
		Precision: 0,
	},
	"ThrottlePosition": {
		Min:       0.0,
		Max:       100.0,
		LowWarn:   0.0,
		HighWarn:  100.0,
		Unit:      "%",
		Precision: 1,
	},
	"EngineTemp": {
		Min:       0.0,
		Max:       150.0,
		LowWarn:   60.0,
		HighWarn:  110.0,
		Unit:      "°C",
		Precision: 1,
	},
	"MAP": {
		Min:       0.0,
		Max:       120.0,
		LowWarn:   20.0,
		HighWarn:  110.0,
		Unit:      "kPa",
		Precision: 1,
	},
	"O2": {
		Min:       0.7,
		Max:       1.3,
		LowWarn:   0.8,
		HighWarn:  1.1,
		Unit:      "λ",
		Precision: 2,
	},
}

// SensorRange defines the operating range for a sensor
type SensorRange struct {
	Min       float64
	Max       float64
	LowWarn   float64
	HighWarn  float64
	Unit      string
	Precision int
}
