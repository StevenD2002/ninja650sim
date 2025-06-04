import { useAppSelector } from '../hooks/redux';

export const ReduxDebugPanel = () => {
  const { transmission, userInput, engineData, ui } = useAppSelector((state: any) => state.motorcycle);
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
          Redux State Debug
        </h3>
        <div style={{ fontSize: '12px', color: '#9ca3af' }}>
          <div><strong>Connection:</strong> {ui.connected ? '✅ Connected' : '❌ Disconnected'}</div>
          <div><strong>User Input:</strong></div>
          <div style={{ marginLeft: '10px' }}>
            Throttle: {userInput.throttlePos}%<br/>
            Clutch: {transmission.clutchPosition.toFixed(2)}<br/>
            Gear: {transmission.currentGear}
          </div>
          {engineData.timestamp > 0 && (
            <div style={{ marginTop: '10px' }}>
              <strong>Engine Data:</strong><br/>
              <div style={{ marginLeft: '10px' }}>
                RPM: {engineData.rpm.toFixed(0)}<br/>
                Power: {engineData.power.toFixed(1)} hp<br/>
                Speed: {engineData.speed.toFixed(1)} km/h
              </div>
            </div>
          )}
          <div style={{ marginTop: '10px' }}>
            <strong>Status:</strong> {ui.statusMessage || 'None'}
          </div>
        </div>
      </div>
  )}
