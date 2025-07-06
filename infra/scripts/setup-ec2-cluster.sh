#!/bin/bash

# Script para configurar instancias EC2 en el cluster de ECS
# Este script verifica si hay instancias EC2 y crea una si es necesario

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Funciones de logging
log() {
    echo -e "${NC}[$(date +'%Y-%m-%d %H:%M:%S')] $1"
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

# Cargar variables de entorno
load_env() {
    if [ -f ".env.aws" ]; then
        log "Cargando variables de entorno de .env.aws"
        set -a
        source .env.aws
        set +a
        success "Variables de entorno cargadas"
    else
        error "Archivo .env.aws no encontrado."
        exit 1
    fi
}

# Verificar si hay instancias EC2 en el cluster
check_ec2_instances() {
    log "Verificando instancias EC2 en el cluster..."
    
    INSTANCE_COUNT=$(aws ecs describe-clusters \
        --clusters $ECS_CLUSTER_NAME \
        --query 'clusters[0].registeredContainerInstancesCount' \
        --output text \
        --region $AWS_REGION)
    
    if [ "$INSTANCE_COUNT" -eq "0" ]; then
        warn "No hay instancias EC2 registradas en el cluster"
        return 1
    else
        success "Hay $INSTANCE_COUNT instancias EC2 registradas en el cluster"
        return 0
    fi
}

# Crear instancia EC2 para ECS
create_ec2_instance() {
    log "Creando instancia EC2 para ECS..."
    
    # Obtener la AMI más reciente de ECS
    ECS_AMI=$(aws ec2 describe-images \
        --owners amazon \
        --filters "Name=name,Values=amzn2-ami-ecs-hvm-*" \
        --query 'Images|sort_by(@, &CreationDate)[-1].ImageId' \
        --output text \
        --region $AWS_REGION)
    
    log "Usando AMI: $ECS_AMI"
    
    # Crear user data script
    cat > user-data.sh << EOF
#!/bin/bash
echo ECS_CLUSTER=$ECS_CLUSTER_NAME >> /etc/ecs/ecs.config
echo ECS_BACKEND_HOST= >> /etc/ecs/ecs.config
yum update -y
yum install -y docker
service docker start
usermod -a -G docker ec2-user
docker pull postgres:15-alpine
docker pull public.ecr.aws/d3a6d4r1/cam/reserve:backend-latest
docker pull public.ecr.aws/d3a6d4r1/cam/reserve:frontend-latest
docker pull public.ecr.aws/d3a6d4r1/cam/reserve:migrator-latest
EOF
    
    # Obtener subnet por defecto
    SUBNET=$(aws ec2 describe-subnets \
        --filters "Name=default-for-az,Values=true" \
        --query 'Subnets[0].SubnetId' \
        --output text \
        --region $AWS_REGION)
    
    # Obtener security group por defecto
    SG=$(aws ec2 describe-security-groups \
        --group-names default \
        --query 'SecurityGroups[0].GroupId' \
        --output text \
        --region $AWS_REGION)
    
    # Crear instancia EC2
    INSTANCE_ID=$(aws ec2 run-instances \
        --image-id $ECS_AMI \
        --count 1 \
        --instance-type t3.medium \
        --key-name "reserve-key" \
        --security-group-ids $SG \
        --subnet-id $SUBNET \
        --user-data file://user-data.sh \
        --iam-instance-profile Name=ecsInstanceRole \
        --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=reserve-ecs-instance}]" \
        --query 'Instances[0].InstanceId' \
        --output text \
        --region $AWS_REGION)
    
    log "Instancia EC2 creada: $INSTANCE_ID"
    
    # Esperar a que la instancia esté ejecutándose
    log "Esperando a que la instancia esté ejecutándose..."
    aws ec2 wait instance-running --instance-ids $INSTANCE_ID --region $AWS_REGION
    
    success "Instancia EC2 está ejecutándose"
    
    # Esperar a que se registre en el cluster
    log "Esperando a que la instancia se registre en el cluster ECS..."
    for i in {1..30}; do
        CURRENT_COUNT=$(aws ecs describe-clusters \
            --clusters $ECS_CLUSTER_NAME \
            --query 'clusters[0].registeredContainerInstancesCount' \
            --output text \
            --region $AWS_REGION)
        
        if [ "$CURRENT_COUNT" -gt "0" ]; then
            success "Instancia registrada en el cluster ECS"
            break
        fi
        
        if [ $i -eq 30 ]; then
            error "Timeout esperando el registro en ECS"
            exit 1
        fi
        
        log "Intento $i/30 - Esperando registro..."
        sleep 10
    done
    
    # Limpiar archivos temporales
    rm -f user-data.sh
    
    success "Instancia EC2 configurada exitosamente"
}

# Función principal
main() {
    log "=================================================="
    log "CONFIGURANDO CLUSTER ECS CON INSTANCIAS EC2"
    log "=================================================="
    
    load_env
    
    case "${1:-setup}" in
        "check")
            check_ec2_instances
            ;;
        "create")
            create_ec2_instance
            ;;
        "setup")
            if ! check_ec2_instances; then
                log "Configurando instancia EC2..."
                create_ec2_instance
            fi
            ;;
        *)
            echo "Uso: $0 [check|create|setup]"
            echo ""
            echo "Comandos:"
            echo "  check  - Verificar instancias EC2 en el cluster"
            echo "  create - Crear nueva instancia EC2"
            echo "  setup  - Configurar instancias EC2 si no existen"
            exit 1
            ;;
    esac
}

# Ejecutar función principal
main "$@" 