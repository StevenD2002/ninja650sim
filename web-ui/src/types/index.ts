// Redux state types for simulator

export interface EngineData {
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

export interface TransmissionState {
  currentGear: number;
  clutchPosition: number;
  clutchPressed: boolean;
}

export interface UIState {
  statusMessage: string;
  statusMessageTime: number;
  statusMessageColor: 'green' | 'yellow' | 'red';
  connected: boolean;
  lastDataUpdate: number;
}

export interface UserInputState {
  throttlePos: number;
}
