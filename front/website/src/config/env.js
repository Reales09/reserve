// Configuración de variables de entorno para el website
export const config = {
  // URL del API - usar variable de entorno o fallback a producción
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://3.220.183.29:3050/api',
  
  // Otras configuraciones
  APP_NAME: 'Trattoria La Bella',
  VERSION: '1.0.0'
};

// Función helper para construir URLs completas
export const buildApiUrl = (endpoint) => {
  const baseUrl = config.API_BASE_URL;
  const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
  return `${baseUrl}${cleanEndpoint}`;
}; 