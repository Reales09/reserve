#!/bin/bash

# Script de despliegue para producción
# Reserve App Frontend - Ambiente de producción

set -e

# Variables específicas para producción
export API_BASE_URL="https://api.tudominio.com"
export APP_NAME="Reserve App"
export APP_VERSION="1.0.0"

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Iniciando despliegue de PRODUCCIÓN${NC}"
echo -e "${YELLOW}Configuración:${NC}"
echo -e "  - API_BASE_URL: ${API_BASE_URL}"
echo -e "  - APP_NAME: ${APP_NAME}"
echo -e "  - APP_VERSION: ${APP_VERSION}"

# Confirmación para producción
echo -e "${RED}⚠️  ATENCIÓN: Vas a desplegar a PRODUCCIÓN${NC}"
read -p "¿Estás seguro? (escribe 'SI' para continuar): " confirm

if [ "$confirm" != "SI" ]; then
    echo -e "${RED}❌ Despliegue cancelado${NC}"
    exit 1
fi

# Obtener versión como parámetro o usar latest
VERSION=${1:-"latest"}

# Llamar al script principal
exec ./scripts/deploy.sh "${VERSION}" 