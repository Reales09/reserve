#!/bin/bash

# =============================================================================
# SCRIPT DE DESPLIEGUE A AWS ECS - RESERVE APP
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

# Verificar dependencias
check_dependencies() {
    log "Verificando dependencias..."
    
    # Verificar AWS CLI
    if ! command -v aws &> /dev/null; then
        error "AWS CLI no está instalado. Instálalo primero."
        exit 1
    fi
    
    # Verificar Docker
    if ! command -v docker &> /dev/null; then
        error "Docker no está instalado."
        exit 1
    fi
    
    # Verificar ECS CLI (opcional)
    if command -v ecs-cli &> /dev/null; then
        success "ECS CLI encontrado"
    else
        warn "ECS CLI no encontrado. Se usará AWS CLI y Docker Compose."
    fi
    
    success "Dependencias verificadas"
}

# Cargar variables de entorno
load_env() {
    if [ -f ".env.aws" ]; then
        log "Cargando variables de entorno de .env.aws"
        set -a
        source .env.aws
        set +a
        success "Variables de entorno cargadas"
    else
        error "Archivo .env.aws no encontrado. Crea el archivo con las variables de AWS."
        exit 1
    fi
}

# Verificar configuración de AWS
check_aws_config() {
    log "Verificando configuración de AWS..."
    
    # Verificar credenciales
    if ! aws sts get-caller-identity &> /dev/null; then
        error "No tienes credenciales de AWS configuradas."
        echo "Ejecuta: aws configure"
        exit 1
    fi
    
    # Verificar región
    CURRENT_REGION=$(aws configure get region)
    if [ "$CURRENT_REGION" != "$AWS_REGION" ]; then
        warn "Región actual ($CURRENT_REGION) diferente a la configurada ($AWS_REGION)"
        log "Usando región configurada en .env.aws: $AWS_REGION"
    fi
    
    success "Configuración de AWS verificada"
}

# Construir y subir imágenes a ECR
build_and_push_images() {
    log "Construyendo y subiendo imágenes a ECR público..."
    
    # Obtener token de ECR público
    aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws
    
    # Repositorio único existente
    ECR_REPO="public.ecr.aws/d3a6d4r1/cam/reserve"
    
    # Lista de imágenes a construir y subir con sus tags
    declare -A images=(
        ["backend"]="../back/central-reserve:docker/Dockerfile"
        ["frontend"]="../front/reserve_app:Dockerfile"
        ["migrator"]="../back/dbpostgres:docker/Dockerfile"
    )
    
    for tag_name in "${!images[@]}"; do
        IFS=':' read -r context dockerfile <<< "${images[$tag_name]}"
        
        log "Construyendo $tag_name..."
        
        # Nombre completo de la imagen
        FULL_IMAGE_NAME="$ECR_REPO:$tag_name-$IMAGE_TAG"
        
        # Construir imagen
        docker build \
            -t "reserve-$tag_name:$IMAGE_TAG" \
            -f "$context/$dockerfile" \
            "$context"
        
        # Tagear para ECR público
        docker tag "reserve-$tag_name:$IMAGE_TAG" "$FULL_IMAGE_NAME"
        
        # Subir a ECR público
        docker push "$FULL_IMAGE_NAME"
        
        success "$tag_name subida exitosamente como $FULL_IMAGE_NAME"
    done
}

# Crear task definition y desplegar a ECS
deploy_to_ecs() {
    log "Desplegando aplicación a ECS usando AWS CLI..."
    
    # Crear task definition JSON
    create_task_definition
    
    # Registrar task definition
    log "Registrando task definition..."
    TASK_DEF_ARN=$(aws ecs register-task-definition \
        --cli-input-json file://task-definition.json \
        --region $AWS_REGION \
        --query 'taskDefinition.taskDefinitionArn' \
        --output text)
    
    log "Task definition registrada: $TASK_DEF_ARN"
    
    # Crear o actualizar servicios
    deploy_services
    
    success "Aplicación desplegada a ECS"
}

# Crear task definition
create_task_definition() {
    log "Creando task definition..."
    
    cat > task-definition.json << EOF
{
    "family": "reserve-app",
    "networkMode": "bridge",
    "requiresCompatibilities": ["EC2"],
    "executionRoleArn": "arn:aws:iam::$AWS_ACCOUNT_ID:role/ecsTaskExecutionRole",
    "containerDefinitions": [
        {
            "name": "postgres",
            "image": "postgres:15-alpine",
            "essential": true,
            "memory": 256,
            "environment": [
                {"name": "POSTGRES_DB", "value": "$DB_NAME"},
                {"name": "POSTGRES_USER", "value": "$DB_USER"},
                {"name": "POSTGRES_PASSWORD", "value": "$DB_PASSWORD"}
            ],
            "portMappings": [
                {
                    "containerPort": 5432,
                    "hostPort": 5432,
                    "protocol": "tcp"
                }
            ],
            "healthCheck": {
                "command": ["CMD-SHELL", "pg_isready -U $DB_USER -d $DB_NAME"],
                "interval": 30,
                "timeout": 5,
                "retries": 3,
                "startPeriod": 60
            },
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": "/ecs/reserve-app",
                    "awslogs-region": "$AWS_REGION",
                    "awslogs-stream-prefix": "postgres"
                }
            }
        },
        {
            "name": "backend",
            "image": "public.ecr.aws/d3a6d4r1/cam/reserve:backend-latest",
            "essential": true,
            "memory": 512,
            "environment": [
                {"name": "DB_HOST", "value": "postgres"},
                {"name": "DB_PORT", "value": "5432"},
                {"name": "DB_NAME", "value": "$DB_NAME"},
                {"name": "DB_USER", "value": "$DB_USER"},
                {"name": "DB_PASS", "value": "$DB_PASSWORD"},
                {"name": "DB_LOG_LEVEL", "value": "info"},
                {"name": "PGSSLMODE", "value": "disable"},
                {"name": "APP_ENV", "value": "production"},
                {"name": "HTTP_PORT", "value": "3050"},
                {"name": "LOG_LEVEL", "value": "info"},
                {"name": "JWT_SECRET", "value": "$JWT_SECRET"},
                {"name": "URL_BASE_SWAGGER", "value": "/swagger/"}
            ],
            "portMappings": [
                {
                    "containerPort": 3050,
                    "hostPort": 3050,
                    "protocol": "tcp"
                }
            ],
            "links": ["postgres"],
            "dependsOn": [
                {
                    "containerName": "postgres",
                    "condition": "HEALTHY"
                }
            ],
            "healthCheck": {
                "command": ["CMD-SHELL", "ps aux | grep '[c]entral_reserve' || exit 1"],
                "interval": 30,
                "timeout": 5,
                "retries": 3,
                "startPeriod": 60
            },
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": "/ecs/reserve-app",
                    "awslogs-region": "$AWS_REGION",
                    "awslogs-stream-prefix": "backend"
                }
            }
        },
                 {
             "name": "migrator",
             "image": "public.ecr.aws/d3a6d4r1/cam/reserve:migrator-latest",
             "essential": false,
             "memory": 128,
             "environment": [
                 {"name": "DB_HOST", "value": "postgres"},
                 {"name": "DB_PORT", "value": "5432"},
                 {"name": "DB_NAME", "value": "$DB_NAME"},
                 {"name": "DB_USER", "value": "$DB_USER"},
                 {"name": "DB_PASS", "value": "$DB_PASSWORD"},
                 {"name": "DB_LOG_LEVEL", "value": "info"},
                 {"name": "PGSSLMODE", "value": "disable"},
                 {"name": "APP_ENV", "value": "production"},
                 {"name": "HTTP_PORT", "value": "8080"},
                 {"name": "LOG_LEVEL", "value": "info"},
                 {"name": "JWT_SECRET", "value": "$JWT_SECRET"}
             ],
             "links": ["postgres"],
             "dependsOn": [
                 {
                     "containerName": "postgres",
                     "condition": "HEALTHY"
                 }
             ],
             "logConfiguration": {
                 "logDriver": "awslogs",
                 "options": {
                     "awslogs-group": "/ecs/reserve-app",
                     "awslogs-region": "$AWS_REGION",
                     "awslogs-stream-prefix": "migrator"
                 }
             }
         },
         {
             "name": "frontend", 
             "image": "public.ecr.aws/d3a6d4r1/cam/reserve:frontend-latest",
             "essential": true,
             "memory": 128,
             "environment": [
                 {"name": "REACT_APP_API_BASE_URL", "value": "$REACT_APP_API_BASE_URL"},
                 {"name": "REACT_APP_NAME", "value": "$REACT_APP_NAME"},
                 {"name": "REACT_APP_VERSION", "value": "$REACT_APP_VERSION"},
                 {"name": "NODE_ENV", "value": "production"}
             ],
             "portMappings": [
                 {
                     "containerPort": 80,
                     "hostPort": 80,
                     "protocol": "tcp"
                 }
             ],
             "links": ["backend"],
             "dependsOn": [
                 {
                     "containerName": "backend",
                     "condition": "HEALTHY"
                 }
             ],
             "logConfiguration": {
                 "logDriver": "awslogs",
                 "options": {
                     "awslogs-group": "/ecs/reserve-app",
                     "awslogs-region": "$AWS_REGION",
                     "awslogs-stream-prefix": "frontend"
                 }
             }
         }
    ]
}
EOF
    
    success "Task definition creada"
}

# Desplegar servicios
deploy_services() {
    log "Creando servicio ECS..."
    
    # Crear o actualizar servicio para EC2
    aws ecs create-service \
        --cluster $ECS_CLUSTER_NAME \
        --service-name reserve-app-ec2 \
        --task-definition reserve-app \
        --desired-count 1 \
        --launch-type EC2 \
        --region $AWS_REGION 2>/dev/null || \
    aws ecs update-service \
        --cluster $ECS_CLUSTER_NAME \
        --service reserve-app-ec2 \
        --task-definition reserve-app \
        --region $AWS_REGION
    
    success "Servicio ECS configurado"
}

# Crear grupo de logs en CloudWatch
create_log_group() {
    log "Creando grupo de logs en CloudWatch..."
    
    aws logs create-log-group \
        --log-group-name /ecs/reserve-app \
        --region $AWS_REGION 2> /dev/null || \
    log "Grupo de logs ya existe"
    
    success "Grupo de logs configurado"
}

# Ejecutar migrador una sola vez
run_migrator() {
    log "Ejecutando migrador de base de datos..."
    
    # Crear task definition para migrador
    cat > migrator-task-definition.json << EOF
{
    "family": "reserve-migrator",
    "networkMode": "awsvpc",
    "requiresCompatibilities": ["FARGATE"],
    "cpu": "256",
    "memory": "512",
    "executionRoleArn": "arn:aws:iam::$AWS_ACCOUNT_ID:role/ecsTaskExecutionRole",
    "containerDefinitions": [
        {
            "name": "migrator",
            "image": "public.ecr.aws/d3a6d4r1/cam/reserve:migrator-latest",
            "essential": true,
            "memory": 512,
            "environment": [
                {"name": "DB_HOST", "value": "postgres"},
                {"name": "DB_PORT", "value": "5432"},
                {"name": "DB_NAME", "value": "$DB_NAME"},
                {"name": "DB_USER", "value": "$DB_USER"},
                {"name": "DB_PASS", "value": "$DB_PASSWORD"},
                {"name": "DB_LOG_LEVEL", "value": "info"},
                {"name": "PGSSLMODE", "value": "disable"},
                {"name": "APP_ENV", "value": "production"},
                {"name": "HTTP_PORT", "value": "8080"},
                {"name": "LOG_LEVEL", "value": "info"},
                {"name": "JWT_SECRET", "value": "$JWT_SECRET"}
            ],
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": "/ecs/reserve-app",
                    "awslogs-region": "$AWS_REGION",
                    "awslogs-stream-prefix": "migrator"
                }
            }
        }
    ]
}
EOF

    # Registrar task definition del migrador
    aws ecs register-task-definition \
        --cli-input-json file://migrator-task-definition.json \
        --region $AWS_REGION > /dev/null
    
    # Obtener subnets y security group
    SUBNETS=$(aws ec2 describe-subnets \
        --filters "Name=default-for-az,Values=true" \
        --query 'Subnets[].SubnetId' \
        --output text \
        --region $AWS_REGION | tr '\t' ',')
    
    SG=$(aws ec2 describe-security-groups \
        --group-names default \
        --query 'SecurityGroups[0].GroupId' \
        --output text \
        --region $AWS_REGION)
    
    # Ejecutar task de migración
    TASK_ARN=$(aws ecs run-task \
        --cluster $ECS_CLUSTER_NAME \
        --task-definition reserve-migrator \
        --launch-type FARGATE \
        --network-configuration "awsvpcConfiguration={subnets=[$SUBNETS],securityGroups=[$SG],assignPublicIp=ENABLED}" \
        --region $AWS_REGION \
        --query 'tasks[0].taskArn' \
        --output text)
    
    log "Esperando a que complete la migración..."
    aws ecs wait tasks-stopped \
        --cluster $ECS_CLUSTER_NAME \
        --tasks $TASK_ARN \
        --region $AWS_REGION
    
    success "Migración de base de datos completada"
}

# Ver estado del despliegue
check_deployment_status() {
    log "Verificando estado del despliegue..."
    
    # Ver estado de servicios
    aws ecs describe-services \
        --cluster $ECS_CLUSTER_NAME \
        --services reserve-app-ec2 \
        --region $AWS_REGION \
        --query 'services[0].{Name:serviceName,Status:status,Running:runningCount,Desired:desiredCount}' \
        --output table
    
    success "Estado del despliegue verificado"
}

# Verificar instancias EC2 en el cluster
check_ec2_cluster() {
    log "Verificando instancias EC2 en el cluster..."
    
    # Llamar al script de setup de EC2
    if [ -f "scripts/setup-ec2-cluster.sh" ]; then
        ./scripts/setup-ec2-cluster.sh setup
    else
        error "Script de setup-ec2-cluster.sh no encontrado"
        exit 1
    fi
    
    success "Cluster EC2 configurado"
}

# Función principal
main() {
    log "=================================================="
    log "DESPLEGANDO RESERVE APP A AWS ECS"
    log "=================================================="
    
    case "${1:-deploy}" in
        "check")
            check_dependencies
            load_env
            log "Cluster: $ECS_CLUSTER_NAME"
            log "Región: $AWS_REGION"
            check_aws_config
            ;;
        "build")
            check_dependencies
            load_env
            log "Cluster: $ECS_CLUSTER_NAME"
            log "Región: $AWS_REGION"
            check_aws_config
            build_and_push_images
            ;;
        "deploy")
            check_dependencies
            load_env
            log "Cluster: $ECS_CLUSTER_NAME"
            log "Región: $AWS_REGION"
            check_aws_config
            check_ec2_cluster
            build_and_push_images
            create_log_group
            deploy_to_ecs
            check_deployment_status
            ;;
        "status")
            load_env
            check_deployment_status
            ;;
        "logs")
            load_env
            log "Viendo logs en CloudWatch..."
            aws logs tail /ecs/reserve-app --follow --region $AWS_REGION
            ;;
        "down")
            load_env
            log "Deteniendo servicio ECS..."
            aws ecs update-service \
                --cluster $ECS_CLUSTER_NAME \
                --service reserve-app-ec2 \
                --desired-count 0 \
                --region $AWS_REGION
            success "Servicio detenido"
            ;;
        "setup-ec2")
            load_env
            log "Configurando instancias EC2..."
            ./scripts/setup-ec2-cluster.sh setup
            ;;
        *)
            echo "Uso: $0 [check|build|deploy|status|logs|down|setup-ec2]"
            echo ""
            echo "Comandos:"
            echo "  check     - Verificar dependencias y configuración"
            echo "  build     - Construir y subir imágenes a ECR"
            echo "  deploy    - Despliegue completo a ECS"
            echo "  status    - Ver estado de los servicios"
            echo "  logs      - Ver logs de los servicios"
            echo "  down      - Detener servicios en ECS"
            echo "  setup-ec2 - Configurar instancias EC2 para el cluster"
            exit 1
            ;;
    esac
}

# Ejecutar función principal
main "$@" 