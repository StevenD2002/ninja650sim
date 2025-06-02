import { motorcycle } from './generated/motorcycle.ts';
function App() {
  const testProtoBufTypes = () => {
    const engineData = new motorcycle.EngineData({
      rpm: 900,
      throttle_position: 0,
      power: 0,
      engine_temp: 90,
      afr_current: 14.7
    });

    const userInput = new motorcycle.UserInput({
      throttle_position: 0,
      clutch_position: 1.0,
      gear: 0,
    });

    return {
      engineData: engineData.toObject(),
      userInput: userInput.toObject(),
      fieldsInEngineData: Object.getOwnPropertyNames(motorcycle.EngineData.prototype),
    };
  } 

  return (
    <div style={{ padding: '20px', fontFamily: 'monospace' }}>
      <h1>Protobuf Test</h1>
      
      <div>
        <h2>Status:</h2>
        <p>✅ React app is running</p>
        <p>✅ Protobuf file generated: motorcycle.ts</p>
        
        <div style={{ 
          backgroundColor: '#f0f0f0', 
          padding: '10px', 
          margin: '10px 0',
          border: '1px solid #ccc'
        }}>
          <h3>Generated Protobuf Types:</h3>
          <ul>
            <li><code>motorcycle.EngineData</code> - Engine telemetry data</li>
            <li><code>motorcycle.UserInput</code> - User controls (throttle, clutch, gear)</li>
            <li><code>motorcycle.ECUMaps</code> - Engine maps (fuel, ignition, AFR)</li>
            <li><code>motorcycle.ECUSettings</code> - ECU configuration</li>
            <li><code>motorcycle.MotorcycleSimulatorClient</code> - gRPC client</li>
          </ul>
        </div>


        <div style={{ 
          backgroundColor: '#e8f5e8', 
          padding: '10px', 
          margin: '10px 0',
          border: '1px solid #4caf50'
        }}>
        </div>

        <div>
          <h3>Test Results:</h3>
          <pre style={{ backgroundColor: '#f5f5f5', padding: '10px' }}>
            {JSON.stringify(testProtoBufTypes(), null, 2)}
          </pre>
        </div>
      </div>
    </div>
  );
}

export default App;
