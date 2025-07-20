import React from 'react';
import Sidebar from './Sidebar.js';
import './Layout.css';

const Layout = ({ children, activeView, onViewChange, userInfo, onLogout }) => {
    return (
        <div className="layout">
            <Sidebar 
                activeView={activeView} 
                onViewChange={onViewChange}
                userInfo={userInfo}
                onLogout={onLogout}
            />
            <main className="main-content">
                <div className="content-wrapper">
                    {children}
                </div>
            </main>
        </div>
    );
};

export default Layout; 