import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from './redux';
import {
  setThrottlePosition,
  toggleClutch,
  shiftUp,
  shiftDown,
  shiftToNeutral,
} from '../store/motorcycleSlice';

export const useKeyboardControls = () => {
  const dispatch = useAppDispatch();
  const { throttlePos } = useAppSelector((state) => state.motorcycle.userInput);

  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      // Prevent default browser behavior for our keys
      const controlKeys = [
        'arrowleft', 'arrowright', 'c', 'u', 'd', 'n',
        '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'
      ];
      
      if (controlKeys.includes(event.key.toLowerCase())) {
        event.preventDefault();
      }

      switch (event.key.toLowerCase()) {
        case 'arrowleft':
          dispatch(setThrottlePosition(throttlePos - 5));
          break;
        case 'arrowright':
          dispatch(setThrottlePosition(throttlePos + 5));
          break;
        case 'c':
          dispatch(toggleClutch());
          break;
        case 'u':
          dispatch(shiftUp());
          break;
        case 'd':
          dispatch(shiftDown());
          break;
        case 'n':
          dispatch(shiftToNeutral());
          break;
        case '0': case '1': case '2': case '3': case '4':
        case '5': case '6': case '7': case '8': case '9':
          const digit = parseInt(event.key);
          if (event.shiftKey && digit === 1) {
            dispatch(setThrottlePosition(100)); // Shift+1 = 100%
          } else {
            dispatch(setThrottlePosition(digit * 10));
          }
          break;
      }
    };

    // Add event listener
    window.addEventListener('keydown', handleKeyPress);
    
    // Cleanup function
    return () => {
      window.removeEventListener('keydown', handleKeyPress);
    };
  }, [dispatch, throttlePos]);

  return
};
