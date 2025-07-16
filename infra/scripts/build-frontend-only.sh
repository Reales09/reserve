#!/bin/bash

# =============================================================================
# SCRIPT DE BUILD SOLO DEL FRONTEND - RESERVE APP
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

# Configurar BuildKit
setup_buildkit() {
    log "Configurando BuildKit..."
    export DOCKER_BUILDKIT=1
    CACHE_DIR=".buildkit/cache"
    mkdir -p "$CACHE_DIR"
    export BUILDX_CACHE_DIR="$CACHE_DIR"
    success "BuildKit configurado"
}

# Verificar dependencias
check_dependencies() {
    log "Verificando dependencias..."
    
    if ! command -v docker &> /dev/null; then
        error "Docker no está instalado."
        exit 1
    fi
    
    if ! docker buildx version &> /dev/null; then
        error "Docker Buildx no está disponible."
        exit 1
    fi
    
    success "Dependencias verificadas"
}

# Configurar API Base URL
setup_api_base_url() {
    log "Configurando API Base URL..."
    
    # Usar la variable de entorno o el valor por defecto
    API_BASE_URL="${REACT_APP_API_BASE_URL:-http://3.220.183.29:3050}"
    
    log "API_BASE_URL configurada: $API_BASE_URL"
    export REACT_APP_API_BASE_URL="$API_BASE_URL"
}

# Construir solo el frontend
build_frontend() {
    log "Construyendo solo el frontend..."
    
    # Obtener token de ECR público
    aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws
    
    # Repositorio
    ECR_REPO="public.ecr.aws/d3a6d4r1/cam/reserve"
    
    # Tag de la imagen
    IMAGE_TAG="${1:-latest}"
    
    log "Tag de imagen: $IMAGE_TAG"
    
    # Contexto y Dockerfile del frontend
    CONTEXT="../../front/reserve_app"
    DOCKERFILE="Dockerfile"
    
    # Nombre completo de la imagen
    FULL_IMAGE_NAME="$ECR_REPO:frontend-$IMAGE_TAG"
    
    log "Construyendo frontend para ARM64..."
    docker buildx build \
        --platform linux/arm64 \
        --cache-from type=local,src="$CACHE_DIR/frontend" \
        --cache-to type=local,dest="$CACHE_DIR/frontend",mode=max \
        --load \
        --build-arg REACT_APP_API_BASE_URL="$REACT_APP_API_BASE_URL" \
        -t "reserve-frontend:$IMAGE_TAG" \
        -f "$CONTEXT/$DOCKERFILE" \
        "$CONTEXT"
    
    # Tagear para ECR público
    docker tag "reserve-frontend:$IMAGE_TAG" "$FULL_IMAGE_NAME"
    
    # Subir a ECR público
    docker push "$FULL_IMAGE_NAME"
    
    success "Frontend subido exitosamente como $FULL_IMAGE_NAME"
}

# Mostrar información
show_info() {
    log "=================================================="
    log "FRONTEND DISPONIBLE EN ECR"
    log "=================================================="
    echo -e "  ${ECR_REPO}:frontend-${IMAGE_TAG}"
    
    log ""
    log "=================================================="
    log "COMANDO DE EJECUCIÓN"
    log "=================================================="
    echo -e "Frontend:  docker run -p 80:80 ${ECR_REPO}:frontend-${IMAGE_TAG}"
    
    log ""
    log "=================================================="
    log "CONFIGURACIÓN API"
    log "=================================================="
    echo -e "API_BASE_URL: $REACT_APP_API_BASE_URL"
}

# Función principal
main() {
    case "${1:-build}" in
        "build")
            log "=================================================="
            log "BUILD SOLO DEL FRONTEND - RESERVE APP"
            log "=================================================="
            
            check_dependencies
            setup_buildkit
            setup_api_base_url
            build_frontend "${2:-latest}"
            show_info
            
            success "=================================================="
            success "FRONTEND CONSTRUIDO Y SUBIDO EXITOSAMENTE"
            success "=================================================="
            ;;
        "clean")
            log "Limpiando cache de BuildKit..."
            rm -rf .buildkit/cache
            success "Cache limpiado"
            ;;
        *)
            echo "Uso: $0 [build|clean] [tag]"
            echo "  build [tag] - Construir y subir solo el frontend (tag opcional, default: latest)"
            echo "  clean       - Limpiar cache de BuildKit"
            echo ""
            echo "Ejemplos:"
            echo "  $0 build           # Construir con tag 'latest'"
            echo "  $0 build v1.0.1    # Construir con tag 'v1.0.1'"
            echo "  $0 clean           # Limpiar cache"
            echo ""
            echo "Variables de entorno:"
            echo "  REACT_APP_API_BASE_URL - URL del backend (default: http://3.220.183.29:3050)"
            exit 1
            ;;
    esac
}

# Ejecutar función principal
main "$@" 