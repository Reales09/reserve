import React from 'react';
import './Sidebar.css';

const Sidebar = ({ activeView, onViewChange }) => {
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
        }
    ];

    return (
        <div className="sidebar">


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


        </div>
    );
};

export default Sidebar; 