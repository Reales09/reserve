#!/bin/bash

# Script de despliegue para staging
# Reserve App Frontend - Ambiente de staging

set -e

# Variables específicas para staging
export API_BASE_URL="https://staging-api.tudominio.com"
export APP_NAME="Reserve App (Staging)"
export APP_VERSION="staging-$(date +%Y%m%d_%H%M%S)"

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Iniciando despliegue de STAGING${NC}"
echo -e "${YELLOW}Configuración:${NC}"
echo -e "  - API_BASE_URL: ${API_BASE_URL}"
echo -e "  - APP_NAME: ${APP_NAME}"
echo -e "  - APP_VERSION: ${APP_VERSION}"

# Llamar al script principal con versión staging
exec ./scripts/deploy.sh "staging" 