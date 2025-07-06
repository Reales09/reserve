# Scripts de Despliegue - Reserve App Frontend

Esta carpeta contiene scripts para desplegar la aplicación React a ECR público de AWS.

## 📁 Archivos

| Script | Descripción |
|--------|-------------|
| `deploy.sh` | Script principal de despliegue |
| `deploy-dev.sh` | Despliegue para desarrollo |
| `deploy-staging.sh` | Despliegue para staging |
| `deploy-prod.sh` | Despliegue para producción |
| `build-local.sh` | Construcción local sin subir a ECR |

## 🚀 Uso Rápido

### Desarrollo
```bash
./scripts/deploy-dev.sh
```

### Staging
```bash
./scripts/deploy-staging.sh
```

### Producción
```bash
./scripts/deploy-prod.sh
```

### Build Local
```bash
./scripts/build-local.sh
```

## 🔧 Configuración

### Variables de Entorno

Los scripts usan las siguientes variables:

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `API_BASE_URL` | URL base de la API | Varía según ambiente |
| `APP_NAME` | Nombre de la aplicación | "Reserve App" |
| `APP_VERSION` | Versión de la aplicación | "1.0.0" |

### Personalización

Puedes sobrescribir variables antes de ejecutar:

```bash
export API_BASE_URL="https://mi-api.com"
export APP_NAME="Mi App"
./scripts/deploy-prod.sh
```

## 📋 Prerrequisitos

1. **Docker** instalado y corriendo
2. **AWS CLI** configurado con permisos para ECR
3. **Node.js** y **npm** instalados
4. Estar en el directorio raíz del proyecto

### Configurar AWS CLI

```bash
# Configurar credenciales
aws configure

# Verificar configuración
aws sts get-caller-identity
```

## 🎯 Ambientes

### Desarrollo
- **API URL**: `https://dev-api.tudominio.com`
- **Tag**: `frontend-dev`
- **Repositorio**: `public.ecr.aws/d3a6d4r1/cam/reserve:frontend-dev`

### Staging
- **API URL**: `https://staging-api.tudominio.com`
- **Tag**: `frontend-staging`
- **Repositorio**: `public.ecr.aws/d3a6d4r1/cam/reserve:frontend-staging`

### Producción
- **API URL**: `https://api.tudominio.com`
- **Tag**: `frontend-latest`
- **Repositorio**: `public.ecr.aws/d3a6d4r1/cam/reserve:frontend-latest`

## 🔍 Ejemplo de Uso Completo

```bash
# 1. Clonar el repositorio
git clone <repo-url>
cd reserve_app

# 2. Instalar dependencias
npm install

# 3. Construir y probar localmente
./scripts/build-local.sh
docker run -p 80:80 reserve-app-frontend:local

# 4. Desplegar a desarrollo
./scripts/deploy-dev.sh

# 5. Desplegar a staging
./scripts/deploy-staging.sh

# 6. Desplegar a producción (con confirmación)
./scripts/deploy-prod.sh
```

## 🐳 Usando las Imágenes

### Desde ECR

```bash
# Desarrollo
docker run -p 80:80 \
  -e REACT_APP_API_BASE_URL=https://dev-api.tudominio.com \
  public.ecr.aws/d3a6d4r1/cam/reserve:frontend-dev

# Producción
docker run -p 80:80 \
  -e REACT_APP_API_BASE_URL=https://api.tudominio.com \
  public.ecr.aws/d3a6d4r1/cam/reserve:frontend-latest
```

### Con Docker Compose

```yaml
version: '3.8'
services:
  frontend:
    image: public.ecr.aws/d3a6d4r1/cam/reserve:frontend-latest
    ports:
      - "80:80"
    environment:
      - REACT_APP_API_BASE_URL=https://api.tudominio.com
      - REACT_APP_NAME=Reserve App
    restart: unless-stopped
```

## 🛠️ Troubleshooting

### Error de Permisos AWS
```bash
# Verificar credenciales
aws sts get-caller-identity

# Hacer login a ECR
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws
```

### Error de Docker
```bash
# Verificar que Docker esté corriendo
docker info

# Limpiar recursos si es necesario
docker system prune -a
```

### Error de Node.js
```bash
# Limpiar cache
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

## 🌐 URLs Útiles

- **ECR Repository**: https://gallery.ecr.aws/d3a6d4r1/cam/reserve
- **AWS Console**: https://console.aws.amazon.com/ecr/repositories
- **Docker Hub**: https://hub.docker.com (alternativa)

## 📝 Logs y Debugging

```bash
# Ver logs durante el despliegue
./scripts/deploy-dev.sh 2>&1 | tee deploy.log

# Verificar imagen construida
docker images | grep reserve-app-frontend

# Inspeccionar imagen
docker inspect reserve-app-frontend:latest
``` 