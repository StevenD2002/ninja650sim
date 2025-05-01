package engine

import (
	"math"
	"math/rand/v2"
	"time"
)

// SensorData represents data from engine sensors
type SensorData struct {
	RPM               float64
	ThrottlePosition  float64
	AirTemperature    float64
	EngineTemperature float64
	MAP               float64 // Manifold Absolute Pressure
	O2                float64 // O2 sensor reading (lambda)
	Speed             float64
	Timestamp         int64
}

// ECUOutputs represents commands from the ECU
type ECUOutputs struct {
	FuelInjectionTime float64 // milliseconds
	IgnitionAdvance   float64 // degrees BTDC
	TargetIdleRPM     float64 // RPM
	LambdaTarget      float64 // Target air/fuel ratio
}

// Engine model representing a Ninja 650 motorcycle
type Engine struct {
	// Basic engine specifications
	Displacement     float64 // cc
	CompressionRatio float64
	MaxRPM           float64
	IdleRPM          float64
	MaxTorque        float64 // Nm
	MaxTorqueRPM     float64
	RedlineRPM       float64

	// Current state
	RPM              float64
	ThrottlePosition float64
	ClutchPosition   float64 // 0.0 = fully engaged, 1.0 = fully disengaged
	EngineTemp       float64 // Celsius
	AirTemp          float64 // Celsius
	MAP              float64 // kPa
	O2Reading        float64 // Lambda value
	Speed            float64 // km/h
	Gear             int     // 0 = neutral, 1-6 = gears
	BrakeApplied     bool

	// Environment settings
	AmbientTemp float64 // Celsius
	Altitude    float64 // meters

	// Configuration
	ExhaustType string // "Stock" or "Yoshimura Alpha 2", etc.

	// Physics parameters
	Responsiveness    float64 // RPM rise rate
	Resistance        float64 // RPM fall rate
	RotationalInertia float64 // kg*m²
	DragCoefficient   float64

	// Internal tracking
	lastUpdateTime time.Time
}

// NewEngine creates a new engine model with default Ninja 650 parameters
func NewEngine() *Engine {
	return &Engine{
		// Basic specifications
		Displacement:     649,   // cc
		CompressionRatio: 10.8,  // Compression ratio
		MaxRPM:           11000, // Max RPM
		IdleRPM:          900,   // Idle RPM
		MaxTorque:        65.7,  // Nm
		MaxTorqueRPM:     6500,  // RPM at max torque
		RedlineRPM:       10500, // Redline

		// Initial state
		RPM:              900,   // Starting at idle
		ThrottlePosition: 0,     // Closed throttle
		ClutchPosition:   1.0,   // Fully disengaged
		EngineTemp:       90,    // Normal operating temp (C)
		AirTemp:          25,    // Ambient temperature (C)
		MAP:              101.3, // kPa (atmospheric at sea level)
		O2Reading:        1.0,   // Lambda = 1.0 (stoichiometric)
		Speed:            0,     // Not moving
		Gear:             0,     // Neutral
		BrakeApplied:     false, // No brakes

		// Environment
		AmbientTemp: 25, // Celsius
		Altitude:    0,  // Meters above sea level

		// Configuration
		ExhaustType: "Yoshimura Alpha 2",

		// Physics parameters
		Responsiveness:    500,  // RPM increase per second at 100% throttle
		Resistance:        200,  // RPM decrease per second at 0% throttle
		RotationalInertia: 0.12, // kg*m² (estimated)
		DragCoefficient:   0.35, // Aerodynamic drag coefficient

		// Initialize the last update time
		lastUpdateTime: time.Now(),
	}
}

// Update engine state based on ECU outputs
func (e *Engine) Update(ecuOutputs ECUOutputs, deltaTime float64) {
	// Calculate forces

	// 1. Engine torque based on throttle, RPM and fuel/ignition settings
	baselineTorque := e.calculateBaselineTorque()
	torqueMultiplier := 1.0

	// Apply ECU fuel enrichment/leaning
	if ecuOutputs.FuelInjectionTime > 0 {
		// Compare to expected stock value
		stockInjectionTime := e.calculateStockInjectionTime()
		torqueMultiplier *= ecuOutputs.FuelInjectionTime / stockInjectionTime

		// Too rich or too lean reduces power
		if torqueMultiplier > 1.3 || torqueMultiplier < 0.8 {
			torqueMultiplier = math.Max(0.5, math.Min(1.1, torqueMultiplier))
		}
	}

	// Apply ignition timing effects
	optimalTiming := e.calculateOptimalTiming()
	timingDifference := ecuOutputs.IgnitionAdvance - optimalTiming

	if timingDifference > 0 {
		// Advanced timing up to a point increases power
		torqueMultiplier *= math.Min(1.1, 1.0+timingDifference*0.01)
	} else {
		// Retarded timing reduces power
		torqueMultiplier *= math.Max(0.7, 1.0+timingDifference*0.03)
	}

	// Apply throttle position
	torqueMultiplier *= e.ThrottlePosition / 100.0

	// Final torque
	engineTorque := baselineTorque * torqueMultiplier

	// 2. Calculate acceleration
	if e.ClutchPosition < 0.5 && e.Gear > 0 {
		// Power is being transmitted to wheels
		// (Simplified - would need gearing ratios, final drive ratio, etc.)
		gearRatio := e.getGearRatio(e.Gear)
		wheelTorque := engineTorque * gearRatio

		// Apply wheel torque to change speed
		// (Very simplified physics)
		acceleration := wheelTorque / 250.0 // Arbitrary mass factor

		// Apply braking
		if e.BrakeApplied {
			acceleration -= 5.0 // Hard braking
		}

		// Update speed
		e.Speed += acceleration * deltaTime
		e.Speed = math.Max(0, e.Speed)

		// Calculate RPM from speed if clutch engaged
		if e.Speed > 0 {
			e.RPM = (e.Speed * 1000 / 3600) * 60 * gearRatio * 100
			e.RPM = math.Min(e.RPM, e.RedlineRPM)
		}
	} else {
		// Clutch disengaged or neutral - engine can rev freely
		rpmChange := (engineTorque - (e.RPM * 0.001)) * 10 / e.RotationalInertia
		e.RPM += rpmChange * deltaTime

		// Apply engine braking and idle control
		if e.ThrottlePosition < 5 {
			targetIdle := ecuOutputs.TargetIdleRPM
			if targetIdle < 800 {
				targetIdle = 800 // Fallback idle
			}

			// Move toward idle RPM
			if e.RPM > targetIdle {
				e.RPM -= math.Min(500*deltaTime, e.RPM-targetIdle)
			} else if e.RPM < targetIdle {
				e.RPM += math.Min(300*deltaTime, targetIdle-e.RPM)
			}
		}

		// Apply speed changes if clutch disengaged
		if e.ClutchPosition > 0.5 && e.Gear > 0 {
			brakeValue := 0.5
			if e.BrakeApplied {
				brakeValue = 5.0
			}
			e.Speed -= math.Min(e.Speed, (e.DragCoefficient*e.Speed*e.Speed*0.001+brakeValue)*deltaTime)
		}
	}

	// Ensure RPM stays in valid range
	e.RPM = math.Max(0, math.Min(e.RPM, e.RedlineRPM*1.05)) // Allow slight overrev

	// 3. Update sensor readings
	e.updateSensorReadings(ecuOutputs, deltaTime)

	// Update last time
	e.lastUpdateTime = time.Now()
}

// Helper methods for the engine model
func (e *Engine) calculateBaselineTorque() float64 {
	// Simplified torque curve based on Ninja 650 characteristics
	// This would ideally be based on dyno data

	// Normalize RPM percentage (0.0 to 1.0)
	rpmPercent := e.RPM / e.MaxRPM

	// Simple torque curve with peak at MaxTorqueRPM
	peakPoint := e.MaxTorqueRPM / e.MaxRPM

	if rpmPercent < 0.1 {
		// Low RPM (idle range)
		return e.MaxTorque * 0.2
	} else if rpmPercent < peakPoint {
		// Building torque
		return e.MaxTorque * (0.4 + 0.6*(rpmPercent-0.1)/(peakPoint-0.1))
	} else {
		// After peak, gradually dropping
		torqueDropoff := (rpmPercent - peakPoint) / (1.0 - peakPoint)
		return e.MaxTorque * (1.0 - torqueDropoff*0.3)
	}
}

// Get current sensor data
func (e *Engine) GetSensorData() SensorData {
	return SensorData{
		RPM:               e.RPM,
		ThrottlePosition:  e.ThrottlePosition,
		AirTemperature:    e.AirTemp,
		EngineTemperature: e.EngineTemp,
		MAP:               e.MAP,
		O2:                e.O2Reading,
		Speed:             e.Speed,
		Timestamp:         time.Now().UnixNano(),
	}
}

// Calculate output metrics
func (e *Engine) CalculatePerformance() (power, torque float64) {
	// Get baseline torque based on RPM curve
	torque = e.calculateBaselineTorque()

	// Apply throttle position
	torque *= (e.ThrottlePosition / 100.0)

	// Apply exhaust effects using the physics model
	torqueMultiplier, powerMultiplier := CalculateExhaustEffect(e.RPM, e.ExhaustType)
	torque *= torqueMultiplier

	// Calculate power
	// Base power calculation: Power = Torque * RPM / 5252 (in HP)
	power = torque * e.RPM / 5252.0

	// Apply any additional power multiplier from the exhaust
	power *= powerMultiplier

	return power, torque
}

// SetThrottle sets the throttle position
func (e *Engine) SetThrottle(position float64) {
	e.ThrottlePosition = math.Max(0, math.Min(100, position))
}

// GetRPM returns the current RPM
func (e *Engine) GetRPM() float64 {
	return e.RPM
}

// GetThrottlePosition returns the current throttle position
func (e *Engine) GetThrottlePosition() float64 {
	return e.ThrottlePosition
}

// Additional helper methods needed by the engine model

// calculateStockInjectionTime calculates the expected stock fuel injection time
func (e *Engine) calculateStockInjectionTime() float64 {
	// This is a simplified model - in reality this would be based on complex maps
	// Base injection time increases with RPM and throttle
	baseTime := 2.0 // Base milliseconds at idle

	// Adjust for RPM (simplified linear relationship)
	rpmFactor := 1.0 + 0.5*(e.RPM-e.IdleRPM)/(e.MaxRPM-e.IdleRPM)

	// Adjust for throttle position
	throttleFactor := 0.5 + 0.5*(e.ThrottlePosition/100.0)

	return baseTime * rpmFactor * throttleFactor
}

// calculateOptimalTiming calculates the optimal ignition timing
func (e *Engine) calculateOptimalTiming() float64 {
	// This is a simplified model - in reality would be based on complex ignition maps

	// Base timing (degrees BTDC)
	baseTiming := 10.0

	// More advance at higher RPM
	rpmAdvance := 20.0 * (e.RPM / e.MaxRPM)

	// Less advance at higher throttle (to prevent knock)
	throttleRetard := 5.0 * (e.ThrottlePosition / 100.0)

	return baseTiming + rpmAdvance - throttleRetard
}

// getGearRatio returns the gear ratio for a given gear
func (e *Engine) getGearRatio(gear int) float64 {
	physics := DefaultNinja650Physics()
	if gear >= 0 && gear < len(physics.GearRatios) {
		return physics.GearRatios[gear]
	}
	return 0.0 // Neutral or invalid gear
}

func (e *Engine) applyExhaustEffects(baseTorque float64) float64 {
	torqueMultiplier, _ := CalculateExhaustEffect(e.RPM, e.ExhaustType)
	return baseTorque * torqueMultiplier
}

// updateSensorReadings updates simulated sensor readings
func (e *Engine) updateSensorReadings(ecuOutputs ECUOutputs, deltaTime float64) {
	// Update O2 sensor based on ECU fuel injection
	stockInjection := e.calculateStockInjectionTime()
	fuelRatio := ecuOutputs.FuelInjectionTime / stockInjection

	// Lambda is inverse of fuel ratio (simplified)
	// Lambda = 1.0 is stoichiometric (ideal mixture)
	// Lambda < 1.0 is rich, Lambda > 1.0 is lean
	e.O2Reading = 1.0 / fuelRatio

	// Add some noise to the sensor readings for realism
	e.O2Reading += (rand.Float64() - 0.5) * 0.05

	// Update MAP based on throttle position and RPM
	baseMAP := 101.3 // kPa (atmospheric)

	// Lower MAP at higher RPM due to intake vacuum
	rpmVacuum := 30.0 * (e.RPM / e.MaxRPM)

	// Higher MAP with more throttle opening
	throttleEffect := rpmVacuum * (e.ThrottlePosition / 100.0)

	e.MAP = baseMAP - rpmVacuum + throttleEffect

	// Update engine temperature
	// Temperature rises with RPM and load, falls with ambient cooling
	heatGeneration := 0.01 * (e.RPM / 1000.0) * (0.5 + 0.5*e.ThrottlePosition/100.0)
	cooling := 0.005 * (e.EngineTemp - e.AmbientTemp)

	e.EngineTemp += (heatGeneration - cooling) * deltaTime
	e.EngineTemp = math.Max(e.AmbientTemp, math.Min(120.0, e.EngineTemp)) // Clamp between ambient and 120C
}
