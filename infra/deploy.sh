#!/bin/bash

# =============================================================================
# SCRIPT DE DESPLIEGUE - RESERVE APP
# =============================================================================

set -e  # Salir si hay error
set -u  # Salir si variable no definida

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

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] ⚠${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ✗${NC} $1"
}

# Verificar que estamos en el directorio correcto
if [ ! -f "docker-compose.yaml" ]; then
    error "No se encontró docker-compose.yaml. Asegúrate de estar en el directorio infra/"
    exit 1
fi

# Verificar que existe el archivo .env
if [ ! -f ".env" ]; then
    error "No se encontró el archivo .env"
    warn "Copia el archivo .env.example a .env y configura tus variables"
    warn "cp .env.example .env"
    exit 1
fi

# Función para mostrar ayuda
show_help() {
    echo "Uso: $0 [COMANDO]"
    echo ""
    echo "Comandos:"
    echo "  build           - Construir todas las imágenes sin cache"
    echo "  build-backend   - Construir solo imagen del backend"
    echo "  build-frontend  - Construir solo imagen del frontend"
    echo "  build-migrator  - Construir solo imagen del migrador"
    echo "  build-all       - Construir todas las imágenes individualmente"
    echo "  up              - Levantar servicios"
    echo "  down            - Bajar servicios"
    echo "  restart         - Reiniciar servicios"
    echo "  logs            - Ver logs de todos los servicios"
    echo "  status          - Ver estado de los servicios"
    echo "  clean           - Limpiar recursos de Docker"
    echo ""
    echo "  aws-check       - Verificar configuración para AWS ECS"
    echo "  aws-build       - Construir y subir imágenes a ECR"
    echo "  aws-deploy      - Despliegue completo a AWS ECS"
    echo "  aws-status      - Ver estado en AWS ECS"
    echo "  aws-logs        - Ver logs en AWS ECS"
    echo "  aws-down        - Detener servicios en AWS ECS"
    echo ""
    echo "  help            - Mostrar esta ayuda"
    echo ""
    echo "Ejemplos:"
    echo "  $0 build && $0 up        # Construir y levantar"
    echo "  $0 build-backend         # Construir solo el backend"
    echo "  $0 build-all latest      # Construir todas las imágenes con tag 'latest'"
    echo "  $0 logs frontend         # Ver logs del frontend"
    echo "  $0 restart central_reserve # Reiniciar solo el backend"
    echo ""
    echo "  $0 aws-check             # Verificar configuración AWS"
    echo "  $0 aws-deploy            # Desplegar a AWS ECS"
    echo "  $0 aws-status            # Ver estado en AWS"
}

# Función para construir imágenes con docker-compose
build_images() {
    log "Construyendo imágenes con docker-compose..."
    docker compose build --no-cache
    success "Imágenes construidas exitosamente"
}

# Función para construir imagen del backend
build_backend() {
    log "Construyendo imagen del backend..."
    if [ -f "scripts/build-backend.sh" ]; then
        ./scripts/build-backend.sh "${2:-latest}"
    else
        error "Script build-backend.sh no encontrado"
        exit 1
    fi
}

# Función para construir imagen del frontend
build_frontend() {
    log "Construyendo imagen del frontend..."
    if [ -f "scripts/build-frontend.sh" ]; then
        ./scripts/build-frontend.sh "${2:-latest}"
    else
        error "Script build-frontend.sh no encontrado"
        exit 1
    fi
}

# Función para construir imagen del migrador
build_migrator() {
    log "Construyendo imagen del migrador..."
    if [ -f "scripts/build-migrator.sh" ]; then
        ./scripts/build-migrator.sh "${2:-latest}"
    else
        error "Script build-migrator.sh no encontrado"
        exit 1
    fi
}

# Función para construir todas las imágenes individualmente
build_all_images() {
    log "Construyendo todas las imágenes individualmente..."
    if [ -f "scripts/build-all.sh" ]; then
        ./scripts/build-all.sh "${2:-latest}" "${3:-false}"
    else
        error "Script build-all.sh no encontrado"
        exit 1
    fi
}

# Función para levantar servicios
start_services() {
    log "Levantando servicios..."
    docker compose up -d
    success "Servicios iniciados"
    
    # Esperar a que los servicios estén listos
    log "Esperando a que los servicios estén listos..."
    sleep 10
    
    # Verificar estado
    docker compose ps
}

# Función para bajar servicios
stop_services() {
    log "Bajando servicios..."
    docker compose down
    success "Servicios detenidos"
}

# Función para reiniciar servicios
restart_services() {
    log "Reiniciando servicios..."
    docker compose restart
    success "Servicios reiniciados"
}

# Función para ver logs
show_logs() {
    if [ -z "${2:-}" ]; then
        docker compose logs -f
    else
        docker compose logs -f "$2"
    fi
}

# Función para ver estado
show_status() {
    log "Estado de los servicios:"
    docker compose ps
    
    log "Uso de recursos:"
    docker stats --no-stream
}

# Función para limpiar recursos
clean_resources() {
    warn "Esto eliminará contenedores, imágenes y volúmenes no utilizados"
    read -p "¿Estás seguro? (y/N): " -r
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log "Limpiando recursos de Docker..."
        docker system prune -f
        docker volume prune -f
        success "Recursos limpiados"
    else
        log "Operación cancelada"
    fi
}

# Función para ejecutar comando AWS
run_aws_command() {
    local aws_command="$1"
    log "Ejecutando comando AWS: $aws_command"
    
    if [ -f "scripts/deploy-to-aws.sh" ]; then
        ./scripts/deploy-to-aws.sh "$aws_command"
    else
        error "Script deploy-to-aws.sh no encontrado"
        exit 1
    fi
}

# Función principal
main() {
    case "${1:-help}" in
        build)
            build_images
            ;;
        build-backend)
            build_backend "$@"
            ;;
        build-frontend)
            build_frontend "$@"
            ;;
        build-migrator)
            build_migrator "$@"
            ;;
        build-all)
            build_all_images "$@"
            ;;
        up)
            start_services
            ;;
        down)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            show_logs "$@"
            ;;
        status)
            show_status
            ;;
        clean)
            clean_resources
            ;;
        aws-check)
            run_aws_command "check"
            ;;
        aws-build)
            run_aws_command "build"
            ;;
        aws-deploy)
            run_aws_command "deploy"
            ;;
        aws-status)
            run_aws_command "status"
            ;;
        aws-logs)
            run_aws_command "logs"
            ;;
        aws-down)
            run_aws_command "down"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "Comando no reconocido: $1"
            show_help
            exit 1
            ;;
    esac
}

# Ejecutar función principal
main "$@" 