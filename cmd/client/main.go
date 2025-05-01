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
	engineData     *pb.EngineData
	gauges         map[string]*ui.Gauge
	layouts        *ui.Layouts
	dataUpdateTime time.Time
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
	c.layouts.Dashboard.AddItem(c.gauges["rpm"].View, 0, 1, true)

	// Create throttle gauge
	c.gauges["throttle"] = ui.NewGauge("Throttle", 0, 100, 0, 0)
	c.layouts.Dashboard.AddItem(c.gauges["throttle"].View, 0, 1, false)

	// Create engine temperature gauge
	c.gauges["engineTemp"] = ui.NewGauge("Engine Temp (°C)", 20, 120, 90, 90)
	c.layouts.Dashboard.AddItem(c.gauges["engineTemp"].View, 0, 1, false)

	// Create AFR gauge
	c.gauges["afr"] = ui.NewGauge("AFR", 10, 16, 14.7, 14.7)
	c.layouts.Dashboard.AddItem(c.gauges["afr"].View, 0, 1, false)

	// Create power/torque panel
	powerPanel := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	powerPanel.SetBorder(true).SetTitle("Power & Torque")
	c.layouts.InfoPanel.AddItem(powerPanel, 0, 1, false)

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

[yellow]Navigation:[-]
[green]Tab[-]: Switch between views
[green]Q[-]: Quit

[yellow]Map Editing:[-]
[green]M[-]: Switch to map view (coming soon)
[green]E[-]: Edit selected map cell (coming soon)
`)
	c.layouts.InfoPanel.AddItem(controlsPanel, 0, 1, false)

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
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				c.app.Stop()
				return nil
			case 'd', 'D':
				c.layouts.SwitchToPage("dashboard")
				return nil
			case 'm', 'M':
				c.layouts.SwitchToPage("maps")
				return nil
			case 'i', 'I':
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
		})
	}
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
			c.gauges["engineTemp"].SetValue(data.EngineTemp)
			c.gauges["afr"].SetValue(data.AfrCurrent)
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
