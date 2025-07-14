import React, { useEffect, useState, useRef } from 'react';

interface CircularGaugeProps {
  value: number;
  min: number;
  max: number;
  size?: number;
  thickness?: number;
  redline?: number;
  unit?: string;
  label: string;
  color?: string;
  animate?: boolean;
}

export const CircularGauge: React.FC<CircularGaugeProps> = ({
  value,
  min,
  max,
  size = 180,
  thickness = 12,
  redline,
  unit = '',
  label,
  color = '#10b981',
  animate = true,
}) => {
  const [animatedValue, setAnimatedValue] = useState(min);
  const [isInRedzone, setIsInRedzone] = useState(false);
  const animationRef = useRef<number>();

  // Smooth animation to new value
  useEffect(() => {
    if (!animate) {
      setAnimatedValue(value);
      return;
    }

    const duration = 300; // ms
    const startValue = animatedValue;
    const difference = value - startValue;
    const startTime = Date.now();

    const animateValue = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      
      // Easing function (ease-out)
      const easeOut = 1 - Math.pow(1 - progress, 3);
      const currentValue = startValue + difference * easeOut;
      
      setAnimatedValue(currentValue);
      
      if (progress < 1) {
        animationRef.current = requestAnimationFrame(animateValue);
      }
    };

    animationRef.current = requestAnimationFrame(animateValue);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [value, animate]);

  // Check if in redzone
  useEffect(() => {
    setIsInRedzone(redline ? animatedValue >= redline : false);
  }, [animatedValue, redline]);

  const normalizedValue = Math.max(min, Math.min(max, animatedValue));
  const percentage = ((normalizedValue - min) / (max - min)) * 100;
  const angle = (percentage / 100) * 270; // 270 degrees for 3/4 circle
  
  const center = size / 2;
  const radius = (size - thickness) / 2;
  const circumference = 2 * Math.PI * radius;
  const strokeDasharray = `${(270 / 360) * circumference} ${circumference}`;
  const strokeDashoffset = circumference - (angle / 360) * circumference;

  // Dynamic colors
  const getGaugeColor = () => {
    if (isInRedzone) return '#ef4444';
    if (redline && normalizedValue > redline * 0.8) return '#f59e0b';
    return color;
  };

  // Needle angle calculation
  const needleAngle = 90 + angle; // Start at -135 degrees

  return (
    <div 
      className="circular-gauge"
      style={{ 
        width: size, 
        height: size,
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center'
      }}
    >
      {/* Background Circle */}
      <svg
        width={size}
        height={size}
        style={{ position: 'absolute', transform: 'rotate(135deg)' }}
      >
        {/* Background arc */}
        <circle
          cx={center}
          cy={center}
          r={radius}
          fill="none"
          stroke="rgba(255, 255, 255, 0.1)"
          strokeWidth={thickness}
          strokeDasharray={strokeDasharray}
          strokeLinecap="round"
        />
        
        {/* Redline indicator */}
        {redline && (
          <circle
            cx={center}
            cy={center}
            r={radius}
            fill="none"
            stroke="rgba(239, 68, 68, 0.3)"
            strokeWidth={thickness}
            strokeDasharray={`${((redline - min) / (max - min)) * (270 / 360) * circumference} ${circumference}`}
            strokeDashoffset={circumference - (((redline - min) / (max - min)) * 270 / 360) * circumference}
            strokeLinecap="round"
          />
        )}

        {/* Value arc */}
        <circle
          cx={center}
          cy={center}
          r={radius}
          fill="none"
          stroke={getGaugeColor()}
          strokeWidth={thickness}
          strokeDasharray={strokeDasharray}
          strokeDashoffset={strokeDashoffset}
          strokeLinecap="round"
          style={{
            transition: animate ? 'stroke-dashoffset 0.3s ease-out, stroke 0.2s ease' : 'none',
            filter: isInRedzone ? 'drop-shadow(0 0 8px #ef4444)' : 'none',
          }}
        />

        {/* Needle */}
        <line
          x1={center}
          y1={center}
          x2={center}
          y2={center - radius + thickness / 2}
          stroke={getGaugeColor()}
          strokeWidth="3"
          strokeLinecap="round"
          style={{
            transformOrigin: `${center}px ${center}px`,
            transform: `rotate(${needleAngle}deg)`,
            transition: animate ? 'transform 0.3s ease-out' : 'none',
          }}
        />

        {/* Center dot */}
        <circle
          cx={center}
          cy={center}
          r="6"
          fill={getGaugeColor()}
        />
      </svg>

      {/* Value display */}
      <div style={{
        position: 'absolute',
        textAlign: 'center',
        color: 'white',
        zIndex: 1,
      }}>
        <div style={{
          fontSize: `${size * 0.12}px`,
          fontWeight: 'bold',
          color: getGaugeColor(),
          textShadow: isInRedzone ? '0 0 10px #ef4444' : 'none',
          transition: 'color 0.2s ease, text-shadow 0.2s ease',
        }}>
          {Math.round(normalizedValue)}{unit}
        </div>
        <div style={{
          fontSize: `${size * 0.08}px`,
          color: 'rgba(255, 255, 255, 0.7)',
          marginTop: `${size * 0.02}px`,
        }}>
          {label}
        </div>
      </div>

      {/* Tick marks */}
      <svg
        width={size}
        height={size}
        style={{ position: 'absolute', transform: 'rotate(-135deg)' }}
      >
        {Array.from({ length: 6 }, (_, i) => {
          const tickAngle = (i / 5) * 270;
          const tickValue = min + (i / 5) * (max - min);
          const isRedlineTick = redline && tickValue >= redline;
          
          return (
            <g key={i}>
              <line
                x1={center + (radius - thickness / 2) * Math.cos((tickAngle * Math.PI) / 180)}
                y1={center + (radius - thickness / 2) * Math.sin((tickAngle * Math.PI) / 180)}
                x2={center + (radius + thickness / 2) * Math.cos((tickAngle * Math.PI) / 180)}
                y2={center + (radius + thickness / 2) * Math.sin((tickAngle * Math.PI) / 180)}
                stroke={isRedlineTick ? '#ef4444' : 'rgba(255, 255, 255, 0.5)'}
                strokeWidth="2"
              />
            </g>
          );
        })}
      </svg>
    </div>
  );
};
