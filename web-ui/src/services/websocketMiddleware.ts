import { WebSocketClient, WSUserInput } from './websocketClient';
import { updateEngineData, setConnectionStatus } from '../store/motorcycleSlice';

// WebSocket middleware that handles connection and message sending
export const createWebSocketMiddleware = () => {
  let client: WebSocketClient | null = null;
  
  return (store: any) => (next: any) => (action: any) => {
    switch (action.type) {
      case 'websocket/connect':
        if (!client) {
          client = new WebSocketClient();
          
          client.connect(
            (data) => {
              store.dispatch(updateEngineData(data));
            },
            (connected) => {
              store.dispatch(setConnectionStatus(connected));
            }
          );
        }
        break;
        
      case 'websocket/disconnect':
        if (client) {
          client.disconnect();
          client = null;
        }
        break;
        
      case 'websocket/sendInput':
        if (client && client.isConnected()) {
          const state = store.getState();
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

export const connectWebSocket = () => ({ type: 'websocket/connect' });
export const disconnectWebSocket = () => ({ type: 'websocket/disconnect' });
export const sendUserInput = () => ({ type: 'websocket/sendInput' });
