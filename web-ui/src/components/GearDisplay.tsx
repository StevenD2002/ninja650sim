import React, { useState, useEffect } from 'react';

interface GearDisplayProps {
  currentGear: number;
  rpm: number;
  clutchPressed: boolean;
  size?: number;
}

export const GearDisplay: React.FC<GearDisplayProps> = ({
  currentGear,
  rpm,
  clutchPressed,
  size = 120,
}) => {
  const [isShifting, setIsShifting] = useState(false);
  const [previousGear, setPreviousGear] = useState(currentGear);

  // Detect gear changes
  useEffect(() => {
    if (currentGear !== previousGear) {
      setIsShifting(true);
      setTimeout(() => setIsShifting(false), 300);
      setPreviousGear(currentGear);
    }
  }, [currentGear, previousGear]);

  const getGearColor = () => {
    if (rpm > 10000) return '#ef4444';
    if (rpm > 8000) return '#f59e0b';
    if (isShifting) return '#3b82f6';
    return '#10b981';
  };

  const getGearText = () => currentGear === 0 ? 'N' : currentGear.toString();

  return (
    <div
      style={{
        width: size,
        height: size,
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        perspective: '1000px',
      }}
    >
      {/* Background glow */}
      <div
        style={{
          position: 'absolute',
          width: '100%',
          height: '100%',
          borderRadius: '50%',
          background: `radial-gradient(circle, ${getGearColor()}22 0%, transparent 70%)`,
          animation: rpm > 8000 ? 'pulse 1s ease-in-out infinite' : 'none',
        }}
      />

      {/* Main gear display */}
      <div
        style={{
          width: size * 0.8,
          height: size * 0.8,
          borderRadius: '50%',
          background: `linear-gradient(135deg, #374151 0%, #1f2937 100%)`,
          border: `3px solid ${getGearColor()}`,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          position: 'relative',
          transform: isShifting ? 'rotateY(180deg) scale(1.1)' : 'rotateY(0deg) scale(1)',
          transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
          boxShadow: `
            inset 0 2px 4px rgba(0, 0, 0, 0.1),
            0 4px 8px rgba(0, 0, 0, 0.2),
            0 0 20px ${getGearColor()}33
          `,
        }}
      >
        {/* Gear number */}
        <div
          style={{
            fontSize: size * 0.4,
            fontWeight: 'bold',
            color: getGearColor(),
            textShadow: `0 0 10px ${getGearColor()}`,
            transform: isShifting ? 'rotateY(180deg)' : 'rotateY(0deg)',
            transition: 'all 0.3s ease',
            filter: clutchPressed ? 'blur(1px)' : 'none',
          }}
        >
          {getGearText()}
        </div>

        {/* Gear shift indicators */}
        <div
          style={{
            position: 'absolute',
            top: -10,
            left: '50%',
            transform: 'translateX(-50%)',
            width: 0,
            height: 0,
            borderLeft: '8px solid transparent',
            borderRight: '8px solid transparent',
            borderBottom: `12px solid ${currentGear < 6 ? getGearColor() : '#374151'}`,
            opacity: currentGear < 6 ? 1 : 0.3,
            transition: 'all 0.2s ease',
          }}
        />
        
        <div
          style={{
            position: 'absolute',
            bottom: -10,
            left: '50%',
            transform: 'translateX(-50%)',
            width: 0,
            height: 0,
            borderLeft: '8px solid transparent',
            borderRight: '8px solid transparent',
            borderTop: `12px solid ${currentGear > 0 ? getGearColor() : '#374151'}`,
            opacity: currentGear > 0 ? 1 : 0.3,
            transition: 'all 0.2s ease',
          }}
        />
      </div>

      {/* Clutch indicator ring */}
      <div
        style={{
          position: 'absolute',
          width: '100%',
          height: '100%',
          borderRadius: '50%',
          border: `2px solid ${clutchPressed ? '#ef4444' : 'transparent'}`,
          transform: clutchPressed ? 'scale(1.1)' : 'scale(1)',
          transition: 'all 0.2s ease',
          animation: clutchPressed ? 'spin 2s linear infinite' : 'none',
        }}
      />

      {/* Gear position indicators */}
      <div style={{ position: 'absolute', width: '100%', height: '100%' }}>
        {[0, 1, 2, 3, 4, 5, 6].map((gear, index) => {
          const angle = (index / 6) * 360 - 90; // Start from top
          const isActive = gear === currentGear;
          const radius = size * 0.45;
          
          return (
            <div
              key={gear}
              style={{
                position: 'absolute',
                width: 8,
                height: 8,
                borderRadius: '50%',
                background: isActive ? getGearColor() : 'rgba(255, 255, 255, 0.3)',
                left: '50%',
                top: '50%',
                transform: `
                  translate(-50%, -50%) 
                  rotate(${angle}deg) 
                  translateY(-${radius}px)
                  ${isActive ? 'scale(1.5)' : 'scale(1)'}
                `,
                transition: 'all 0.3s ease',
                boxShadow: isActive ? `0 0 10px ${getGearColor()}` : 'none',
              }}
            />
          );
        })}
      </div>

      {/* CSS-in-JS animations */}
      <style>{`
        @keyframes pulse {
          0%, 100% { opacity: 1; }
          50% { opacity: 0.7; }
        }
        
        @keyframes spin {
          from { transform: scale(1.1) rotate(0deg); }
          to { transform: scale(1.1) rotate(360deg); }
        }
      `}</style>
    </div>
  );
};
