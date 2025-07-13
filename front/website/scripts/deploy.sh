#!/bin/bash

# =============================================================================
# SCRIPT DE DESPLIEGUE PARA WEBSITE (ASTRO) - RESERVE APP
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

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] ⚠${NC} $1"
}

# Configuración
ECR_REPO="public.ecr.aws/d3a6d4r1/cam/reserve"
IMAGE_NAME="reserve-website"
VERSION="${1:-latest}"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
DESCRIPTIVE_TAG="website-${VERSION}"
DATED_TAG="website-${TIMESTAMP}"

log "=================================================="
log "DESPLEGANDO WEBSITE (ASTRO) - RESERVE APP"
log "=================================================="
log "Repositorio: $ECR_REPO"
log "Imagen: $IMAGE_NAME"
log "Versión: $VERSION"
log "Tag descriptivo: $DESCRIPTIVE_TAG"
log "Tag con fecha: $DATED_TAG"

# Verificar que estamos en el directorio correcto
if [ ! -f "Dockerfile" ]; then
    error "No se encontró Dockerfile. Asegúrate de estar en el directorio del website."
    exit 1
fi

# Verificar que Docker está instalado
if ! command -v docker &> /dev/null; then
    error "Docker no está instalado."
    exit 1
fi

# Verificar que AWS CLI está instalado
if ! command -v aws &> /dev/null; then
    error "AWS CLI no está instalado."
    exit 1
fi

# Login a ECR público
log "Haciendo login a ECR público..."
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

if [ $? -eq 0 ]; then
    success "Login exitoso a ECR público"
else
    error "Error en login a ECR público"
    exit 1
fi

# Construir imagen
log "Construyendo imagen Docker para ARM64..."
docker build \
    --no-cache \
    --platform linux/arm64 \
    --tag "${IMAGE_NAME}:${VERSION}" \
    .

if [ $? -eq 0 ]; then
    success "Imagen construida exitosamente: ${IMAGE_NAME}:${VERSION}"
else
    error "Error construyendo imagen"
    exit 1
fi

# Tagear para ECR
log "Tageando imagen para ECR..."
docker tag "${IMAGE_NAME}:${VERSION}" "${ECR_REPO}:${DESCRIPTIVE_TAG}"
docker tag "${IMAGE_NAME}:${VERSION}" "${ECR_REPO}:${DATED_TAG}"

if [ "$VERSION" != "latest" ]; then
    docker tag "${IMAGE_NAME}:${VERSION}" "${ECR_REPO}:website-${VERSION}"
fi

success "Imágenes tageadas para ECR"

# Subir a ECR
log "Subiendo imágenes a ECR..."
docker push "${ECR_REPO}:${DESCRIPTIVE_TAG}"
docker push "${ECR_REPO}:${DATED_TAG}"

if [ "$VERSION" != "latest" ]; then
    docker push "${ECR_REPO}:website-${VERSION}"
fi

success "Imágenes subidas exitosamente a ECR"

# Mostrar información de las imágenes
log "=================================================="
log "IMÁGENES DISPONIBLES EN ECR"
log "=================================================="
echo -e "  ${ECR_REPO}:${DESCRIPTIVE_TAG}"
echo -e "  ${ECR_REPO}:${DATED_TAG}"
if [ "$VERSION" != "latest" ]; then
    echo -e "  ${ECR_REPO}:website-${VERSION}"
fi

log ""
log "=================================================="
log "COMANDOS DE EJECUCIÓN"
log "=================================================="
echo -e "docker run -p 80:80 ${ECR_REPO}:${DESCRIPTIVE_TAG}"
echo -e ""
echo -e "Imágenes disponibles:"
echo -e "  - ${ECR_REPO}:${DESCRIPTIVE_TAG}"
echo -e "  - ${ECR_REPO}:${DATED_TAG}"
if [ "$VERSION" != "latest" ]; then
    echo -e "  - ${ECR_REPO}:website-${VERSION}"
fi

success "=================================================="
success "DESPLIEGUE COMPLETADO EXITOSAMENTE"
success "==================================================" 