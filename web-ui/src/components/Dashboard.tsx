import React, { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import {
  setThrottlePosition,
  toggleClutch,
  shiftUp,
  shiftDown,
  shiftToNeutral,
  clearStatusMessage,
} from '../store/motorcycleSlice';
import {
  connectWebSocket,
  disconnectWebSocket,
  sendUserInput,
} from '../services/websocketMiddleware';

const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  
  // Get state from Redux store
  const {
    engineData,
    transmission,
    ui,
    userInput,
  } = useAppSelector((state: any) => state.motorcycle);

  // Initialize WebSocket connection on mount
  useEffect(() => {
    dispatch(connectWebSocket());
    
    // Cleanup on unmount
    return () => {
      dispatch(disconnectWebSocket());
    };
  }, [dispatch]);

  // Send user input whenever relevant state changes
  useEffect(() => {
    if (ui.connected) {
      dispatch(sendUserInput());
    }
  }, [
    userInput.throttlePos,
    transmission.clutchPosition,
    transmission.currentGear,
    ui.connected,
    dispatch,
  ]);

  // Keyboard controls
  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      switch (event.key.toLowerCase()) {
        case 'arrowleft':
          event.preventDefault();
          dispatch(setThrottlePosition(userInput.throttlePos - 5));
          break;
        case 'arrowright':
          event.preventDefault();
          dispatch(setThrottlePosition(userInput.throttlePos + 5));
          break;
        case 'c':
          event.preventDefault();
          dispatch(toggleClutch());
          break;
        case 'u':
          event.preventDefault();
          dispatch(shiftUp());
          break;
        case 'd':
          event.preventDefault();
          dispatch(shiftDown());
          break;
        case 'n':
          event.preventDefault();
          dispatch(shiftToNeutral());
          break;
        case '0': case '1': case '2': case '3': case '4':
        case '5': case '6': case '7': case '8': case '9':
          event.preventDefault();
          const digit = parseInt(event.key);
          if (event.shiftKey && digit === 1) {
            dispatch(setThrottlePosition(100)); // Shift+1 = 100%
          } else {
            dispatch(setThrottlePosition(digit * 10));
          }
          break;
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [dispatch, userInput.throttlePos]);

  // Clear status message after 3 seconds
  useEffect(() => {
    if (ui.statusMessage) {
      const timer = setTimeout(() => {
        if (Date.now() - ui.statusMessageTime >= 3000) {
          dispatch(clearStatusMessage());
        }
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [ui.statusMessage, ui.statusMessageTime, dispatch]);

  // Helper functions
  const getGearText = () => transmission.currentGear === 0 ? 'N' : transmission.currentGear.toString();
  const getGearColor = () => {
    if (engineData.rpm > 10000) return '#ef4444';
    if (engineData.rpm > 8000) return '#eab308';
    return '#10b981';
  };
  const getClutchStatus = () => {
    if (transmission.clutchPosition > 0.8) return { text: 'DISENGAGED', color: '#ef4444' };
    if (transmission.clutchPosition > 0.2) return { text: 'SLIPPING', color: '#eab308' };
    return { text: 'ENGAGED', color: '#10b981' };
  };

  const clutchStatus = getClutchStatus();
  const lastUpdate = new Date(ui.lastDataUpdate).toLocaleTimeString();

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
            backgroundColor: ui.connected ? '#10b981' : '#ef4444'
          }}></div>
          <span style={{ fontSize: '14px', color: '#9ca3af' }}>
            {ui.connected ? 'Connected to Go Server' : 'Disconnected'}
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
          {ui.connected && engineData.timestamp > 0 ? (
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
                    {lastUpdate}
                  </span>
                </div>
              </div>
            </div>
          ) : (
            <div style={{ textAlign: 'center', color: '#6b7280' }}>
              {ui.connected ? 'Waiting for data...' : 'Not connected'}
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
              Throttle: {userInput.throttlePos.toFixed(0)}%
            </label>
            <input
              type="range"
              min="0"
              max="100"
              value={userInput.throttlePos}
              onChange={(e) => dispatch(setThrottlePosition(Number(e.target.value)))}
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
              onClick={() => dispatch(toggleClutch())}
              style={{
                padding: '8px 16px',
                borderRadius: '4px',
                border: 'none',
                fontWeight: '600',
                backgroundColor: transmission.clutchPressed ? '#dc2626' : '#16a34a',
                color: 'white',
                cursor: 'pointer'
              }}
            >
              Clutch
            </button>
            <button
              onClick={() => dispatch(shiftToNeutral())}
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
              onClick={() => dispatch(shiftUp())}
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
              onClick={() => dispatch(shiftDown())}
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

        {/* Redux Debug Panel */}
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
            Redux State Debug
          </h3>
          <div style={{ fontSize: '12px', color: '#9ca3af' }}>
            <div><strong>Connection:</strong> {ui.connected ? '✅ Connected' : '❌ Disconnected'}</div>
            <div><strong>User Input:</strong></div>
            <div style={{ marginLeft: '10px' }}>
              Throttle: {userInput.throttlePos}%<br/>
              Clutch: {transmission.clutchPosition.toFixed(2)}<br/>
              Gear: {transmission.currentGear}
            </div>
            {engineData.timestamp > 0 && (
              <div style={{ marginTop: '10px' }}>
                <strong>Engine Data:</strong><br/>
                <div style={{ marginLeft: '10px' }}>
                  RPM: {engineData.rpm.toFixed(0)}<br/>
                  Power: {engineData.power.toFixed(1)} hp<br/>
                  Speed: {engineData.speed.toFixed(1)} km/h
                </div>
              </div>
            )}
            <div style={{ marginTop: '10px' }}>
              <strong>Status:</strong> {ui.statusMessage || 'None'}
            </div>
          </div>
        </div>
      </div>

      {/* Status Message */}
      {ui.statusMessage && Date.now() - ui.statusMessageTime < 3000 && (
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
          {ui.statusMessage}
        </div>
      )}
    </div>
  );
};

export default Dashboard;
