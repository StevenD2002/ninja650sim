import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { EngineData, TransmissionState, UIState, UserInputState } from '../types';

interface MotorcycleState {
  engineData: EngineData;
  transmission: TransmissionState;
  ui: UIState;
  userInput: UserInputState;
}

const initialState: MotorcycleState = {
  engineData: {
    rpm: 0,
    throttle_position: 0,
    timestamp: 0,
    power: 0,
    torque: 0,
    speed: 0,
    engine_temp: 90,
    afr_current: 14.7,
    afr_target: 14.7,
    fuel_injection_ms: 0,
    ignition_advance: 0,
    gear: 0,
    clutch_position: 1.0,
  },
  transmission: {
    currentGear: 0,        // Start in neutral
    clutchPosition: 1.0,   // Start with clutch disengaged
    clutchPressed: true,
  },
  ui: {
    statusMessage: '',
    statusMessageTime: 0,
    statusMessageColor: 'green',
    connected: false,
    lastDataUpdate: 0,
  },
  userInput: {
    throttlePos: 0,
  },
};

const motorcycleSlice = createSlice({
  name: 'motorcycle',
  initialState,
  reducers: {
    // Engine data updates from WebSocket
    updateEngineData: (state, action: PayloadAction<EngineData>) => {
      state.engineData = { ...action.payload };
      state.ui.lastDataUpdate = Date.now();
    },

    // User input actions
    setThrottlePosition: (state, action: PayloadAction<number>) => {
      const newThrottle = Math.max(0, Math.min(100, action.payload));
      state.userInput.throttlePos = newThrottle;
    },

    // Transmission actions
    setClutchPosition: (state, action: PayloadAction<number>) => {
      state.transmission.clutchPosition = action.payload;
      state.transmission.clutchPressed = action.payload > 0.8;
    },

    toggleClutch: (state) => {
      if (state.transmission.clutchPressed) {
        // Release clutch (engage)
        state.transmission.clutchPosition = 0.0;
        state.transmission.clutchPressed = false;
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Clutch ENGAGED', color: 'green' }
        });
      } else {
        // Press clutch (disengage)
        state.transmission.clutchPosition = 1.0;
        state.transmission.clutchPressed = true;
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Clutch DISENGAGED', color: 'red' }
        });
      }
    },

    setGear: (state, action: PayloadAction<number>) => {
      state.transmission.currentGear = action.payload;
    },

    shiftUp: (state) => {
      // Check clutch requirement
      if (state.transmission.clutchPosition < 0.8) {
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Press clutch to shift gears', color: 'red' }
        });
        return;
      }
      
      // Don't exceed 6th gear
      if (state.transmission.currentGear < 6) {
        const newGear = state.transmission.currentGear + 1;
        state.transmission.currentGear = newGear;
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: `Shifted to ${newGear} gear`, color: 'green' }
        });
      } else {
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Already in top gear', color: 'yellow' }
        });
      }
    },

    shiftDown: (state) => {
      // Check clutch requirement
      if (state.transmission.clutchPosition < 0.8) {
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Press clutch to shift gears', color: 'red' }
        });
        return;
      }
      
      // Don't go below neutral
      if (state.transmission.currentGear > 0) {
        const newGear = state.transmission.currentGear - 1;
        state.transmission.currentGear = newGear;
        const message = newGear === 0 ? 'Shifted to Neutral' : `Shifted to ${newGear} gear`;
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message, color: 'green' }
        });
      } else {
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Already in neutral', color: 'yellow' }
        });
      }
    },

    shiftToNeutral: (state) => {
      // Check clutch requirement
      if (state.transmission.clutchPosition < 0.8) {
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Press clutch to shift gears', color: 'red' }
        });
        return;
      }
      
      state.transmission.currentGear = 0;
      motorcycleSlice.caseReducers.showStatusMessage(state, {
        type: 'showStatusMessage',
        payload: { message: 'Shifted to Neutral', color: 'green' }
      });
    },

    // UI actions
    showStatusMessage: (state, action: PayloadAction<{ message: string; color: 'green' | 'yellow' | 'red' }>) => {
      state.ui.statusMessage = action.payload.message;
      state.ui.statusMessageColor = action.payload.color;
      state.ui.statusMessageTime = Date.now();
    },

    clearStatusMessage: (state) => {
      state.ui.statusMessage = '';
    },

    setConnectionStatus: (state, action: PayloadAction<boolean>) => {
      state.ui.connected = action.payload;
      if (action.payload) {
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Connected to server', color: 'green' }
        });
      } else {
        motorcycleSlice.caseReducers.showStatusMessage(state, {
          type: 'showStatusMessage',
          payload: { message: 'Disconnected from server', color: 'red' }
        });
      }
    },
  },
});

export const {
  updateEngineData,
  setThrottlePosition,
  setClutchPosition,
  toggleClutch,
  setGear,
  shiftUp,
  shiftDown,
  shiftToNeutral,
  showStatusMessage,
  clearStatusMessage,
  setConnectionStatus,
} = motorcycleSlice.actions;

export default motorcycleSlice.reducer;
