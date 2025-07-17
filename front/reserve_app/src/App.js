import React, { useState } from 'react';
import Layout from './presentation/components/Layout.js';
import GestionReservas from './presentation/pages/GestionReservas.js';
import CalendarioPage from './presentation/pages/CalendarioPage.js';
import './App.css';

function App() {
  const [activeView, setActiveView] = useState('calendario'); // Calendario como página principal

  const handleViewChange = (view) => {
    setActiveView(view);
  };

  const renderContent = () => {
    switch (activeView) {
      case 'calendario':
        return <CalendarioPage />;
      case 'reservas':
        return <GestionReservas />;
      default:
        return <CalendarioPage />;
    }
  };

  return (
    <div className="App">
      <Layout activeView={activeView} onViewChange={handleViewChange}>
        {renderContent()}
      </Layout>
    </div>
  );
}

export default App;
