#!/bin/bash

# =============================================================================
# SCRIPT PARA CONSTRUIR IMAGEN DEL BACKEND (CENTRAL-RESERVE)
# =============================================================================

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para logging
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] ✓${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ✗${NC} $1"
}

# Configuración
PROJECT_NAME="central-reserve"
PROJECT_PATH="../back/central-reserve"
DOCKERFILE_PATH="docker/Dockerfile"
IMAGE_NAME="reserve-backend"
TAG="${1:-latest}"

log "Construyendo imagen del backend..."
log "Proyecto: $PROJECT_NAME"
log "Ruta: $PROJECT_PATH"
log "Imagen: $IMAGE_NAME:$TAG"

# Verificar que el directorio existe
if [ ! -d "$PROJECT_PATH" ]; then
    error "No se encontró el directorio del proyecto: $PROJECT_PATH"
    exit 1
fi

# Verificar que el Dockerfile existe
if [ ! -f "$PROJECT_PATH/$DOCKERFILE_PATH" ]; then
    error "No se encontró el Dockerfile en: $PROJECT_PATH/$DOCKERFILE_PATH"
    exit 1
fi

# Construir la imagen
log "Construyendo imagen Docker..."
docker build \
    --no-cache \
    --tag "$IMAGE_NAME:$TAG" \
    --file "$PROJECT_PATH/$DOCKERFILE_PATH" \
    "$PROJECT_PATH"

success "Imagen del backend construida exitosamente: $IMAGE_NAME:$TAG"

# Mostrar información de la imagen
log "Información de la imagen:"
docker images "$IMAGE_NAME:$TAG" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

log "Para usar esta imagen, puedes ejecutar:"
echo "  docker run -p 3050:3050 $IMAGE_NAME:$TAG" 