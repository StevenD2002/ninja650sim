export const transmissionPanel = ({getGearColor, getGearText}) => {
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
