import { useAppSelector, useAppDispatch } from '../hooks/redux';
import {
  setThrottlePosition,
  toggleClutch,
  shiftUp,
  shiftDown,
  shiftToNeutral,
} from '../store/motorcycleSlice';
export const ControlsPanel = () => {
  // Get state from Redux store
  const {
    transmission,
    userInput,
  } = useAppSelector((state: any) => state.motorcycle);

  const dispatch = useAppDispatch();

  return (
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


  )
}
