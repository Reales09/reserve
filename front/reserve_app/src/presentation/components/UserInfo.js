import React, { useState } from 'react';
import { useAuth } from '../hooks/useAuth.js';
import './UserInfo.css';

const UserInfo = () => {
  const { userInfo, userRolesPermissions, isSuperAdmin, getUserRoles, getUserPermissions } = useAuth();
  const [showDetails, setShowDetails] = useState(false);

  const toggleDetails = () => {
    setShowDetails(!showDetails);
  };

  if (!userInfo) {
    return null;
  }

  const roles = getUserRoles();
  const permissions = getUserPermissions();

  return (
    <div className="user-info-container">
      <div className="user-info-header" onClick={toggleDetails}>
        <div className="user-avatar">
          {userInfo.avatar_url ? (
            <img src={userInfo.avatar_url} alt="Avatar" />
          ) : (
            <span className="user-avatar-placeholder">
              {userInfo.name?.charAt(0)?.toUpperCase() || 'U'}
            </span>
          )}
        </div>
        <div className="user-details">
          <div className="user-name">{userInfo.name || 'Usuario'}</div>
          <div className="user-email">{userInfo.email || 'usuario@ejemplo.com'}</div>
          {isSuperAdmin() && (
            <div className="super-admin-badge">👑 Super Admin</div>
          )}
        </div>
        <button className="toggle-details-btn">
          {showDetails ? '▼' : '▶'}
        </button>
      </div>

      {showDetails && (
        <div className="user-details-expanded">
          <div className="roles-section">
            <h4>🎭 Roles ({roles.length})</h4>
            <div className="roles-list">
              {roles.map((role) => (
                <div key={role.id} className="role-item">
                  <span className="role-name">{role.name}</span>
                  <span className="role-code">({role.code})</span>
                  <span className="role-level">Nivel {role.level}</span>
                </div>
              ))}
            </div>
          </div>

          <div className="permissions-section">
            <h4>🔐 Permisos ({permissions.length})</h4>
            <div className="permissions-list">
              {permissions.map((permission) => (
                <div key={permission.id} className="permission-item">
                  <div className="permission-header">
                    <span className="permission-name">{permission.name}</span>
                    <span className={`permission-action ${permission.action}`}>
                      {permission.action}
                    </span>
                  </div>
                  <div className="permission-resource">
                    📁 {permission.resource}
                  </div>
                  <div className="permission-description">
                    {permission.description}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserInfo; 