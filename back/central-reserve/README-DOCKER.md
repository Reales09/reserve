# 🐳 Docker - Central Reserve Backend

Guía completa para ejecutar el backend con Docker.

## 🚀 Inicio Rápido

### Opción 1: Script Automatizado (Recomendado)
```bash
# Desarrollo completo con todos los servicios
./scripts/build-docker.sh dev

# Solo build para producción
./scripts/build-docker.sh prod
```

### Opción 2: Makefile
```bash
# Ver todos los comandos disponibles
make help

# Entorno de desarrollo
make docker-dev

# Entorno de producción
make docker-prod

# Ver logs
make docker-logs
```

### Opción 3: Docker Compose Directo
```bash
# Desarrollo
cd docker
docker-compose -f docker-compose.dev.yml up -d

# Producción
docker-compose -f docker-compose.prod.yml up -d
```

## 📋 Servicios Incluidos

### 🏗️ Entorno de Desarrollo
- **central_reserve**: API Backend (puerto 3050)
- **postgres**: Base de datos PostgreSQL (puerto 5432)
- **redis**: Cache Redis (puerto 6379)
- **nats**: Mensajería NATS (puerto 4222)
- **nats_dashboard**: Dashboard NATS (puerto 8111)
- **adminer**: Gestor de base de datos (puerto 8080)

### 🚀 Entorno de Producción
- **central_reserve**: API Backend optimizado
- **postgres**: Base de datos PostgreSQL
- **redis**: Cache Redis
- **nats**: Mensajería NATS

## 🔧 Configuración

### Variables de Entorno
Crea un archivo `.env` en la raíz del proyecto:

```env
# Configuración de la aplicación
APP_ENV=development
HTTP_PORT=3050
LOG_LEVEL=debug
JWT_SECRET=tu-jwt-secret-aqui

# Base de datos
DB_HOST=postgres
DB_USER=postgres
DB_PASS=password
DB_PORT=5432
DB_NAME=central_reserve
DB_LOG_LEVEL=info
PGSSLMODE=disable

# Swagger
URL_BASE_SWAGGER=http://localhost:3050

# Email (opcional)
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USER=tu-email@outlook.com
SMTP_PASS=tu-contraseña
FROM_EMAIL=tu-email@outlook.com
SMTP_USE_STARTTLS=true
SMTP_USE_TLS=false
```

### Puertos Utilizados
- **3050**: API Backend
- **5432**: PostgreSQL
- **6379**: Redis
- **4222**: NATS
- **8111**: NATS Dashboard
- **8080**: Adminer

## 🛠️ Comandos Útiles

### Gestión de Contenedores
```bash
# Ver contenedores activos
docker ps

# Ver logs en tiempo real
make docker-logs

# Ver logs de todos los servicios
make docker-logs-all

# Detener todos los servicios
make docker-stop

# Reiniciar servicios
docker-compose -f docker/docker-compose.dev.yml restart
```

### Base de Datos
```bash
# Acceder a PostgreSQL
docker exec -it postgres_dev psql -U postgres -d central_reserve

# Resetear base de datos
make db-reset

# Ver logs de PostgreSQL
docker logs postgres_dev
```

### Desarrollo
```bash
# Rebuild de la imagen
make docker-build

# Ejecutar tests
make test

# Verificar salud de servicios
make health

# Probar envío de emails
make test-email
```

## 🔍 Troubleshooting

### Problemas Comunes

#### 1. Puerto ya en uso
```bash
# Ver qué está usando el puerto
sudo lsof -i :3050

# Matar proceso
sudo kill -9 <PID>
```

#### 2. Contenedor no inicia
```bash
# Ver logs detallados
docker logs central_reserve_dev

# Verificar variables de entorno
docker exec central_reserve_dev env | grep -E "(DB_|SMTP_)"
```

#### 3. Base de datos no conecta
```bash
# Verificar que PostgreSQL esté corriendo
docker ps | grep postgres

# Ver logs de PostgreSQL
docker logs postgres_dev

# Probar conexión
docker exec -it postgres_dev pg_isready -U postgres
```

#### 4. Permisos de archivos
```bash
# Dar permisos al script
chmod +x scripts/build-docker.sh

# Si hay problemas con volúmenes
sudo chown -R $USER:$USER .
```

### Limpieza
```bash
# Limpieza completa
make clean-all

# Solo contenedores
docker-compose -f docker/docker-compose.dev.yml down -v

# Solo imágenes
docker rmi central-reserve:latest
```

## 📊 Monitoreo

### Health Checks
```bash
# Verificar salud de la API
curl http://localhost:3050/health

# Verificar todos los servicios
make health
```

### Logs Estructurados
Los logs incluyen:
- Timestamp
- Nivel de log (INFO, ERROR, WARN)
- Contexto de la operación
- Métricas de rendimiento

### Métricas Disponibles
- Latencia de requests HTTP
- Estado de conexiones a BD
- Uso de memoria y CPU
- Errores y excepciones

## 🚀 Despliegue a Producción

### 1. Build de Producción
```bash
./scripts/build-docker.sh prod v1.0.0
```

### 2. Configurar Variables de Producción
```env
APP_ENV=production
LOG_LEVEL=info
# Configurar credenciales reales de BD y email
```

### 3. Desplegar
```bash
docker-compose -f docker/docker-compose.prod.yml up -d
```

### 4. Verificar
```bash
# Health check
curl http://tu-servidor:3050/health

# Logs
docker-compose -f docker/docker-compose.prod.yml logs -f
```

## 🔒 Seguridad

### Buenas Prácticas Implementadas
- ✅ Usuario no-root en contenedores
- ✅ Imagen minimalista (Alpine)
- ✅ Variables de entorno para secretos
- ✅ Health checks configurados
- ✅ Volúmenes persistentes para datos
- ✅ Red aislada para servicios

### Recomendaciones de Producción
- Usar secrets management (Docker Secrets, AWS Secrets Manager)
- Configurar backup automático de PostgreSQL
- Implementar rate limiting
- Configurar SSL/TLS para HTTPS
- Monitoreo con Prometheus/Grafana

## 📚 Recursos Adicionales

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [PostgreSQL Docker](https://hub.docker.com/_/postgres)
- [Redis Docker](https://hub.docker.com/_/redis)
- [NATS Docker](https://hub.docker.com/_/nats) 