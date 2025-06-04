import { useAppSelector } from '../hooks/redux';

export const TransmissionPanel = () => {
  const {
    engineData,
    transmission,
  } = useAppSelector((state: any) => state.motorcycle);


  // Helper functions
  const getGearText = () => transmission.currentGear === 0 ? 'N' : transmission.currentGear.toString();
  const getGearColor = () => {
    if (engineData.rpm > 10000) return '#ef4444';
    if (engineData.rpm > 8000) return '#eab308';
    return '#10b981';
  };
  const getClutchStatus = () => {
      if (transmission.clutchPosition > 0.8) return { text: 'DISENGAGED', color: '#ef4444' };
      if (transmission.clutchPosition > 0.2) return { text: 'SLIPPING', color: '#eab308' };
      return { text: 'ENGAGED', color: '#10b981' };
  };

  const clutchStatus = getClutchStatus();

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
            Transmission
          </h3>
          <div style={{ textAlign: 'center' }}>
            <div style={{ marginBottom: '20px' }}>
              <div style={{ color: '#9ca3af', marginBottom: '8px' }}>Current Gear:</div>
              <div style={{ 
                fontSize: '3rem', 
                fontWeight: 'bold', 
                color: getGearColor() 
              }}>
                {getGearText()}
              </div>
            </div>
            <div>
              <div style={{ color: '#9ca3af', marginBottom: '8px' }}>Clutch:</div>
              <div style={{ 
                fontSize: '1.25rem', 
                fontWeight: 'bold', 
                color: clutchStatus.color 
              }}>
                {clutchStatus.text}
              </div>
            </div>
          </div>
        </div>

  )
}
