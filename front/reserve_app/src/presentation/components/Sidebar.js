import React from 'react';
import UserInfo from './UserInfo.js';
import './Sidebar.css';

const Sidebar = ({ activeView, onViewChange, userInfo, onLogout }) => {
    const menuItems = [
        {
            id: 'calendario',
            icon: '▦',
            label: 'Calendario',
            path: '/calendario'
        },
        {
            id: 'reservas',
            icon: '≡',
            label: 'Reservas',
            path: '/reservas'
        },
        {
            id: 'auth-test',
            icon: '🔐',
            label: 'Prueba Auth',
            path: '/auth-test'
        }
    ];

    return (
        <div className="sidebar">
            {/* Header con información del usuario */}
            <div className="sidebar-header">
                <div className="user-info">
                    <div className="user-avatar">
                        {userInfo?.avatar_url ? (
                            <img src={userInfo.avatar_url} alt="Avatar" />
                        ) : (
                            <span className="user-avatar-placeholder">
                                {userInfo?.name?.charAt(0)?.toUpperCase() || 'U'}
                            </span>
                        )}
                    </div>
                    <div className="user-details">
                        <div className="user-name">{userInfo?.name || 'Usuario'}</div>
                        <div className="user-email">{userInfo?.email || 'usuario@ejemplo.com'}</div>
                    </div>
                </div>
                <button 
                    className="logout-button"
                    onClick={onLogout}
                    title="Cerrar sesión"
                >
                    <span className="logout-icon">🚪</span>
                    <span className="logout-label">Salir</span>
                </button>
            </div>

            <nav className="sidebar-nav">
                <ul className="nav-list">
                    {menuItems.map((item) => (
                        <li
                            key={item.id}
                            className="nav-item"
                            data-tooltip={item.label}
                        >
                            <button
                                className={`nav-button ${activeView === item.id ? 'active' : ''}`}
                                onClick={() => onViewChange(item.id)}
                                title={item.label}
                            >
                                <span className="nav-icon">{item.icon}</span>
                                <span className="nav-label">{item.label}</span>
                            </button>
                        </li>
                    ))}
                </ul>
            </nav>

            {/* Componente UserInfo expandido */}
            <div className="sidebar-user-info">
                <UserInfo />
            </div>
        </div>
    );
};

export default Sidebar; 