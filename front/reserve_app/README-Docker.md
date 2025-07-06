# Docker Setup para Reserve App

Este proyecto incluye configuración completa de Docker para despliegue en producción.

## 🚀 Características

- **Multi-stage build** optimizado para React
- **Variables de entorno** configurables en tiempo de ejecución
- **Nginx** optimizado para SPA (Single Page Application)
- **Compresión gzip** habilitada
- **Cache headers** para archivos estáticos
- **Health checks** incluidos

## 📋 Variables de Entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `REACT_APP_API_BASE_URL` | URL base de la API | `http://localhost:3050` |
| `REACT_APP_NAME` | Nombre de la aplicación | `Reserve App` |
| `REACT_APP_VERSION` | Versión de la aplicación | `1.0.0` |

## 🛠️ Uso

### Desarrollo Local

```bash
# Construir la imagen
docker build -t reserve-app .

# Ejecutar con variables de entorno
docker run -p 80:80 \
  -e REACT_APP_API_BASE_URL=http://localhost:3050 \
  reserve-app
```

### Con Docker Compose

```bash
# Desarrollo
docker-compose up -d

# Producción
docker-compose -f docker-compose.prod.yml up -d
```

### Despliegue en Producción

```bash
# Construir imagen
docker build -t reserve-app:latest .

# Ejecutar con URL de producción
docker run -p 80:80 \
  -e REACT_APP_API_BASE_URL=https://api.tudominio.com \
  -e REACT_APP_NAME="Reserve App Pro" \
  -e REACT_APP_VERSION="1.0.0" \
  reserve-app:latest
```

## 🔧 Configuración Avanzada

### Variables de Entorno en Tiempo de Ejecución

Las variables de entorno se pueden cambiar sin reconstruir la imagen:

```bash
# Ejemplo con diferentes APIs
docker run -p 80:80 \
  -e REACT_APP_API_BASE_URL=https://staging-api.tudominio.com \
  reserve-app

docker run -p 80:80 \
  -e REACT_APP_API_BASE_URL=https://prod-api.tudominio.com \
  reserve-app
```

### Usar con Docker Swarm

```bash
# Deploy en swarm
docker stack deploy -c docker-compose.prod.yml reserve-stack
```

### Usar con Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: reserve-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: reserve-app
  template:
    metadata:
      labels:
        app: reserve-app
    spec:
      containers:
      - name: frontend
        image: reserve-app:latest
        ports:
        - containerPort: 80
        env:
        - name: REACT_APP_API_BASE_URL
          value: "https://api.tudominio.com"
```

## 🧪 Testing

```bash
# Verificar que la aplicación funciona
curl -f http://localhost:80/

# Verificar variables de entorno
docker logs [container-id]
```

## 🔍 Troubleshooting

### Variables de Entorno No Funcionan

1. Verifica que uses el prefijo `REACT_APP_`
2. Asegúrate de que no haya espacios en las variables
3. Revisa los logs del contenedor

### Problemas de CORS

Si tienes problemas de CORS, asegúrate de que la URL de la API sea correcta:

```bash
# Verificar configuración
docker exec [container-id] grep -r "REACT_APP_API_BASE_URL" /usr/share/nginx/html/
```

### Optimización

Para mejorar el rendimiento:

1. Usa un registry privado para imágenes
2. Implementa cache layers en el build
3. Usa CDN para archivos estáticos
4. Configura load balancer si es necesario

## 📁 Estructura de Archivos Docker

```
├── Dockerfile              # Configuración principal
├── docker-compose.yml      # Para desarrollo
├── docker-compose.prod.yml # Para producción  
├── nginx.conf              # Configuración de Nginx
├── docker-entrypoint.sh    # Script de inicio
└── .dockerignore           # Archivos ignorados
``` 