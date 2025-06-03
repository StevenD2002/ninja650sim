package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// WebSocket message types, just matches the protobuf definitions
type WSUserInput struct {
	ThrottlePosition float64 `json:"throttle_position"`
	ClutchPosition   float64 `json:"clutch_position"`
	Gear             int     `json:"gear"`
}

type WSEngineData struct {
	RPM              float64 `json:"rpm"`
	ThrottlePosition float64 `json:"throttle_position"`
	Timestamp        int64   `json:"timestamp"`
	Power            float64 `json:"power"`
	Torque           float64 `json:"torque"`
	Speed            float64 `json:"speed"`
	EngineTemp       float64 `json:"engine_temp"`
	AFRCurrent       float64 `json:"afr_current"`
	AFRTarget        float64 `json:"afr_target"`
	FuelInjectionMs  float64 `json:"fuel_injection_ms"`
	IgnitionAdvance  float64 `json:"ignition_advance"`
	Gear             int     `json:"gear"`
	ClutchPosition   float64 `json:"clutch_position"`
}

// WebSocket client connection
type WSClient struct {
	conn   *websocket.Conn
	server *server
	send   chan WSEngineData
	input  chan WSUserInput
}

// Handle WebSocket connections
func (s *server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket client connected")

	// Create client
	client := &WSClient{
		conn:   conn,
		server: s,
		send:   make(chan WSEngineData, 10),
		input:  make(chan WSUserInput, 10),
	}

	// Start goroutines for handling messages
	go client.readMessages()
	go client.writeMessages()
	go client.simulationLoop()

	// Keep connection alive
	select {}
}

// Read messages from WebSocket client
func (c *WSClient) readMessages() {
	defer c.conn.Close()

	for {
		var input WSUserInput
		err := c.conn.ReadJSON(&input)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Send input to simulation
		select {
		case c.input <- input:
		default:
			// Channel full, skip this input
		}
	}
}

// Write messages to WebSocket client
func (c *WSClient) writeMessages() {
	defer c.conn.Close()

	for {
		select {
		case data := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteJSON(data); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// Simulation loop for WebSocket client
func (c *WSClient) simulationLoop() {
	ticker := time.NewTicker(50 * time.Millisecond) // 20Hz
	defer ticker.Stop()

	for {
		select {
		case input := <-c.input:
			// Apply user input to engine
			c.server.engine.SetThrottle(input.ThrottlePosition)
			c.server.engine.ClutchPosition = input.ClutchPosition
			c.server.engine.Gear = input.Gear

			log.Printf("WS Input - Throttle: %.1f%%, Clutch: %.2f, Gear: %d",
				input.ThrottlePosition, input.ClutchPosition, input.Gear)

		case <-ticker.C:
			// Update simulation
			sensorData := c.server.engine.GetSensorData()
			ecuOutputs := c.server.ecu.ProcessSensorData(sensorData)
			c.server.engine.Update(ecuOutputs, 0.05)

			// Calculate performance
			power, torque := c.server.engine.CalculatePerformance()

			// Create WebSocket response (convert from protobuf format)
			wsData := WSEngineData{
				RPM:              c.server.engine.GetRPM(),
				ThrottlePosition: c.server.engine.GetThrottlePosition(),
				Timestamp:        time.Now().UnixNano(),
				Power:            power,
				Torque:           torque,
				Speed:            sensorData.Speed,
				EngineTemp:       sensorData.EngineTemperature,
				AFRCurrent:       sensorData.O2 * 14.7,
				AFRTarget:        ecuOutputs.LambdaTarget * 14.7,
				FuelInjectionMs:  ecuOutputs.FuelInjectionTime,
				IgnitionAdvance:  ecuOutputs.IgnitionAdvance,
				Gear:             c.server.engine.Gear,
				ClutchPosition:   c.server.engine.ClutchPosition,
			}

			// Send to client
			select {
			case c.send <- wsData:
			default:
				// Channel full, skip this update
			}
		}
	}
}

// Add this to your main() function:
func SetupWebSocketServer(s *server) {
	// WebSocket endpoint
	http.HandleFunc("/ws", s.handleWebSocket)

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Starting WebSocket server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()
}
