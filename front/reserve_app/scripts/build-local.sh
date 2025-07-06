#!/bin/bash

# Script para construir imagen localmente
# Reserve App Frontend - Build local

set -e

# Variables
IMAGE_NAME="reserve-app-frontend"
VERSION=${1:-"local"}
API_BASE_URL=${API_BASE_URL:-"http://localhost:3050"}
APP_NAME=${APP_NAME:-"Reserve App (Local)"}
APP_VERSION=${APP_VERSION:-"local-dev"}

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔨 Construyendo imagen Docker localmente${NC}"
echo -e "${YELLOW}Configuración:${NC}"
echo -e "  - IMAGE_NAME: ${IMAGE_NAME}:${VERSION}"
echo -e "  - API_BASE_URL: ${API_BASE_URL}"
echo -e "  - APP_NAME: ${APP_NAME}"
echo -e "  - APP_VERSION: ${APP_VERSION}"

# Verificar que estamos en el directorio correcto
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ Error: No se encontró package.json. Ejecuta desde el directorio raíz del proyecto${NC}"
    exit 1
fi

# Verificar que Docker esté corriendo
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Docker no está corriendo${NC}"
    exit 1
fi

# Construir la imagen
echo -e "${YELLOW}🔨 Construyendo imagen...${NC}"
docker build \
  --build-arg REACT_APP_API_BASE_URL="${API_BASE_URL}" \
  --build-arg REACT_APP_NAME="${APP_NAME}" \
  --build-arg REACT_APP_VERSION="${APP_VERSION}" \
  -t ${IMAGE_NAME}:${VERSION} \
  .

echo -e "${GREEN}✅ Imagen construida exitosamente!${NC}"
echo -e "${YELLOW}📋 Para ejecutar la imagen:${NC}"
echo -e "docker run -p 80:80 \\"
echo -e "  -e REACT_APP_API_BASE_URL=${API_BASE_URL} \\"
echo -e "  ${IMAGE_NAME}:${VERSION}"
echo -e ""
echo -e "${BLUE}🌐 La aplicación estará disponible en: http://localhost:80${NC}" 