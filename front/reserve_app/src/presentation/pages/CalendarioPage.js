import React from 'react';
import './CalendarioPage.css';

const CalendarioPage = () => {
    return (
        <div className="calendario-page">
            <header className="calendario-header">
                <div className="header-content">
                    <h1>📅 Calendario de Reservas</h1>
                    <p>Vista de calendario para gestionar reservas por fecha</p>
                </div>
            </header>

            <div className="calendario-content">
                <div className="calendario-placeholder">
                    <div className="placeholder-icon">📅</div>
                    <h2>Calendario en Desarrollo</h2>
                    <p>Esta vista mostrará un calendario interactivo para gestionar las reservas por fecha.</p>
                    <div className="placeholder-features">
                        <div className="feature-item">
                            <span className="feature-icon">📊</span>
                            <span>Vista mensual y semanal</span>
                        </div>
                        <div className="feature-item">
                            <span className="feature-icon">✨</span>
                            <span>Arrastrar y soltar reservas</span>
                        </div>
                        <div className="feature-item">
                            <span className="feature-icon">🔔</span>
                            <span>Notificaciones de cambios</span>
                        </div>
                        <div className="feature-item">
                            <span className="feature-icon">📱</span>
                            <span>Vista responsive</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default CalendarioPage; 