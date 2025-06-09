import { Provider } from 'react-redux';
import { store } from './store';
import AdvancedDashboard from './components/AdvancedDashboard';

function App() {
  return (
    <Provider store={store}>
      <AdvancedDashboard />
    </Provider>
  );
}

export default App;
