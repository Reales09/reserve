// Configuration - Environment Variables
export const config = {
  // API Base URL - se puede sobrescribir con variable de entorno
  API_BASE_URL: process.env.REACT_APP_API_BASE_URL || 'REACT_APP_API_BASE_URL_PLACEHOLDER',
  
  // Otros configuraciones que puedas necesitar
  APP_NAME: process.env.REACT_APP_NAME || 'Reserve App',
  APP_VERSION: process.env.REACT_APP_VERSION || '1.0.0',
  
  // Configuración de desarrollo
  isDevelopment: process.env.NODE_ENV === 'development',
  isProduction: process.env.NODE_ENV === 'production',
};

// Función para validar configuración
export const validateConfig = () => {
  const requiredVars = [
    'API_BASE_URL'
  ];
  
  const missing = requiredVars.filter(varName => !config[varName] || config[varName] === 'REACT_APP_API_BASE_URL_PLACEHOLDER');
  
  if (missing.length > 0) {
    console.warn('⚠️  Variables de entorno faltantes o no configuradas:', missing);
    
    // En desarrollo, usar localhost por defecto
    if (config.isDevelopment) {
      config.API_BASE_URL = 'http://localhost:3050';
      console.log('🔧 Usando URL base por defecto para desarrollo:', config.API_BASE_URL);
    }
  }
  
  console.log('🔧 Configuración cargada:', {
    API_BASE_URL: config.API_BASE_URL,
    APP_NAME: config.APP_NAME,
    NODE_ENV: process.env.NODE_ENV,
    isDevelopment: config.isDevelopment,
    isProduction: config.isProduction
  });
  
  return config;
};

// Validar configuración al cargar
validateConfig(); 