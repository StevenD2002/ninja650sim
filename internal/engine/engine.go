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

	// Advanced Environmental Factors:
	FuelOctane           float64 // 87, 91, 93, 100+ for race fuel
	AirFilterRestriction float64 // 0.0-1.0, 0=clean, 1=completely blocked
	AtmosphericPressure  float64 // kPa
	Humidity             float64 // 0.0-1.0
	FuelQuality          float64 // 0.0-1.0, accounts for ethanol content, age, etc.

	// Engine condition factors
	EngineWear    float64 // 0.0-1.0, affects compression and efficiency
	CarbonBuildup float64 // 0.0-1.0, affects timing and efficiency

	// Additional physics properties
	FrontalArea float64 // m² for better drag calculation
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

		// New environmental factors
		FuelOctane:           91.0,  // Premium fuel
		AirFilterRestriction: 0.0,   // Clean filter
		AtmosphericPressure:  101.3, // kPa at sea level
		Humidity:             0.5,   // 50% humidity
		FuelQuality:          1.0,   // Perfect fuel

		// Engine condition
		EngineWear:    0.0, // New engine
		CarbonBuildup: 0.0, // Clean engine

		// Physics
		FrontalArea: 0.7, // m² this is more of an estimate
	}
}

func (e *Engine) Update(ecuOutputs ECUOutputs, deltaTime float64) {
	// Get physics constants and motorcycle specs
	physics := DefaultNinja650Physics()

	// Calculate proper air density based on environment
	airDensity := CalculateAirDensity(e.Altitude, e.AirTemp)

	// Existing torque calculations...
	baselineTorque := e.calculateBaselineTorque()
	torqueMultiplier := 1.0

	// Apply advanced environmental effects
	torqueMultiplier = e.applyEnvironmentalEffects(torqueMultiplier, airDensity)

	// Apply ECU fuel enrichment/leaning
	if ecuOutputs.FuelInjectionTime > 0 {
		stockInjectionTime := e.calculateStockInjectionTime()
		torqueMultiplier *= ecuOutputs.FuelInjectionTime / stockInjectionTime

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

	// Apply octane effects on ignition timing
	_, powerMultiplier := e.calculateOctaneEffects(ecuOutputs.IgnitionAdvance)
	torqueMultiplier *= powerMultiplier

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

		// Handle vehicle deceleration using proper physics
		dragForce := CalculateAerodynamicDrag(e.Speed, e.DragCoefficient, e.FrontalArea, airDensity)
		rollingForce := CalculateRollingResistance(e.Speed, physics.Mass, physics.RollingResistance)

		brakeForce := 0.0
		if e.BrakeApplied {
			brakeForce = 1000.0 // N
		}

		totalResistance := (dragForce + rollingForce + brakeForce) / physics.Mass
		e.Speed = math.Max(0, e.Speed-totalResistance*deltaTime)

	} else if clutchSlipping {
		// CASE 2: Clutch partially engaged - complex model with slip

		// Use proper physics calculation for wheel torque
		wheelTorque := CalculateWheelTorqueFromEngineTorque(
			transferTorque, e.Gear, e.GearRatios, e.FinalDriveRatio, 0.9)

		// Calculate wheel force and acceleration
		wheelRadius := e.WheelCircumference / (2 * math.Pi)
		wheelForce := wheelTorque / wheelRadius

		// Calculate resistance forces using proper physics
		dragForce := CalculateAerodynamicDrag(e.Speed, e.DragCoefficient, e.FrontalArea, airDensity)
		rollingForce := CalculateRollingResistance(e.Speed, physics.Mass, physics.RollingResistance)

		brakeForce := 0.0
		if e.BrakeApplied {
			brakeForce = 1000.0 // N
		}

		// Net force and acceleration
		netForce := wheelForce - dragForce - rollingForce - brakeForce
		acceleration := netForce / physics.Mass

		// Apply to speed
		e.Speed += acceleration * deltaTime
		e.Speed = math.Max(0, e.Speed)

		// Calculate engine RPM changes due to torque and slip
		// Engine is pulled down by transmission but also pushed by throttle
		rpmPulldown := e.ClutchSlip * (1.0 - e.ClutchPosition) * 0.1
		rpmFromTorque := (engineTorque - (e.RPM * 0.001)) * 5 / e.RotationalInertia

		// Combine effects
		e.RPM += (rpmFromTorque - rpmPulldown) * deltaTime

	} else {
		// CASE 3: Clutch fully engaged, in gear - direct connection

		// Use proper physics calculation for wheel torque
		wheelTorque := CalculateWheelTorqueFromEngineTorque(
			engineTorque, e.Gear, e.GearRatios, e.FinalDriveRatio, 0.9)

		// Calculate wheel force
		wheelRadius := e.WheelCircumference / (2 * math.Pi)
		wheelForce := wheelTorque / wheelRadius

		// Calculate resistance forces using proper physics
		dragForce := CalculateAerodynamicDrag(e.Speed, e.DragCoefficient, e.FrontalArea, airDensity)
		rollingForce := CalculateRollingResistance(e.Speed, physics.Mass, physics.RollingResistance)

		// Road gradient (simplified - flat road)
		roadGradient := 0.0
		gravityComponent := 9.81 * math.Sin(roadGradient) * physics.Mass

		// Apply braking
		brakeForce := 0.0
		if e.BrakeApplied {
			brakeForce = 1000.0 // N
		}

		// Net force and acceleration
		netForce := wheelForce - dragForce - rollingForce - gravityComponent - brakeForce
		acceleration := netForce / physics.Mass

		// Update speed
		e.Speed += acceleration * deltaTime
		e.Speed = math.Max(0, e.Speed)

		// Calculate RPM from speed using proper physics
		if e.Speed > 0 {
			e.RPM = CalculateRPMFromSpeed(e.Speed, e.Gear, e.GearRatios,
				e.FinalDriveRatio, physics.WheelDiameter)
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

	// Simulate engine wear over time
	e.SimulateEngineWear(deltaTime)

	// Update sensor readings
	e.updateSensorReadings(ecuOutputs, deltaTime)

	// Update last time
	e.lastUpdateTime = time.Now()
}

// SetEnvironmentalConditions allows setting multiple environmental factors
func (e *Engine) SetEnvironmentalConditions(octane, airFilterRestriction, altitude, humidity float64) {
	e.FuelOctane = math.Max(80, math.Min(110, octane))
	e.AirFilterRestriction = math.Max(0, math.Min(1, airFilterRestriction))
	e.Altitude = math.Max(0, math.Min(5000, altitude)) // Up to 5000m altitude
	e.Humidity = math.Max(0, math.Min(1, humidity))
}

// SimulateEngineWear advances engine wear over time
func (e *Engine) SimulateEngineWear(deltaTime float64) {
	// Wear rate depends on RPM, temperature, and load
	baseWearRate := 0.000001 // Very slow base rate

	// Higher RPM increases wear
	rpmFactor := math.Pow(e.RPM/e.MaxRPM, 2)

	// Higher temperature increases wear
	tempFactor := math.Max(1.0, (e.EngineTemp-90.0)*0.02)

	// Higher load increases wear
	loadFactor := 1.0 + (e.ThrottlePosition / 100.0)

	// Poor fuel quality increases wear
	fuelFactor := 2.0 - e.FuelQuality

	wearRate := baseWearRate * rpmFactor * tempFactor * loadFactor * fuelFactor

	e.EngineWear += wearRate * deltaTime
	e.EngineWear = math.Min(1.0, e.EngineWear)

	// Carbon buildup over time (especially with rich mixtures)
	carbonRate := 0.00001
	if e.O2Reading < 0.95 { // Rich mixture
		carbonRate *= 2.0
	}

	e.CarbonBuildup += carbonRate * deltaTime
	e.CarbonBuildup = math.Min(1.0, e.CarbonBuildup)
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

// Apply environmental effects to torque multiplier
func (e *Engine) applyEnvironmentalEffects(torqueMultiplier, airDensity float64) float64 {
	// Air filter restriction reduces airflow
	airflowReduction := 1.0 - (e.AirFilterRestriction * 0.25) // Max 25% reduction
	torqueMultiplier *= airflowReduction

	// Air density effects (altitude, temperature, humidity)
	standardAirDensity := 1.225 // kg/m³ at sea level, 15°C
	densityFactor := airDensity / standardAirDensity
	torqueMultiplier *= (0.7 + 0.3*densityFactor) // Partial density dependence

	// Fuel quality effects
	torqueMultiplier *= (0.9 + 0.1*e.FuelQuality) // 10% swing for fuel quality

	// Engine wear effects
	wearFactor := 1.0 - (e.EngineWear * 0.15) // Max 15% loss from wear
	torqueMultiplier *= wearFactor

	// Carbon buildup effects (reduces efficiency)
	carbonFactor := 1.0 - (e.CarbonBuildup * 0.08) // Max 8% loss
	torqueMultiplier *= carbonFactor

	return torqueMultiplier
}

// Calculate octane effects on knock limit and power
func (e *Engine) calculateOctaneEffects(ignitionAdvance float64) (knockLimit, powerMultiplier float64) {
	// Base knock threshold for different octanes
	var baseKnockThreshold float64
	switch {
	case e.FuelOctane >= 100:
		baseKnockThreshold = 45.0 // Race fuel
	case e.FuelOctane >= 93:
		baseKnockThreshold = 35.0 // Premium
	case e.FuelOctane >= 91:
		baseKnockThreshold = 30.0 // Mid-grade
	case e.FuelOctane >= 87:
		baseKnockThreshold = 25.0 // Regular
	default:
		baseKnockThreshold = 20.0 // Low quality fuel
	}

	// Adjust for engine temperature and load
	tempAdjustment := (e.EngineTemp - 90.0) * 0.15 // Hotter = more knock prone
	loadAdjustment := e.ThrottlePosition * 0.1     // Higher load = more knock prone
	carbonAdjustment := e.CarbonBuildup * 5.0      // Carbon increases knock tendency

	knockLimit = baseKnockThreshold - tempAdjustment - loadAdjustment - carbonAdjustment

	// Power multiplier based on knock proximity
	if ignitionAdvance > knockLimit {
		// Knock occurring - major power loss and potential damage
		powerMultiplier = 0.6
		// In a real implementation, you might track engine damage here
	} else if ignitionAdvance > knockLimit*0.95 {
		// Very close to knock - slight power loss
		powerMultiplier = 0.95
	} else if ignitionAdvance > knockLimit*0.8 {
		// Getting close to knock - minimal loss
		powerMultiplier = 0.98
	} else {
		// Safe operation
		powerMultiplier = 1.0
	}

	return knockLimit, powerMultiplier
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
// the formula for injection time is based on volumetric efficiency, RPM, throttle position, and environmental factors (this is a simplified model)
func (e *Engine) calculateStockInjectionTime() float64 {
	// Base injection time
	baseTime := 2.0

	// Calculate volumetric efficiency with environmental effects
	baseVE := 0.85

	// Air filter restriction
	airflowFactor := 1.0 - (e.AirFilterRestriction * 0.3)

	// Altitude effects (less dense air)
	altitudeFactor := math.Exp(-e.Altitude / 8400.0)

	// Temperature effects (colder air is denser)
	tempFactor := (273.15 + 15.0) / (273.15 + e.AirTemp)

	// Humidity reduces air density
	humidityFactor := 1.0 - (e.Humidity * 0.02)

	// Engine wear affects breathing efficiency
	wearFactor := 1.0 - (e.EngineWear * 0.1)

	// Carbon buildup affects valve operation
	carbonFactor := 1.0 - (e.CarbonBuildup * 0.05)

	// Combined volumetric efficiency
	totalVE := baseVE * airflowFactor * altitudeFactor * tempFactor * humidityFactor * wearFactor * carbonFactor
	totalVE = math.Max(0.3, math.Min(1.1, totalVE)) // Reasonable bounds

	// RPM factor
	rpmFactor := 1.0 + 0.5*(e.RPM-e.IdleRPM)/(e.MaxRPM-e.IdleRPM)

	// Throttle factor
	throttleFactor := 0.5 + 0.5*(e.ThrottlePosition/100.0)

	return baseTime * totalVE * rpmFactor * throttleFactor
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
