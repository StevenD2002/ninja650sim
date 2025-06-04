import React, { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { useKeyboardControls } from '../hooks/useKeyboardControls';
import {
  clearStatusMessage,
} from '../store/motorcycleSlice';
import {
  connectWebSocket,
  disconnectWebSocket,
  sendUserInput,
} from '../services/websocketMiddleware';

//TODO: eventually need to add these to an index and export them for cleanliness
import { EngineDataPanel } from './EngineDataPanel';
import { TransmissionPanel } from './TransmissionPanel';
import { ControlsPanel } from './ControlsPanel';
import { ReduxDebugPanel } from './ReduxDebugPanel';
import { StatusMessage } from './StatusMessage';

const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  
  // Get state from Redux store
  const {
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

  useKeyboardControls();

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
        
        <EngineDataPanel />
        <TransmissionPanel />
        <ControlsPanel />
        <ReduxDebugPanel />

      </div>
      <StatusMessage />
    </div>
  );
};

export default Dashboard;
