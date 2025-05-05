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
	EngineTemp       float64 // Celsius
	AirTemp          float64 // Celsius
	MAP              float64 // kPa
	O2Reading        float64 // Lambda value
	Speed            float64 // km/h
	Gear             int     // 0 = neutral, 1-6 = gears
	BrakeApplied     bool
	RevLimit         float64 // RPM at which rev limiter kicks in

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

	// Transmission properties
	CurrentGear        int       // 0 = neutral, 1-6 = gears
	GearRatios         []float64 // Gear ratios for each gear
	FinalDriveRatio    float64   // Chain/sprocket ratio
	WheelCircumference float64   // Wheel circumference in meters

	// Clutch properties
	ClutchSlip     float64 // How much power gets through partially engaged clutch
	ClutchPosition float64 // 0.0 = fully engaged, 1.0 = fully disengaged
	ShiftTimer     float64 // Timer for shift animation/physics
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
		RevLimit:         0,     // No rev limit

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

		ClutchSlip: 0.0, // start with no slip

		GearRatios: []float64{
			0.0,   // Neutral
			2.438, // 1st gear
			1.714, // 2nd gear
			1.333, // 3rd gear
			1.111, // 4th gear
			0.966, // 5th gear
			0.852, // 6th gear
		},
		FinalDriveRatio:    3.067, // Chain drive ratio
		WheelCircumference: 1.95,  // meters (650cc sport bike)
	}
}

func (e *Engine) Update(ecuOutputs ECUOutputs, deltaTime float64) {
	// Calculate baseline engine torque
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

	// Final engine torque
	engineTorque := baselineTorque * torqueMultiplier

	// Calculate clutch transfer torque
	transferTorque := engineTorque
	clutchSlipping := false

	if e.ClutchPosition > 0.0 {
		// Reduce torque transfer based on clutch position
		// 0.0 = fully engaged, 1.0 = fully disengaged
		clutchFactor := 1.0 - e.ClutchPosition
		transferTorque = engineTorque * clutchFactor

		// Calculate clutch slip when partially engaged
		if e.ClutchPosition > 0.0 && e.ClutchPosition < 1.0 {
			clutchSlipping = true

			// Calculate RPM difference between engine and transmission input
			transmissionInputRPM := e.calculateTransmissionInputRPM()
			rpmDiff := e.RPM - transmissionInputRPM

			// Store slip value for potential display/analysis
			e.ClutchSlip = rpmDiff * (1.0 - e.ClutchPosition)
		}
	}

	// Handle vehicle physics based on clutch and gear state
	if e.ClutchPosition >= 0.95 || e.Gear == 0 {
		// CASE 1: Clutch fully disengaged or in neutral - engine runs free

		// Calculate RPM change based only on engine torque and internal friction
		rpmChange := (engineTorque - (e.RPM * 0.001)) * 10 / e.RotationalInertia
		e.RPM += rpmChange * deltaTime

		// Apply idle control
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

		// Handle vehicle deceleration (no engine braking)
		brakeValue := 0.5
		if e.BrakeApplied {
			brakeValue = 5.0
		}

		// Apply aerodynamic drag and rolling resistance
		dragForce := e.DragCoefficient*e.Speed*e.Speed*0.001 + brakeValue
		e.Speed = math.Max(0, e.Speed-dragForce*deltaTime)

	} else if clutchSlipping {
		// CASE 2: Clutch partially engaged - complex model with slip

		// Calculate torque at wheels
		gearRatio := e.getGearRatio(e.Gear)
		finalDriveRatio := e.FinalDriveRatio
		wheelTorque := transferTorque * gearRatio * finalDriveRatio

		// Calculate wheel force and acceleration
		wheelRadius := e.WheelCircumference / (2 * math.Pi)
		wheelForce := wheelTorque / wheelRadius

		// Mass factor - simplified physics
		vehicleMass := 200.0 // kg, approximate mass of Ninja 650 with rider
		acceleration := wheelForce / vehicleMass

		// Apply to speed
		e.Speed += acceleration * deltaTime

		// Calculate engine RPM changes due to torque and slip
		// Engine is pulled down by transmission but also pushed by throttle
		rpmPulldown := e.ClutchSlip * (1.0 - e.ClutchPosition) * 0.1
		rpmFromTorque := (engineTorque - (e.RPM * 0.001)) * 5 / e.RotationalInertia

		// Combine effects
		e.RPM += (rpmFromTorque - rpmPulldown) * deltaTime

		// Apply braking
		if e.BrakeApplied {
			brakeDecel := 5.0 // m/s²
			e.Speed = math.Max(0, e.Speed-brakeDecel*deltaTime)
		}

	} else {
		// CASE 3: Clutch fully engaged, in gear - direct connection

		// Calculate wheel torque through drivetrain
		gearRatio := e.getGearRatio(e.Gear)
		finalDriveRatio := e.FinalDriveRatio
		wheelTorque := engineTorque * gearRatio * finalDriveRatio

		// Apply drivetrain efficiency
		drivetrainEfficiency := 0.9 // 90% efficiency
		wheelTorque *= drivetrainEfficiency

		// Calculate wheel force and acceleration
		wheelRadius := e.WheelCircumference / (2 * math.Pi)
		wheelForce := wheelTorque / wheelRadius

		// Vehicle mass and load simulation
		vehicleMass := 200.0 // kg
		roadGradient := 0.0  // flat road

		// Factor in road gradient (simplified)
		gravityComponent := 9.81 * math.Sin(roadGradient) * vehicleMass

		// Total force = wheel force - drag - rolling resistance - gravity
		brakeValue := 0.0
		if e.BrakeApplied {
			brakeValue = 1000.0 // Brake force in N
		}

		dragForce := e.DragCoefficient * e.Speed * e.Speed * 0.2
		rollingResistance := 0.015 * vehicleMass * 9.81 // Rolling resistance coefficient * normal force

		netForce := wheelForce - dragForce - rollingResistance - gravityComponent - brakeValue
		acceleration := netForce / vehicleMass

		// Update speed
		e.Speed += acceleration * deltaTime
		e.Speed = math.Max(0, e.Speed)

		// Calculate RPM directly from wheel speed
		if e.Speed > 0 {
			// Speed to wheel RPM to engine RPM
			wheelRPM := (e.Speed * 1000 / 3600) * 60 / e.WheelCircumference
			e.RPM = wheelRPM * gearRatio * finalDriveRatio

			// Limit to redline (with slight overrev allowed)
			e.RPM = math.Min(e.RPM, e.RedlineRPM*1.05)
		} else if e.ThrottlePosition < 5 {
			// If stopped with throttle closed, engine may stall
			if e.RPM < 1000 {
				stallProbability := (1000 - e.RPM) / 1000.0

				// Simple stalling model
				if rand.Float64() < stallProbability*deltaTime*0.5 {
					e.RPM = 0 // Engine stalled
				}
			}
		}
	}

	// Engine rev limiter
	if e.RPM > e.RevLimit && e.RevLimit > 0 {
		// Cut spark/fuel when hitting rev limiter
		e.RPM = e.RevLimit - (rand.Float64() * 200) // Bouncing on limiter
	}

	// Ensure RPM stays in valid range and above 0
	e.RPM = math.Max(0, math.Min(e.RPM, e.RedlineRPM*1.05))

	// Update sensor readings
	e.updateSensorReadings(ecuOutputs, deltaTime)

	// Update last time
	e.lastUpdateTime = time.Now()
}

// Calculate transmission input RPM from wheel speed
func (e *Engine) calculateTransmissionInputRPM() float64 {
	if e.Gear == 0 {
		return 0.0 // Neutral
	}

	// Calculate wheel RPM from speed (km/h)
	// Speed (km/h) = Wheel RPM * Circumference (m) * 60 / 1000
	wheelRPM := e.Speed * 1000.0 / (e.WheelCircumference * 60.0)

	// Calculate transmission input RPM
	return wheelRPM * e.FinalDriveRatio * e.getGearRatio(e.Gear)
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

// Calculate vehicle speed from engine RPM
func (e *Engine) calculateSpeedFromRPM() float64 {
	if e.Gear == 0 || e.ClutchPosition >= 1.0 {
		return e.Speed // No change if clutch disengaged or in neutral
	}

	// Calculate wheel RPM
	wheelRPM := e.RPM / (e.GearRatios[e.Gear] * e.FinalDriveRatio)

	// Calculate speed
	return wheelRPM * e.WheelCircumference * 60.0 / 1000.0
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
	if gear >= 0 && gear < len(e.GearRatios) {
		return e.GearRatios[gear]
	}
	return 0.0 // Neutral or invalid gear
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
