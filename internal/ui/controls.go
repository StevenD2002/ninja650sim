package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Control represents an interactive UI control
type Control struct {
	View    tview.Primitive
	OnFocus func()
	OnBlur  func()
}

// MapEditor represents a control for editing ECU maps
type MapEditor struct {
	Control
	Table       *tview.Table
	RPMValues   []float64
	LoadValues  []float64
	Values      [][]float64
	SelectedRow int
	SelectedCol int
	MapType     string
	OnCellEdit  func(row, col int, value float64) bool
}

// NewMapEditor creates a new map editor control
func NewMapEditor(mapType string) *MapEditor {
	table := tview.NewTable().SetBorders(true)

	editor := &MapEditor{
		Control: Control{
			View: table,
		},
		Table:       table,
		SelectedRow: 0,
		SelectedCol: 0,
		MapType:     mapType,
	}

	// Configure table
	table.SetBorder(true)
	table.SetTitle(mapType + " Map Editor")

	// Set up key handling
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			// Edit selected cell
			editor.EditSelectedCell()
			return nil
		}
		return event
	})

	return editor
}

// SetData sets the map data to display
func (m *MapEditor) SetData(rpmValues, loadValues []float64, values [][]float64) {
	m.RPMValues = rpmValues
	m.LoadValues = loadValues
	m.Values = values

	// Clear existing cells
	m.Table.Clear()

	// Add header row (RPM values)
	m.Table.SetCell(0, 0, tview.NewTableCell("Load\\RPM").SetTextColor(tcell.ColorYellow))
	for i, rpm := range rpmValues {
		m.Table.SetCell(0, i+1, tview.NewTableCell(formatFloat(rpm, 0)).SetTextColor(tcell.ColorYellow))
	}

	// Add load column and values
	for i, load := range loadValues {
		m.Table.SetCell(i+1, 0, tview.NewTableCell(formatFloat(load, 0)).SetTextColor(tcell.ColorYellow))

		// Add map values
		for j, value := range values[i] {
			cell := tview.NewTableCell(formatFloat(value, 2))
			m.Table.SetCell(i+1, j+1, cell)
		}
	}

	// Select first cell
	m.Table.Select(1, 1)
}

// EditSelectedCell opens a form to edit the selected cell
func (m *MapEditor) EditSelectedCell() {
	row, col := m.Table.GetSelection()
	if row == 0 || col == 0 {
		// Don't edit headers
		return
	}

	// Get current value
	// dataRow, dataCol := row-1, col-1
	// currentValue := m.Values[dataRow][dataCol]

	// Create modal
	// modal := tview.NewModal().
	// 	SetText(formatFloat(m.RPMValues[dataCol], 0) + " RPM, " + formatFloat(m.LoadValues[dataRow], 0) + "% Load\nCurrent value: " + formatFloat(currentValue, 2)).
	// 	AddButtons([]string{"OK", "Cancel"}).
	// 	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	// 		// Handle button press
	// 	})
	//
	// TODO: Complete editing functionality
}

// Helper function to format float with specified precision
func formatFloat(value float64, precision int) string {
	format := "%." + string(rune('0'+precision)) + "f"
	return fmt.Sprintf(format, value)
}
