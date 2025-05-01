package ui

import (
	"fmt"
	"math"

	"github.com/rivo/tview"
)

// Gauge represents a visual gauge for a numeric value
type Gauge struct {
	View       *tview.TextView
	title      string
	min        float64
	max        float64
	warningMin float64
	warningMax float64
	value      float64
	width      int
	precision  int
}

// NewGauge creates a new gauge
func NewGauge(title string, min, max, warningMin, warningMax float64) *Gauge {
	gauge := &Gauge{
		View:       tview.NewTextView().SetDynamicColors(true),
		title:      title,
		min:        min,
		max:        max,
		warningMin: warningMin,
		warningMax: warningMax,
		value:      min,
		width:      50,
		precision:  1,
	}

	gauge.View.SetBorder(true)
	gauge.View.SetTitle(title)
	gauge.SetValue(min)

	return gauge
}

// SetValue updates the gauge with a new value
func (g *Gauge) SetValue(value float64) {
	// Clamp value to range
	g.value = math.Max(g.min, math.Min(g.max, value))

	// Calculate percentage
	percentage := (g.value - g.min) / (g.max - g.min)
	filledWidth := int(float64(g.width) * percentage)

	// Choose color based on value range
	color := "green"
	if g.value < g.warningMin || g.value > g.warningMax {
		color = "red"
	} else if math.Abs(g.value-g.warningMin) < (g.warningMax-g.warningMin)*0.2 ||
		math.Abs(g.value-g.warningMax) < (g.warningMax-g.warningMin)*0.2 {
		color = "yellow"
	}

	// Create gauge text
	g.View.Clear()
	fmt.Fprintf(g.View, "[%s]", color)
	for i := 0; i < filledWidth; i++ {
		fmt.Fprintf(g.View, "█")
	}
	fmt.Fprintf(g.View, "[-:-]")
	for i := filledWidth; i < g.width; i++ {
		fmt.Fprintf(g.View, "░")
	}

	// Add value text
	fmt.Fprintf(g.View, "\n\nValue: %.*f", g.precision, g.value)

	// Add visual scale markers
	g.drawScale()
}

// drawScale adds scale markers to the gauge
func (g *Gauge) drawScale() {
	fmt.Fprintf(g.View, "\n")

	// Draw scale markers
	numMarkers := 5
	for i := 0; i <= numMarkers; i++ {
		position := int(float64(g.width) * float64(i) / float64(numMarkers))
		value := g.min + (g.max-g.min)*float64(i)/float64(numMarkers)

		// Print spaces up to position
		for j := 0; j < position; j++ {
			if j == position-1 {
				fmt.Fprintf(g.View, "|")
			} else if j == 0 || j == g.width-1 {
				fmt.Fprintf(g.View, "|")
			} else {
				fmt.Fprintf(g.View, " ")
			}
		}

		// Print value under marker
		fmt.Fprintf(g.View, "\n%*.*f", position, g.precision, value)
	}
}

// SetPrecision sets the decimal precision for displayed values
func (g *Gauge) SetPrecision(precision int) {
	g.precision = precision
	g.SetValue(g.value) // Redraw
}
