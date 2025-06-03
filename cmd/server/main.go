package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/StevenD2002/ninja650sim/internal/ecu"
	"github.com/StevenD2002/ninja650sim/internal/engine"
	pb "github.com/StevenD2002/ninja650sim/proto"
	"google.golang.org/grpc"
)

// Server implements the gRPC MotorcycleSimulator service
type server struct {
	pb.UnimplementedMotorcycleSimulatorServer
	engine *engine.Engine
	ecu    *ecu.ECU

	// Current motorcycle state
	running bool
}

// NewServer creates a new simulator server
func NewServer() *server {
	return &server{
		engine:  engine.NewEngine(),
		ecu:     ecu.NewECU(),
		running: false,
	}
}

// StreamEngine implements the gRPC service method for streaming engine data
func (s *server) StreamEngine(stream pb.MotorcycleSimulator_StreamEngineServer) error {
	// Mark the simulation as running
	s.running = true
	defer func() { s.running = false }()

	// Create ticker for simulation updates
	ticker := time.NewTicker(50 * time.Millisecond) // 20Hz simulation rate
	defer ticker.Stop()

	// Channel for user input
	inputChan := make(chan *pb.UserInput)

	// Start goroutine to receive user input
	go func() {
		for {
			input, err := stream.Recv()
			if err == io.EOF {
				close(inputChan)
				return
			}
			if err != nil {
				log.Printf("Error receiving user input: %v", err)
				return
			}
			inputChan <- input
		}
	}()

	// Main simulation loop
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case input, ok := <-inputChan:
			if ok {
				// Update throttle position
				s.engine.SetThrottle(input.ThrottlePosition)
				// Update clutch position
				s.engine.ClutchPosition = input.ClutchPosition

				// Update gear
				s.engine.Gear = int(input.Gear)

				// For debugging
				log.Printf("Input received - Throttle: %.1f%%, Clutch: %.2f, Gear: %d",
					input.ThrottlePosition, input.ClutchPosition, input.Gear)
			}
		case <-ticker.C:
			// Get sensor data from engine
			sensorData := s.engine.GetSensorData()

			// Process sensor data through ECU
			ecuOutputs := s.ecu.ProcessSensorData(sensorData)

			// Update engine based on ECU outputs
			s.engine.Update(ecuOutputs, 0.05) // 50ms

			log.Printf("Engine state - RPM: %.1f, Speed: %.1f, Throttle: %.1f%%, Gear: %d, Clutch: %.2f",
				s.engine.RPM, s.engine.Speed, s.engine.ThrottlePosition, s.engine.Gear, s.engine.ClutchPosition)

			// Calculate performance metrics
			power, torque := s.engine.CalculatePerformance()

			// Create response message
			response := &pb.EngineData{
				Rpm:              s.engine.GetRPM(),
				ThrottlePosition: s.engine.GetThrottlePosition(),
				Timestamp:        time.Now().UnixNano(),
				// Add additional fields if you've extended your proto definition
				// For example:
				Power:           power,
				Torque:          torque,
				Speed:           sensorData.Speed,
				EngineTemp:      sensorData.EngineTemperature,
				AfrCurrent:      sensorData.O2 * 14.7, // Convert lambda to AFR
				AfrTarget:       ecuOutputs.LambdaTarget * 14.7,
				FuelInjectionMs: ecuOutputs.FuelInjectionTime,
				IgnitionAdvance: ecuOutputs.IgnitionAdvance,
			}

			// Send update to client
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

// GetECUMaps returns the current ECU maps
func (s *server) GetECUMaps(ctx context.Context, req *pb.MapsRequest) (*pb.ECUMaps, error) {
	// Get the current maps from the ECU
	fuelMap := s.ecu.FuelMap
	ignitionMap := s.ecu.IgnitionMap
	afrMap := s.ecu.TargetAFRMap

	// Convert to protobuf format
	response := &pb.ECUMaps{
		FuelMap:     convertMap2DToProto(fuelMap.Map2D, "fuel"),
		IgnitionMap: convertMap2DToProto(ignitionMap.Map2D, "ignition"),
		AfrMap:      convertMap2DToProto(afrMap.Map2D, "afr"),
	}

	return response, nil
}

// UpdateECUMap updates a specific ECU map
func (s *server) UpdateECUMap(ctx context.Context, req *pb.MapUpdateRequest) (*pb.UpdateStatus, error) {
	// Check which map to update
	switch req.MapType {
	case "fuel":
		// Update a single cell in the fuel map
		s.ecu.FuelMap.SetValue(req.Rpm, req.Load, req.Value)
		return &pb.UpdateStatus{Success: true, Message: "Fuel map updated"}, nil
	case "ignition":
		// Update a single cell in the ignition map
		s.ecu.IgnitionMap.SetValue(req.Rpm, req.Load, req.Value)
		return &pb.UpdateStatus{Success: true, Message: "Ignition map updated"}, nil
	case "afr":
		// Update a single cell in the AFR map
		s.ecu.TargetAFRMap.SetValue(req.Rpm, req.Load, req.Value)
		return &pb.UpdateStatus{Success: true, Message: "AFR map updated"}, nil
	default:
		return &pb.UpdateStatus{Success: false, Message: "Unknown map type"}, nil
	}
}

// SetECUSettings updates the ECU settings
func (s *server) SetECUSettings(ctx context.Context, req *pb.ECUSettings) (*pb.UpdateStatus, error) {
	// Update ECU settings
	s.ecu.FuelTrim = req.FuelTrim
	s.ecu.IgnitionTrim = req.IgnitionTrim
	s.ecu.IdleRPM = req.IdleRpm
	s.ecu.RevLimit = req.RevLimit
	s.ecu.TempCompensation = req.TempCompensation

	return &pb.UpdateStatus{Success: true, Message: "ECU settings updated"}, nil
}

// Helper function to convert Map2D to protobuf format
func convertMap2DToProto(m ecu.Map2D, mapType string) *pb.Map2D {
	protoMap := &pb.Map2D{
		Type:            mapType,
		RpmBreakpoints:  m.RPMBreakpoints,
		LoadBreakpoints: m.LoadBreakpoints,
		Values:          make([]*pb.MapRow, len(m.Values)),
	}

	for i, row := range m.Values {
		protoMap.Values[i] = &pb.MapRow{
			Values: row,
		}
	}

	return protoMap
}

func main() {
	// Create a TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()
	simulatorServer := NewServer()

	// Register our implementation
	pb.RegisterMotorcycleSimulatorServer(s, simulatorServer)

	// setup the websocket server
	SetupWebSocketServer(simulatorServer)

	// Start the server
	log.Println("Starting Motorcycle Simulator gRPC server on :50051")
	log.Println("WebSocket server running on :8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
