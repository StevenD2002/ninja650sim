export interface WSUserInput {
  throttle_position: number;
  clutch_position: number;
  gear: number;
}

export interface WSEngineData {
  rpm: number;
  throttle_position: number;
  timestamp: number;
  power: number;
  torque: number;
  speed: number;
  engine_temp: number;
  afr_current: number;
  afr_target: number;
  fuel_injection_ms: number;
  ignition_advance: number;
  gear: number;
  clutch_position: number;
}

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private onDataCallback?: (data: WSEngineData) => void;
  private onConnectionCallback?: (connected: boolean) => void;

  //TODO: need to move the url to an environment variable or config file
  constructor(private serverUrl: string = 'ws://localhost:8080/ws') {}

  connect(
    onData: (data: WSEngineData) => void,
    onConnection: (connected: boolean) => void
  ): void {
    this.onDataCallback = onData;
    this.onConnectionCallback = onConnection;
    this.connectWebSocket();
  }

  private connectWebSocket(): void {
    try {
      this.ws = new WebSocket(this.serverUrl);

      this.ws.onopen = () => {
        console.log('✅ WebSocket connected to Go server');
        this.reconnectAttempts = 0;
        this.onConnectionCallback?.(true);
      };

      this.ws.onmessage = (event) => {
        try {
          const data: WSEngineData = JSON.parse(event.data);
          this.onDataCallback?.(data);
        } catch (error) {
          console.error('❌ Error parsing engine data:', error);
        }
      };

      this.ws.onclose = (event) => {
        console.log(`📡 WebSocket disconnected (code: ${event.code})`);
        this.onConnectionCallback?.(false);
        this.handleReconnect();
      };

      this.ws.onerror = (error) => {
        console.error('❌ WebSocket error:', error);
        this.onConnectionCallback?.(false);
      };

    } catch (error) {
      console.error('❌ Error creating WebSocket:', error);
      this.handleReconnect();
    }
  }

  private handleReconnect(): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`🔄 Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
      
      setTimeout(() => {
        this.connectWebSocket();
      }, this.reconnectDelay * this.reconnectAttempts);
    } else {
      console.error('💔 Max reconnection attempts reached');
    }
  }

  sendUserInput(input: WSUserInput): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(input));
    } else {
      console.warn('⚠️  WebSocket is not connected, cannot send input');
    }
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}
