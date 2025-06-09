import React, { useEffect, useState, useRef } from 'react';

interface EngineVisualizerProps {
  rpm: number;
  throttle: number;
  engineTemp: number;
  power: number;
  size?: number;
}

export const EngineVisualizer: React.FC<EngineVisualizerProps> = ({
  rpm,
  throttle,
  engineTemp,
  power,
  size = 200,
}) => {
  const [pistonOffset, setPistonOffset] = useState(0);
  const animationRef = useRef<number>();
  const startTimeRef = useRef<number>(Date.now());

  // Animate pistons based on RPM
  useEffect(() => {
    if (rpm > 0) {
      const animate = () => {
        const elapsed = Date.now() - startTimeRef.current;
        const rpmFrequency = (rpm / 60) * 4; // Convert RPM to Hz, multiply by 4 for 4-stroke simulation
        const offset = Math.sin(elapsed * 0.01 * rpmFrequency) * 10;
        setPistonOffset(offset);
        animationRef.current = requestAnimationFrame(animate);
      };
      
      animationRef.current = requestAnimationFrame(animate);
    } else {
      setPistonOffset(0);
    }

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [rpm]);

  const getEngineColor = () => {
    if (engineTemp > 110) return '#ef4444';
    if (engineTemp > 100) return '#f59e0b';
    if (rpm > 8000) return '#3b82f6';
    return '#10b981';
  };

  const getFlameIntensity = () => {
    return Math.max(0, (throttle / 100) * (rpm / 11000));
  };

  return (
    <div
      style={{
        width: size,
        height: size,
        position: 'relative',
        background: 'radial-gradient(circle, #1f2937 0%, #111827 100%)',
        borderRadius: '12px',
        padding: '20px',
        border: `2px solid ${getEngineColor()}`,
        boxShadow: `0 0 20px ${getEngineColor()}33`,
        overflow: 'hidden',
      }}
    >
      {/* Engine block */}
      <svg
        width="100%"
        height="100%"
        viewBox="0 0 160 160"
        style={{ position: 'absolute', top: 0, left: 0 }}
      >
        {/* Cylinder block */}
        <rect
          x="40"
          y="60"
          width="80"
          height="60"
          fill="url(#engineGradient)"
          stroke={getEngineColor()}
          strokeWidth="2"
          rx="4"
        />

        {/* Cylinder heads */}
        <rect
          x="50"
          y="50"
          width="25"
          height="20"
          fill="url(#cylinderGradient)"
          stroke={getEngineColor()}
          strokeWidth="1"
          rx="2"
        />
        <rect
          x="85"
          y="50"
          width="25"
          height="20"
          fill="url(#cylinderGradient)"
          stroke={getEngineColor()}
          strokeWidth="1"
          rx="2"
        />

        {/* Pistons (animated) */}
        <rect
          x="55"
          y={65 + pistonOffset}
          width="15"
          height="10"
          fill={getEngineColor()}
          rx="2"
        />
        <rect
          x="90"
          y={65 - pistonOffset}
          width="15"
          height="10"
          fill={getEngineColor()}
          rx="2"
        />

        {/* Connecting rods */}
        <line
          x1="62.5"
          y1={75 + pistonOffset}
          x2="80"
          y2="100"
          stroke={getEngineColor()}
          strokeWidth="3"
          strokeLinecap="round"
        />
        <line
          x1="97.5"
          y1={75 - pistonOffset}
          x2="80"
          y2="100"
          stroke={getEngineColor()}
          strokeWidth="3"
          strokeLinecap="round"
        />

        {/* Crankshaft */}
        <circle
          cx="80"
          cy="100"
          r="8"
          fill={getEngineColor()}
          stroke="#000"
          strokeWidth="2"
        />

        {/* Exhaust flames (based on throttle and RPM) */}
        {getFlameIntensity() > 0.1 && (
          <>
            <ellipse
              cx="130"
              cy="80"
              rx={5 + getFlameIntensity() * 10}
              ry={3 + getFlameIntensity() * 5}
              fill="url(#flameGradient)"
              opacity={0.3 + getFlameIntensity() * 0.7}
              style={{
                animation: rpm > 0 ? 'flicker 0.1s ease-in-out infinite alternate' : 'none',
              }}
            />
            <ellipse
              cx="135"
              cy="80"
              rx={3 + getFlameIntensity() * 7}
              ry={2 + getFlameIntensity() * 3}
              fill="url(#flameGradient2)"
              opacity={0.5 + getFlameIntensity() * 0.5}
              style={{
                animation: rpm > 0 ? 'flicker 0.15s ease-in-out infinite alternate-reverse' : 'none',
              }}
            />
          </>
        )}

        {/* Gradients */}
        <defs>
          <linearGradient id="engineGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#4b5563" />
            <stop offset="100%" stopColor="#1f2937" />
          </linearGradient>
          
          <linearGradient id="cylinderGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#6b7280" />
            <stop offset="100%" stopColor="#374151" />
          </linearGradient>
          
          <radialGradient id="flameGradient" cx="50%" cy="50%" r="50%">
            <stop offset="0%" stopColor="#fbbf24" stopOpacity="1" />
            <stop offset="50%" stopColor="#f59e0b" stopOpacity="0.8" />
            <stop offset="100%" stopColor="#ef4444" stopOpacity="0.4" />
          </radialGradient>
          
          <radialGradient id="flameGradient2" cx="50%" cy="50%" r="50%">
            <stop offset="0%" stopColor="#fef3c7" stopOpacity="1" />
            <stop offset="100%" stopColor="#f59e0b" stopOpacity="0.6" />
          </radialGradient>
        </defs>
      </svg>

      {/* Engine stats overlay */}
      <div
        style={{
          position: 'absolute',
          bottom: '10px',
          left: '10px',
          right: '10px',
          background: 'rgba(0, 0, 0, 0.8)',
          borderRadius: '6px',
          padding: '8px',
          fontSize: '10px',
          color: 'white',
        }}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <span>RPM: {Math.round(rpm)}</span>
          <span style={{ color: getEngineColor() }}>
            {Math.round(power)}hp
          </span>
        </div>
      </div>

      {/* Heat waves effect */}
      {engineTemp > 100 && (
        <div
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: 'linear-gradient(45deg, transparent 48%, rgba(239, 68, 68, 0.1) 50%, transparent 52%)',
            animation: 'heatWave 2s ease-in-out infinite',
            pointerEvents: 'none',
          }}
        />
      )}

      {/* Vibration effect for high RPM */}
      {rpm > 9000 && (
        <div
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            animation: 'vibrate 0.1s linear infinite',
            pointerEvents: 'none',
          }}
        />
      )}

      <style>{`
        @keyframes flicker {
          0% { opacity: 0.8; transform: scale(1); }
          100% { opacity: 1; transform: scale(1.1); }
        }
        
        @keyframes heatWave {
          0%, 100% { transform: translateX(-2px); }
          50% { transform: translateX(2px); }
        }
        
        @keyframes vibrate {
          0% { transform: translate(0); }
          25% { transform: translate(0.5px, 0.5px); }
          50% { transform: translate(-0.5px, 0.5px); }
          75% { transform: translate(0.5px, -0.5px); }
          100% { transform: translate(-0.5px, -0.5px); }
        }
      `}</style>
    </div>
  );
};
