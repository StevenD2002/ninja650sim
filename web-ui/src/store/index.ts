import { configureStore } from '@reduxjs/toolkit';
import motorcycleReducer from './motorcycleSlice';
import { createWebSocketMiddleware } from '../services/websocketMiddleware';

const websocketMiddleware = createWebSocketMiddleware();

export const store = configureStore({
  reducer: {
    motorcycle: motorcycleReducer,
  },
  // use our middleware so we have a centralized way to handle WebSocket connections
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(websocketMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
