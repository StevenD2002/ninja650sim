import React, { useEffect, useMemo } from 'react';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { useKeyboardControls } from '../hooks/useKeyboardControls';
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
import { CircularGauge } from './CircularGauge';
import { GearDisplay } from './GearDisplay';
import { EngineVisualizer } from './EngineVisualizer';

const AdvancedDashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  
  // Get state from Redux store
  const {
    engineData,
    transmission,
    ui,
    userInput,
  } = useAppSelector((state) => state.motorcycle);

  // Initialize keyboard controls
  useKeyboardControls();

  // Initialize WebSocket connection
  useEffect(() => {
    dispatch(connectWebSocket());
    return () => {
      dispatch(disconnectWebSocket());
    };
  }, [dispatch]);

  // Send user input when state changes
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

  // Auto-clear status messages
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

  // Memoized computed values for performance
  const computedValues = useMemo(() => ({
    isHighRpm: engineData.rpm > 8000,
    isRedline: engineData.rpm > 10000,
    isOverheating: engineData.engine_temp > 110,
    powerPercentage: (engineData.power / 67) * 100, // 67hp max for Ninja 650
    lastUpdate: new Date(ui.lastDataUpdate).toLocaleTimeString(),
  }), [engineData.rpm, engineData.engine_temp, engineData.power, ui.lastDataUpdate]);

  return (
    <div
      style={{
        padding: '20px',
        fontFamily: 'monospace',
        background: 'linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)',
        color: 'white',
        minHeight: '100vh',
        position: 'relative',
        overflow: 'hidden',
      }}
    >
      {/* Background grid pattern */}
      <div
        style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundImage: `
            linear-gradient(rgba(59, 130, 246, 0.1) 1px, transparent 1px),
            linear-gradient(90deg, rgba(59, 130, 246, 0.1) 1px, transparent 1px)
          `,
          backgroundSize: '50px 50px',
          animation: ui.connected ? 'gridMove 20s linear infinite' : 'none',
          zIndex: -1,
        }}
      />

      {/* Header with connection status */}
      <div style={{ textAlign: 'center', marginBottom: '30px', position: 'relative' }}>
        <h1
          style={{
            color: '#fbbf24',
            fontSize: '2.5rem',
            margin: '0 0 10px 0',
            textShadow: '0 0 20px #fbbf2444',
            animation: computedValues.isRedline ? 'redlineFlash 0.5s ease-in-out infinite alternate' : 'none',
          }}
        >
          Ninja 650 ECU Simulator
        </h1>
        
        {/* Connection indicator with pulse effect */}
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            gap: '10px',
          }}
        >
          <div
            style={{
              width: '16px',
              height: '16px',
              borderRadius: '50%',
              backgroundColor: ui.connected ? '#10b981' : '#ef4444',
              boxShadow: ui.connected 
                ? '0 0 20px #10b981, inset 0 0 10px rgba(255, 255, 255, 0.2)' 
                : '0 0 20px #ef4444',
              animation: ui.connected ? 'pulse 2s ease-in-out infinite' : 'none',
            }}
          />
          <span style={{ fontSize: '14px', color: '#9ca3af' }}>
            {ui.connected ? 'Connected to Go Server' : 'Disconnected'}
          </span>
        </div>
      </div>

      {/* Main dashboard grid */}
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
          gap: '25px',
          marginBottom: '30px',
        }}
      >
        {/* RPM Gauge */}
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <CircularGauge
            value={engineData.rpm}
            min={0}
            max={11000}
            redline={9000}
            label="RPM"
            size={200}
            color="#3b82f6"
            animate={true}
          />
        </div>

        {/* Speed Gauge */}
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <CircularGauge
            value={engineData.speed}
            min={0}
            max={200}
            label="Speed"
            unit=" km/h"
            size={180}
            color="#10b981"
            animate={true}
          />
        </div>

        {/* Engine Temperature Gauge */}
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <CircularGauge
            value={engineData.engine_temp}
            min={20}
            max={120}
            redline={110}
            label="Engine Temp"
            unit="°C"
            size={160}
            color="#f59e0b"
            animate={true}
          />
        </div>

        {/* Power Gauge */}
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <CircularGauge
            value={engineData.power}
            min={0}
            max={67}
            label="Power"
            unit=" hp"
            size={160}
            color="#8b5cf6"
            animate={true}
          />
        </div>
      </div>

      {/* Central control area */}
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
          gap: '25px',
          marginBottom: '30px',
        }}
      >
        {/* Gear Display */}
        <div
          style={{
            background: 'rgba(30, 41, 59, 0.8)',
            borderRadius: '16px',
            padding: '30px',
            border: '1px solid rgba(59, 130, 246, 0.3)',
            backdropFilter: 'blur(10px)',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '20px',
          }}
        >
          <h3 style={{ color: '#fbbf24', margin: 0, fontSize: '1.5rem' }}>Transmission</h3>
          <GearDisplay
            currentGear={transmission.currentGear}
            rpm={engineData.rpm}
            clutchPressed={transmission.clutchPressed}
            size={140}
          />
          
          {/* Clutch status */}
          <div style={{ textAlign: 'center' }}>
            <div style={{ color: '#9ca3af', fontSize: '14px', marginBottom: '8px' }}>
              Clutch Status
            </div>
            <div
              style={{
                fontSize: '18px',
                fontWeight: 'bold',
                color: transmission.clutchPressed ? '#ef4444' : '#10b981',
                textShadow: `0 0 10px ${transmission.clutchPressed ? '#ef4444' : '#10b981'}`,
              }}
            >
              {transmission.clutchPressed ? 'DISENGAGED' : 'ENGAGED'}
            </div>
          </div>
        </div>

        {/* Engine Visualizer */}
        <div
          style={{
            background: 'rgba(30, 41, 59, 0.8)',
            borderRadius: '16px',
            padding: '30px',
            border: '1px solid rgba(59, 130, 246, 0.3)',
            backdropFilter: 'blur(10px)',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '20px',
          }}
        >
          <h3 style={{ color: '#fbbf24', margin: 0, fontSize: '1.5rem' }}>Engine</h3>
          <EngineVisualizer
            rpm={engineData.rpm}
            throttle={userInput.throttlePos}
            engineTemp={engineData.engine_temp}
            power={engineData.power}
            size={200}
          />
        </div>

        {/* Controls Panel */}
        <div
          style={{
            background: 'rgba(30, 41, 59, 0.8)',
            borderRadius: '16px',
            padding: '30px',
            border: '1px solid rgba(59, 130, 246, 0.3)',
            backdropFilter: 'blur(10px)',
          }}
        >
          <h3 style={{ color: '#fbbf24', marginBottom: '20px', fontSize: '1.5rem', textAlign: 'center' }}>
            Controls
          </h3>

          {/* Throttle Control */}
          <div style={{ marginBottom: '25px' }}>
            <div
              style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: '10px',
              }}
            >
              <label style={{ fontSize: '16px', fontWeight: '500', color: '#d1d5db' }}>
                Throttle
              </label>
              <span
                style={{
                  fontSize: '18px',
                  fontWeight: 'bold',
                  color: '#3b82f6',
                  textShadow: '0 0 10px #3b82f6',
                }}
              >
                {userInput.throttlePos.toFixed(0)}%
              </span>
            </div>
            
            {/* Custom styled range input */}
            <div style={{ position: 'relative' }}>
              <input
                type="range"
                min="0"
                max="100"
                value={userInput.throttlePos}
                onChange={(e) => dispatch(setThrottlePosition(Number(e.target.value)))}
                style={{
                  width: '100%',
                  height: '8px',
                  borderRadius: '4px',
                  background: `linear-gradient(to right, #3b82f6 0%, #3b82f6 ${userInput.throttlePos}%, #374151 ${userInput.throttlePos}%, #374151 100%)`,
                  outline: 'none',
                  appearance: 'none',
                  cursor: 'pointer',
                }}
              />
            </div>
          </div>

          {/* Control Buttons */}
          <div
            style={{
              display: 'grid',
              gridTemplateColumns: '1fr 1fr',
              gap: '12px',
              marginBottom: '20px',
            }}
          >
            {[
              {
                label: 'Clutch',
                onClick: () => dispatch(toggleClutch()),
                color: transmission.clutchPressed ? '#ef4444' : '#10b981',
              },
              {
                label: 'Neutral',
                onClick: () => dispatch(shiftToNeutral()),
                color: '#3b82f6',
              },
              {
                label: 'Shift Up',
                onClick: () => dispatch(shiftUp()),
                color: '#6b7280',
              },
              {
                label: 'Shift Down',
                onClick: () => dispatch(shiftDown()),
                color: '#6b7280',
              },
            ].map((button, index) => (
              <button
                key={index}
                onClick={button.onClick}
                style={{
                  padding: '12px 16px',
                  borderRadius: '8px',
                  border: `2px solid ${button.color}`,
                  background: `rgba(${button.color === '#ef4444' ? '239, 68, 68' : 
                              button.color === '#10b981' ? '16, 185, 129' :
                              button.color === '#3b82f6' ? '59, 130, 246' : '107, 114, 128'}, 0.2)`,
                  color: button.color,
                  fontWeight: '600',
                  cursor: 'pointer',
                  transition: 'all 0.2s ease',
                  textShadow: `0 0 10px ${button.color}`,
                  boxShadow: `0 0 20px ${button.color}33`,
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.transform = 'translateY(-2px)';
                  e.currentTarget.style.boxShadow = `0 4px 20px ${button.color}66`;
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = 'translateY(0)';
                  e.currentTarget.style.boxShadow = `0 0 20px ${button.color}33`;
                }}
              >
                {button.label}
              </button>
            ))}
          </div>

          {/* Keyboard help */}
          <div
            style={{
              fontSize: '12px',
              color: '#6b7280',
              background: 'rgba(0, 0, 0, 0.3)',
              padding: '12px',
              borderRadius: '6px',
              border: '1px solid rgba(107, 114, 128, 0.3)',
            }}
          >
            <div style={{ marginBottom: '8px', color: '#9ca3af', fontWeight: 'bold' }}>
              Keyboard Shortcuts:
            </div>
            <div>← → : Throttle | C: Clutch | U/D: Shift | N: Neutral</div>
            <div>0-9: Set throttle | Shift+1: 100%</div>
          </div>
        </div>

        {/* Engine Data Panel */}
        <div
          style={{
            background: 'rgba(30, 41, 59, 0.8)',
            borderRadius: '16px',
            padding: '30px',
            border: '1px solid rgba(59, 130, 246, 0.3)',
            backdropFilter: 'blur(10px)',
          }}
        >
          <h3 style={{ color: '#fbbf24', marginBottom: '20px', fontSize: '1.5rem', textAlign: 'center' }}>
            Engine Data
          </h3>
          
          {ui.connected && engineData.timestamp > 0 ? (
            <div style={{ display: 'grid', gap: '12px', fontSize: '14px' }}>
              {[
                { label: 'AFR', value: engineData.afr_current.toFixed(1), unit: '' },
                { label: 'Torque', value: engineData.torque.toFixed(1), unit: ' Nm' },
                { label: 'Ignition', value: engineData.ignition_advance.toFixed(1), unit: '° BTDC' },
                { label: 'Fuel Injection', value: engineData.fuel_injection_ms.toFixed(2), unit: ' ms' },
              ].map((item, index) => (
                <div
                  key={index}
                  style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '8px 12px',
                    background: 'rgba(0, 0, 0, 0.3)',
                    borderRadius: '6px',
                    border: '1px solid rgba(59, 130, 246, 0.2)',
                  }}
                >
                  <span style={{ color: '#9ca3af' }}>{item.label}:</span>
                  <span
                    style={{
                      color: '#fbbf24',
                      fontWeight: 'bold',
                      textShadow: '0 0 10px #fbbf24',
                    }}
                  >
                    {item.value}{item.unit}
                  </span>
                </div>
              ))}
              
              <div
                style={{
                  marginTop: '15px',
                  padding: '8px 12px',
                  background: 'rgba(59, 130, 246, 0.1)',
                  borderRadius: '6px',
                  border: '1px solid rgba(59, 130, 246, 0.3)',
                  textAlign: 'center',
                }}
              >
                <div style={{ color: '#6b7280', fontSize: '12px' }}>Last Update</div>
                <div style={{ color: '#60a5fa', fontSize: '14px', fontWeight: 'bold' }}>
                  {computedValues.lastUpdate}
                </div>
              </div>
            </div>
          ) : (
            <div
              style={{
                textAlign: 'center',
                color: '#6b7280',
                padding: '40px 20px',
                background: 'rgba(0, 0, 0, 0.3)',
                borderRadius: '6px',
                border: '1px dashed rgba(107, 114, 128, 0.3)',
              }}
            >
              {ui.connected ? 'Waiting for data...' : 'Not connected to server'}
            </div>
          )}
        </div>
      </div>

      {/* Status message with advanced styling */}
      {ui.statusMessage && Date.now() - ui.statusMessageTime < 3000 && (
        <div
          style={{
            position: 'fixed',
            bottom: '30px',
            left: '50%',
            transform: 'translateX(-50%)',
            padding: '16px 32px',
            borderRadius: '12px',
            fontWeight: '600',
            background: 'rgba(30, 41, 59, 0.95)',
            border: `2px solid ${
              ui.statusMessageColor === 'red' ? '#ef4444' :
              ui.statusMessageColor === 'yellow' ? '#f59e0b' : '#10b981'
            }`,
            color: ui.statusMessageColor === 'red' ? '#ef4444' :
                   ui.statusMessageColor === 'yellow' ? '#f59e0b' : '#10b981',
            backdropFilter: 'blur(10px)',
            boxShadow: `0 8px 32px rgba(0, 0, 0, 0.4), 0 0 20px ${
              ui.statusMessageColor === 'red' ? '#ef444433' :
              ui.statusMessageColor === 'yellow' ? '#f59e0b33' : '#10b98133'
            }`,
            zIndex: 1000,
            animation: 'slideUp 0.3s ease-out',
            textShadow: '0 0 10px currentColor',
          }}
        >
          {ui.statusMessage}
        </div>
      )}

      {/* CSS animations */}
      <style>{`
        @keyframes pulse {
          0%, 100% { opacity: 1; transform: scale(1); }
          50% { opacity: 0.8; transform: scale(1.05); }
        }
        
        @keyframes redlineFlash {
          0% { text-shadow: 0 0 20px #fbbf2444; }
          100% { text-shadow: 0 0 30px #ef4444, 0 0 40px #ef4444; }
        }
        
        @keyframes gridMove {
          0% { transform: translate(0, 0); }
          100% { transform: translate(50px, 50px); }
        }
        
        @keyframes slideUp {
          from {
            transform: translateX(-50%) translateY(100px);
            opacity: 0;
          }
          to {
            transform: translateX(-50%) translateY(0);
            opacity: 1;
          }
        }
        
        input[type="range"]::-webkit-slider-thumb {
          appearance: none;
          height: 20px;
          width: 20px;
          border-radius: 50%;
          background: #3b82f6;
          cursor: pointer;
          box-shadow: 0 0 10px #3b82f6, inset 0 0 5px rgba(255, 255, 255, 0.3);
          border: 2px solid #1e40af;
        }
        
        input[type="range"]::-moz-range-thumb {
          height: 20px;
          width: 20px;
          border-radius: 50%;
          background: #3b82f6;
          cursor: pointer;
          box-shadow: 0 0 10px #3b82f6, inset 0 0 5px rgba(255, 255, 255, 0.3);
          border: 2px solid #1e40af;
        }
      `}</style>
    </div>
  );
};

export default AdvancedDashboard;
