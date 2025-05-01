package ecu

import (
	"math"
)

// Map2D represents a 2D lookup table with RPM and load breakpoints
type Map2D struct {
	RPMBreakpoints  []float64
	LoadBreakpoints []float64
	Values          [][]float64
}

// FuelMap represents the fuel map of the ECU
type FuelMap struct {
	Map2D
}

// IgnitionMap represents the ignition timing map of the ECU
type IgnitionMap struct {
	Map2D
}

// AFRMap represents the target air-fuel ratio map of the ECU
type AFRMap struct {
	Map2D
}

// GetValue retrieves an interpolated value from a 2D map
func (m *Map2D) GetValue(rpm, load float64) float64 {
	// Ensure RPM and load are within bounds
	rpm = math.Max(m.RPMBreakpoints[0], math.Min(rpm, m.RPMBreakpoints[len(m.RPMBreakpoints)-1]))
	load = math.Max(m.LoadBreakpoints[0], math.Min(load, m.LoadBreakpoints[len(m.LoadBreakpoints)-1]))

	// Find the indices for interpolation
	rpmLowIdx, rpmHighIdx := 0, 0
	loadLowIdx, loadHighIdx := 0, 0

	// Find RPM indices
	for i := 0; i < len(m.RPMBreakpoints)-1; i++ {
		if rpm >= m.RPMBreakpoints[i] && rpm <= m.RPMBreakpoints[i+1] {
			rpmLowIdx = i
			rpmHighIdx = i + 1
			break
		}
	}

	// Find load indices
	for i := 0; i < len(m.LoadBreakpoints)-1; i++ {
		if load >= m.LoadBreakpoints[i] && load <= m.LoadBreakpoints[i+1] {
			loadLowIdx = i
			loadHighIdx = i + 1
			break
		}
	}

	// Get the four corner values
	v1 := m.Values[rpmLowIdx][loadLowIdx]
	v2 := m.Values[rpmHighIdx][loadLowIdx]
	v3 := m.Values[rpmLowIdx][loadHighIdx]
	v4 := m.Values[rpmHighIdx][loadHighIdx]

	// Calculate interpolation factors
	rpmFactor := 0.0
	if m.RPMBreakpoints[rpmHighIdx] != m.RPMBreakpoints[rpmLowIdx] {
		rpmFactor = (rpm - m.RPMBreakpoints[rpmLowIdx]) / (m.RPMBreakpoints[rpmHighIdx] - m.RPMBreakpoints[rpmLowIdx])
	}

	loadFactor := 0.0
	if m.LoadBreakpoints[loadHighIdx] != m.LoadBreakpoints[loadLowIdx] {
		loadFactor = (load - m.LoadBreakpoints[loadLowIdx]) / (m.LoadBreakpoints[loadHighIdx] - m.LoadBreakpoints[loadLowIdx])
	}

	// Bilinear interpolation
	v12 := v1 + rpmFactor*(v2-v1)
	v34 := v3 + rpmFactor*(v4-v3)

	return v12 + loadFactor*(v34-v12)
}

// SetValue sets a value in the map at the nearest breakpoints
func (m *Map2D) SetValue(rpm, load, value float64) {
	// Find the nearest RPM and load breakpoints
	nearestRPMIdx := 0
	nearestLoadIdx := 0

	smallestRPMDiff := math.Abs(rpm - m.RPMBreakpoints[0])
	smallestLoadDiff := math.Abs(load - m.LoadBreakpoints[0])

	for i, breakpoint := range m.RPMBreakpoints {
		diff := math.Abs(rpm - breakpoint)
		if diff < smallestRPMDiff {
			smallestRPMDiff = diff
			nearestRPMIdx = i
		}
	}

	for i, breakpoint := range m.LoadBreakpoints {
		diff := math.Abs(load - breakpoint)
		if diff < smallestLoadDiff {
			smallestLoadDiff = diff
			nearestLoadIdx = i
		}
	}

	// Set the value
	m.Values[nearestRPMIdx][nearestLoadIdx] = value
}

// ModifyRegion modifies values in a region of the map by a percentage or fixed amount
func (m *Map2D) ModifyRegion(startRPM, endRPM, startLoad, endLoad, modificationPercent float64) {
	for i, rpm := range m.RPMBreakpoints {
		if rpm >= startRPM && rpm <= endRPM {
			for j, load := range m.LoadBreakpoints {
				if load >= startLoad && load <= endLoad {
					// Apply percentage change
					m.Values[i][j] *= (1.0 + modificationPercent/100.0)
				}
			}
		}
	}
}

// DefaultNinja650FuelMap returns a default fuel map for a Ninja 650
func DefaultNinja650FuelMap() FuelMap {
	// RPM breakpoints (1000 to 11000, 1000 RPM steps)
	rpmBreakpoints := []float64{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 11000}

	// Load breakpoints (0 to 100, 10% steps)
	loadBreakpoints := []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

	// Create a default fuel map (values are multipliers, 1.0 = baseline)
	// Rows are RPM, columns are load
	values := [][]float64{
		{0.85, 0.88, 0.90, 0.92, 0.94, 0.96, 0.98, 1.00, 1.02, 1.04, 1.06}, // 1000 RPM
		{0.86, 0.89, 0.91, 0.93, 0.95, 0.97, 0.99, 1.01, 1.03, 1.05, 1.07}, // 2000 RPM
		{0.87, 0.90, 0.92, 0.94, 0.96, 0.98, 1.00, 1.02, 1.04, 1.06, 1.08}, // 3000 RPM
		{0.88, 0.91, 0.93, 0.95, 0.97, 0.99, 1.01, 1.03, 1.05, 1.07, 1.09}, // 4000 RPM
		{0.89, 0.92, 0.94, 0.96, 0.98, 1.00, 1.02, 1.04, 1.06, 1.08, 1.10}, // 5000 RPM
		{0.90, 0.93, 0.95, 0.97, 0.99, 1.01, 1.03, 1.05, 1.07, 1.09, 1.11}, // 6000 RPM
		{0.91, 0.94, 0.96, 0.98, 1.00, 1.02, 1.04, 1.06, 1.08, 1.10, 1.12}, // 7000 RPM
		{0.92, 0.95, 0.97, 0.99, 1.01, 1.03, 1.05, 1.07, 1.09, 1.11, 1.13}, // 8000 RPM
		{0.93, 0.96, 0.98, 1.00, 1.02, 1.04, 1.06, 1.08, 1.10, 1.12, 1.14}, // 9000 RPM
		{0.94, 0.97, 0.99, 1.01, 1.03, 1.05, 1.07, 1.09, 1.11, 1.13, 1.15}, // 10000 RPM
		{0.95, 0.98, 1.00, 1.02, 1.04, 1.06, 1.08, 1.10, 1.12, 1.14, 1.16}, // 11000 RPM
	}

	return FuelMap{
		Map2D{
			RPMBreakpoints:  rpmBreakpoints,
			LoadBreakpoints: loadBreakpoints,
			Values:          values,
		},
	}
}

// DefaultNinja650IgnitionMap returns a default ignition map for a Ninja 650
func DefaultNinja650IgnitionMap() IgnitionMap {
	// RPM breakpoints (1000 to 11000, 1000 RPM steps)
	rpmBreakpoints := []float64{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 11000}

	// Load breakpoints (0 to 100, 10% steps)
	loadBreakpoints := []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

	// Create a default ignition map (values are degrees BTDC)
	// Rows are RPM, columns are load
	values := [][]float64{
		{10, 15, 20, 25, 28, 30, 32, 33, 32, 30, 28}, // 1000 RPM
		{12, 18, 23, 28, 30, 32, 34, 35, 34, 32, 30}, // 2000 RPM
		{15, 20, 25, 30, 32, 34, 36, 37, 36, 34, 32}, // 3000 RPM
		{18, 23, 28, 32, 34, 36, 38, 39, 38, 36, 34}, // 4000 RPM
		{20, 25, 30, 34, 36, 38, 40, 40, 39, 37, 35}, // 5000 RPM
		{22, 27, 32, 36, 38, 40, 41, 41, 40, 38, 36}, // 6000 RPM
		{24, 29, 34, 38, 40, 41, 42, 42, 41, 39, 37}, // 7000 RPM
		{25, 30, 35, 39, 41, 42, 43, 43, 42, 40, 38}, // 8000 RPM
		{25, 30, 35, 39, 41, 42, 43, 43, 42, 40, 38}, // 9000 RPM
		{24, 29, 34, 38, 40, 41, 42, 42, 41, 39, 37}, // 10000 RPM
		{22, 27, 32, 36, 38, 40, 41, 41, 40, 38, 36}, // 11000 RPM
	}

	return IgnitionMap{
		Map2D{
			RPMBreakpoints:  rpmBreakpoints,
			LoadBreakpoints: loadBreakpoints,
			Values:          values,
		},
	}
}

// DefaultNinja650AFRMap returns a default target air-fuel ratio map for a Ninja 650
func DefaultNinja650AFRMap() AFRMap {
	// RPM breakpoints (1000 to 11000, 1000 RPM steps)
	rpmBreakpoints := []float64{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000, 11000}

	// Load breakpoints (0 to 100, 10% steps)
	loadBreakpoints := []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

	// Create a default AFR map (values are AFR, 14.7 is stoichiometric for gasoline)
	// Rows are RPM, columns are load
	values := [][]float64{
		{14.7, 14.7, 14.7, 14.5, 14.3, 14.1, 13.8, 13.5, 13.2, 12.9, 12.6}, // 1000 RPM
		{14.7, 14.7, 14.7, 14.5, 14.3, 14.1, 13.8, 13.5, 13.2, 12.9, 12.6}, // 2000 RPM
		{14.7, 14.7, 14.7, 14.5, 14.3, 14.1, 13.8, 13.5, 13.2, 12.9, 12.6}, // 3000 RPM
		{14.7, 14.7, 14.7, 14.5, 14.3, 14.1, 13.8, 13.5, 13.2, 12.9, 12.6}, // 4000 RPM
		{14.7, 14.7, 14.6, 14.4, 14.2, 14.0, 13.7, 13.4, 13.1, 12.8, 12.5}, // 5000 RPM
		{14.7, 14.6, 14.5, 14.3, 14.1, 13.9, 13.6, 13.3, 13.0, 12.7, 12.4}, // 6000 RPM
		{14.6, 14.5, 14.4, 14.2, 14.0, 13.8, 13.5, 13.2, 12.9, 12.6, 12.3}, // 7000 RPM
		{14.5, 14.4, 14.3, 14.1, 13.9, 13.7, 13.4, 13.1, 12.8, 12.5, 12.2}, // 8000 RPM
		{14.4, 14.3, 14.2, 14.0, 13.8, 13.6, 13.3, 13.0, 12.7, 12.4, 12.1}, // 9000 RPM
		{14.3, 14.2, 14.1, 13.9, 13.7, 13.5, 13.2, 12.9, 12.6, 12.3, 12.0}, // 10000 RPM
		{14.2, 14.1, 14.0, 13.8, 13.6, 13.4, 13.1, 12.8, 12.5, 12.2, 11.9}, // 11000 RPM
	}

	return AFRMap{
		Map2D{
			RPMBreakpoints:  rpmBreakpoints,
			LoadBreakpoints: loadBreakpoints,
			Values:          values,
		},
	}
}
