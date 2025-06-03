import { Middleware } from '@reduxjs/toolkit';
import { WebSocketClient, WSUserInput } from './websocketClient';
import { updateEngineData, setConnectionStatus } from '../store/motorcycleSlice';
import type { RootState } from '../store';

// WebSocket middleware that handles connection and message sending
// This way, each component can just focus on ui and not have to manage connecting to the 
// web socket. Its a shared connection
export const createWebSocketMiddleware = (): Middleware => {
  let client: WebSocketClient | null = null;
  
  return (store) => (next) => (action) => {
    switch (action.type) {
      case 'websocket/connect':
        // Initialize WebSocket connection
        if (!client) {
          client = new WebSocketClient('ws://localhost:8080/ws');
          
          client.connect(
            (data) => {
              // Dispatch engine data updates
              store.dispatch(updateEngineData(data));
            },
            (connected) => {
              // Dispatch connection status updates
              store.dispatch(setConnectionStatus(connected));
            }
          );
        }
        break;
        
      case 'websocket/disconnect':
        // Disconnect WebSocket
        if (client) {
          client.disconnect();
          client = null;
        }
        break;
        
      case 'websocket/sendInput':
        // Send user input to server
        if (client && client.isConnected()) {
          const state = store.getState() as RootState;
          const input: WSUserInput = {
            throttle_position: state.motorcycle.userInput.throttlePos,
            clutch_position: state.motorcycle.transmission.clutchPosition,
            gear: state.motorcycle.transmission.currentGear,
          };
          client.sendUserInput(input);
        }
        break;
    }
    
    return next(action);
  };
};

// Action creators for WebSocket middleware
export const connectWebSocket = () => ({ type: 'websocket/connect' });
export const disconnectWebSocket = () => ({ type: 'websocket/disconnect' });
export const sendUserInput = () => ({ type: 'websocket/sendInput' });
