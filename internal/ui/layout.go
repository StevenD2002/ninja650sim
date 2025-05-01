package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// Layouts contains the UI layout components
type Layouts struct {
	Root      *tview.Flex
	Pages     *tview.Pages
	Dashboard *tview.Flex
	InfoPanel *tview.Flex
	MapPanel  *tview.Flex
	StatusBar *tview.TextView
}

// NewLayouts creates the UI layout structure
func NewLayouts() *Layouts {
	layouts := &Layouts{
		Root:      tview.NewFlex().SetDirection(tview.FlexRow),
		Pages:     tview.NewPages(),
		Dashboard: tview.NewFlex().SetDirection(tview.FlexRow),
		InfoPanel: tview.NewFlex().SetDirection(tview.FlexRow),
		MapPanel:  tview.NewFlex().SetDirection(tview.FlexRow),
		StatusBar: tview.NewTextView().SetDynamicColors(true),
	}

	// Setup status bar
	layouts.StatusBar.SetTextAlign(tview.AlignCenter)
	layouts.StatusBar.SetText("[yellow]Ninja 650 ECU Simulator[white] | [green]Connected to localhost:50051[white] | Press [blue]Q[white] to quit")

	// Create pages for different screens
	layouts.Pages.AddPage("dashboard", layouts.Dashboard, true, true)
	layouts.Pages.AddPage("maps", layouts.MapPanel, true, false)
	layouts.Pages.AddPage("info", layouts.InfoPanel, true, false)

	// Compose root layout
	layouts.Root.AddItem(layouts.Pages, 0, 1, true)
	layouts.Root.AddItem(layouts.StatusBar, 1, 0, false)

	return layouts
}

// SwitchToPage changes the active page
// SwitchToPage changes the active page
func (l *Layouts) SwitchToPage(name string) {
	l.Pages.SwitchToPage(name)
	// Update status bar to show current page
	l.StatusBar.SetText(fmt.Sprintf("[yellow]Ninja 650 ECU Simulator[white] | [blue]%s View[white] | Press [green]Q[white] to quit", strings.Title(name)))
}

// GetCurrentPage returns the name of the current page
func (l *Layouts) GetCurrentPage() string {
	frontPageName, _ := l.Pages.GetFrontPage()
	if frontPageName != "" {
		return frontPageName
	}
	return "dashboard" // Default
}
