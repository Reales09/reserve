import React from 'react';
import Sidebar from './Sidebar.js';
import './Layout.css';

const Layout = ({ children, activeView, onViewChange }) => {
    return (
        <div className="layout">
            <Sidebar activeView={activeView} onViewChange={onViewChange} />
            <main className="main-content">
                <div className="content-wrapper">
                    {children}
                </div>
            </main>
        </div>
    );
};

export default Layout; 