package engine

import (
	"math"
)

// PhysicsConstants contains universal constants used in calculations
type PhysicsConstants struct {
	GravityAcceleration float64 // m/s²
	AirDensityAtSea     float64 // kg/m³
	StandardPressure    float64 // kPa
	GasConstant         float64 // J/(kg·K)
}

// MotorcyclePhysics contains motorcycle-specific physical properties
type MotorcyclePhysics struct {
	Mass                   float64 // kg
	FrontalArea            float64 // m²
	WheelDiameter          float64 // m
	FinalDriveRatio        float64
	GearRatios             []float64
	WheelInertia           float64 // kg·m²
	EngineMomentOfInertia  float64 // kg·m²
	TransmissionEfficiency float64 // 0-1
	RollingResistance      float64 // coefficient
}

// DefaultPhysicsConstants returns standard physics constants
func DefaultPhysicsConstants() PhysicsConstants {
	return PhysicsConstants{
		GravityAcceleration: 9.81,    // m/s²
		AirDensityAtSea:     1.225,   // kg/m³ at sea level, 15°C
		StandardPressure:    101.325, // kPa at sea level
		GasConstant:         287.058, // J/(kg·K) for dry air
	}
}

// DefaultNinja650Physics returns physics parameters for a Ninja 650
func DefaultNinja650Physics() MotorcyclePhysics {
	return MotorcyclePhysics{
		Mass:            196.0, // kg (wet weight)
		FrontalArea:     0.7,   // m² (approximate)
		WheelDiameter:   0.43,  // m (17 inch wheel)
		FinalDriveRatio: 3.067, // Chain drive ratio
		GearRatios: []float64{
			0.0,   // Neutral
			2.438, // 1st gear
			1.714, // 2nd gear
			1.333, // 3rd gear
			1.111, // 4th gear
			0.966, // 5th gear
			0.852, // 6th gear
		},
		WheelInertia:           0.8,   // kg·m² (approximate)
		EngineMomentOfInertia:  0.12,  // kg·m² (approximate)
		TransmissionEfficiency: 0.9,   // 90% efficiency
		RollingResistance:      0.015, // typical motorcycle tire
	}
}

// CalculateAirDensity calculates air density based on altitude and temperature
func CalculateAirDensity(altitude, temperature float64) float64 {
	// Standard physics calculations for air density
	constants := DefaultPhysicsConstants()

	// Temperature in Kelvin
	tempK := temperature + 273.15

	// Barometric pressure approximation based on altitude
	// Using the barometric formula: P = P0 * exp(-g*h/(R*T))
	pressure := constants.StandardPressure * math.Exp(-constants.GravityAcceleration*altitude/(constants.GasConstant*tempK))

	// Density calculation: ρ = P/(R*T)
	density := pressure / (constants.GasConstant * tempK)

	return density
}

// CalculateAerodynamicDrag calculates aerodynamic drag force
func CalculateAerodynamicDrag(speed, dragCoefficient, frontalArea, airDensity float64) float64 {
	// Drag Force = 0.5 * ρ * v² * Cd * A
	// speed should be in m/s
	speedMS := speed / 3.6 // Convert km/h to m/s
	return 0.5 * airDensity * speedMS * speedMS * dragCoefficient * frontalArea
}

// CalculateRollingResistance calculates rolling resistance force
func CalculateRollingResistance(speed, mass, rollingCoefficient float64) float64 {
	// Simple rolling resistance model
	// F = Crr * m * g
	// With slight speed dependency
	constants := DefaultPhysicsConstants()
	speedFactor := 1.0 + (speed/100.0)*0.1 // Slight increase with speed

	return rollingCoefficient * mass * constants.GravityAcceleration * speedFactor
}

// CalculateWheelTorqueFromEngineTorque converts engine torque to wheel torque
func CalculateWheelTorqueFromEngineTorque(
	engineTorque float64,
	gear int,
	gearRatios []float64,
	finalDriveRatio float64,
	efficiency float64,
) float64 {
	if gear <= 0 || gear >= len(gearRatios) {
		return 0.0 // Neutral or invalid gear
	}

	// Wheel torque = Engine torque * Gear ratio * Final drive ratio * Efficiency
	return engineTorque * gearRatios[gear] * finalDriveRatio * efficiency
}

// CalculateRPMFromSpeed calculates engine RPM based on vehicle speed
func CalculateRPMFromSpeed(
	speedKmh float64,
	gear int,
	gearRatios []float64,
	finalDriveRatio float64,
	wheelDiameter float64,
) float64 {
	if gear <= 0 || gear >= len(gearRatios) {
		return 0.0 // Neutral or invalid gear
	}

	// Convert speed to wheel RPM
	wheelCircumference := math.Pi * wheelDiameter // meters
	speedMS := speedKmh / 3.6                     // Convert km/h to m/s
	wheelRPM := (speedMS * 60.0) / wheelCircumference

	// Convert wheel RPM to engine RPM
	return wheelRPM * gearRatios[gear] * finalDriveRatio
}

// CalculateSpeedFromRPM calculates vehicle speed based on engine RPM
func CalculateSpeedFromRPM(
	rpm float64,
	gear int,
	gearRatios []float64,
	finalDriveRatio float64,
	wheelDiameter float64,
) float64 {
	if gear <= 0 || gear >= len(gearRatios) {
		return 0.0 // Neutral or invalid gear
	}

	// Convert engine RPM to wheel RPM
	wheelRPM := rpm / (gearRatios[gear] * finalDriveRatio)

	// Convert wheel RPM to speed
	wheelCircumference := math.Pi * wheelDiameter // meters
	speedMS := (wheelRPM * wheelCircumference) / 60.0
	return speedMS * 3.6 // Convert m/s to km/h
}

// CalculateEngineBraking calculates engine braking torque
func CalculateEngineBraking(rpm, displacement, compressionRatio float64) float64 {
	// Simple engine braking model based on engine parameters
	// - Higher RPM = more braking
	// - Higher displacement = more braking
	// - Higher compression ratio = more braking

	// Base torque is proportional to displacement
	baseTorque := displacement * 0.01

	// RPM factor (more braking at higher RPM)
	rpmFactor := math.Min(1.0, rpm/3000.0)

	// Compression factor
	compressionFactor := compressionRatio / 10.0

	return baseTorque * rpmFactor * compressionFactor
}

// CalculateEngineInertia calculates approximate engine rotational inertia
func CalculateEngineInertia(displacement float64) float64 {
	// Very rough approximation based on engine displacement
	// Actual values would depend on specific engine design
	return displacement * 0.0002
}

// CalculateOptimalShiftPoints calculates optimal RPM for upshifting gears
func CalculateOptimalShiftPoints(maxTorqueRPM, redlineRPM float64) []float64 {
	// Simple shift point calculation
	// For maximum acceleration, usually shift just after max torque RPM
	// For maximum power, shift closer to redline

	// Calculate a point between max torque and redline for each gear
	powerBand := redlineRPM - maxTorqueRPM
	shiftPoints := make([]float64, 7) // Neutral + 6 gears

	shiftPoints[0] = 0 // Neutral

	// Lower gears - shift closer to redline for maximum acceleration
	shiftPoints[1] = redlineRPM - powerBand*0.1  // 1st gear
	shiftPoints[2] = redlineRPM - powerBand*0.15 // 2nd gear
	shiftPoints[3] = redlineRPM - powerBand*0.2  // 3rd gear

	// Higher gears - shift closer to max torque for better efficiency
	shiftPoints[4] = maxTorqueRPM + powerBand*0.4 // 4th gear
	shiftPoints[5] = maxTorqueRPM + powerBand*0.3 // 5th gear
	shiftPoints[6] = maxTorqueRPM + powerBand*0.2 // 6th gear

	return shiftPoints
}

// CalculateThrottleResponse models the non-linear response of throttle openings
func CalculateThrottleResponse(throttlePosition float64) float64 {
	// Convert linear throttle position (0-100%) to non-linear power delivery
	// This accounts for the fact that most bikes deliver power non-linearly
	// Typically the first 25% of throttle might only deliver 5-10% of power

	if throttlePosition <= 0 {
		return 0
	}

	// Normalize to 0-1 range
	normalizedPosition := throttlePosition / 100.0

	// Apply a power curve for more realistic throttle response
	// Using a cubic curve for more sensitive mid-range
	return math.Pow(normalizedPosition, 3)
}

// CalculateExhaustEffect calculates the effect of an aftermarket exhaust
func CalculateExhaustEffect(rpm float64, exhaustType string) (torqueMultiplier, powerMultiplier float64) {
	// Default - no change
	torqueMultiplier = 1.0
	powerMultiplier = 1.0

	// RPM normalized to 0-1 range (assuming max RPM of 11000)
	normalizedRPM := math.Min(1.0, rpm/11000.0)

	switch exhaustType {
	case "Stock":
		// No change
		return

	case "Yoshimura Alpha 2":
		// Yoshimura typically adds power in mid-high RPM range
		// with less improvement at very low and very high RPM

		if normalizedRPM < 0.2 {
			// Low RPM range - slight improvement
			torqueMultiplier = 1.02
		} else if normalizedRPM < 0.4 {
			// Low-mid RPM range - moderate improvement
			torqueMultiplier = 1.03 + 0.02*(normalizedRPM-0.2)/0.2
		} else if normalizedRPM < 0.7 {
			// Mid-high RPM range - best improvement
			torqueMultiplier = 1.05 + 0.05*(normalizedRPM-0.4)/0.3
		} else {
			// High RPM range - good improvement but tapering off
			torqueMultiplier = 1.10 - 0.03*(normalizedRPM-0.7)/0.3
		}

		// Power is affected similarly
		powerMultiplier = torqueMultiplier

	case "Akrapovic Full System":
		// Akrapovic typically offers more significant gains
		if normalizedRPM < 0.2 {
			torqueMultiplier = 1.03
		} else if normalizedRPM < 0.4 {
			torqueMultiplier = 1.04 + 0.04*(normalizedRPM-0.2)/0.2
		} else if normalizedRPM < 0.7 {
			torqueMultiplier = 1.08 + 0.06*(normalizedRPM-0.4)/0.3
		} else {
			torqueMultiplier = 1.14 - 0.02*(normalizedRPM-0.7)/0.3
		}

		powerMultiplier = torqueMultiplier * 1.01 // Slightly more power gain

	default:
		// Unknown exhaust type
		return
	}

	return
}

// SimulateAcceleration simulates acceleration time from 0 to a target speed
func SimulateAcceleration(
	engine *Engine,
	physics MotorcyclePhysics,
	targetSpeedKmh float64,
) float64 {
	// Create a copy of the engine to avoid modifying the original
	simulationEngine := *engine

	// Reset engine to idle in neutral
	simulationEngine.RPM = simulationEngine.IdleRPM
	simulationEngine.Speed = 0
	simulationEngine.Gear = 0
	simulationEngine.ThrottlePosition = 0

	// Prepare simulation
	totalTime := 0.0
	timeStep := 0.05 // 50ms simulation steps

	// ECU outputs for simulation (full power)
	ecuOutputs := ECUOutputs{
		FuelInjectionTime: simulationEngine.calculateStockInjectionTime() * 1.0, // Stock fueling
		IgnitionAdvance:   simulationEngine.calculateOptimalTiming(),            // Optimal timing
		TargetIdleRPM:     simulationEngine.IdleRPM,
		LambdaTarget:      1.0, // Stoichiometric
	}

	// Calculate optimal shift points
	shiftPoints := CalculateOptimalShiftPoints(simulationEngine.MaxTorqueRPM, simulationEngine.RedlineRPM)

	// Acceleration simulation
	for simulationEngine.Speed < targetSpeedKmh {
		// Full throttle
		simulationEngine.ThrottlePosition = 100.0

		// Determine optimal gear
		if simulationEngine.Speed > 5.0 && simulationEngine.Gear == 0 {
			// Start moving - engage 1st gear
			simulationEngine.Gear = 1
		} else if simulationEngine.Gear > 0 && simulationEngine.Gear < len(shiftPoints)-1 {
			// Check if we should shift up
			if simulationEngine.RPM >= shiftPoints[simulationEngine.Gear] {
				// Simulate clutch operation and shift
				simulationEngine.Gear++
				simulationEngine.RPM = CalculateRPMFromSpeed(
					simulationEngine.Speed,
					simulationEngine.Gear,
					physics.GearRatios,
					physics.FinalDriveRatio,
					physics.WheelDiameter,
				)
			}
		}

		// Update engine state
		simulationEngine.Update(ecuOutputs, timeStep)

		// Advance time
		totalTime += timeStep

		// Safety break to prevent infinite loops
		if totalTime > 30.0 {
			break // 30 seconds should be enough to reach any reasonable speed
		}
	}

	return totalTime
}
