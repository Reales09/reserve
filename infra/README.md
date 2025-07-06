# 🚀 Despliegue Centralizado - Reserve App

Este directorio contiene la configuración centralizada para el despliegue de toda la aplicación Reserve App en producción.

## 📋 Estructura del Monorepo

```
reserve/
├── back/
│   ├── central-reserve/    # API principal (Go)
│   └── dbpostgres/        # Servicio de migración de BD (Go)
├── front/
│   └── reserve_app/       # Frontend React
└── infra/                 # ⭐ Configuración de despliegue
    ├── docker-compose.yaml
    ├── .env.example
    ├── init.sql
    ├── deploy.sh
    └── README.md
```

## 🛠️ Servicios Configurados

### 1. **Frontend** (`frontend`)
- **Imagen**: React + Nginx
- **Puerto**: 80
- **Función**: Interfaz de usuario
- **Configuración**: Variables de entorno dinámicas

### 2. **Backend API** (`central_reserve`)
- **Imagen**: Go aplicación principal
- **Puerto**: 3050
- **Función**: API REST principal
- **Dependencias**: PostgreSQL, Migración de BD

### 3. **Migración de BD** (`db_migrator`)
- **Imagen**: Go servicio de migración
- **Función**: Crear y actualizar esquema de BD
- **Comportamiento**: ⚠️ **Se ejecuta UNA sola vez y se apaga**
- **Dependencias**: PostgreSQL

### 4. **Base de Datos** (`postgres`)
- **Imagen**: PostgreSQL 15
- **Puerto**: 5432
- **Función**: Almacenamiento de datos
- **Volumen**: Datos persistentes

## 🔧 Configuración Inicial

### 1. Crear archivo de variables de entorno

```bash
# Copiar archivo de ejemplo
cp .env.example .env

# Editar con tus valores
nano .env
```

### 2. Variables de entorno principales

```env
# Base de datos
DB_NAME=central_reserve
DB_USER=postgres
DB_PASSWORD=tu_password_super_secreto

# Frontend
REACT_APP_API_BASE_URL=https://api.tudominio.com
DOMAIN=tudominio.com

# Aplicación
APP_ENV=production
```

## 🚀 Despliegue

### Opción 1: Script automatizado (Recomendado)

```bash
# Hacer el script ejecutable
chmod +x deploy.sh

# Ver ayuda
./deploy.sh help

# Construir y desplegar
./deploy.sh build
./deploy.sh up

# Ver logs
./deploy.sh logs

# Ver estado
./deploy.sh status
```

### Opción 2: Docker Compose manual

```bash
# Construir imágenes
docker-compose build

# Levantar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Ver estado
docker-compose ps
```

## 📊 Orden de Ejecución

El sistema está configurado para ejecutarse en el siguiente orden:

1. **PostgreSQL** se inicia primero
2. **Migrador de BD** se ejecuta y configura el esquema
3. **Backend API** se inicia y se conecta a la BD
4. **Frontend** se inicia y se conecta al backend

## 🔍 Monitoreo

### Ver logs de un servicio específico
```bash
docker-compose logs -f frontend
docker-compose logs -f central_reserve
docker-compose logs -f db_migrator
docker-compose logs -f postgres
```

### Verificar estado de servicios
```bash
docker-compose ps
```

### Verificar recursos del sistema
```bash
docker stats
```

## 🛡️ Características de Seguridad

- **Usuarios no root**: Todos los contenedores usan usuarios no privilegiados
- **Certificados SSL**: Incluidos en las imágenes
- **Logging estructurado**: Logs en formato JSON con rotación
- **Health checks**: Verificación automática de servicios

## 🔄 Gestión de Servicios

### Reiniciar un servicio
```bash
docker-compose restart central_reserve
```

### Reconstruir y actualizar
```bash
docker-compose build --no-cache
docker-compose up -d
```

### Detener todo
```bash
docker-compose down
```

### Limpiar recursos
```bash
docker system prune -f
docker volume prune -f
```

## 📁 Volúmenes de Datos

- **postgres_data**: Datos persistentes de PostgreSQL
- Los logs se almacenan en el contenedor con rotación automática

## 🌐 Configuración de Dominio

Para usar con tu dominio:

1. Configurar DNS apuntando a tu servidor
2. Actualizar `DOMAIN` en `.env`
3. Configurar SSL/TLS (Traefik labels incluidos)

## 🆘 Troubleshooting

### El migrador no se ejecuta
- Verificar que PostgreSQL esté corriendo
- Revisar logs: `docker-compose logs db_migrator`
- Verificar variables de entorno de BD

### Backend no se conecta a BD
- Verificar que la migración se completó
- Revisar logs: `docker-compose logs central_reserve`
- Verificar conectividad: `docker-compose exec central_reserve ping postgres`

### Frontend no se conecta al backend
- Verificar `REACT_APP_API_BASE_URL` en `.env`
- Revisar logs: `docker-compose logs frontend`
- Verificar que el backend esté corriendo en puerto 3050

## 🔗 URLs de Acceso

Una vez desplegado:

- **Frontend**: http://localhost (puerto 80)
- **API Backend**: http://localhost:3050
- **PostgreSQL**: localhost:5432 (solo para debugging)

## 📝 Notas Importantes

1. **Migración de BD**: El servicio `db_migrator` se ejecuta una sola vez y se apaga automáticamente
2. **Persistencia**: Los datos de PostgreSQL se mantienen en volúmenes Docker
3. **SSL**: Las etiquetas de Traefik están configuradas para SSL automático
4. **Logs**: Se rotan automáticamente para evitar llenar disco

## 🤝 Contribuir

Para modificar la configuración:

1. Editar `docker-compose.yaml` para cambios de servicios
2. Actualizar `.env.example` para nuevas variables
3. Modificar `deploy.sh` para nueva funcionalidad
4. Actualizar este README para cambios importantes 