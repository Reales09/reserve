#!/bin/bash

# =============================================================================
# SCRIPT MAESTRO PARA CONSTRUIR TODAS LAS IMÁGENES
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
TAG="${1:-latest}"
PARALLEL="${2:-false}"

log "=================================================="
log "CONSTRUYENDO TODAS LAS IMÁGENES - RESERVE APP"
log "=================================================="
log "Tag: $TAG"
log "Construcción paralela: $PARALLEL"

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.yaml" ]; then
    error "No se encontró docker-compose.yaml. Asegúrate de estar en el directorio infra/"
    exit 1
fi

# Verificar que los scripts existen
SCRIPTS_DIR="scripts"
if [ ! -d "$SCRIPTS_DIR" ]; then
    error "No se encontró el directorio de scripts: $SCRIPTS_DIR"
    exit 1
fi

# Lista de scripts a ejecutar
SCRIPTS=(
    "build-migrator.sh"
    "build-backend.sh"
    "build-frontend.sh"
    "build-website.sh"
)

# Función para ejecutar un script
execute_script() {
    local script_name="$1"
    local script_path="$SCRIPTS_DIR/$script_name"
    
    if [ -f "$script_path" ]; then
        log "Ejecutando: $script_name"
        chmod +x "$script_path"
        if "$script_path" "$TAG"; then
            success "✓ $script_name completado"
        else
            error "✗ $script_name falló"
            return 1
        fi
    else
        error "Script no encontrado: $script_path"
        return 1
    fi
}

# Construir imágenes
if [ "$PARALLEL" = "true" ]; then
    log "Construyendo imágenes en paralelo..."
    
    # Ejecutar todos los scripts en paralelo
    for script in "${SCRIPTS[@]}"; do
        execute_script "$script" &
    done
    
    # Esperar a que todos terminen
    wait
    
    # Verificar que todos terminaron exitosamente
    for script in "${SCRIPTS[@]}"; do
        if ! wait; then
            error "Error en construcción paralela"
            exit 1
        fi
    done
else
    log "Construyendo imágenes secuencialmente..."
    
    # Ejecutar scripts uno por uno
    for script in "${SCRIPTS[@]}"; do
        execute_script "$script" || exit 1
    done
fi

success "=================================================="
success "TODAS LAS IMÁGENES CONSTRUIDAS EXITOSAMENTE"
success "=================================================="

# Mostrar resumen de imágenes
log "Resumen de imágenes construidas:"
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}" | grep -E "(reserve-|REPOSITORY)"

log ""
log "Para desplegar todas las imágenes:"
echo "  docker compose up -d"
echo ""
log "Para construir usando docker-compose:"
echo "  docker compose build"
