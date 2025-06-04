import { useAppSelector } from '../hooks/redux'
  // Get state from Redux store
export const EngineDataPanel = () => {
  const {
    engineData,
    ui,
  } = useAppSelector((state: any) => state.motorcycle);
  const lastUpdate = new Date(ui.lastDataUpdate).toLocaleTimeString();

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
            Engine Data
          </h3>
          {ui.connected && engineData.timestamp > 0 ? (
            <div style={{ display: 'grid', gap: '8px', fontSize: '14px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>RPM:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.rpm.toFixed(0)}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Power:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.power.toFixed(1)} hp
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Torque:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.torque.toFixed(1)} Nm
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Speed:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.speed.toFixed(1)} km/h
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Engine Temp:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.engine_temp.toFixed(0)}°C
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>AFR:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.afr_current.toFixed(1)}
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Ignition:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.ignition_advance.toFixed(1)}° BTDC
                </span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#9ca3af' }}>Fuel:</span>
                <span style={{ color: '#fbbf24', fontWeight: 'bold' }}>
                  {engineData.fuel_injection_ms.toFixed(2)} ms
                </span>
              </div>
              <div style={{ 
                borderTop: '1px solid #4b5563', 
                paddingTop: '8px', 
                marginTop: '8px' 
              }}>
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: '#6b7280', fontSize: '12px' }}>Last Update:</span>
                  <span style={{ color: '#60a5fa', fontSize: '12px' }}>
                    {lastUpdate}
                  </span>
                </div>
              </div>
            </div>
          ) : (
            <div style={{ textAlign: 'center', color: '#6b7280' }}>
              {ui.connected ? 'Waiting for data...' : 'Not connected'}
            </div>
          )}
        </div>
  )
}
