package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/StevenD2002/ninja650sim/internal/ui"
	pb "github.com/StevenD2002/ninja650sim/proto"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents the motorcycle simulator client
type Client struct {
	app            *tview.Application
	stream         pb.MotorcycleSimulator_StreamEngineClient
	ecuClient      pb.MotorcycleSimulatorClient
	conn           *grpc.ClientConn
	ctx            context.Context
	cancel         context.CancelFunc
	throttlePos    float64
	clutchPos      float64
	currentGear    int
	clutchPressed  bool
	engineData     *pb.EngineData
	gauges         map[string]*ui.Gauge
	layouts        *ui.Layouts
	dataUpdateTime time.Time
	statusMsg      string
	statusMsgTime  time.Time
}

// NewClient creates a new client
func NewClient(serverAddr string) (*Client, error) {
	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Connect to the gRPC server
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		return nil, err
	}

	// Create gRPC client
	grpcClient := pb.NewMotorcycleSimulatorClient(conn)

	// Create tview application
	app := tview.NewApplication()

	// Create client
	client := &Client{
		app:            app,
		ecuClient:      grpcClient,
		conn:           conn,
		ctx:            ctx,
		cancel:         cancel,
		throttlePos:    0.0,
		clutchPos:      1.0, // Start with clutch disengaged
		currentGear:    0,   // Start in neutral
		clutchPressed:  true,
		engineData:     &pb.EngineData{},
		gauges:         make(map[string]*ui.Gauge),
		dataUpdateTime: time.Now(),
	}

	// Setup UI
	client.setupUI()

	return client, nil
}

// setupUI initializes the terminal UI components
func (c *Client) setupUI() {
	// Create layouts
	c.layouts = ui.NewLayouts()

	// Create RPM gauge
	c.gauges["rpm"] = ui.NewGauge("RPM", 0, 11000, 0, 900)
	c.gauges["rpm"].SetPrecision(0)
	c.layouts.Dashboard.AddItem(c.gauges["rpm"].View, 0, 1, true)

	// Create throttle gauge
	c.gauges["throttle"] = ui.NewGauge("Throttle", 0, 100, 0, 0)
	c.gauges["throttle"].SetPrecision(0)
	c.layouts.Dashboard.AddItem(c.gauges["throttle"].View, 0, 1, false)

	// Create speed gauge
	c.gauges["speed"] = ui.NewGauge("Speed (km/h)", 0, 200, 0, 0)
	c.gauges["speed"].SetPrecision(1)
	c.layouts.Dashboard.AddItem(c.gauges["speed"].View, 0, 1, false)

	// Create engine temperature gauge
	c.gauges["engineTemp"] = ui.NewGauge("Engine Temp (°C)", 20, 120, 90, 90)
	c.layouts.Dashboard.AddItem(c.gauges["engineTemp"].View, 0, 1, false)

	// Create AFR gauge
	c.gauges["afr"] = ui.NewGauge("AFR", 10, 16, 14.7, 14.7)
	c.gauges["afr"].SetPrecision(1)
	c.layouts.Dashboard.AddItem(c.gauges["afr"].View, 0, 1, false)

	// Create power/torque panel
	powerPanel := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	powerPanel.SetBorder(true).SetTitle("Power & Torque")
	c.layouts.InfoPanel.AddItem(powerPanel, 0, 1, false)

	// Create transmission panel
	transmissionPanel := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	transmissionPanel.SetBorder(true).SetTitle("Transmission")
	c.layouts.Dashboard.AddItem(transmissionPanel, 0, 1, false)

	// Create controls info panel
	controlsPanel := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	controlsPanel.SetBorder(true).SetTitle("Controls")
	controlsPanel.SetText(`
[yellow]Throttle Controls:[-]
[green]Left/Right Arrows[-]: Adjust throttle position
[green]0-9[-]: Set throttle to 0-90%
[green]Shift+1[-]: Set throttle to 100%

[yellow]Transmission Controls:[-]
[green]C[-]: Toggle clutch (press/release)
[green]U[-]: Shift up
[green]D[-]: Shift down
[green]N[-]: Shift to neutral

[yellow]Navigation:[-]
[green]Tab[-]: Switch between views
[green]Q[-]: Quit

[yellow]Map Editing:[-]
[green]M[-]: Switch to map view
[green]E[-]: Edit selected map cell
`)
	c.layouts.InfoPanel.AddItem(controlsPanel, 0, 1, false)

	// Status bar for messages
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusBar.SetBorder(false)
	c.layouts.Root.AddItem(statusBar, 1, 0, false)

	// Update power panel periodically
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if c.engineData != nil {
				c.app.QueueUpdateDraw(func() {
					powerPanel.Clear()
					fmt.Fprintf(powerPanel, "[yellow]Power:[-] %.1f hp\n", c.engineData.Power)
					fmt.Fprintf(powerPanel, "[yellow]Torque:[-] %.1f Nm\n", c.engineData.Torque)
					fmt.Fprintf(powerPanel, "[yellow]Ignition:[-] %.1f° BTDC\n", c.engineData.IgnitionAdvance)
					fmt.Fprintf(powerPanel, "[yellow]Fuel:[-] %.2f ms\n", c.engineData.FuelInjectionMs)
					fmt.Fprintf(powerPanel, "\n[blue]Last Update:[-] %s", c.dataUpdateTime.Format("15:04:05.000"))
				})
			}
		}
	}()

	// Update transmission panel periodically
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			c.app.QueueUpdateDraw(func() {
				transmissionPanel.Clear()

				// Display gear
				gearText := "N"
				if c.currentGear > 0 {
					gearText = fmt.Sprintf("%d", c.currentGear)
				}

				// Colorize gear text
				gearColor := "green"
				if c.engineData != nil && c.engineData.Rpm > 10000 {
					gearColor = "red"
				} else if c.engineData != nil && c.engineData.Rpm > 8000 {
					gearColor = "yellow"
				}

				fmt.Fprintf(transmissionPanel, "\n[white]Current Gear: [%s]%s[-]\n\n", gearColor, gearText)

				// Display clutch status
				clutchText := "ENGAGED"
				clutchColor := "green"
				if c.clutchPos > 0.8 {
					clutchText = "DISENGAGED"
					clutchColor = "red"
				} else if c.clutchPos > 0.2 {
					clutchText = "SLIPPING"
					clutchColor = "yellow"
				}

				fmt.Fprintf(transmissionPanel, "[white]Clutch: [%s]%s[-]", clutchColor, clutchText)
			})
		}
	}()

	// Update status bar periodically
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			c.app.QueueUpdateDraw(func() {
				statusBar.Clear()

				// Show status message if it's recent (less than 3 seconds old)
				if c.statusMsg != "" && time.Since(c.statusMsgTime) < 3*time.Second {
					statusBar.SetText(c.statusMsg)
				} else {
					// Default status bar
					currentPage := c.layouts.GetCurrentPage()
					statusBar.SetText(fmt.Sprintf("[yellow]Ninja 650 ECU Simulator[white] | [blue]%s View[white] | Press [green]Q[white] to quit",
						title(currentPage)))
				}
			})
		}
	}()

	// Set root and handle input
	c.app.SetRoot(c.layouts.Root, true)
	c.setupInputHandling()
}

// setupInputHandling configures keyboard input handlers
func (c *Client) setupInputHandling() {
	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Handle global keys
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			c.app.Stop()
			return nil
		case tcell.KeyTab:
			// Cycle through pages
			currentPage := c.layouts.GetCurrentPage()
			switch currentPage {
			case "dashboard":
				c.layouts.SwitchToPage("maps")
			case "maps":
				c.layouts.SwitchToPage("info")
			case "info":
				c.layouts.SwitchToPage("dashboard")
			default:
				c.layouts.SwitchToPage("dashboard")
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				c.app.Stop()
				return nil
			case 'm':
				c.layouts.SwitchToPage("maps")
				return nil
			case 'i':
				c.layouts.SwitchToPage("info")
				return nil
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// Set throttle position (0-9 = 0-90%)
				if event.Modifiers()&tcell.ModShift != 0 && event.Rune() == '1' {
					// Shift+1 = 100%
					c.setThrottle(100.0)
				} else {
					digit := int(event.Rune() - '0')
					c.setThrottle(float64(digit * 10))
				}
				return nil
			case 'c', 'C':
				// Toggle clutch
				if c.clutchPressed {
					c.setClutch(0.0) // Release clutch
					c.clutchPressed = false
					c.showStatusMessage("Clutch ENGAGED", "green")
				} else {
					c.setClutch(1.0) // Press clutch
					c.clutchPressed = true
					c.showStatusMessage("Clutch DISENGAGED", "red")
				}
				return nil
			case 'u', 'U':
				// Shift up
				c.shiftUp()
				return nil
			case 'd', 'D':
				// Shift down
				c.shiftDown()
				return nil
			case 'n', 'N':
				// Shift to neutral
				c.shiftToNeutral()
				return nil
			}
		case tcell.KeyLeft:
			// Decrease throttle
			c.setThrottle(math.Max(0, c.throttlePos-5))
			return nil
		case tcell.KeyRight:
			// Increase throttle
			c.setThrottle(math.Min(100, c.throttlePos+5))
			return nil
		}
		return event
	})
}

// setThrottle changes the throttle position and sends it to the server
func (c *Client) setThrottle(position float64) {
	c.throttlePos = position

	// Update throttle gauge
	c.gauges["throttle"].SetValue(position)

	// Send to server if stream is active
	if c.stream != nil {
		c.stream.Send(&pb.UserInput{
			ThrottlePosition: position,
			ClutchPosition:   c.clutchPos,
			Gear:             int32(c.currentGear),
		})
	}
}

// setClutch changes the clutch position and sends it to the server
func (c *Client) setClutch(position float64) {
	c.clutchPos = position

	// Send to server if stream is active
	if c.stream != nil {
		c.stream.Send(&pb.UserInput{
			ThrottlePosition: c.throttlePos,
			ClutchPosition:   position,
			Gear:             int32(c.currentGear),
		})
	}
}

// shiftUp shifts to the next higher gear
func (c *Client) shiftUp() {
	// Only shift if clutch is pressed
	if c.clutchPos < 0.8 {
		c.showStatusMessage("Press clutch to shift gears", "red")
		return
	}

	// Don't exceed top gear
	if c.currentGear < 6 {
		newGear := c.currentGear + 1
		c.setGear(newGear)
	} else {
		c.showStatusMessage("Already in top gear", "yellow")
	}
}

// shiftDown shifts to the next lower gear
func (c *Client) shiftDown() {
	// Only shift if clutch is pressed
	if c.clutchPos < 0.8 {
		c.showStatusMessage("Press clutch to shift gears", "red")
		return
	}

	// Don't go below neutral
	if c.currentGear > 0 {
		newGear := c.currentGear - 1
		c.setGear(newGear)
	} else {
		c.showStatusMessage("Already in neutral", "yellow")
	}
}

// shiftToNeutral shifts directly to neutral
func (c *Client) shiftToNeutral() {
	// Only shift if clutch is pressed
	if c.clutchPos < 0.8 {
		c.showStatusMessage("Press clutch to shift gears", "red")
		return
	}

	c.setGear(0)
}

// setGear changes the current gear and sends it to the server
func (c *Client) setGear(gear int) {
	c.currentGear = gear

	// Send to server if stream is active
	if c.stream != nil {
		c.stream.Send(&pb.UserInput{
			ThrottlePosition: c.throttlePos,
			ClutchPosition:   c.clutchPos,
			Gear:             int32(gear),
		})
	}

	// Show gear change message
	if gear == 0 {
		c.showStatusMessage("Shifted to Neutral", "green")
	} else {
		c.showStatusMessage(fmt.Sprintf("Shifted to %d gear", gear), "green")
	}
}

// showStatusMessage displays a message in the status bar
func (c *Client) showStatusMessage(message string, color string) {
	c.statusMsg = fmt.Sprintf("[%s]%s[-]", color, message)
	c.statusMsgTime = time.Now()
}

// Start begins the connection and UI
func (c *Client) Start() error {
	// Start stream
	stream, err := c.ecuClient.StreamEngine(c.ctx)
	if err != nil {
		return err
	}
	c.stream = stream

	// Start receive goroutine
	go c.receiveEngineData()

	// Run the UI
	return c.app.Run()
}

// receiveEngineData handles incoming data from the server
func (c *Client) receiveEngineData() {
	for {
		data, err := c.stream.Recv()
		if err != nil {
			log.Printf("Error receiving data: %v", err)
			c.app.Stop()
			return
		}

		// Update client data
		c.engineData = data
		c.dataUpdateTime = time.Now()

		// Update gauges
		c.app.QueueUpdateDraw(func() {
			c.gauges["rpm"].SetValue(data.Rpm)
			c.gauges["throttle"].SetValue(data.ThrottlePosition)
			if c.gauges["speed"] != nil && data.Speed > 0 {
				c.gauges["speed"].SetValue(data.Speed)
			}
			if c.gauges["engineTemp"] != nil && data.EngineTemp > 0 {
				c.gauges["engineTemp"].SetValue(data.EngineTemp)
			}
			if c.gauges["afr"] != nil && data.AfrCurrent > 0 {
				c.gauges["afr"].SetValue(data.AfrCurrent)
			}
		})
	}
}

// Cleanup closes connections and frees resources
func (c *Client) Cleanup() {
	if c.conn != nil {
		c.conn.Close()
	}
	c.cancel()
}

// Helper function to capitalize the first letter of a string
func title(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

func main() {
	// Create client
	client, err := NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Cleanup()

	// Start client
	log.Println("Starting Motorcycle Simulator client")
	if err := client.Start(); err != nil {
		log.Fatalf("Error running client: %v", err)
	}
}

