import { useState, useEffect, useRef } from 'react';
import { WebSocketClient, WSEngineData, WSUserInput } from './services/websocketClient';

function App() {
  // Connection state
  const [connected, setConnected] = useState(false);
  const [engineData, setEngineData] = useState<WSEngineData | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);
  
  // User input state
  const [throttlePos, setThrottlePos] = useState(0);
  const [clutchPos, setClutchPos] = useState(1.0); // Start disengaged
  const [currentGear, setCurrentGear] = useState(0); // Start in neutral
  const [clutchPressed, setClutchPressed] = useState(true);
  
  // Status messages
  const [statusMessage, setStatusMessage] = useState('');
  const [statusTime, setStatusTime] = useState(0);
  
  // WebSocket client ref
  const clientRef = useRef<WebSocketClient | null>(null);

  // Initialize WebSocket connection
  useEffect(() => {
    console.log('🚀 Initializing WebSocket connection...');
    clientRef.current = new WebSocketClient('ws://localhost:8080/ws');
    
    clientRef.current.connect(
      (data: WSEngineData) => {
        console.log('📊 Received engine data:', data);
        setEngineData(data);
        setLastUpdate(new Date());
      },
      (isConnected: boolean) => {
        console.log('🔌 Connection status:', isConnected);
        setConnected(isConnected);
        if (isConnected) {
          showStatusMessage('Connected to server');
        } else {
          showStatusMessage('Disconnected from server');
        }
      }
    );

    // Cleanup on unmount
    return () => {
      console.log('🧹 Cleaning up WebSocket connection');
      clientRef.current?.disconnect();
    };
  }, []);

  // Send user input to server
  const sendUserInput = () => {
    if (clientRef.current && connected) {
      const input: WSUserInput = {
        throttle_position: throttlePos,
        clutch_position: clutchPos,
        gear: currentGear
      };
      console.log('📤 Sending user input:', input);
      clientRef.current.sendUserInput(input);
    }
  };

  // Send input whenever state changes
  useEffect(() => {
    sendUserInput();
  }, [throttlePos, clutchPos, currentGear, connected]);

  // these are just the keboard controls
  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      switch (event.key.toLowerCase()) {
        case 'arrowleft':
          event.preventDefault();
          setThrottlePos(prev => Math.max(0, prev - 5));
          break;
        case 'arrowright':
          event.preventDefault();
          setThrottlePos(prev => Math.min(100, prev + 5));
          break;
        case 'c':
          event.preventDefault();
          toggleClutch();
          break;
        case 'u':
          event.preventDefault();
          shiftUp();
          break;
        case 'd':
          event.preventDefault();
          shiftDown();
          break;
        case 'n':
          event.preventDefault();
          shiftToNeutral();
          break;
        case '0': case '1': case '2': case '3': case '4':
        case '5': case '6': case '7': case '8': case '9':
          event.preventDefault();
          const digit = parseInt(event.key);
          if (event.shiftKey && digit === 1) {
            setThrottlePos(100); // Shift+1 = 100%
          } else {
            setThrottlePos(digit * 10);
          }
          break;
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [clutchPos, currentGear]);

  // Control functions 
  const toggleClutch = () => {
    if (clutchPressed) {
      setClutchPos(0.0); // Release clutch
      setClutchPressed(false);
      showStatusMessage('Clutch ENGAGED');
    } else {
      setClutchPos(1.0); // Press clutch
      setClutchPressed(true);
      showStatusMessage('Clutch DISENGAGED');
    }
  };

  const shiftUp = () => {
    if (clutchPos < 0.8) {
      showStatusMessage('Press clutch to shift gears');
      return;
    }
    if (currentGear < 6) {
      const newGear = currentGear + 1;
      setCurrentGear(newGear);
      showStatusMessage(`Shifted to ${newGear} gear`);
    } else {
      showStatusMessage('Already in top gear');
    }
  };

  const shiftDown = () => {
    if (clutchPos < 0.8) {
      showStatusMessage('Press clutch to shift gears');
      return;
    }
    if (currentGear > 0) {
      const newGear = currentGear - 1;
      setCurrentGear(newGear);
      const message = newGear === 0 ? 'Shifted to Neutral' : `Shifted to ${newGear} gear`;
      showStatusMessage(message);
    } else {
      showStatusMessage('Already in neutral');
    }
  };

  const shiftToNeutral = () => {
    if (clutchPos < 0.8) {
      showStatusMessage('Press clutch to shift gears');
      return;
    }
    setCurrentGear(0);
    showStatusMessage('Shifted to Neutral');
  };

  const showStatusMessage = (message: string) => {
    setStatusMessage(message);
    setStatusTime(Date.now());
  };

  // Clear status message after 3 seconds
  useEffect(() => {
    if (statusMessage) {
      const timer = setTimeout(() => {
        if (Date.now() - statusTime >= 3000) {
          setStatusMessage('');
        }
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [statusMessage, statusTime]);

  // Helper functions
  const getGearText = () => currentGear === 0 ? 'N' : currentGear.toString();
  const getGearColor = () => {
    if (engineData?.rpm && engineData.rpm > 10000) return '#ef4444';
    if (engineData?.rpm && engineData.rpm > 8000) return '#eab308';
    return '#10b981';
  };
  const getClutchStatus = () => {
    if (clutchPos > 0.8) return { text: 'DISENGAGED', color: '#ef4444' };
    if (clutchPos > 0.2) return { text: 'SLIPPING', color: '#eab308' };
    return { text: 'ENGAGED', color: '#10b981' };
  };

  const clutchStatus = getClutchStatus();

  return (
    <div style={{ 
      padding: '20px', 
      fontFamily: 'monospace', 
      backgroundColor: '#1f2937', 
      color: 'white', 
      minHeight: '100vh' 
    }}>
      <div style={{ textAlign: 'center', marginBottom: '30px' }}>
        <h1 style={{ 
          color: '#fbbf24', 
          fontSize: '2rem', 
          margin: '0 0 10px 0' 
        }}>
          Ninja 650 ECU Simulator
        </h1>
        <div style={{ 
          display: 'flex', 
          alignItems: 'center', 
          justifyContent: 'center', 
          gap: '10px' 
        }}>
          <div style={{
            width: '12px',
            height: '12px',
            borderRadius: '50%',
            backgroundColor: connected ? '#10b981' : '#ef4444'
          }}></div>
          <span style={{ fontSize: '14px', color: '#9ca3af' }}>
            {connected ? 'Connected to Go Server' : 'Disconnected'}
          </span>
        </div>
      </div>

      <div style={{ 
        display: 'grid', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', 
        gap: '20px', 
        marginBottom: '20px' 
      }}>
        
        {/* Engine Data Panel */}
        <div style={{ 
          backgroundColor: '#374151', 
          padding: '20px', 
          borderRadius: '8px', 
          border: '1px solid #4b5563' 
        }}>
          <h3 style={{ 
            color: '#fbbf24', 
            textAlign: 'center', 
            marginBottom: '15px' 
          }}>
            Engine Data
          </h3>
          {engineData ? (
            <div style={{ display: 'grid', gap: '8px', fontSize: '14px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>RPM:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.rpm.toFixed(0)}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Power:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.power.toFixed(1)} hp
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Torque:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.torque.toFixed(1)} Nm
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Speed:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.speed.toFixed(1)} km/h
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Engine Temp:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.engine_temp.toFixed(0)}°C
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>AFR:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.afr_current.toFixed(1)}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Ignition:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.ignition_advance.toFixed(1)}° BTDC
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Fuel:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.fuel_injection_ms.toFixed(2)} ms
                </span>
              </div>
              <div style={{ 
                borderTop: '1px solid #4b5563', 
                paddingTop: '8px', 
                marginTop: '8px' 
              }}>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: '#6b7280', fontSize: '12px' }}>Last Update:</span>
                  <span style={{ color: '#60a5fa', fontSize: '12px' }}>
                    {lastUpdate?.toLocaleTimeString()}
                  </span>
                </div>
              </div>
            </div>
          ) : (
            <div style={{ textAlign: 'center', color: '#6b7280' }}>
              {connected ? 'Waiting for data...' : 'Not connected'}
            </div>
          )}
        </div>

        {/* Transmission Panel */}
        <div style={{ 
          backgroundColor: '#374151', 
          padding: '20px', 
          borderRadius: '8px', 
          border: '1px solid #4b5563' 
        }}>
          <h3 style={{ 
            color: '#fbbf24', 
            textAlign: 'center', 
            marginBottom: '15px' 
          }}>
            Transmission
          </h3>
          <div style={{ textAlign: 'center' }}>
            <div style={{ marginBottom: '20px' }}>
              <div style={{ color: '#9ca3af', marginBottom: '8px' }}>Current Gear:</div>
              <div style={{ 
                fontSize: '3rem', 
                fontWeight: 'bold', 
                color: getGearColor() 
              }}>
                {getGearText()}
              </div>
            </div>
            <div>
              <div style={{ color: '#9ca3af', marginBottom: '8px' }}>Clutch:</div>
              <div style={{ 
                fontSize: '1.25rem', 
                fontWeight: 'bold', 
                color: clutchStatus.color 
              }}>
                {clutchStatus.text}
              </div>
            </div>
          </div>
        </div>

        {/* Controls Panel */}
        <div style={{ 
          backgroundColor: '#374151', 
          padding: '20px', 
          borderRadius: '8px', 
          border: '1px solid #4b5563' 
        }}>
          <h3 style={{ 
            color: '#fbbf24', 
            textAlign: 'center', 
            marginBottom: '15px' 
          }}>
            Controls
          </h3>
          
          {/* Throttle Slider */}
          <div style={{ marginBottom: '20px' }}>
            <label style={{ 
              display: 'block', 
              fontSize: '14px', 
              fontWeight: '500', 
              color: '#d1d5db', 
              marginBottom: '8px' 
            }}>
              Throttle: {throttlePos.toFixed(0)}%
            </label>
            <input
              type="range"
              min="0"
              max="100"
              value={throttlePos}
              onChange={(e) => setThrottlePos(Number(e.target.value))}
              style={{ width: '100%' }}
            />
          </div>

          {/* Control Buttons */}
          <div style={{ 
            display: 'grid', 
            gridTemplateColumns: '1fr 1fr', 
            gap: '8px' 
          }}>
            <button
              onClick={toggleClutch}
              style={{
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                fontWeight: '600',
                backgroundColor: clutchPressed ? '#dc2626' : '#16a34a',
                color: 'white',
                cursor: 'pointer'
              }}
            >
              Clutch
            </button>
            <button
              onClick={shiftToNeutral}
              style={{
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                fontWeight: '600',
                backgroundColor: '#2563eb',
                color: 'white',
                cursor: 'pointer'
              }}
            >
              Neutral
            </button>
            <button
              onClick={shiftUp}
              style={{
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                fontWeight: '600',
                backgroundColor: '#4b5563',
                color: 'white',
                cursor: 'pointer'
              }}
            >
              Shift Up
            </button>
            <button
              onClick={shiftDown}
              style={{
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                fontWeight: '600',
                backgroundColor: '#4b5563',
                color: 'white',
                cursor: 'pointer'
              }}
            >
              Shift Down
            </button>
          </div>

          {/* Keyboard Help */}
          <div style={{ 
            marginTop: '20px', 
            fontSize: '12px', 
            color: '#6b7280' 
          }}>
            <div><strong>Keyboard Controls:</strong></div>
            <div>Left/Right: Throttle ±5%</div>
            <div>C: Toggle Clutch</div>
            <div>U: Shift Up | D: Shift Down</div>
            <div>N: Neutral | 0-9: Set throttle</div>
            <div>Shift+1: 100% throttle</div>
          </div>
        </div>

        {/* Connection Debug Panel */}
        <div style={{ 
          backgroundColor: '#374151', 
          padding: '20px', 
          borderRadius: '8px', 
          border: '1px solid #4b5563' 
        }}>
          <h3 style={{ 
            color: '#fbbf24', 
            textAlign: 'center', 
            marginBottom: '15px' 
          }}>
            Debug Info
          </h3>
          <div style={{ fontSize: '12px', color: '#9ca3af' }}>
            <div><strong>WebSocket URL:</strong> ws://localhost:8080/ws</div>
            <div><strong>Status:</strong> {connected ? '✅ Connected' : '❌ Disconnected'}</div>
            <div><strong>Current Input:</strong></div>
            <div style={{ marginLeft: '10px' }}>
              Throttle: {throttlePos}%<br/>
              Clutch: {clutchPos.toFixed(2)}<br/>
              Gear: {currentGear}
            </div>
            {engineData && (
              <div style={{ marginTop: '10px' }}>
                <strong>Last Data Received:</strong><br/>
                <div style={{ marginLeft: '10px' }}>
                  RPM: {engineData.rpm.toFixed(0)}<br/>
                  Power: {engineData.power.toFixed(1)} hp
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Status Message */}
      {statusMessage && Date.now() - statusTime < 3000 && (
        <div style={{
          position: 'fixed',
          bottom: '20px',
          left: '50%',
          transform: 'translateX(-50%)',
          padding: '12px 24px',
          borderRadius: '8px',
          fontWeight: '600',
          backgroundColor: '#374151',
          border: '1px solid #4b5563',
          color: 'white',
          zIndex: 1000
        }}>
          {statusMessage}
        </div>
      )}
    </div>
  );
}

export default App;
