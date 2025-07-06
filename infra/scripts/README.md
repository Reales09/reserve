# 🐳 Scripts de Construcción de Imágenes Docker

Este directorio contiene scripts especializados para construir las imágenes Docker de cada componente del proyecto Reserve App.

## 📁 Estructura de Scripts

```
scripts/
├── build-all.sh        # Script maestro - construye todas las imágenes
├── build-backend.sh    # Construye imagen del backend (Go)
├── build-frontend.sh   # Construye imagen del frontend (React)
├── build-migrator.sh   # Construye imagen del migrador de BD (Go)
└── README.md          # Este archivo
```

## 🚀 Uso de los Scripts

### 1. Script Maestro - Construir Todas las Imágenes

```bash
# Construir todas las imágenes secuencialmente
./scripts/build-all.sh

# Construir con tag específico
./scripts/build-all.sh v1.0.0

# Construir en paralelo (más rápido)
./scripts/build-all.sh latest true
```

### 2. Scripts Individuales

#### Backend (central-reserve)
```bash
# Construir imagen del backend
./scripts/build-backend.sh

# Con tag específico
./scripts/build-backend.sh v1.0.0
```

#### Frontend (reserve_app)
```bash
# Construir imagen del frontend
./scripts/build-frontend.sh

# Con tag específico
./scripts/build-frontend.sh v1.0.0
```

#### Migrador de BD (dbpostgres)
```bash
# Construir imagen del migrador
./scripts/build-migrator.sh

# Con tag específico
./scripts/build-migrator.sh v1.0.0
```

## 🛠️ Usando desde el Script Principal

También puedes usar estos scripts desde el script principal de despliegue:

```bash
# Construir todas las imágenes
./deploy.sh build-all

# Construir solo el backend
./deploy.sh build-backend

# Construir solo el frontend
./deploy.sh build-frontend

# Construir solo el migrador
./deploy.sh build-migrator

# Construir con tag específico
./deploy.sh build-all v1.0.0
```

## 📊 Comparación de Métodos

| Método | Comando | Ventajas | Desventajas |
|--------|---------|----------|-------------|
| **Docker Compose** | `docker compose build` | Rápido, integrado | Menos control individual |
| **Scripts Individuales** | `./scripts/build-*.sh` | Control granular | Más manual |
| **Script Maestro** | `./scripts/build-all.sh` | Automatizado, flexible | Más complejo |

## 🏗️ Imágenes Generadas

| Script | Imagen Generada | Tamaño Aprox. | Propósito |
|--------|-----------------|---------------|-----------|
| `build-backend.sh` | `reserve-backend:latest` | ~50MB | API Go principal |
| `build-frontend.sh` | `reserve-frontend:latest` | ~25MB | React + Nginx |
| `build-migrator.sh` | `reserve-migrator:latest` | ~40MB | Migración BD |

## 🔧 Configuración Avanzada

### Variables de Entorno para Frontend

Los scripts del frontend pueden usar estas variables:

```bash
# Configurar antes de construir
export REACT_APP_API_BASE_URL=https://api.midominio.com
export REACT_APP_NAME="Mi App"
export REACT_APP_VERSION="2.0.0"

# Construir con variables personalizadas
./scripts/build-frontend.sh
```

### Construcción Paralela

Para acelerar el proceso, puedes construir en paralelo:

```bash
# Construir todas las imágenes en paralelo
./scripts/build-all.sh latest true
```

## 🚨 Solución de Problemas

### Error: "Script no encontrado"
```bash
# Verificar que estás en el directorio correcto
pwd  # Debería mostrar .../reserve/infra

# Hacer scripts ejecutables
chmod +x scripts/*.sh
```

### Error: "Dockerfile no encontrado"
```bash
# Verificar estructura de proyectos
ls -la ../back/central-reserve/docker/
ls -la ../back/dbpostgres/docker/
ls -la ../front/reserve_app/
```

### Error: "Docker no encontrado"
```bash
# Verificar que Docker está instalado
docker --version
docker compose version
```

## 📝 Personalización

Para modificar los scripts:

1. **Cambiar nombres de imágenes**: Edita la variable `IMAGE_NAME` en cada script
2. **Modificar rutas**: Actualiza `PROJECT_PATH` si cambias estructura
3. **Agregar build args**: Modifica los parámetros de `docker build`

## 🔄 Integración con CI/CD

Estos scripts están diseñados para ser usados en pipelines:

```yaml
# Ejemplo para GitHub Actions
- name: Build Backend
  run: ./infra/scripts/build-backend.sh ${{ github.sha }}

- name: Build Frontend
  run: ./infra/scripts/build-frontend.sh ${{ github.sha }}

- name: Build Migrator
  run: ./infra/scripts/build-migrator.sh ${{ github.sha }}
```

## 🏷️ Gestión de Tags

### Estrategia de Tagging Recomendada

```bash
# Desarrollo
./scripts/build-all.sh dev

# Staging
./scripts/build-all.sh staging

# Producción
./scripts/build-all.sh $(git rev-parse --short HEAD)
./scripts/build-all.sh v1.0.0
```

### Tags Múltiples

```bash
# Construir con múltiples tags
./scripts/build-backend.sh latest
docker tag reserve-backend:latest reserve-backend:v1.0.0
docker tag reserve-backend:latest reserve-backend:stable
```

## 🎯 Mejores Prácticas

1. **Usa tags específicos** para producción
2. **Construye en paralelo** para desarrollo
3. **Verifica imágenes** antes de desplegar
4. **Limpia imágenes viejas** regularmente
5. **Documenta cambios** en cada build

## 🤝 Contribuir

Para agregar nuevos scripts:

1. Crea el script siguiendo el patrón existente
2. Hazlo ejecutable: `chmod +x scripts/nuevo-script.sh`
3. Agrégalo a `build-all.sh`
4. Actualiza la documentación 