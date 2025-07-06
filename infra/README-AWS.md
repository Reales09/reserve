# 🚀 Despliegue en AWS ECS - Reserve App

Esta guía te ayudará a desplegar tu aplicación Reserve App en Amazon ECS (Elastic Container Service).

## 📋 Requisitos Previos

### 1. **Infraestructura AWS necesaria:**
- ✅ **Cluster ECS**: `central-reserve` (ya existe)
- ✅ **RDS PostgreSQL**: Base de datos gestionada
- ✅ **ECR Repositories**: Para almacenar imágenes Docker
- ✅ **VPC y Subnets**: Red configurada
- ✅ **Security Groups**: Puertos 80, 3050, 5432 abiertos
- ⚠️ **Application Load Balancer** (opcional, recomendado)

### 2. **Herramientas locales:**
```bash
# Verificar que tienes las herramientas instaladas
aws --version          # AWS CLI v2
docker --version       # Docker Engine
```

### 3. **Credenciales AWS:**
```bash
# Configurar credenciales
aws configure
# O usar variables de entorno:
export AWS_ACCESS_KEY_ID="tu_access_key"
export AWS_SECRET_ACCESS_KEY="tu_secret_key"
export AWS_DEFAULT_REGION="us-east-1"
```

## 🔧 Configuración Inicial

### 1. **Configurar variables de entorno para AWS:**

```bash
# Copiar plantilla de AWS
cp .env.aws .env.aws.local

# Editar con tus valores reales
nano .env.aws.local
```

### 2. **Variables críticas a configurar:**

```bash
# En .env.aws.local
AWS_ACCOUNT_ID=123456789012                    # Tu Account ID
AWS_REGION=us-east-1                           # Tu región
ECS_CLUSTER_NAME=central-reserve               # Tu cluster
RDS_ENDPOINT=tu-rds.us-east-1.rds.amazonaws.com # Tu RDS endpoint
DB_PASSWORD=tu_password_super_secreto          # Password de RDS
JWT_SECRET=tu_jwt_secret_super_largo           # JWT para producción
REACT_APP_API_BASE_URL=https://api.tudominio.com # Tu dominio público
```

## 🚀 Proceso de Despliegue

### **Paso 1: Verificar configuración**
```bash
./scripts/deploy-to-aws.sh check
```

### **Paso 2: Construir y subir imágenes a ECR**
```bash
./scripts/deploy-to-aws.sh build
```

### **Paso 3: Despliegue completo**
```bash
./scripts/deploy-to-aws.sh deploy
```

### **Verificar estado**
```bash
./scripts/deploy-to-aws.sh status
```

## 📊 Comandos Disponibles

| Comando | Descripción | Ejemplo |
|---------|-------------|---------|
| `check` | Verificar dependencias y configuración | `./scripts/deploy-to-aws.sh check` |
| `build` | Construir y subir imágenes a ECR | `./scripts/deploy-to-aws.sh build` |
| `deploy` | Despliegue completo a ECS | `./scripts/deploy-to-aws.sh deploy` |
| `status` | Ver estado de servicios en ECS | `./scripts/deploy-to-aws.sh status` |
| `logs` | Ver logs de servicios | `./scripts/deploy-to-aws.sh logs` |
| `down` | Detener servicios en ECS | `./scripts/deploy-to-aws.sh down` |

## 🏗️ Arquitectura en AWS

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │     Backend     │    │   DB Migrator   │
│   (React+Nginx) │    │   (Go API)      │    │   (Go - Once)   │
│   Port: 80      │    │   Port: 3050    │    │   (Ejecuta 1x)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   RDS PostgreSQL│
                    │   Port: 5432    │
                    └─────────────────┘
```

## 🔒 Configuración de Seguridad

### **Security Groups necesarios:**

1. **ECS Tasks Security Group:**
   - Inbound: 80 (Frontend), 3050 (Backend)
   - Outbound: 443 (HTTPS), 5432 (RDS)

2. **RDS Security Group:**
   - Inbound: 5432 desde ECS Tasks

3. **Load Balancer Security Group:**
   - Inbound: 80, 443 desde 0.0.0.0/0
   - Outbound: 80, 3050 hacia ECS Tasks

## 🌐 Configuración de Dominio

### **Con Application Load Balancer:**

1. **Target Groups:**
   - Frontend: Puerto 80
   - Backend: Puerto 3050

2. **Listeners:**
   - HTTP (80) → HTTPS redirect
   - HTTPS (443) → Target Groups

3. **DNS:**
   - `tudominio.com` → ALB
   - `api.tudominio.com` → ALB (Backend)

## 📝 Variables de Entorno Críticas

### **Base de datos (RDS):**
```bash
RDS_ENDPOINT=your-rds-endpoint.region.rds.amazonaws.com
DB_NAME=central_reserve
DB_USER=postgres
DB_PASSWORD=super_secure_password
DB_SSLMODE=require  # IMPORTANTE: Usar SSL en producción
```

### **Aplicación:**
```bash
APP_ENV=production
LOG_LEVEL=info
JWT_SECRET=very_long_and_secure_jwt_secret_for_production
```

### **Frontend:**
```bash
REACT_APP_API_BASE_URL=https://api.yourdomain.com
```

## 🔍 Monitoreo y Logs

### **CloudWatch Logs:**
```bash
# Ver logs en tiempo real
aws logs tail /ecs/reserve-app --follow

# Logs específicos por servicio
aws logs tail /ecs/reserve-app --follow --filter-pattern="frontend"
aws logs tail /ecs/reserve-app --follow --filter-pattern="backend"
```

### **Métricas ECS:**
- CPU y memoria por servicio
- Health checks
- Task count

## 🚨 Troubleshooting

### **Problema: Imágenes no se suben a ECR**
```bash
# Verificar permisos ECR
aws ecr get-authorization-token

# Recrear repositorios
aws ecr delete-repository --repository-name reserve-backend --force
aws ecr create-repository --repository-name reserve-backend
```

### **Problema: Servicios no inician**
```bash
# Ver logs de ECS
./scripts/deploy-to-aws.sh logs

# Verificar task definition
aws ecs describe-task-definition --task-definition your-task-def
```

### **Problema: Base de datos no conecta**
```bash
# Verificar connectivity desde ECS task
aws ecs execute-command --cluster central-reserve \
  --task your-task-id \
  --container central_reserve \
  --interactive \
  --command "/bin/sh"
```

## 🔄 Actualización de la Aplicación

### **Deploy de nueva versión:**
```bash
# 1. Cambiar tag en .env.aws
IMAGE_TAG=v1.1.0

# 2. Construir y desplegar
./scripts/deploy-to-aws.sh build
./scripts/deploy-to-aws.sh deploy
```

### **Rollback:**
```bash
# 1. Cambiar a tag anterior
IMAGE_TAG=v1.0.0

# 2. Redesplegar
./scripts/deploy-to-aws.sh deploy
```

## 💰 Consideraciones de Costos

### **Recursos que generan costo:**
- **ECS Tasks**: Según CPU/memoria asignada
- **RDS**: Instancia db.t3.micro (mínimo)
- **ECR**: Almacenamiento de imágenes
- **CloudWatch Logs**: Ingesta y almacenamiento
- **Data Transfer**: Tráfico saliente

### **Optimización:**
- Usar Fargate Spot para development
- Configurar log retention (7-30 días)
- Limpiar imágenes ECR antiguas

## 📈 Escalamiento

### **Auto Scaling:**
```bash
# Configurar en ECS Service
aws ecs put-scaling-policy \
  --service-name reserve-backend \
  --cluster central-reserve \
  --scalable-dimension ecs:service:DesiredCount \
  --target-tracking-configuration file://scaling-config.json
```

### **Load Balancer Health Checks:**
- Path: `/health` (Backend)
- Interval: 30s
- Timeout: 5s
- Healthy threshold: 2

## 🤝 Contribuir

Para agregar nueva funcionalidad al despliegue:

1. Modificar `docker-compose.aws.yaml`
2. Actualizar variables en `.env.aws`
3. Ajustar script `deploy-to-aws.sh`
4. Documentar cambios en este README

## 📞 Soporte

Si tienes problemas:

1. Verificar logs: `./scripts/deploy-to-aws.sh logs`
2. Revisar estado: `./scripts/deploy-to-aws.sh status`
3. Verificar configuración: `./scripts/deploy-to-aws.sh check` 