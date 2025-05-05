package ecu

import (
	"math"
	"time"

	"github.com/StevenD2002/ninja650sim/internal/engine"
)

// ECU represents the Engine Control Unit that manages engine parameters
type ECU struct {
	// Maps
	FuelMap      FuelMap
	IgnitionMap  IgnitionMap
	TargetAFRMap AFRMap

	// ECU Settings
	IdleRPM  float64
	RevLimit float64

	// Tuning parameters
	FuelTrim     float64 // Global fuel adjustment (-100% to +100%)
	IgnitionTrim float64 // Global ignition adjustment (-10 to +10 degrees)

	// Exhaust settings
	ExhaustType string

	// Engine temperature compensation
	TempCompensation bool

	// Statistics for analysis
	KnockCount   int
	AFRDeviation float64 // How far from target AFR

	// Sensors - current values read from the engine
	ThrottlePosition float64
	RPM              float64
	EngineTemp       float64
	AirTemp          float64
	MAP              float64
	O2Reading        float64

	// Last update time for internal timing
	lastUpdateTime time.Time
}

// NewECU creates a new ECU with default maps for a Ninja 650
func NewECU() *ECU {
	return &ECU{
		FuelMap:      DefaultNinja650FuelMap(),
		IgnitionMap:  DefaultNinja650IgnitionMap(),
		TargetAFRMap: DefaultNinja650AFRMap(),

		IdleRPM:  900,
		RevLimit: 11000,

		FuelTrim:     0.0, // No adjustment
		IgnitionTrim: 0.0, // No adjustment

		ExhaustType: "Yoshimura Alpha 2",

		TempCompensation: true,

		KnockCount:   0,
		AFRDeviation: 0.0,

		ThrottlePosition: 0.0,
		RPM:              0.0,
		EngineTemp:       0.0,
		AirTemp:          0.0,
		MAP:              0.0,
		O2Reading:        0.0,

		lastUpdateTime: time.Now(),
	}
}

// UpdateSensors updates the ECU's internal sensor readings from engine state
func (e *ECU) UpdateSensors(engineState engine.SensorData) {
	e.ThrottlePosition = engineState.ThrottlePosition
	e.RPM = engineState.RPM
	e.EngineTemp = engineState.EngineTemperature
	e.AirTemp = engineState.AirTemperature
	e.MAP = engineState.MAP
	e.O2Reading = engineState.O2
}

// ProcessSensorData processes current sensor readings and returns ECU outputs
func (e *ECU) ProcessSensorData(sensors engine.SensorData) engine.ECUOutputs {
	// Update internal sensor state
	e.UpdateSensors(sensors)

	// Calculate load (simplified - using throttle position as load)
	load := e.ThrottlePosition

	// Get base values from maps
	baseFuel := e.FuelMap.GetValue(e.RPM, load)
	baseIgnition := e.IgnitionMap.GetValue(e.RPM, load)
	targetAFR := e.TargetAFRMap.GetValue(e.RPM, load)

	// Apply global trims
	fuelMultiplier := baseFuel * (1.0 + (e.FuelTrim / 100.0))
	ignitionAdjusted := baseIgnition + e.IgnitionTrim

	// Apply temperature compensation if enabled
	if e.TempCompensation {
		// Cold engine needs more fuel
		if e.EngineTemp < 80.0 {
			coldFactor := 1.0 + (80.0-e.EngineTemp)*0.01
			fuelMultiplier *= coldFactor
		}

		// Cold engine needs less timing advance
		if e.EngineTemp < 60.0 {
			coldRetard := (60.0 - e.EngineTemp) * 0.1
			ignitionAdjusted -= coldRetard
		}
	}

	// Apply modifications for aftermarket exhaust
	if e.ExhaustType != "Stock" {
		// Aftermarket exhaust generally runs leaner, so add fuel
		if e.ExhaustType == "Yoshimura Alpha 2" {
			fuelMultiplier *= 1.03 // 3% more fuel
		}
	}

	// Calculate final fuel injection time (ms)
	// This is simplified - real ECUs use complex algorithms
	// Base injection time (arbitrary value for simulation)
	baseInjectionTime := 2.5 // ms at stoichiometric mixture, 100% VE

	// Adjust for RPM and load
	// Higher RPM = less time for injection
	rpmFactor := math.Sqrt(5000.0 / math.Max(1000.0, e.RPM))

	// Calculate volumetric efficiency (simplified)
	volumetricEfficiency := calculateVolumetricEfficiency(e.RPM, load)

	// Calculate final injection time
	fuelInjectionTime := baseInjectionTime * volumetricEfficiency * fuelMultiplier * rpmFactor

	// Convert target AFR to lambda (lambda = AFR / 14.7 for gasoline)
	lambdaTarget := targetAFR / 14.7

	// Calculate AFR deviation for statistics
	currentAFR := 14.7 * e.O2Reading // Convert lambda to AFR
	e.AFRDeviation = math.Abs(currentAFR - targetAFR)

	// Check for potential knock conditions (simplified)
	knockRisk := checkKnockRisk(e.RPM, load, ignitionAdjusted, e.EngineTemp)
	if knockRisk > 0.7 {
		// Reduce timing if knock risk is high
		ignitionAdjusted -= (knockRisk - 0.7) * 10.0
		e.KnockCount++
	}

	// Create and return ECU outputs
	return engine.ECUOutputs{
		FuelInjectionTime: fuelInjectionTime,
		IgnitionAdvance:   ignitionAdjusted,
		TargetIdleRPM:     e.IdleRPM,
		LambdaTarget:      lambdaTarget,
	}
}

// calculateVolumetricEfficiency calculates how efficiently the engine breathes
// This is a simplified model - real engines have complex VE curves
func calculateVolumetricEfficiency(rpm, load float64) float64 {
	// Baseline efficiency
	baseVE := 0.85

	// VE peaks in mid-RPM range for most engines
	rpmFactor := 1.0 - 0.3*math.Pow((rpm-5500)/5500, 2)

	// VE generally increases with load (throttle opening)
	loadFactor := 0.6 + 0.4*(load/100.0)

	return baseVE * rpmFactor * loadFactor
}

// checkKnockRisk evaluates the risk of engine knock
func checkKnockRisk(rpm, load, ignitionAdvance, engineTemp float64) float64 {
	// Higher load = higher knock risk
	loadFactor := math.Pow(load/100.0, 2)

	// Higher RPM generally means less knock risk up to a point
	rpmFactor := 0.0
	if rpm < 2500 {
		rpmFactor = 0.5 - 0.5*(rpm/2500.0) // Higher risk at very low RPM
	} else if rpm > 7500 {
		rpmFactor = (rpm - 7500) / 3500 // Higher risk at very high RPM
	}

	// More ignition advance = more knock risk
	timingFactor := 0.0
	if ignitionAdvance > 30 {
		timingFactor = (ignitionAdvance - 30) / 10.0
	}

	// Higher engine temp = more knock risk
	tempFactor := 0.0
	if engineTemp > 90 {
		tempFactor = (engineTemp - 90) / 30.0
	}

	// Combine factors (weighted)
	knockRisk := 0.5*loadFactor + 0.15*rpmFactor + 0.25*timingFactor + 0.1*tempFactor

	return math.Min(1.0, math.Max(0.0, knockRisk))
}

// ResetStatistics resets the ECU statistics
func (e *ECU) ResetStatistics() {
	e.KnockCount = 0
	e.AFRDeviation = 0.0
}

// GetPerformanceStats returns information about the ECU performance
func (e *ECU) GetPerformanceStats() map[string]float64 {
	return map[string]float64{
		"KnockCount":   float64(e.KnockCount),
		"AFRDeviation": e.AFRDeviation,
	}
}
